terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

# The Foundation Central configuration APIs are served by the Foundation Central
# (FCVM) service, NOT Prism Central. Configure the provider with the Foundation
# endpoint so these resources talk to the correct host.
provider "nutanix" {
  username            = var.nutanix_username
  password            = var.nutanix_password
  endpoint            = var.nutanix_endpoint
  port                = 9440
  insecure            = true
  foundation_endpoint = var.foundation_endpoint
  foundation_port     = var.foundation_port
}

# Manage the Foundation Central service configuration (singleton).
# Only the timeout knobs are user-writable; the remaining attributes are computed.
resource "nutanix_foundation_central_config_v2" "example" {
  ahv_installation_timeout_minutes = var.ahv_installation_timeout_minutes
  aos_download_timeout_minutes     = var.aos_download_timeout_minutes
}

# Read back the current Foundation Central configuration.
data "nutanix_foundation_central_config_v2" "current" {
  depends_on = [nutanix_foundation_central_config_v2.example]
}

output "foundation_central_version" {
  value = data.nutanix_foundation_central_config_v2.current.version
}
