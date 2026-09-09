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

# Register a new installer image from an external URL in Foundation Central.
# `source = "REMOTE_URL"` downloads the image from `url`. `type` is one of
# AOS, AHV or ESX. For AHV images a `checksum` block (sha256 or md5) may be set.
resource "nutanix_installer_image_v2" "example" {
  name    = var.image_name
  type    = var.image_type
  source  = "REMOTE_URL"
  url     = var.image_url
  version = var.image_version

  # Optional checksum (applicable to AHV images). Uncomment ONE of the blocks.
  # checksum {
  #   sha256 {
  #     hex_digest = "0000000000000000000000000000000000000000000000000000000000000000"
  #   }
  # }
}

# Read a single installer image back by its external ID.
data "nutanix_installer_image_v2" "get" {
  ext_id = nutanix_installer_image_v2.example.ext_id
}

# List all installer images registered in Foundation Central.
data "nutanix_installer_images_v2" "list" {
  depends_on = [nutanix_installer_image_v2.example]
}

# List with paging / filtering options.
data "nutanix_installer_images_v2" "filtered" {
  limit      = 10
  depends_on = [nutanix_installer_image_v2.example]
}
