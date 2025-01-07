terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.7.0"
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

# Add a User group to the system.

resource "nutanix_user_groups_v2" "example" {
  # Type of the User Group. LDAP, SAML
  group_type         = "<group Type>"
  idp_id             = "<idp uuid of user group>"
  name               = "<group name>"
  distinguished_name = "<distinguished name of the user group>"
}


# List all the user groups in the system.
data "nutanix_user_groups_v2" "example" {}

# List user groups with a filter.
data "nutanix_user_groups_v2" "example" {
  filter = "name eq '<group name>'"
}

# Get the details of a user group.
data "nutanix_user_group_v2" "example" {
  ext_id = nutanix_user_groups_v2.example.id
}
