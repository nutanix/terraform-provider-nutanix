terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

# Service Provider Admin — allow a registered domain for tenancy.
provider "nutanix" {
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  port                     = 443
}

# Adds an existing domain to the tenancy allowed-domain list.
resource "nutanix_add_domain_v2" "example" {
  domain_ext_id = var.domain_ext_id
}

# Lists all domains that are allowed for tenancy.
data "nutanix_domains_v2" "example" {
  depends_on = [nutanix_add_domain_v2.example]
}

output "allowed_domains" {
  value = data.nutanix_domains_v2.example.domains
}
