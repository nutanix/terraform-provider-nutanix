terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
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

# Create a tenant with a local admin and aggregated quota.
resource "nutanix_tenant_v2" "example" {
  name        = var.tenant_name
  description = "Tenant created via Terraform"
  url         = "https://example.com/${var.tenant_name}"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = var.tenant_admin_password
  }

  aggregated_quota {
    vcpu    = 8
    memory  = 17179869184
    storage = 107374182400
  }
}

# Fetch a single tenant by its external identifier.
data "nutanix_tenant_v2" "get-tenant" {
  ext_id = nutanix_tenant_v2.example.ext_id
}

# List all tenants.
data "nutanix_tenants_v2" "tenants-list" {
  depends_on = [nutanix_tenant_v2.example]
}
