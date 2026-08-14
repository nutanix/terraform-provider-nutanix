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

# Register a Prism Central (PC) domain with Nutanix Central (NC).
resource "nutanix_registered_domain_v2" "example" {
  name                  = var.registered_domain_name
  fqdn                  = var.registered_domain_fqdn
  dns_name              = var.registered_domain_fqdn
  domain_ext_id         = "00000000-0000-0000-0000-000000000000"
  location_ext_id       = "00000000-0000-0000-0000-000000000000"
  credential_issuer     = "nutanix-central"
  created_by            = "admin"
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
    domain_ext_id = "00000000-0000-0000-0000-000000000000"
    target_url    = "https://nc.example.com"
    tenant_ext_id = "00000000-0000-0000-0000-000000000000"
    credentials {
      api_key = var.registered_domain_api_key
      key_id  = "example-key-id"
    }
  }

  tunnel {
    endpoint = "https://tunnel.example.com"
    provider = "IKAT"
  }
}

# Sync (change) the registration state of a registered domain.
resource "nutanix_registered_domain_change_state_v2" "sync" {
  ext_id             = nutanix_registered_domain_v2.example.id
  registration_state = "REGISTERED"
}

# Establish trust for a registered domain by exchanging a client-signed request.
resource "nutanix_registered_domain_establish_trust_v2" "trust" {
  ext_id                 = nutanix_registered_domain_v2.example.id
  client_signing_request = var.client_signing_request
}

# Refresh the API credentials used to communicate with a registered domain.
resource "nutanix_registered_domain_refresh_api_credentials_v2" "refresh" {
  ext_id = nutanix_registered_domain_v2.example.id
}

# Fetch a single registered domain by its external identifier.
data "nutanix_registered_domain_v2" "get" {
  ext_id = nutanix_registered_domain_v2.example.id
}

# List all registered domains.
data "nutanix_registered_domains_v2" "list" {}

# List registered domains filtered by name.
data "nutanix_registered_domains_v2" "filtered" {
  limit  = 5
  filter = "name eq '${nutanix_registered_domain_v2.example.name}'"
}
