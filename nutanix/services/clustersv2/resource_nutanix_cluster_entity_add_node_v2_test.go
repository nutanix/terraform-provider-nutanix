package clustersv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameDiscoverUnconfiguredNode         = "nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node"
	resourceNameDiscoverUnconfiguredClusterNodes = "nutanix_clusters_discover_unconfigured_nodes_v2.cluster-nodes"
	resourceNameFetchUnconfiguredNodeNetwork     = "nutanix_clusters_unconfigured_node_networks_v2.node-network-info"
	resourceNameAddNodeToCluster                 = "nutanix_cluster_add_node_v2.test"
	resourceName3NodesCluster                    = "nutanix_cluster_v2.cluster-3nodes"
)

func TestAccV2NutanixClusterAddNodeResource_Basic(t *testing.T) {
	if testVars.Clusters.Nodes[1].CvmIP == "" &&
		testVars.Clusters.Nodes[2].CvmIP == "" &&
		testVars.Clusters.Nodes[3].CvmIP == "" {
		t.Skip("Skipping test as No available nodes to be used for testing")
	}
	r := acctest.RandInt()
	clusterName := fmt.Sprintf("tf-3node-cluster-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClustersConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					// add node to cluster check
					// unconfigured Nodes check
					resource.TestCheckResourceAttr(resourceNameDiscoverUnconfiguredClusterNodes, "unconfigured_nodes.#", "3"),

					//Cluster Check
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "nodes.0.node_list.#", "3"),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_function.0", testVars.Clusters.Config.ClusterFunctions[0]),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
				),
			},
			{
				Config: testAccClustersConfig(clusterName) + testAccAddNodeToClusterConfig(),
				Check: resource.ComposeTestCheckFunc(
					// unconfigured Node to be added check
					resource.TestCheckResourceAttr(resourceNameDiscoverUnconfiguredNode, "unconfigured_nodes.#", "1"),
					resource.TestCheckResourceAttr(resourceNameDiscoverUnconfiguredNode, "unconfigured_nodes.0.cvm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[3].CvmIP),
					resource.TestCheckResourceAttrSet(resourceNameDiscoverUnconfiguredNode, "unconfigured_nodes.0.nos_version"),
					resource.TestCheckResourceAttrSet(resourceNameDiscoverUnconfiguredNode, "unconfigured_nodes.0.node_uuid"),

					// fetch network info for unconfigured node check
					resource.TestCheckResourceAttr(resourceNameFetchUnconfiguredNodeNetwork, "nodes_networking_details.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameFetchUnconfiguredNodeNetwork, "nodes_networking_details.0.network_info.#"),
					resource.TestCheckResourceAttrSet(resourceNameFetchUnconfiguredNodeNetwork, "nodes_networking_details.0.uplinks.#"),
				),
			},
			{
				PreConfig: func() {
					t.Log("Sleeping for 10 Minute before removing the node")
					time.Sleep(10 * time.Minute)
				},
				Config: testAccClustersConfig(clusterName) + testAccAddNodeToClusterConfig(),
				Check: resource.ComposeTestCheckFunc(
					// add node to cluster check
					resource.TestCheckResourceAttr(resourceNameAddNodeToCluster, "node_params.0.node_list.0.cvm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[3].CvmIP),
				),
			},
		},
	})
}

func testAccClustersConfig(clusterName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]
  config   = (jsondecode(file("%[2]s")))
  clusters = local.config.clusters
}


############################ cluster with 3 nodes

## check if the nodes is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-nodes" {
  ext_id       = local.cluster_ext_id
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = local.clusters.nodes[0].cvm_ip
    }
  }
  ip_filter_list {
    ipv4 {
      value = local.clusters.nodes[1].cvm_ip
    }
  }
  ip_filter_list {
    ipv4 {
      value = local.clusters.nodes[2].cvm_ip
    }
  }
  depends_on = [data.nutanix_clusters_v2.clusters]

  ## check if the 3 nodes are un configured or not
  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 3
      error_message = "The nodes are not unconfigured"
    }
  }
}


resource "nutanix_cluster_v2" "cluster-3nodes" {
  name   = "%[1]s"
  dryrun = false
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = local.clusters.nodes[0].cvm_ip
        }
      }
    }
    node_list {
      controller_vm_ip {
        ipv4 {
          value = local.clusters.nodes[1].cvm_ip
        }
      }
    }
    node_list {
      controller_vm_ip {
        ipv4 {
          value = local.clusters.nodes[2].cvm_ip
        }
      }
    }
  }
  config {
    cluster_function = local.clusters.config.cluster_functions
    cluster_arch     = local.clusters.config.cluster_arch
    fault_tolerance_state {
      domain_awareness_level          = "NODE"
    }
  }

  provisioner "local-exec" {
    command = "ssh-keygen -f ~/.ssh/known_hosts -R ${local.clusters.nodes[1].cvm_ip};   sshpass -p '${local.clusters.pe_password}' ssh -o StrictHostKeyChecking=no ${local.clusters.pe_username}@${local.clusters.nodes[1].cvm_ip} '/home/nutanix/prism/cli/ncli user reset-password user-name=${local.clusters.nodes[1].username} password=${local.clusters.nodes[1].password}'"

    on_failure = continue
  }

  lifecycle {
    ignore_changes = [nodes.0.node_list, links, categories, config.0.cluster_function]
  }

  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.cluster-nodes]
}



## we need only to rgister on of 3 nodes tp pc
resource "nutanix_pc_registration_v2" "nodes-registration" {
  pc_ext_id = local.cluster_ext_id
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = local.clusters.nodes[1].cvm_ip
          }
        }
        credentials {
          authentication {
            username = local.clusters.nodes[1].username
            password = local.clusters.nodes[1].password
          }
        }
      }
    }
  }
  depends_on = [nutanix_cluster_v2.cluster-3nodes]

  provisioner "local-exec" {
    command    = " sleep 5s"
    on_failure = continue
  }

}


`, clusterName, filepath)
}

func testAccAddNodeToClusterConfig() string {
	return `


################################# add node

## check if the node to add is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-node" {
  ext_id = nutanix_cluster_v2.cluster-3nodes.id
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

## fetch Network info for unconfigured node
resource "nutanix_clusters_unconfigured_node_networks_v2" "node-network-info" {
  ext_id       = nutanix_cluster_v2.cluster-3nodes.id
  request_type = "expand_cluster"
  node_list {
    cvm_ip {
      ipv4 {
        value = local.clusters.nodes[3].cvm_ip
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
resource "nutanix_cluster_add_node_v2" "test" {
  cluster_ext_id = nutanix_cluster_v2.cluster-3nodes.id

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
      hypervisor_hostname       = "test"
      current_network_interface = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].name
      # required for adding a node
      hypervisor_ip {
        ipv4 {
          value = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_ip.0.ipv4.0.value
        }
      }
      cvm_ip {
        ipv4 {
          value = local.clusters.nodes[3].cvm_ip
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


`
}
