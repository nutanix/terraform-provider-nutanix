---
layout: "nutanix"
page_title: "Migration Guide: V1 to V2 Resources"
sidebar_current: "docs-nutanix-migration-v1-to-v2"
description: |-
  This guide provides step-by-step instructions for migrating from Nutanix Terraform Provider v1 resources to v2 resources, which are built on v4 APIs/SDKs.
---

# Migration Guide: V1 to V2 Resources

This guide walks you through the process of migrating your existing Terraform-managed infrastructure from Nutanix Provider v1 resources to v2 resources. The v2 resources are built on Nutanix v4 APIs/SDKs and provide enhanced capabilities and improved performance.

~> **Important Notice:** Nutanix strongly encourages migrating to v2 resources built on v4 APIs/SDKs.

## Overview

The migration process involves:
1. Adding v2 resource configurations alongside existing v1 resources.
2. Importing existing infrastructure into v2 resource state using entity UUID.
3. Removing v1 resource configurations.
4. Managing resources using v2 modules going forward.

This process allows you to migrate without recreating or disrupting your existing infrastructure.

This guide covers migration for:
- **Resources:** Individual Terraform resources
- **Modules:** Terraform modules that use v1 resources internally

## Resource Migration

### Step 1: Identify Your V1 Resources

First, identify the v1 resources you want to migrate. For example, if you have a virtual machine created using the v1 module:

```hcl
resource "nutanix_virtual_machine" "dev_vm" {
  name = "dev_vm"
  # ... other configuration
}
```

### Step 2: Add V2 Resource Configuration

Add a corresponding v2 resource configuration to your Terraform file:

```hcl
resource "nutanix_virtual_machine_v2" "import_dev_vm" {
  # Configuration will be imported from existing resource
  # You can add minimal required attributes if needed
}
```

~> **Note:** At this stage, you only need to declare the resource block.

### Step 3: Get the UUID of the Existing Resource

Retrieve the UUID of the resource created with the v1 module using Terraform state:

```bash
terraform state show nutanix_virtual_machine.dev_vm
```

Look for the `id` or `uuid` field in the output.

### Step 4: Import the Resource into V2 State

Execute the terraform import command to import the existing resource into the v2 resource state:

```bash
terraform import nutanix_virtual_machine_v2.import_dev_vm <UUID_OF_ENTITY>
```

~> **Note:** Replace `UUID_OF_ENTITY` with the actual UUID obtained from Step 3.

### Step 5: Verify Import is Successful

Verify that the resource has been imported successfully by listing all resources in the Terraform state:

```bash
terraform state list
```
You should now see both the v1 and v2 resources in the state:
```
nutanix_virtual_machine.dev_vm
nutanix_virtual_machine_v2.import_dev_vm
```
This confirms that the import was successful and both resources are now tracked in your Terraform state.


### Step 6: Remove the V1 Resource Config and Remove it from State

Once you've verified the import is successful, remove the v1 resource configuration and state:

1. **Remove the v1 resource from your Terraform configuration file:**
   
   Delete or comment out the v1 resource block:
   
   ```hcl
   # resource "nutanix_virtual_machine" "dev_vm" {
   #   name = "dev_vm"
   #   # ... other configuration
   # }
   ```

2. **Remove the v1 resource from Terraform state:**
   
   ```bash
   terraform state rm nutanix_virtual_machine.dev_vm
   ```

~> **Important:** Ensure the v2 resource import was successful (verified in Step 5) before removing the v1 resource from state. Once removed, the v1 resource will no longer be managed by Terraform.

### Step 7: Use V2 Resources from Now On

Congratulations! Your migration is complete. From now on, you can manage the imported resource using the v2 module:

```hcl
resource "nutanix_virtual_machine_v2" "import_dev_vm" {
  name = "dev_vm"
  # All v2 resource attributes can now be managed here
  # ... other v2 configuration
}
```

You can now:
- Update the resource configuration using v2 attributes
- Apply changes using `terraform apply`
- Manage the resource lifecycle with v2 resources
- Benefit from the enhanced capabilities of v4 APIs/SDKs


## Module Migration

If you're using Terraform modules that reference v1 resources internally, you'll need to migrate those as well. The process is similar to resource migration, but involves identifying and migrating all resources managed by the module.

### Step 1: Identify Modules Using V1 Resources

First, identify modules in your configuration that use v1 resources. For example, in your main configuration file:

``hcl
module "vm_module" {
  source = "./modules/vm-module"
  
  vm_name = "dev_vm"
  # ... other module inputs
}
```
Then, check the module source code to see if it uses v1 resources. Inside the module directory (e.g., `./modules/vm-module/main.tf`), you might see:

```hcl
# modules/vm-module/main.tf

variable "vm_name" {
  description = "Name of the virtual machine"
  type        = string
}

resource "nutanix_virtual_machine" "vm" {
  name = var.vm_name
  # ... other v1 resource configuration
}
```
If the module source code uses v1 resources like `nutanix_virtual_machine`, it needs migration.

### Step 2: List Resources Created by the Module

Identify all resources created by the module using Terraform state:

```bash
terraform state list | grep module.vm_module
```
This will show all resources managed by the module
```
module.vm_module.nutanix_virtual_machine.vm
# ... other resources created by the module.
```
### Step 3: Add V2 Resource Configuration in Module Source
Update your module source code to include v2 resource configurations. Add v2 resources alongside v1 resources in the module source.

Edit your module source file (e.g., `modules/vm-module/main.tf`) to add v2 resources:

```hcl
# modules/vm-module/main.tf

variable "vm_name" {
  description = "Name of the virtual machine"
  type        = string
}

# Existing v1 resources (will be removed after import)
resource "nutanix_virtual_machine" "vm" {
  name = var.vm_name
  # ... other v1 resource configuration
}

# Add v2 resource configuration for import
resource "nutanix_virtual_machine_v2" "vm_v2" {
  # Configuration will be imported from existing resource
  # You can add minimal required attributes if needed
}
```
~> **Note:** At this stage, you only need to declare the v2 resource blocks. The actual configuration will be imported in the further step.

### Step 4: Get the UUID of the Existing Resource
For resource in the module, retrieve its UUID using Terraform state:

```bash
terraform state show module.vm_module.nutanix_virtual_machine.vm
```
Look for the `id` or `uuid` field in the output for resource.

### Step 5: Import Module Resources into V2 State
Execute the terraform import command to import each resource managed by the module into v2 resource state. Use the UUIDs obtained from Step 4:
```bash
terraform import module.vm_module.nutanix_virtual_machine_v2.vm_v2 <UUID_OF_ENTITY>
```

### Step 6: Verify Import is Successful
Verify that all module resources have been imported successfully by listing all resources in the Terraform state:

```bash
terraform state list
```

You should now see both the v1 module resources and v2 resources in the state:
```
module.vm_module.nutanix_virtual_machine.vm
module.vm_module.nutanix_virtual_machine_v2.vm_v2
```

### Step 7: Remove the V1 Resource Config and Remove it from State
Once you've verified that all imports are successful, remove the v1 resource configurations and state:

1. **Remove the v1 resources from your module source code:**

   remove or comment out the v1 resource blocks in your module file (e.g., `modules/vm-module/main.tf`):
   
   ```hcl
   # modules/vm-module/main.tf
   
   # Existing v1 resources (removed after migration)
   # resource "nutanix_virtual_machine" "vm" {
   #   name = var.vm_name
   #   # ... other v1 resource configuration
   # }
   ```
2. **Remove v1 module resources from Terraform state:**

   Remove v1 resource managed by the module from Terraform state:
   
   ```bash
   terraform state rm module.vm_module.nutanix_virtual_machine.vm
   ```
~> **Important:** Ensure all v2 resource imports were successful (verified in Step 6) before removing the v1 module resources from state. Once removed, the v1 resources will no longer be managed by Terraform.

### Step 8: Continue to Use Modules Which Internally Use V2 Resources
Congratulations! Your module migration is complete. From now on, you can manage the resources using modules that internally use v2 resources.

Your module source now uses v2 resources (e.g., `modules/vm-module/main.tf`):
```hcl
# modules/vm-module/main.tf

variable "vm_name" {
  description = "Name of the virtual machine"
  type        = string
}

# Using v2 resources
resource "nutanix_virtual_machine_v2" "vm_v2" {
  name = var.vm_name
  # All v2 resource attributes can now be managed here
  # ... other v2 configuration
}
```

You can now:
- Continue using your module with v2 resources going forward
- Update the resource configuration using v2 attributes
- Apply changes using `terraform apply`
- Manage the resource lifecycle with v2 resources
- Benefit from the enhanced capabilities of v4 APIs/SDKs

## Important Considerations

- **No Downtime:** The migration process does not recreate resources, so there's no service interruption
- **Gradual Migration:** You can migrate resources one at a time - you don't need to migrate everything at once
- **Dependencies:** Pay attention to resource dependencies and migrate them in the correct order

### Configuration Mismatches
After import, if `terraform plan` shows changes:
- Review the differences carefully
- Some attribute changes are expected due to API differences
- Update your configuration to match the actual resource state


## Troubleshooting

### Import Fails
If the import fails, verify:
- The UUID is correct
- The resource still exists in Nutanix
- You have proper permissions to access the resource