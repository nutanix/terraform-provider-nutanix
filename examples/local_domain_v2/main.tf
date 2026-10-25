terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.1"
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

# Register this Prism Central as a local domain with Nutanix Central.
# LocalDomainService is a singleton: there is exactly one local domain per
# Prism Central. The api_key/key_id credentials are issued by Nutanix Central.
resource "nutanix_local_domain_v2" "example" {
  credentials {
    api_key = var.nc_api_key
    key_id  = var.nc_key_id
  }
  domain_ext_id            = var.nc_domain_ext_id
  identity_provider_ext_id = var.nc_identity_provider_ext_id
  target_url               = var.nc_target_url
  tenant_ext_id            = var.nc_tenant_ext_id
}

# Read the registered local domain entity (singleton — no ext_id input).
data "nutanix_local_domain_v2" "get-local-domain" {
  depends_on = [nutanix_local_domain_v2.example]
}
