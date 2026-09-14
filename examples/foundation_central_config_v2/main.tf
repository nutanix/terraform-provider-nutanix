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

# Manage the (singleton) Foundation Central service configuration.
# Only the two timeout fields are writable; commit_id, tls_certificate_fingerprint
# and version are populated by the server.
resource "nutanix_foundation_central_config_v2" "example" {
  ahv_installation_timeout_minutes = var.ahv_installation_timeout_minutes
  aos_download_timeout_minutes     = var.aos_download_timeout_minutes
}

# Read the current Foundation Central configuration (singleton — no id input).
data "nutanix_foundation_central_config_v2" "current" {
  depends_on = [nutanix_foundation_central_config_v2.example]
}
