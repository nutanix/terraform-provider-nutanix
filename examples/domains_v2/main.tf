terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
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

# Lists all domains that are allowed for tenancy.
data "nutanix_domains_v2" "example" {
}

# Output the list of domains for verification
output "domains" {
  value = data.nutanix_domains_v2.example.domains
}