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



# Discover unconfigured nodes in a cluster

resource "nutanix_clusters_discover_unconfigured_nodes_v2" "test" {
  ext_id       = "<Cluster-UUID>"
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = "<IP-Address>"
    }
  }
}
