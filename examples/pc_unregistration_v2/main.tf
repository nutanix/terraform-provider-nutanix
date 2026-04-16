terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1.0"
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


resource "nutanix_pc_unregistration_v2 " "unregister_pc" {
  pc_ext_id = var.local_pc_ext_id
  ext_id    = var.remote_pc_ext_id
}
