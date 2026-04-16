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

#pull the storage container stats data
data "nutanix_storage_container_stats_info_v2" "test" {
  ext_id     = "<storage_container_ext_id>"
  start_time = "<start_time>"
  end_time   = "<end_time>"
}


