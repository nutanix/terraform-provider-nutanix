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

# Create a tenant to demonstrate the data source
resource "nutanix_tenant_v2" "example_tenant" {
  name        = var.tenant_name
  description = "Tenant for demonstrating domain allocations data source"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = var.tenant_admin_password
  }
}

# Lists all domain allocations for a tenant.
data "nutanix_tenant_domain_allocations_v2" "example" {
  tenant_ext_id = nutanix_tenant_v2.example_tenant.ext_id
}

# Output the list of domain allocations for verification
output "domain_allocations" {
  value = data.nutanix_tenant_domain_allocations_v2.example.domain_allocations
}