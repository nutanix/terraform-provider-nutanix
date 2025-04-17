terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.2.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  insecure = true
  port     = var.nutanix_port
}

// Example 1: Update VM Configuration

// Provision Application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = var.blueprint_name
    app_name        = var.app_name
    app_description = var.app_description
}

// Run patch config (update config)
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = var.patch_name
    config_name = var.config_name // same as patch name
}

// Example 2: Update VM Configuration with runtime editable

# Run patch config with runtime editable on above application
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = var.patch_name
    config_name = var.config_name // same as patch name
    vm_config {
        memory_size_mib = var.memory_size_mib
        num_sockets = var.num_sockets
        num_vcpus_per_socket = var.num_vcpus_per_socket
    }
}

## Example 3: Add Category

# Run patch config to add category in above application
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = var.patch_name
    config_name = var.config_name // same as patch name
    categories {
        value = var.category_value
        operation = var.add_operation
    }
}

## Example 4: Delete Category

# Run patch config to delete category in above application
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = var.patch_name
    config_name = var.config_name // same as patch name
    categories {
        value = var.category_value
        operation = var.delete_operation
    }
}

// Example 5: Add Disk

// To add disk without runtime editable. 
// Execute the patch action (similar to example 1) having disk configured

// Run patch config to add disk (runtime editable) in above application
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = var.patch_name
    config_name = var.config_name // same as patch name
    disks {
        disk_size_mib = var.disk_size_mib
        operation = var.add_operation
    }
}

// Example 6: Add Nic

// To add nic without runtime editable. 
// Execute the patch action (similar to example 1) having nic configured

// Run patch config to add nic (runtime editable) in above application
resource "nutanix_self_service_app_patch" "test" {
    app_uuid = nutanix_self_service_app_provision.test.id
    patch_name = var.patch_name
    config_name = var.config_name // same as patch name
    nics {
        index = var.index
        operation = var.add_operation
        subnet_uuid = var.subnet_uuid
    }
}



