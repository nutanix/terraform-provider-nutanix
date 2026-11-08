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

# Create a couple of tenants to demonstrate the data source
resource "nutanix_tenant_v2" "example_tenant_1" {
  name        = var.tenant_name_1
  description = "First example tenant"

  admin_info {
    username   = "tenant-admin-1"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin-1@example.com"
    password   = var.tenant_admin_password
  }
}

resource "nutanix_tenant_v2" "example_tenant_2" {
  name        = var.tenant_name_2
  description = "Second example tenant"

  admin_info {
    username   = "tenant-admin-2"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin-2@example.com"
    password   = var.tenant_admin_password
  }
}

# Lists all tenants.
data "nutanix_tenants_v2" "example" {
  # Optional: You can add filters, pagination, etc.
  # filter = "name eq 'example*'"
  # page   = 1
  # limit  = 10
}

# Output the list of tenants for verification
output "tenants" {
  value = data.nutanix_tenants_v2.example.tenants
}