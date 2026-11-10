terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

# Nutanix Central Admin — NC-side domain registration and location discovery.
provider "nutanix" {
  alias                    = "nc"
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port                     = 443
}

# Prism Central — local-domain registration runs against the PC endpoint.
provider "nutanix" {
  alias    = "pc"
  username = var.nutanix_pc_username
  password = var.nutanix_pc_password
  endpoint = var.nutanix_pc_endpoint
  port     = 9440
  insecure = true
}

# ---------------------------------------------------------------------------
# Locations available for domain registration (NC Admin)
# ---------------------------------------------------------------------------
data "nutanix_locations_v2" "all" {
  provider = nutanix.nc
}

data "nutanix_locations_v2" "paged" {
  provider = nutanix.nc
  page     = 0
  limit    = 20
}

data "nutanix_locations_v2" "filtered" {
  provider = nutanix.nc
  limit    = 10
  filter   = "startswith(country, 'India')"
  order_by = "name desc"
}

# ---------------------------------------------------------------------------
# Register a Prism Central domain with Nutanix Central (NC side)
# ---------------------------------------------------------------------------
resource "nutanix_registered_domain_v2" "example" {
  provider              = nutanix.nc
  name                  = var.registered_domain_name
  fqdn                  = var.registered_domain_fqdn
  dns_name              = var.registered_domain_fqdn
  domain_ext_id         = var.domain_ext_id
  location_ext_id       = var.location_ext_id
  credential_issuer     = "nutanix-central"
  created_by            = var.nutanix_central_username
  api_credential_status = "VALID"
  registration_state    = "REGISTERED"
  registered_time       = "2024-01-01T00:00:00Z"

  platform_data {
    pc_version          = "pc.7.6"
    connectivity_status = "REACHABLE"
    pc_env_type         = "ONPREM"
    pc_provider_type    = "VM"
    cluster_count       = 1
  }

  registration_config {
    domain_ext_id = var.domain_ext_id
    target_url    = var.nc_target_url
    tenant_ext_id = var.nc_tenant_ext_id
    credentials {
      api_key = var.registered_domain_api_key
      key_id  = var.registered_domain_key_id
    }
  }

  tunnel {
    endpoint = var.tunnel_endpoint
    provider = "IKAT"
  }
}

resource "nutanix_registered_domain_change_state_v2" "sync" {
  provider           = nutanix.nc
  ext_id             = nutanix_registered_domain_v2.example.id
  registration_state = "REGISTERED"
}

resource "nutanix_registered_domain_establish_trust_v2" "trust" {
  provider               = nutanix.nc
  ext_id                 = nutanix_registered_domain_v2.example.id
  client_signing_request = var.client_signing_request
}

resource "nutanix_registered_domain_refresh_api_credentials_v2" "refresh" {
  provider = nutanix.nc
  ext_id   = nutanix_registered_domain_v2.example.id
}

data "nutanix_registered_domain_v2" "get" {
  provider = nutanix.nc
  ext_id   = nutanix_registered_domain_v2.example.id
}

data "nutanix_registered_domains_v2" "list" {
  provider = nutanix.nc
}

data "nutanix_registered_domains_v2" "filtered" {
  provider = nutanix.nc
  limit    = 5
  filter   = "name eq '${nutanix_registered_domain_v2.example.name}'"
}

# ---------------------------------------------------------------------------
# Register this Prism Central as a local domain with Nutanix Central (PC side)
# LocalDomainService is a singleton: one local domain per Prism Central.
# ---------------------------------------------------------------------------
resource "nutanix_local_domain_v2" "example" {
  provider = nutanix.pc

  credentials {
    api_key = var.nc_api_key
    key_id  = var.nc_key_id
  }

  domain_ext_id            = var.domain_ext_id
  identity_provider_ext_id = var.nc_identity_provider_ext_id
  target_url               = var.nc_target_url
  tenant_ext_id            = var.nc_tenant_ext_id
}

data "nutanix_local_domain_v2" "get_local_domain" {
  provider   = nutanix.pc
  depends_on = [nutanix_local_domain_v2.example]
}

output "registered_domain_ext_id" {
  value = nutanix_registered_domain_v2.example.id
}

output "locations" {
  value = data.nutanix_locations_v2.all
}
