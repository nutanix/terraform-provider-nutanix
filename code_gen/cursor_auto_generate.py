#!/usr/bin/env python3
"""
Automated Cursor IDE Command Execution
Sends commands to Cursor Composer without manual intervention
"""

import sys
import subprocess
import time
import os
import json
import re
from pathlib import Path

def load_sdk_info(sdk_info_path):
    """Load and parse sdk_info.json file"""
    try:
        with open(sdk_info_path, 'r') as f:
            return json.load(f)
    except FileNotFoundError:
        print(f"❌ Error: sdk_info.json not found at {sdk_info_path}")
        return None
    except json.JSONDecodeError as e:
        print(f"❌ Error: Invalid JSON in sdk_info.json: {e}")
        return None
    except Exception as e:
        print(f"❌ Error loading sdk_info.json: {e}")
        return None

def extract_namespace_from_package(package_path):
    """Extract namespace from package path
    e.g., 'monitoring-go-client/v4' -> 'monitoring'
    """
    match = re.search(r'/([a-z0-9-]+)-go-client/', package_path)
    if match:
        return match.group(1)
    return "unknown"

def camel_to_snake(name):
    """Convert camelCase or PascalCase to snake_case.
    e.g., 'identityExtId' -> 'identity_ext_id', 'scopeTemplateName' -> 'scope_template_name'
    """
    # Insert underscore before uppercase letters that follow lowercase letters or digits
    s1 = re.sub(r'([a-z0-9])([A-Z])', r'\1_\2', name)
    # Insert underscore before uppercase letters followed by lowercase (handles acronyms like 'ExtID')
    s2 = re.sub(r'([A-Z]+)([A-Z][a-z])', r'\1_\2', s1)
    return s2.lower()


def extract_required_fields_from_body(api_entry):
    """Extract required field info from the Body field's description in a Create request_struct.

    The Body description follows the pattern:
      "(required) ... It requires the <field1>, <field2>, ..., and <fieldN> attributes."

    The description names are informal shorthand that may NOT match struct field names exactly
    (e.g. "role" for struct field "RoleExtId", "name" for "DisplayName").

    We resolve them against actual Body struct fields using multiple strategies:
      1. Exact match on camelCase field name (e.g. "identityExtId" -> "IdentityExtId")
      2. Exact match on JSON tag name (e.g. "identityExtId" -> json:"identityExtId")
      3. Prefix match on JSON tag (e.g. "role" matches "roleExtId" -> "RoleExtId")
      4. Fallback: keep the raw description name (let the AI agent resolve it against the struct)

    Returns a dict with:
      - "resolved": list of snake_case field names successfully mapped to struct fields
      - "unresolved": list of raw description names that could NOT be mapped
      - "body_description": the original Body description text (authoritative source)
    Returns None if no Body field or no required-fields pattern is found.
    """
    request_struct = api_entry.get('request_struct', {})
    fields = request_struct.get('fields', [])

    for field in fields:
        if field.get('name') != 'Body':
            continue

        description = field.get('description', '')
        match = re.search(r'[Ii]t requires the (.+?) attributes?\.', description)
        if not match:
            return None

        raw = match.group(1)
        raw = re.sub(r',?\s+and\s+', ', ', raw)
        desc_names = [name.strip() for name in raw.split(',') if name.strip()]

        body_fields = field.get('fields', [])

        # Build lookup tables from the actual Body struct fields
        camel_to_field = {}
        json_tag_to_field = {}

        for bf in body_fields:
            bf_name = bf.get('name', '')
            if not bf_name:
                continue
            camel = bf_name[0].lower() + bf_name[1:]
            camel_to_field[camel] = bf_name

            tag = bf.get('tag', '')
            tag_match = re.search(r'json:"([^",]+)', tag)
            if tag_match:
                json_tag_to_field[tag_match.group(1)] = bf_name

        resolved = []
        unresolved = []
        for desc_name in desc_names:
            field_name = _resolve_field_name(desc_name, camel_to_field, json_tag_to_field)
            if field_name != desc_name:
                resolved.append(camel_to_snake(field_name))
            else:
                unresolved.append(desc_name)

        return {
            "resolved": resolved,
            "unresolved": unresolved,
            "body_description": description,
        }

    return None


def _resolve_field_name(desc_name, camel_to_field, json_tag_to_field):
    """Resolve an informal description name to the actual struct field name.

    Strategies in order:
      1. Exact camelCase match  (e.g. "identityExtId" -> "IdentityExtId")
      2. Exact JSON tag match   (e.g. "identityExtId" -> json:"identityExtId")
      3. Prefix match on JSON tags — pick the shortest tag that starts with desc_name
         (e.g. "role" matches "roleExtId" -> "RoleExtId")
      4. Fallback: return desc_name as-is (caller treats as unresolved)
    """
    if desc_name in camel_to_field:
        return camel_to_field[desc_name]

    if desc_name in json_tag_to_field:
        return json_tag_to_field[desc_name]

    candidates = []
    for tag_name, field_name in json_tag_to_field.items():
        if tag_name.lower().startswith(desc_name.lower()) and tag_name != desc_name:
            candidates.append((len(tag_name), tag_name, field_name))
    if candidates:
        candidates.sort()
        return candidates[0][2]

    return desc_name


def get_datasource_methods(api_list):
    """Extract datasource methods (Get*ById) and List* methods"""
    datasources = []
    for api in api_list:
        method_name = api.get('api_method', {}).get('name', '')
        if method_name.startswith('List'):
            datasources.append({
                'name': method_name,
                'resource': method_name[4:],
                'receiver': api.get('api_method', {}).get('receiver', '')
            })
        if method_name.startswith('Get') and method_name.endswith('ById'):
            # Extract resource name (e.g., GetAlertById -> Alert)
            resource_name = method_name[3:-2]  # Remove 'Get' and 'ById'
            datasources.append({
                'name': method_name,
                'resource': resource_name,
                'receiver': api.get('api_method', {}).get('receiver', '')
            })
    return datasources

def get_resource_methods(api_list):
    """Extract resource methods grouped by receiver.
    
    Rules:
    1. If a receiver has Create, Update, Delete methods, they form a single CRUD resource
    2. The resource will have: create (Create API), read (GetById API), update (Update API), delete (Delete API)
    3. Other methods (not Get/List/Create/Update/Delete) are also considered resources
    """
    # Group methods by receiver
    receiver_methods = {}
    for api in api_list:
        receiver = api.get('api_method', {}).get('receiver', '')
        method_name = api.get('api_method', {}).get('name', '')
        
        if receiver not in receiver_methods:
            receiver_methods[receiver] = []
        receiver_methods[receiver].append({
            'name': method_name,
            'api': api
        })
    
    resources = []
    processed_receivers = set()
    
    # Process each receiver
    for receiver, methods in receiver_methods.items():
        # Check for CRUD operations
        create_method = None
        update_method = None
        delete_method = None
        read_method = None  # GetById for Read context
        other_methods = []
        
        for method in methods:
            method_name = method['name']
            
            if method_name.startswith('Create'):
                create_method = method
            elif method_name.startswith('Update'):
                update_method = method
            elif method_name.startswith('Delete'):
                delete_method = method
            elif method_name.startswith('Get') and method_name.endswith('ById'):
                read_method = method
            elif not (method_name.startswith('Get') or method_name.startswith('List')):
                # Other methods (not Get/List) are also resources
                other_methods.append(method)
        
        # If receiver has Create, Update, or Delete, create a CRUD resource
        if create_method or update_method or delete_method:
            # Extract resource name from Create/Update/Delete method
            resource_name = None
            if create_method:
                resource_name = create_method['name'][6:]  # Remove 'Create'
            elif update_method:
                resource_name = update_method['name'][6:]  # Remove 'Update'
            elif delete_method:
                resource_name = delete_method['name'][6:]  # Remove 'Delete'
            
            if resource_name:
                resource = {
                    'name': resource_name,
                    'receiver': receiver,
                    'create': create_method['name'] if create_method else None,
                    'read': read_method['name'] if read_method else None,
                    'update': update_method['name'] if update_method else None,
                    'delete': delete_method['name'] if delete_method else None,
                }
                resources.append(resource)
                processed_receivers.add(receiver)
        
        # Process other methods as separate resources
        # Note: Other methods are always separate resources, even if receiver has CRUD operations
        for method in other_methods:
            # Create a resource for this method
            # Extract a meaningful resource name from the method name
            method_name = method['name']
            resource_name = method_name  # Use method name as resource name
            
            resource = {
                'name': resource_name,
                'receiver': receiver,
                'method': method_name,
            }
            resources.append(resource)
    
    return resources

def find_cursor_app():
    """Find Cursor.app on macOS"""
    possible_paths = [
        "/Applications/Cursor.app",
        os.path.expanduser("~/Applications/Cursor.app"),
    ]
    
    for path in possible_paths:
        if os.path.exists(path):
            return path
    return None

def send_command_via_cli(command):
    """Send command via Cursor CLI"""
    try:
        # Try using cursor CLI directly
        result = subprocess.run(
            ["cursor", "--command", command],
            capture_output=True,
            text=True,
            timeout=30
        )
        return result.returncode == 0, result.stdout, result.stderr
    except (subprocess.TimeoutExpired, FileNotFoundError):
        return False, "", "CLI not available"

def send_command_via_applescript(command):
    """Send command to Cursor using AppleScript (macOS)"""
    # First, copy command to clipboard
    try:
        subprocess.run(
            ["pbcopy"],
            input=command.encode(),
            check=True,
            timeout=5
        )
    except Exception as e:
        return False, "", f"Failed to copy to clipboard: {e}"
    
    # Use AppleScript to automate Cursor
    applescript = '''
    tell application "Cursor"
        activate
    end tell
    
    delay 0.5
    
    tell application "System Events"
        tell process "Cursor"
            -- Open Composer with Cmd+I
            key code 34 using {command down}  -- Cmd+I
            delay 1.5
            
            -- Clear any existing text
            keystroke "a" using {command down}
            delay 0.2
            
            -- Paste from clipboard
            keystroke "v" using {command down}
            delay 0.5
            
            -- Press Enter to execute
            key code 36  -- Enter key
        end tell
    end tell
    '''
    
    try:
        result = subprocess.run(
            ["osascript", "-e", applescript],
            capture_output=True,
            text=True,
            timeout=60
        )
        if result.returncode == 0:
            return True, "Command sent via AppleScript (copied to clipboard and pasted)", ""
        else:
            error_msg = result.stderr or result.stdout or "AppleScript execution failed"
            if "not allowed assistive access" in error_msg.lower():
                return False, "", "Accessibility permissions required. Grant Terminal/Python access in System Preferences > Security & Privacy > Accessibility"
            return False, "", error_msg
    except subprocess.TimeoutExpired:
        return False, "", "AppleScript timeout"
    except Exception as e:
        return False, "", str(e)

def send_command_via_debugger(command, port=9222):
    """Send command via Cursor's debugger protocol"""
    import socket
    
    try:
        # Connect to debugger port
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(5)
        sock.connect(("localhost", port))
        
        # Send command (simplified - actual protocol is more complex)
        message = json.dumps({"method": "Runtime.evaluate", "params": {"expression": command}})
        sock.sendall(message.encode())
        
        sock.close()
        return True, "Command sent", ""
    except Exception as e:
        return False, "", str(e)

def send_command_via_file(command, workspace_path):
    """Send command by writing to a file and using clipboard + AppleScript"""
    cursor_commands_dir = Path(workspace_path) / ".cursor" / "commands"
    cursor_commands_dir.mkdir(parents=True, exist_ok=True)
    
    # Save command to a file for reference
    command_file = cursor_commands_dir / f"auto_generate_{int(time.time())}.md"
    command_file.write_text(f"# Auto-Generated Command\n\n{command}\n")
    
    # Copy command to clipboard
    try:
        if sys.platform == "darwin":  # macOS
            subprocess.run(
                ["pbcopy"],
                input=command.encode(),
                check=True
            )
            clipboard_ready = True
        else:
            clipboard_ready = False
    except:
        clipboard_ready = False
    
    return True, f"Command saved to {command_file}. Use AppleScript method to send it automatically.", ""

def generate_terraform_command_from_sdk(sdk_info, workspace_path):
    """Generate a simple, clear Terraform generation command from sdk_info.json"""
    if not sdk_info:
        return "Generate Terraform provider code. Please check sdk_info.json file."
    
    package_path = sdk_info.get('package_path', '')
    namespace = extract_namespace_from_package(package_path)
    api_list = sdk_info.get('api_request_response_struct', [])

    # Extract datasources (Get*ById methods)
    datasources = get_datasource_methods(api_list)
    
    # Extract resources (Create/Update/Delete methods)
    resources = get_resource_methods(api_list)
    
    # Build a simple, actionable prompt
    namespace_camel = namespace.capitalize()
    datasource_list = ', '.join([ds['name'] for ds in datasources]) if datasources else 'N/A'
    sdk_info_file = Path(workspace_path) / 'code_gen' / 'sdk_extract_output' / 'sdk_info.json'
    
    # Build a lookup from method name to api_entry for required-field extraction
    method_to_api = {}
    for api in api_list:
        method_name = api.get('api_method', {}).get('name', '')
        if method_name:
            method_to_api[method_name] = api

    # Build resource details
    resource_details = []
    if resources:
        for resource in resources:
            if 'create' in resource or 'update' in resource or 'delete' in resource:
                # CRUD resource
                details = f"Resource: {resource['name']} (CRUD)"
                if resource.get('create'):
                    details += f"\n  - Create context: {resource['create']}"
                if resource.get('read'):
                    details += f"\n  - Read context: {resource['read']}"
                if resource.get('update'):
                    details += f"\n  - Update context: {resource['update']}"
                if resource.get('delete'):
                    details += f"\n  - Delete context: {resource['delete']}"

                # Extract required fields from the Create method's Body description
                create_method_name = resource.get('create')
                if create_method_name and create_method_name in method_to_api:
                    req_info = extract_required_fields_from_body(method_to_api[create_method_name])
                    if req_info:
                        if req_info["resolved"]:
                            details += f"\n  - Required schema fields (mark as Required: true): {', '.join(req_info['resolved'])}"
                        if req_info["unresolved"]:
                            details += f"\n  - Unresolved required attrs from description (match to nearest Body struct field): {', '.join(req_info['unresolved'])}"
                        details += f"\n  - Body description (authoritative): \"{req_info['body_description']}\""

                resource_details.append(details)
            elif 'method' in resource:
                # Other resource (non-CRUD)
                resource_details.append(f"Resource: {resource['name']} (method: {resource['method']})")
    
    resource_section = '\n'.join(resource_details) if resource_details else 'N/A'

    # Build API method summary from sdk_info.json description fields
    api_method_summary_lines = []
    for api in api_list:
        method = api.get('api_method', {})
        method_name = method.get('name', '')
        description = method.get('description', '')
        uri = method.get('uri', '')
        receiver = method.get('receiver', '')
        if method_name:
            parts = [f"  - {method_name} (receiver: {receiver})"]
            if description:
                parts.append(f"    Description: {description}")
            if uri:
                parts.append(f"    URI: {uri}")
            api_method_summary_lines.append('\n'.join(parts))
    api_method_summary = '\n'.join(api_method_summary_lines) if api_method_summary_lines else '  Refer to sdk_info.json'

    # Build datasource details with descriptions
    datasource_details = []
    for api in api_list:
        method = api.get('api_method', {})
        method_name = method.get('name', '')
        description = method.get('description', '')
        if method_name.startswith('Get') and method_name.endswith('ById'):
            datasource_details.append(f"  - {method_name} (singular datasource): {description}")
        elif method_name.startswith('List'):
            datasource_details.append(f"  - {method_name} (plural datasource): {description}")
    datasource_section = '\n'.join(datasource_details) if datasource_details else '  N/A'

    prompt = f"""# Terraform Provider Code Generation — {namespace} namespace

## Source of Truth
- SDK info file: {sdk_info_file}
- Reference existing namespaces (datapoliciesv2, networkingv2) ONLY for file/folder structure patterns.
- Use sdk_info.json for ALL field mappings, import paths, types, and request/response structs, descriptions, and URIs.

## API Methods Available
{api_method_summary}

---

## Step 1 — SDK Client
- File: nutanix/sdks/v4/{namespace}/{namespace}.go (follow the pattern in networking.go).
- If the client file already exists, update it with any new API methods or receivers.

## Step 2 — Provider Config
- Update nutanix/config.go: add a {namespace_camel}API client field and initialize it.
- if {namespace_camel}API client field is already present, update it with any new API methods or receivers.

## Step 3 — Datasources
- Directory: nutanix/services/{namespace}v2/
- Methods to implement:
{datasource_section}
- Rules:
  - Get*ById methods → singular datasource (fetches one resource by ID). Build schema from response_struct in sdk_info.json.
  - List* methods → plural datasource (fetches a list of resources). Build schema from response_struct in sdk_info.json.

## Step 4 — Resources
- Directory: nutanix/services/{namespace}v2/
{resource_section}
- Rules:
  - If a receiver has Create + Update + Delete methods, combine them into ONE resource file with Create, Read, Update, Delete contexts.
    - Create context → Create API method
    - Read context   → GetById API method (if available)
    - Update context → Update API method (if available)
    - Delete context → Delete API method (if available)
  - Build the schema from request_struct (for inputs) and response_struct (for computed outputs) in sdk_info.json.
  - **Required fields**: The Create method's request_struct has a Body field whose description lists the required attributes
    (e.g., "It requires the role, identityExtId, ... attributes."). These description names are **informal shorthand** that
    may not match struct field names exactly (e.g., "role" refers to the struct field "RoleExtId" with json tag "roleExtId",
    "name" may refer to "DisplayName" or a similarly named field). When "Required schema fields" are listed above, they have
    been resolved to actual struct field names. When "Unresolved required attrs" are listed, you MUST manually match them to
    the closest Body struct field by inspecting the json tags and field names in sdk_info.json. The Body description is also
    provided verbatim as the authoritative source. Mark resolved + matched fields as `Required: true` in the resource schema.
    All other Body fields that are NOT listed as required should be `Optional: true` (and `Computed: true` if they also appear in the response).
    Fields that only appear in the response_struct (not in Body) should be `Computed: true`.
  - Methods that are not Get/List/Create/Update/Delete → implement as separate action resources.

## Step 5 — Tests
- Build test files for every datasource and resource generated.
- Derive test assertions from the response_struct and request_struct fields in sdk_info.json.
- Dry run the tests to ensure each and every attribute is covered during resource tests, each and every attribute should be validated, present in state file.

## Step 6 — Examples
- Directory: examples/{namespace}_v2/
- Provide working .tf examples for each datasource and resource.

## Step 7 — Documentation
- Datasource docs: website/docs/d/
- Resource docs: website/docs/r/
- Rules: All descriptions MUST come from sdk_info.json.
  - API-level description (page subtitle, resource/datasource summary) → use the "description" field from api_method in sdk_info.json.
  - Attribute-level description (each argument/attribute row) → use the "description" field from each field inside request_struct and response_struct in sdk_info.json.

## Step 8 — Registration
- Register all new datasources and resources in nutanix/provider/provider.go.

## Step 9 - Review the generated code
- Review Resource and DataSource files, navigate through schema and validate against sdk_info.json.
- Review test files, validate against sdk_info.json. Validate against new resources and datasources.
- Review examples files, validate against sdk_info.json. Validate against new resources and datasources.
- Review documentation files, validate against sdk_info.json. Validate against new resources and datasources.
- During reviews, if you found any issues, please fix them and re-run the review process until it matched with sdk_info.json.
---

## OneOf Type Handling (CRITICAL)
OneOfTypeX fields (e.g., oneOfType0, oneOfType1) are PRIVATE and cannot be accessed directly.
Always use GetValue() and switch on ObjectType_:

```go
func flattenOneOfValue(oneOfValue *import1.OneOfSomeValue) []map[string]interface{{}} {{
    if oneOfValue != nil && oneOfValue.ObjectType_ != nil {{
        valueMap := make(map[string]interface{{}})
        value := oneOfValue.GetValue()
        if value != nil {{
            switch *oneOfValue.ObjectType_ {{
            case "monitoring.v4.common.StringValue":
                if strVal, ok := value.(import1.StringValue); ok && strVal.StringValue != nil {{
                    valueMap["string_value"] = utils.StringValue(strVal.StringValue)
                }}
            case "monitoring.v4.common.BoolValue":
                if boolVal, ok := value.(import1.BoolValue); ok && boolVal.BoolValue != nil {{
                    valueMap["bool_value"] = utils.BoolValue(boolVal.BoolValue)
                }}
            case "monitoring.v4.common.IntValue":
                if intVal, ok := value.(import1.IntValue); ok && intVal.IntValue != nil {{
                    valueMap["int_value"] = utils.Int64Value(intVal.IntValue)
                }}
            case "monitoring.v4.common.DoubleValue":
                if doubleVal, ok := value.(import1.DoubleValue); ok && doubleVal.DoubleValue != nil {{
                    valueMap["double_value"] = utils.Float64Value(doubleVal.DoubleValue)
                }}
            }}
        }}
        return []map[string]interface{{}}{{valueMap}}
    }}
    return nil
}}
```

Do NOT access oneOfType0, oneOfType1, etc. directly — they are unexported. Always use GetValue() + ObjectType_ switch.
"""

    return prompt

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 cursor_auto_generate.py <sdk_info.json_path> [workspace_path]")
        print("\nExample:")
        print('  python3 cursor_auto_generate.py code_gen/sdk_extract_output/sdk_info.json')
        print('  python3 cursor_auto_generate.py code_gen/sdk_extract_output/sdk_info.json /path/to/workspace')
        sys.exit(1)
    
    sdk_info_path = sys.argv[1]
    workspace_path = sys.argv[2] if len(sys.argv) > 2 else os.getcwd()
    
    # Resolve paths
    script_dir = Path(__file__).parent
    if not os.path.isabs(sdk_info_path):
        sdk_info_path = str(script_dir / sdk_info_path)
    
    print("🚀 Automated Terraform Code Generation")
    print("=" * 80)
    print(f"📋 SDK Info: {sdk_info_path}")
    print(f"📁 Workspace: {workspace_path}")
    print()
    
    # Load SDK info
    print("📖 Loading sdk_info.json...")
    sdk_info = load_sdk_info(sdk_info_path)
    if not sdk_info:
        print("❌ Failed to load sdk_info.json")
        sys.exit(1)
    
    namespace = extract_namespace_from_package(sdk_info.get('package_path', ''))
    api_list = sdk_info.get('api_request_response_struct', [])
    datasources = get_datasource_methods(api_list)
    resources = get_resource_methods(api_list)

    print(f"✅ Loaded SDK info for namespace: {namespace}")
    print(f"   API version: {sdk_info.get('api_version', '') or 'Not detected'}")
    print(f"   Internal SDK: {sdk_info.get('is_internal', False)}")
    print(f"   Found {len(api_list)} API methods")
    print(f"   Found {len(datasources)} datasource(s): {', '.join([ds['name'] for ds in datasources]) if datasources else 'None'}")
    print(f"   Found {len(resources)} resource(s): {', '.join([r['name'] for r in resources]) if resources else 'None'}")
    
    desc_count = sum(1 for api in api_list if api.get('api_method', {}).get('description'))
    uri_count = sum(1 for api in api_list if api.get('api_method', {}).get('uri'))
    print(f"   📄 Descriptions: {desc_count}/{len(api_list)} methods have descriptions")
    print(f"   🔗 URIs: {uri_count}/{len(api_list)} methods have URIs")
    print()
    
    # Generate the command
    print("🔍 Generating Terraform generation prompt...")
    command = generate_terraform_command_from_sdk(sdk_info, workspace_path)
    print("✅ Prompt generated!")
    print()
    print("📝 Generated prompt:")
    print("-" * 80)
    print(command)
    print("-" * 80)
    print()
    
    # Try different methods to send the command
    # AppleScript is most reliable on macOS for UI automation
    methods = [
        ("AppleScript (UI Automation)", lambda: send_command_via_applescript(command)),
        ("CLI", lambda: send_command_via_cli(command)),
        ("File-based (Fallback)", lambda: send_command_via_file(command, workspace_path)),
    ]
    
    success = False
    for method_name, method_func in methods:
        print(f"🤖 Trying {method_name} method...")
        try:
            ok, stdout, stderr = method_func()
            if ok:
                print(f"✅ {method_name} method succeeded!")
                if stdout:
                    print(f"   Output: {stdout[:100]}")
                success = True
                break
            else:
                print(f"❌ {method_name} method failed: {stderr}")
        except Exception as e:
            print(f"❌ {method_name} method error: {e}")
        print()
    
    if not success:
        print("⚠️  All automated methods failed.")
        print()
        print("📝 Alternative: The command has been saved to:")
        print(f"   {Path(workspace_path) / '.cursor' / 'commands' / 'auto_generate_manual.md'}")
        print()
        print("   You can:")
        print("   1. Open Cursor Composer (Cmd+I)")
        print("   2. Type: /auto_generate_manual")
        print("   3. Or manually paste this command:")
        print()
        print(command)
        print()
        sys.exit(1)
    
    print("=" * 80)
    print("✨ Command sent successfully!")
    print("⏳ Cursor should now be processing the command...")
    print("=" * 80)

if __name__ == "__main__":
    main()

