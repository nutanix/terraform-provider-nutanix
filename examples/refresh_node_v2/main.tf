terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
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

# Refresh a node's information by fetching the latest data from the hardware provider.
resource "nutanix_refresh_node_v2" "example" {
  ext_id = var.node_ext_id
}
