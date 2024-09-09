terraform{
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = "1.3.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}

#pull all clusters data
data "nutanix_clusters" "clusters"{}

#create local variable pointing to desired cluster
locals {
  cluster1 = [
    for cluster in data.nutanix_clusters.clusters.entities :
    cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
  ][0]
}

#pull the storage container stats data
data "nutanix_storage_container_stats_info_v2" "test" {
  ext_id = "<storage_container_ext_id>"
  start_time = "<start_time>"
  end_time = "<end_time>"
}


