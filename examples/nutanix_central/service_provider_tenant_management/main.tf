# ==============================================================================
# Service Provider Admin: Tenant Management
# ==============================================================================
#
# What this file does:
#   Shows the Service Provider Admin workflow to onboard a tenant: allow a
#   domain for tenancy, create the tenant, and allocate that domain (capacity
#   + optional networking) to the tenant.
#
# Flow:
#   1. Authenticate as Service Provider Admin via nutanix_central_endpoint
#      (no tenant_name / domain_name).
#   2. Shortlist (allow) an existing registered domain for tenancy
#      (nutanix_add_domain_v2) and list allowed domains.
#   3. Create a tenant with a local Tenant Admin user (nutanix_tenant_v2).
#   4. Allocate the domain to the tenant with project quota, placement targets,
#      and NAT / no-NAT external network config
#      (nutanix_tenant_domain_allocation_v2).
#   5. Read tenant and domain-allocation data sources for verification.
#
# Prerequisites:
#   1. Nutanix Central with Service Provider Admin credentials
#   2. A registered PC domain (domain_ext_id) available to shortlist
#   3. Cluster / storage container / transit VPC / overlay subnet IDs for the
#      allocation (and networking blocks if used)
#
# ==============================================================================

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.1"
    }
  }
}

# Service Provider Admin — NC endpoint only (no tenant_name / domain_name).
provider "nutanix" {
  username                 = var.nutanix_username
  password                 = var.nutanix_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port                     = 443
}

// Shortlist the domain. Which can be later assigned to tenant
resource "nutanix_add_domain_v2" "add_domain" {
  domain_ext_id = var.domain_ext_id
}

# Validate the domains shortlisted in SP.
data "nutanix_domains_v2" "registered_domain" {
  depends_on = [nutanix_add_domain_v2.add_domain]
}

# Prerequisite tenant that owns the domain allocation.
resource "nutanix_tenant_v2" "tenant" {
  name        = var.tenant_name
  description = "Tenant that owns the domain allocation"
  admin_info {
    username   = var.tenant_admin_username
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = var.tenant_admin_password
  }
  lifecycle {
    ignore_changes = [admin_info]
  }
}

# Allocate the domain to the tenant.
resource "nutanix_tenant_domain_allocation_v2" "tenant_domain_allocation" {
  tenant_ext_id = nutanix_tenant_v2.tenant.ext_id
  domain_ext_id = var.domain_ext_id
  projects {
    quota {
      vcpu    = 50
      memory  = 75  # value in GB
      storage = 200 # value in GB
    }
    resource_group {
      placement_targets {
        cluster_uuid            = var.cluster_ext_id
        storage_container_uuids = [var.cpaas_storage_container_ext_id, var.database_storage_container_ext_id, var.self_service_storage_container_ext_id]
      }
    }
  }
  # Configure the allocation's external network connection directly. Only one
  # of nat or no_nat may be configured.
  external_network_config {
    config {
      nat {
        transit_vpc_ext_id    = var.nat_transit_vpc_ext_id
        overlay_subnet_ext_id = var.nat_overlay_subnet_ext_id
        floating_ip_count     = 25
      }
      no_nat {
        transit_vpc_ext_id           = var.no_nat_transit_vpc_ext_id
        overlay_subnet_ext_id        = var.no_nat_overlay_subnet_ext_id
        externally_routable_prefixes = var.externally_routable_prefixes
      }
    }
  }
  depends_on = [nutanix_tenant_v2.tenant,nutanix_add_domain_v2.add_domain]
}

data "nutanix_tenant_domain_allocation_v2" "get_domain_allocation" {
  ext_id        = nutanix_tenant_domain_allocation_v2.tenant_domain_allocation.id
  tenant_ext_id = nutanix_tenant_v2.tenant.ext_id
}

data "nutanix_tenant_domain_allocations_v2" "get_all_domain_allocations" {
  tenant_ext_id = nutanix_tenant_v2.tenant.ext_id
  depends_on    = [nutanix_tenant_domain_allocation_v2.tenant_domain_allocation]
}

data "nutanix_tenants_v2" "get_all_tenants" {
  depends_on = [nutanix_tenant_v2.tenant]
}

data "nutanix_tenant_v2" "get_tenant"{
  ext_id = nutanix_tenant_v2.tenant.ext_id
}