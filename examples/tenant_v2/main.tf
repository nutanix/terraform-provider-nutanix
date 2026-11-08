terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
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

# Creates a new tenant
resource "nutanix_tenant_v2" "example" {
  name        = var.tenant_name
  description = var.tenant_description

  admin_info {
    username   = var.admin_username
    first_name = var.admin_first_name
    last_name  = var.admin_last_name
    email      = var.admin_email
    password   = var.admin_password
  }

  # Optional: Add aggregated quota for the tenant
  aggregated_quota {
    memory = var.quota_memory
    storage = var.quota_storage
    vcpu   = var.quota_vcpu
  }
}

# Output the tenant details for verification
output "tenant_ext_id" {
  value = nutanix_tenant_v2.example.ext_id
}

output "tenant_name" {
  value = nutanix_tenant_v2.example.name
}

output "tenant_created_time" {
  value = nutanix_tenant_v2.example.created_time
}