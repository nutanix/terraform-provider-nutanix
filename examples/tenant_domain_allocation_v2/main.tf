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

# Prerequisite tenant that owns the domain allocation.
resource "nutanix_tenant_v2" "example" {
  name        = var.tenant_name
  description = "Tenant that owns the domain allocation"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = var.tenant_admin_password
  }
}

# Create a domain allocation for the tenant with a project + placement targets.
resource "nutanix_tenant_domain_allocation_v2" "example" {
  tenant_ext_id = nutanix_tenant_v2.example.ext_id
  domain_ext_id = var.domain_ext_id

  projects {
    name = "example-project"
    quota {
      vcpu    = 4
      memory  = 8589934592
      storage = 53687091200
    }
    resource_group {
      placement_targets {
        cluster_uuid            = var.cluster_uuid
        storage_container_uuids = [var.storage_container_uuid]
      }
    }
  }
}

# Fetch a single domain allocation by its external identifier.
data "nutanix_tenant_domain_allocation_v2" "get-domain-allocation" {
  tenant_ext_id = nutanix_tenant_v2.example.ext_id
  ext_id        = nutanix_tenant_domain_allocation_v2.example.ext_id
}

# List all domain allocations for the tenant.
data "nutanix_tenant_domain_allocations_v2" "domain-allocations-list" {
  tenant_ext_id = nutanix_tenant_v2.example.ext_id
  depends_on    = [nutanix_tenant_domain_allocation_v2.example]
}
