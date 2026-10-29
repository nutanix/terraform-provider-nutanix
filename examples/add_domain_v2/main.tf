terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_central_username # User with Service Provider Admin access
  password = var.nutanix_central_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port     = 443
}

# Adds an existing domain to the tenancy allowed-domain list.
resource "nutanix_add_domain_v2" "example" {
  domain_ext_id = var.domain_ext_id
}

# Lists all domains that are allowed for tenancy.
data "nutanix_domains_v2" "example" {
  depends_on = [nutanix_add_domain_v2.example]
}
