terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

# Service Provider Admin — create a tenant under Nutanix Central.
provider "nutanix" {
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port                     = 443
}

# Creates a new tenant with a local admin user and optional aggregated quota.
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

  aggregated_quota {
    memory  = var.quota_memory
    storage = var.quota_storage
    vcpu    = var.quota_vcpu
  }
}

# Read the created tenant by ext_id.
data "nutanix_tenant_v2" "example" {
  ext_id = nutanix_tenant_v2.example.ext_id
}

# List tenants (optionally filter by name).
data "nutanix_tenants_v2" "example" {
  depends_on = [nutanix_tenant_v2.example]
}

output "tenant_ext_id" {
  value = nutanix_tenant_v2.example.ext_id
}

output "tenant_name" {
  value = nutanix_tenant_v2.example.name
}

output "tenants" {
  value = data.nutanix_tenants_v2.example.tenants
}
