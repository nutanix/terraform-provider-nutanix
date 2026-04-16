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

# Get scope template by id
data "nutanix_scope_template_v2" "get-by-id" {
  ext_id = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
}

# List all scope templates
data "nutanix_scope_templates_v2" "list-all" {}

# List scope templates with filter
data "nutanix_scope_templates_v2" "filtered" {
  filter = "displayName eq '${nutanix_scope_template_v2.example.display_name}'"
}
