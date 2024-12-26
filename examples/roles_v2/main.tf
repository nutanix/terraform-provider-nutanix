terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
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


# Create role
resource "nutanix_roles_v2" "test" {
  display_name = "test_role"
  description  = "creat a test role using terraform"
  operations   = var.operations
}

# list Roles
data "nutanix_roles_v2" "test" {}

# list Roles with filter
data "nutanix_roles_v2" "test" {
  filter = "displayName eq 'test_role'"
}
# get a specific role by id
data "nutanix_role_v2" "test" {
  ext_id = nutanix_roles_v2.test.id
}

