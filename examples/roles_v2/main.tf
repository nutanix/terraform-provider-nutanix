terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}


# filtered list operation
data "nutanix_operations_v2" "operations-filtered-list" {
  filter = "startswith(displayName, 'Create_')"
}

# Create role
resource "nutanix_roles_v2" "example-role" {
  display_name = "example_role"
  description  = "create example role"
  operations = [
    data.nutanix_operations_v2.operations-filtered-list.operations[0].ext_id,
    data.nutanix_operations_v2.operations-filtered-list.operations[1].ext_id,
    data.nutanix_operations_v2.operations-filtered-list.operations[2].ext_id,
    data.nutanix_operations_v2.operations-filtered-list.operations[3].ext_id
  ]
}

# List all Roles
data "nutanix_roles_v2" "roles" {}

# List Roles with filter
data "nutanix_roles_v2" "filtered-roles" {
  filter = "displayName eq '${nutanix_roles_v2.example-role.display_name}'"
}

# List Roles with filter and orderby
data "nutanix_roles_v2" "ordered-roles" {
  order_by = "createdTime desc"
}

# List Roles with filter and orderby
data "nutanix_roles_v2" "filtered-ordered-roles"{
  filter = "displayName eq '${nutanix_roles_v2.example-role.display_name}'"
  order_by = "createdTime desc"
}
