terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}



# ## check if the node to add is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-node" {
  ext_id = "<Cluster UUID>"
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = "<Node CVM IPV4 Address>"
    }
  }

  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 1
      error_message = "The node is not unconfigured"
    }
  }
}

# ## fetch Network info for unconfigured node
resource "nutanix_clusters_unconfigured_node_networks_v2" "node-network-info" {
    ext_id = "<Cluster UUID>"
  request_type = "expand_cluster"
  node_list {
    cvm_ip {
      ipv4 {
        value = "<Node CVM IPV4 Address>"
      }
    }
    hypervisor_ip {
      ipv4 {
        value = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_ip.0.ipv4.0.value
      }
    }
  }
  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node]

}
