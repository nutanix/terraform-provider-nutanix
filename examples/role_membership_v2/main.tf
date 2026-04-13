terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=1.0.0"
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
  role_ext_id         = var.role_ext_id
  identity_type       = "USER"
  identity_ext_id     = var.identity_ext_id
  idp_ext_id          = var.idp_ext_id
  scope_template_name = var.scope_template_name
  project_ext_id      = var.project_ext_id

  key_value_pairs {
    key   = "projectId"
    value = var.project_ext_id
  }

  scope_template_name_values {
    name  = "projectId"
    value = var.project_ext_id
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

output "role_membership_id" {
  value = nutanix_role_membership_v2.example.ext_id
}

output "role_memberships_count" {
  value = length(data.nutanix_role_memberships_v2.list.role_memberships)
}
