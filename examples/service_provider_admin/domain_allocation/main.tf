terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

# Service Provider Admin — allocate a domain (capacity + networking) to a tenant.
provider "nutanix" {
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port                     = 443
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

  # Configure the allocation's external network connection directly. Only one
  # of nat or no_nat may be configured.
  external_network_config {
    config {
      nat {
        transit_vpc_ext_id    = var.transit_vpc_ext_id
        overlay_subnet_ext_id = var.overlay_subnet_ext_id
        floating_ip_count     = 2
      }
    }
  }

  projects {
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
data "nutanix_tenant_domain_allocation_v2" "get_domain_allocation" {
  tenant_ext_id = nutanix_tenant_v2.example.ext_id
  ext_id        = nutanix_tenant_domain_allocation_v2.example.ext_id
}

# List all domain allocations for the tenant.
data "nutanix_tenant_domain_allocations_v2" "domain_allocations_list" {
  tenant_ext_id = nutanix_tenant_v2.example.ext_id
  depends_on    = [nutanix_tenant_domain_allocation_v2.example]
}

output "domain_allocation_ext_id" {
  value = nutanix_tenant_domain_allocation_v2.example.ext_id
}

output "domain_allocations" {
  value = data.nutanix_tenant_domain_allocations_v2.domain_allocations_list.domain_allocations
}
