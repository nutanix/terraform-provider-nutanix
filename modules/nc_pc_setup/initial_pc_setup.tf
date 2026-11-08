# ==============================================================================
# STAGE 1: Prism Central Setup
# ==============================================================================
#
# This stage prepares the Prism Central instance for NC registration:
#   1. Discover PC, cluster, and storage container information
#   2. Configure Active Directory integration
#   3. Create LDAP user with Prism Admin privileges
#   4. Create Transit VPCs and Overlay Subnets for tenant networking
#
# ==============================================================================

# ------------------------------------------------------------------------------
# DATA SOURCES - Discover PC Infrastructure
# ------------------------------------------------------------------------------

data "nutanix_pcs_v2" "pc" {
  provider = nutanix.pc
}

data "nutanix_clusters_v2" "clusters" {
  provider = nutanix.pc
}

# data "nutanix_storage_containers_v2" "storage_container" {
#   provider = nutanix.pc
# }

data "nutanix_roles_v2" "pc" {
  provider = nutanix.pc
}

# ------------------------------------------------------------------------------
# LOCAL VALUES - Extract Infrastructure IDs
# ------------------------------------------------------------------------------

locals {
  # Find PC by name
  pc_ext_id = one([
    for pc in data.nutanix_pcs_v2.pc.pcs :
    pc.ext_id if lower(pc.config[0].name) == lower(var.pc_name)
  ])

  # Get first non-PC cluster
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]

  # Find Prism Admin role
  pc_admin_role_ext_id = one([
    for role in data.nutanix_roles_v2.pc.roles :
    role.ext_id if role.display_name == var.pc_admin_role_name
  ])
}

# Storage Container Creation
resource "nutanix_storage_containers_v2" "example" {
  provider = nutanix.pc
  name                                     = "cpaas_storage_container"
  cluster_ext_id                           = local.cluster_ext_id
  logical_advertised_capacity_bytes        = 1073741824000
  logical_explicit_reserved_capacity_bytes = 32
  replication_factor                       = 1
  nfs_whitelist_addresses {
    ipv4 {
      value         = "192.168.15.0"
      prefix_length = 32
    }
  }
  erasure_code                          = "OFF"
  is_inline_ec_enabled                  = false
  has_higher_ec_fault_domain_preference = false
  cache_deduplication                   = "OFF"
  on_disk_dedup                         = "OFF"
  is_compression_enabled                = true
  is_internal                           = false
  is_software_encryption_enabled        = false
}

# Storage Container Creation
resource "nutanix_storage_containers_v2" "example_2" {
  provider = nutanix.pc
  name                                     = "database_storage_container"
  cluster_ext_id                           = local.cluster_ext_id
  logical_advertised_capacity_bytes        = 1073741824000
  logical_explicit_reserved_capacity_bytes = 32
  replication_factor                       = 1
  nfs_whitelist_addresses {
    ipv4 {
      value         = "192.168.15.0"
      prefix_length = 32
    }
  }
  erasure_code                          = "OFF"
  is_inline_ec_enabled                  = false
  has_higher_ec_fault_domain_preference = false
  cache_deduplication                   = "OFF"
  on_disk_dedup                         = "OFF"
  is_compression_enabled                = true
  is_internal                           = false
  is_software_encryption_enabled        = false
}

data "nutanix_storage_containers_v2" "get_storage_containers" {
  provider = nutanix.pc
  filter   = "name eq 'SelfServiceContainer'"
}

locals {
  cpaas_storage_container_ext_id = nutanix_storage_containers_v2.example.container_ext_id
  database_storage_container_ext_id = nutanix_storage_containers_v2.example_2.container_ext_id
  storage_container_self_service_ext_id = data.nutanix_storage_containers_v2.get_storage_containers.storage_containers[0].container_ext_id
}
# ------------------------------------------------------------------------------
# ACTIVE DIRECTORY INTEGRATION
# ------------------------------------------------------------------------------

resource "nutanix_directory_services_v2" "pc_active_directory" {
  provider       = nutanix.pc
  name           = var.pc_directory_service.name
  url            = var.pc_directory_service.url
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = var.pc_directory_service.domain_name

  service_account {
    username = var.pc_directory_service.service_account_username
    password = var.pc_directory_service.service_account_password
  }

  lifecycle {
    ignore_changes = [service_account[0].password]
  }
}

# ------------------------------------------------------------------------------
# LDAP USER & AUTHORIZATION POLICY
# ------------------------------------------------------------------------------

resource "nutanix_users_v2" "pc_ldap_user" {
  provider  = nutanix.pc
  username  = var.pc_ldap_username
  user_type = "LDAP"
  idp_id    = nutanix_directory_services_v2.pc_active_directory.id
}

resource "nutanix_authorization_policy_v2" "pc_admin" {
  provider                  = nutanix.pc
  role                      = local.pc_admin_role_ext_id
  display_name              = "Terraform Prism Admin Policy"
  description               = "Grants Prism Admin access to the configured LDAP user."
  authorization_policy_type = "USER_DEFINED"

  identities {
    reserved = jsonencode({
      user = { uuid = { anyof = [nutanix_users_v2.pc_ldap_user.ext_id] } }
    })
  }

  entities {
    reserved = jsonencode({ "*" = { "*" = { eq = "*" } } })
  }
}

# ------------------------------------------------------------------------------
# TRANSIT VPC - NAT ENABLED
# ------------------------------------------------------------------------------

resource "nutanix_vpc_v2" "nat_transit" {
  provider    = nutanix.pc
  name        = "nat-tenant-transit-vpc-akhil"
  description = "Transit VPC for tenant external connectivity (NAT enabled)"
  vpc_type    = "TRANSIT"
}

resource "nutanix_subnet_v2" "nat_overlay_external" {
  provider       = nutanix.pc
  name           = "tenant-nat-overlay-external-akhil"
  description    = "NAT overlay external network for tenant floating IPs"
  vpc_reference  = nutanix_vpc_v2.nat_transit.id
  subnet_type    = "OVERLAY"
  is_external    = true
  is_nat_enabled = true

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "10.57.183.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "10.57.183.1"
      }
      pool_list {
        start_ip { value = "10.57.183.2" }
        end_ip { value = "10.57.183.100" }
      }
    }
  }
}

# ------------------------------------------------------------------------------
# TRANSIT VPC - NO-NAT (Direct Routing)
# ------------------------------------------------------------------------------

resource "nutanix_vpc_v2" "no_nat_transit" {
  provider    = nutanix.pc
  name        = "no-nat-tenant-transit-vpc-akhil"
  description = "Transit VPC for tenant external connectivity (No-NAT / Direct Routing)"
  vpc_type    = "TRANSIT"

  # Advertise workload subnet CIDR without NAT
  externally_routable_prefixes {
    ipv4 {
      ip {
        value         = "10.50.0.0"
        prefix_length = 32
      }
      prefix_length = 24
    }
  }
}

resource "nutanix_subnet_v2" "no_nat_overlay_external" {
  provider       = nutanix.pc
  name           = "no-nat-tenant-overlay-external-akhil"
  description    = "No-NAT overlay external network for direct routing"
  vpc_reference  = nutanix_vpc_v2.no_nat_transit.id
  subnet_type    = "OVERLAY"
  is_external    = true
  is_nat_enabled = false

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "10.128.42.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "10.128.42.1"
      }
      pool_list {
        start_ip { value = "10.128.42.2" }
        end_ip { value = "10.128.42.100" }
      }
    }
  }
}

# ------------------------------------------------------------------------------
# LOCAL VALUES - Network Infrastructure IDs
# ------------------------------------------------------------------------------

locals {
  nat_transit_vpc_ext_id       = nutanix_vpc_v2.nat_transit.ext_id
  nat_overlay_subnet_ext_id    = nutanix_subnet_v2.nat_overlay_external.ext_id
  no_nat_transit_vpc_ext_id    = nutanix_vpc_v2.no_nat_transit.ext_id
  no_nat_overlay_subnet_ext_id = nutanix_subnet_v2.no_nat_overlay_external.ext_id
}
