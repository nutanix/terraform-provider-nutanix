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



data "nutanix_clusters_v2" "clusters" {}
locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]

}

# Discover unconfigured nodes in a cluster
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "test" {
  ext_id       = local.cluster_ext_id
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = var.cvm_ip
    }
  }
}
