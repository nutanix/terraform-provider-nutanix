terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

# Create a role membership for a user
resource "nutanix_role_membership_v2" "example" {
  role_ext_id         = "ca386756-e45f-5555-8625-5b68ae17393b" # ProjectAdmin
  identity_type       = "USER"
  identity_ext_id     = "8a49f561-6bd7-5d26-b53e-661d63e7bdb8" # User UUID
  idp_ext_id          = "d711f713-cdf5-5ee9-9936-4e67373eb842" # Identity Provider UUID
  scope_template_name = "ProjectsScopeTemplate"
  project_ext_id      = "78198f9c-063d-590e-9e9a-939b51829a39" # Project UUID

  scope_template_name_values {
    name  = "projectExtId"
    value = "78198f9c-063d-590e-9e9a-939b51829a39" # Project UUID
  }
}


# Create a role membership for a user group
resource "nutanix_role_membership_v2" "example" {
  role_ext_id         = "ca386756-e45f-5555-8625-5b68ae17393b" # ProjectAdmin
  identity_type       = "GROUP"
  identity_ext_id     = "8a49f561-6bd7-5d26-b53e-661d63e7bdb8" # UserGroup UUID
  idp_ext_id          = "d711f713-cdf5-5ee9-9936-4e67373eb842" # Identity Provider UUID
  scope_template_name = "ProjectsScopeTemplate"
  project_ext_id      = "78198f9c-063d-590e-9e9a-939b51829a39" # Project UUID

  scope_template_name_values {
    name  = "projectExtId"
    value = "78198f9c-063d-590e-9e9a-939b51829a39" # Project UUID
  }
}

# Data source to fetch a single role membership by ID
data "nutanix_role_membership_v2" "by_id" {
  ext_id = nutanix_role_membership_v2.example.ext_id
}

# Data source to list all role memberships
data "nutanix_role_memberships_v2" "list" {}

# Data source to list role membership summaries
data "nutanix_role_membership_summary_v2" "summaries" {}
