# ==============================================================================
# NC Admin: Prism Central Registration with Nutanix Central
# ==============================================================================
#
# What this file does:
#   Registers a Prism Central (PC) instance with Nutanix Central (NC) so the PC
#   becomes a managed domain for NC Admin / Service Provider / Tenant workflows.
#
# Flow:
#   1. Configure two providers — PC (direct endpoint) and NC (nutanix_central_endpoint).
#   2. Discover the Prism Central and cluster ext_ids from the PC APIs.
#   3. Look up an NC location to attach to the registered domain.
#   4. Register the PC domain in Nutanix Central (nutanix_registered_domain_v2).
#   5. Complete the handshake on PC by creating the local domain
#      (nutanix_local_domain_v2) with the API credentials issued by NC.
#
# Prerequisites:
#   - NC Admin credentials and PC Admin credentials
#   - Directory service / IAM setup on PC and NC (authorization policies) must
#     already exist; this registration depends on those policies
#
# ==============================================================================

# ==============================================================================
# PROVIDERS
# ==============================================================================

# Prism Central - Admin access for PC configuration
provider "nutanix" {
  alias    = "pc"
  endpoint = var.pc_endpoint
  username = var.pc_username
  password = var.pc_password
  port     = 9440
}

# Nutanix Central
provider "nutanix" {
  alias                    = "nc"
  nutanix_central_endpoint = var.nutanix_central_endpoint
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  port                     = 443
}

# ------------------------------------------------------------------------------
# DATA SOURCES - Discover PC Infrastructure
# ------------------------------------------------------------------------------

data "nutanix_pcs_v2" "pc" {
  provider = nutanix.pc
}

data "nutanix_clusters_v2" "clusters" {
  provider = nutanix.pc
}

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
}

# ------------------------------------------------------------------------------
# LOCATION LOOKUP (Nutanix Central)
# ------------------------------------------------------------------------------

data "nutanix_locations_v2" "registered_domain" {
  provider = nutanix.nc
  limit    = 10
  filter   = "startswith(country, 'India')"
  order_by = "name desc"
}

# ------------------------------------------------------------------------------
# REGISTER PC DOMAIN IN NUTANIX CENTRAL
# ------------------------------------------------------------------------------

resource "nutanix_registered_domain_v2" "registered_domain" {
  provider = nutanix.nc

  # recovered from state: domain_ext_id == the discovered Prism Central ext_id
  domain_ext_id   = local.pc_ext_id
  name            = var.registered_domain_name
  location_ext_id = data.nutanix_locations_v2.registered_domain.locations[1].ext_id

  # Domain registration must run only after AD, LDAP users and the admin
  # authorization policies on both PC and NC have been configured.
  depends_on = [
    nutanix_authorization_policy_v2.pc_admin,
    nutanix_authorization_policy_v2.nc_admin,
    nutanix_authorization_policy_v2.service_provider_admin,
  ]

  # API-issued; read via data.nutanix_registered_domain_v2.registration_creds.
  lifecycle {
    ignore_changes = [registration_config]
  }
}

# ------------------------------------------------------------------------------
# LOCAL DOMAIN ON PRISM CENTRAL (consumes NC registration credentials)
# ------------------------------------------------------------------------------
#
# Resource registration_config credentials are WriteOnly and omitted from state.
# Non-secret fields are also empty at plan time (attribute is not Computed), so
# read issued config/secrets from the registered-domain data source after create.

data "nutanix_registered_domain_v2" "registration_creds" {
  provider = nutanix.nc
  ext_id   = nutanix_registered_domain_v2.registered_domain.id
}

resource "nutanix_local_domain_v2" "local_domain" {
  provider = nutanix.pc

  domain_ext_id            = data.nutanix_registered_domain_v2.registration_creds.registration_config[0].domain_ext_id
  identity_provider_ext_id = nutanix_directory_services_v2.pc_active_directory.ext_id
  target_url               = data.nutanix_registered_domain_v2.registration_creds.registration_config[0].target_url
  tenant_ext_id            = data.nutanix_registered_domain_v2.registration_creds.registration_config[0].tenant_ext_id

  credentials {
    api_key = data.nutanix_registered_domain_v2.registration_creds.registration_config[0].credentials[0].api_key
    key_id  = data.nutanix_registered_domain_v2.registration_creds.registration_config[0].credentials[0].key_id
  }
}
