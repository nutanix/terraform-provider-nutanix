terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
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

# Fetch the configured welcome banner (singleton datasource — no input required).
data "nutanix_welcome_banner_v2" "welcome-banner" {}

output "welcome_banner" {
  value = data.nutanix_welcome_banner_v2.welcome-banner
}
