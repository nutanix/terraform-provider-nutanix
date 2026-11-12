# ==============================================================================
# NC Admin: Domain Infrastructure Management
# ==============================================================================
#
# This configuration demonstrates how an NC Admin can manage infrastructure
# within a registered domain:
#   • Create projects for organizing resources
#   • Create Category
#
# PREREQUISITES:
#   1. NC Admin credentials with domain management permissions
# ==============================================================================


# ==============================================================================
# TERRAFORM CONFIGURATION
# ==============================================================================

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.1"
    }
  }
}

# ==============================================================================
# PROVIDER CONFIGURATION FOR NC ADMIN
# ==============================================================================

provider "nutanix" {
  username                 = var.nutanix_username
  password                 = var.nutanix_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  domain_name              = var.domain_name
  port                     = 443
}

# ==============================================================================
# PROJECT
# ==============================================================================

resource "nutanix_project_v2" "project" {
  name        = "cpass catgeory"
  description = "Project created by NC Admin for domain infrastructure"
}

resource "nutanix_category_v2" "category" {
  description = "Category created by NC Admin for domain infrastructure"
  key = "cpass"
  value = "cpass"
  project_ext_id = nutanix_project_v2.project.ext_id
}