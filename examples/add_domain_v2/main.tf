terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Adds an existing domain to the tenancy allowed-domain list.
resource "nutanix_add_domain_v2" "example" {
  domain_ext_id = var.domain_ext_id
}

# Lists all domains that are allowed for tenancy.
data "nutanix_domains_v2" "example" {
  depends_on = [nutanix_add_domain_v2.example]
}
