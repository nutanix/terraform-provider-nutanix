terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1"
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


resource "nutanix_unregister_cluster_v2 " "pc"{
  pc_ext_id = "<PC_UUID>"
  ext_id = "<cluster uuid>"
}