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
# ## check if the node to add is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-node" {
  ext_id       = local.cluster_ext_id
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = var.cvm_ip
    }
  }

  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 1
      error_message = "The node is not unconfigured"
    }
  }
}

## check if the node to add is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-node" {
  ext_id       = nutanix_cluster_v2.cluster-3nodes.id
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = local.clusters.nodes[3].cvm_ip
    }
  }

  ## check if the 3 nodes are un configured or not
  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 1
      error_message = "The node ${local.clusters.nodes[3].cvm_ip} is configured"
    }
  }
  depends_on = [nutanix_pc_registration_v2.nodes-registration]
}


# ## fetch Network info for unconfigured node
resource "nutanix_clusters_unconfigured_node_networks_v2" "node-network-info" {
  ext_id       = local.cluster_ext_id
  request_type = "expand_cluster"
  node_list {
    cvm_ip {
      ipv4 {
        value = var.cvm_ip
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
