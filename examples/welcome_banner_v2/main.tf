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

# Configure the cluster-wide welcome banner (iam / UpdateWelcomeBanner).
# The welcome banner is a singleton configuration; there is exactly one per PC.
resource "nutanix_welcome_banner_v2" "example" {
  content    = var.welcome_banner_content
  is_enabled = var.welcome_banner_is_enabled
}
