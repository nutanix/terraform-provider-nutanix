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
                resource_details.append(details)
            elif 'method' in resource:
                # Other resource (non-CRUD)
                resource_details.append(f"Resource: {resource['name']} (method: {resource['method']})")
    
    resource_section = '\n'.join(resource_details) if resource_details else 'N/A'
    
    prompt = f"""Generate complete Terraform provider code for {namespace} namespace using {sdk_info_file}.
Generate:
1. SDK client: nutanix/sdks/v4/{namespace}/{namespace}.go (pattern: networking.go), if client is already present please update the client with the new api methods or recivers.
2. Update nutanix/config.go: add {namespace_camel}API client init + field
3. Datasources in nutanix/services/{namespace}v2/: {datasource_list}: 
IMPORTANT for DataSources:
1. If Get*ById methods are present, please build the datasource using the response_struct from the sdk_info.json file for the api method, it's a module for fetching a single resource by id.
2. If List* methods are present, please build the datasource using the response_struct from the sdk_info.json file for the api method, it's a module for fetching a list of resources.
4. Resources in nutanix/services/{namespace}v2/:
{resource_section}
IMPORTANT for Resources: 
- If a receiver has Create, Update, Delete methods, they should all be in the SAME resource file with Create, Read, Update, Delete contexts. Build the schema for the resource using the request_struct and response_struct from the sdk_info.json file for the api method.
- Create context should use the Create API method
- Read context should use the GetById API method (if available)
- Update context should use the Update API method
- Delete context should use the Delete API method
- Other methods (not Get/List/Create/Update/Delete) should be considered as separate resources
5. Tests: Build the test files for the datasource and resource files by analyzing the api method and the response_struct or Request struct from the sdk_info.json file.
6. Examples: examples/{namespace}_v2/
7. Docs: website/docs/d/ and website/docs/r/
8. Register in nutanix/provider/provider.go

Reference: datapoliciesv2, networkingv2 just for file structures. Use sdk_info.json for exact field mappings, import paths and request/response structs.

IMPORTANT - OneOf Type Handling:
When flattening OneOf types (e.g., OneOfMetricDetailMetricValue, OneOfParameterParamValue), the OneOfTypeX fields are PRIVATE and cannot be accessed directly. Use this pattern:

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

DO NOT try to access OneOfType0, OneOfType1, etc. directly - they are private fields. Always use GetValue() and check ObjectType_ field.
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
    print(f"   Found {len(api_list)} API methods")
    print(f"   Found {len(datasources)} datasource(s): {', '.join([ds['name'] for ds in datasources]) if datasources else 'None'}")
    print(f"   Found {len(resources)} resource(s): {', '.join([r['name'] for r in resources]) if resources else 'None'}")
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

