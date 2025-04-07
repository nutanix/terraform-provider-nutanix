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

#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  clusters_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

## check if the node to add is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-node" {
  ext_id       = local.clusters_ext_id
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = var.node_ip
    }
  }

  ## check if the 3 nodes are un configured or not
  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 1
      error_message = "The node ${var.node_ip} is configured"
    }
  }
}

## fetch Network info for unconfigured node
resource "nutanix_clusters_unconfigured_node_networks_v2" "node-network-info" {
  ext_id       = local.clusters_ext_id
  request_type = "expand_cluster"
  node_list {
    cvm_ip {
      ipv4 {
        value = var.node_ip
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

## add node to the cluster
resource "nutanix_cluster_add_node_v2" "add-node" {
  cluster_ext_id = local.clusters_ext_id

  should_skip_add_node          = false
  should_skip_pre_expand_checks = false

  node_params {
    should_skip_host_networking = false
    hypervisor_isos {
      type = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_type
    }
    node_list {
      node_uuid                 = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].node_uuid
      model                     = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].rackable_unit_model
      block_id                  = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].rackable_unit_serial
      hypervisor_type           = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_type
      hypervisor_version        = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_version
      node_position             = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].node_position
      nos_version               = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].nos_version
      hypervisor_hostname       = "example"
      current_network_interface = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].name
      # required for adding a node
      hypervisor_ip {
        ipv4 {
          value = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_ip.0.ipv4.0.value
        }
      }
      cvm_ip {
        ipv4 {
          value = var.node_ip
        }
      }
      ipmi_ip {
        ipv4 {
          value = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].ipmi_ip.0.ipv4.0.value
        }
      }

      is_robo_mixed_hypervisor = true
      networks {
        name     = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].network_info[0].hci[0].name
        networks = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].network_info[0].hci[0].networks
        uplinks {
          active {
            name  = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].name
            mac   = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].mac
            value = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].name
          }
          standby {
            name  = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[1].name
            mac   = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[1].mac
            value = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[1].name
          }
        }
      }
    }

  }

  config_params {
    should_skip_imaging = true
    target_hypervisor   = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_type
  }

  remove_node_params {
    extra_params {
      should_skip_upgrade_check = false
      skip_space_check          = false
      should_skip_add_check     = false
    }
    should_skip_remove    = false
    should_skip_prechecks = false
  }

  depends_on = [nutanix_clusters_unconfigured_node_networks_v2.node-network-info]
}
