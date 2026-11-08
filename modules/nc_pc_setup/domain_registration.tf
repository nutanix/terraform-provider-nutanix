# ==============================================================================
# STAGE 3: Domain Registration (Nutanix Central <-> Prism Central)
# ==============================================================================
#
# Registers the Prism Central domain with Nutanix Central and then establishes
# the reciprocal local-domain record on Prism Central using the NC-issued
# registration credentials.
#
# NOTE: This file was reconstructed from terraform.tfstate after the original
#       source was lost. The resource wiring (attribute sources) was recovered
#       from the recorded state and dependency graph; review the lines marked
#       "recovered from state" if behaviour differs from the original.
# ==============================================================================

# ------------------------------------------------------------------------------
# LOCATION LOOKUP (Nutanix Central)
# ------------------------------------------------------------------------------

data "nutanix_locations_v2" "registered_domain" {
  provider = nutanix.nc
  limit    = 10
  filter   = "startswith(country, 'India')"
  order_by = "name desc"

  # LDAP NC auth only works after AD + LDAP user + admin policies exist.
  depends_on = [
    nutanix_authorization_policy_v2.nc_admin,
    nutanix_authorization_policy_v2.service_provider_admin,
  ]
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
