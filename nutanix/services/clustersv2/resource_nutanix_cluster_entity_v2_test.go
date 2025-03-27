package clustersv2_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameCluster              = "nutanix_cluster_v2.test"
	resourceNameDiscoverUnConfigNode = "nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node"
	resourceNameClusterRegistration  = "nutanix_pc_registration_v2.node-registration"
)

func TestAccV2NutanixClusterResource_CreateClusterWithMinimumConfig(t *testing.T) {
	if testVars.Clusters.Nodes[0].CvmIP == "" {
		t.Skip("Skipping test as No available node to be used for testing")
	}
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-cluster-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccClusterResourceMinimumConfig(name),
				PlanOnly: false,
			},
			{
				PreConfig: func() {
					time.Sleep(10 * time.Second) // 10-second delay
				},
				Config: testAccClusterResourceMinimumConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCluster, "name", name),
					resource.TestCheckResourceAttr(resourceNameCluster, "dryrun", "false"),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
				),
			},
		},
	})
}

func TestAccV2NutanixClusterResource_CreateClusterWithAllConfig(t *testing.T) {
	if testVars.Clusters.Nodes[0].CvmIP == "" {
		t.Skip("Skipping test as No available node to be used for testing")
	}
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-cluster-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:   testAccClusterResourceAllConfig(name),
				PlanOnly: false,
			},
			{
				PreConfig: func() {
					time.Sleep(10 * time.Second) // 10-second delay
				},
				Config: testAccClusterResourceAllConfig(name),
				Check: resource.ComposeTestCheckFunc(
					// check the unconfigured node is discovered or not
					resource.TestCheckResourceAttr(resourceNameDiscoverUnConfigNode, "address_type", "IPV4"),
					resource.TestCheckResourceAttr(resourceNameDiscoverUnConfigNode, "ip_filter_list.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(resourceNameDiscoverUnConfigNode, "unconfigured_nodes.#", "1"),
					resource.TestCheckResourceAttr(resourceNameDiscoverUnConfigNode, "unconfigured_nodes.0.cvm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttrSet(resourceNameDiscoverUnConfigNode, "unconfigured_nodes.0.nos_version"),
					resource.TestCheckResourceAttrSet(resourceNameDiscoverUnConfigNode, "unconfigured_nodes.0.hypervisor_type"),

					// check the cluster is created with minimum config
					resource.TestCheckResourceAttr(resourceNameCluster, "name", name),
					resource.TestCheckResourceAttr(resourceNameCluster, "dryrun", "false"),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.external_address.0.ipv4.0.value", testVars.Clusters.Network.VirtualIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.0.fqdn.0.value", testVars.Clusters.Network.NTPServers[0]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.1.fqdn.0.value", testVars.Clusters.Network.NTPServers[1]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.2.fqdn.0.value", testVars.Clusters.Network.NTPServers[2]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.3.fqdn.0.value", testVars.Clusters.Network.NTPServers[3]),

					resource.TestCheckResourceAttrSet(resourceNameClusterRegistration, "pc_ext_id"),
					resource.TestCheckResourceAttr(resourceNameClusterRegistration, "remote_cluster.0.aos_remote_cluster_spec.0.remote_cluster.0.address.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(resourceNameClusterRegistration, "remote_cluster.0.aos_remote_cluster_spec.0.remote_cluster.0.credentials.0.authentication.0.username", testVars.Clusters.Nodes[0].Username),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(10 * time.Second) // 10-second delay
				},
				Config: testAccClusterResourceUpdateConfig(name+"-updated", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCluster, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceNameCluster, "dryrun", "false"),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.pulse_status.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.pulse_status.0.pii_scrubbing_level", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.external_address.0.ipv4.0.value", testVars.Clusters.Network.VirtualIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.external_data_services_ip.0.ipv4.0.value", testVars.Clusters.Network.IscsiIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.0.fqdn.0.value", testVars.Clusters.Network.NTPServers[0]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.1.fqdn.0.value", testVars.Clusters.Network.NTPServers[1]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.2.fqdn.0.value", testVars.Clusters.Network.NTPServers[2]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.3.fqdn.0.value", testVars.Clusters.Network.NTPServers[3]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.email_address", testVars.Clusters.Network.SMTPServer.EmailAddress),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.server.0.ip_address.0.ipv4.0.value", testVars.Clusters.Network.SMTPServer.IP),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.server.0.port", strconv.Itoa(testVars.Clusters.Network.SMTPServer.Port)),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.server.0.username", testVars.Clusters.Network.SMTPServer.Username),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.type", testVars.Clusters.Network.SMTPServer.Type),
				),
			},
			// Disable the cluster pulse status
			{
				PreConfig: func() {
					time.Sleep(10 * time.Second) // 10-second delay
				},
				Config: testAccClusterResourceUpdateConfig(name+"-updated", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCluster, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceNameCluster, "dryrun", "false"),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.pulse_status.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.external_address.0.ipv4.0.value", testVars.Clusters.Network.VirtualIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.external_data_services_ip.0.ipv4.0.value", testVars.Clusters.Network.IscsiIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.0.fqdn.0.value", testVars.Clusters.Network.NTPServers[0]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.1.fqdn.0.value", testVars.Clusters.Network.NTPServers[1]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.2.fqdn.0.value", testVars.Clusters.Network.NTPServers[2]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.ntp_server_ip_list.3.fqdn.0.value", testVars.Clusters.Network.NTPServers[3]),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.email_address", testVars.Clusters.Network.SMTPServer.EmailAddress),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.server.0.ip_address.0.ipv4.0.value", testVars.Clusters.Network.SMTPServer.IP),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.server.0.port", strconv.Itoa(testVars.Clusters.Network.SMTPServer.Port)),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.server.0.username", testVars.Clusters.Network.SMTPServer.Username),
					resource.TestCheckResourceAttr(resourceNameCluster, "network.0.smtp_server.0.type", testVars.Clusters.Network.SMTPServer.Type),
				),
			},
		},
	})
}

var clusterConfig = fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
	  cluster_ext_id = [
		for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
	  ][0]
	  config   = (jsondecode(file("%[1]s")))
	  clusters = local.config.clusters
	}`, filepath)

func testAccClusterResourceMinimumConfig(name string) string {
	return fmt.Sprintf(`
		# cluster config
		%[1]s

		# check if the nodes is un configured or not
		resource "nutanix_clusters_discover_unconfigured_nodes_v2" "test-discover-cluster-node" {
		  ext_id       = local.cluster_ext_id
		  address_type = "IPV4"
		  ip_filter_list {
			ipv4 {
			  value = local.clusters.nodes[0].cvm_ip
			}
		  }
		  depends_on = [data.nutanix_clusters_v2.clusters]

		  ## check if the node is  un configured or not
		  lifecycle {
			postcondition {
			  condition     = length(self.unconfigured_nodes) == 1
			  error_message = "The node ${local.clusters.nodes[0].cvm_ip} are not unconfigured"
			}
		  }
		}

		# create a new cluster
		resource "nutanix_cluster_v2" "test" {
		  name   = "%[2]s"
		  dryrun = false
		  nodes {
			node_list {
			  controller_vm_ip {
				ipv4 {
				  value = local.clusters.nodes[0].cvm_ip
				}
			  }
			}
		  }
		  config {
			cluster_function = local.clusters.config.cluster_functions
			cluster_arch     = local.clusters.config.cluster_arch
			fault_tolerance_state {
			  domain_awareness_level          = local.clusters.config.fault_tolerance_state.domain_awareness_level
			}
		  }

		  provisioner "local-exec" {
			command = "ssh-keygen -f '~/.ssh/known_hosts' -R '${local.clusters.nodes[0].cvm_ip}';  sshpass -p '${local.clusters.pe_password}' ssh -o StrictHostKeyChecking=no ${local.clusters.pe_username}@${local.clusters.nodes[0].cvm_ip} '/home/nutanix/prism/cli/ncli user reset-password user-name=${local.clusters.nodes[0].username} password=${local.clusters.nodes[0].password}' "

			on_failure = continue
		  }
          # Set lifecycle to ignore changes
		  lifecycle {
			ignore_changes = [network.0.smtp_server.0.server.0.password,  links, categories, config.0.cluster_function]
		  }
		  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node]
		}


		# register the cluster to pc
		resource "nutanix_pc_registration_v2" "node-registration" {
		  pc_ext_id = local.cluster_ext_id
		  remote_cluster {
			aos_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = local.clusters.nodes[0].cvm_ip
				  }
				}
				credentials {
				  authentication {
					username = local.clusters.nodes[0].username
					password = local.clusters.nodes[0].password
				  }
				}
			  }
			}
		  }
		  depends_on = [nutanix_cluster_v2.test]
		}

`, clusterConfig, name)
}

func testAccClusterResourceAllConfig(name string) string {
	return fmt.Sprintf(`
		%[1]s

		# check if the nodes is un configured or not
		resource "nutanix_clusters_discover_unconfigured_nodes_v2" "test-discover-cluster-node" {
		  ext_id       = local.cluster_ext_id
		  address_type = "IPV4"
		  ip_filter_list {
			ipv4 {
			  value = local.clusters.nodes[0].cvm_ip
			}
		  }
		  depends_on = [data.nutanix_clusters_v2.clusters]

		  ## check if the node is  un configured or not
		  lifecycle {
			postcondition {
			  condition     = length(self.unconfigured_nodes) == 1
			  error_message = "The node ${local.clusters.nodes[0].cvm_ip} are not unconfigured"
			}
		  }
		}


		resource "nutanix_cluster_v2" "test" {
		  name   = "%[2]s"
		  dryrun = false
		  nodes {
			node_list {
			  controller_vm_ip {
				ipv4 {
				  value = local.clusters.nodes[0].cvm_ip
				}
			  }
			}
		  }
		  config {
			cluster_function = local.clusters.config.cluster_functions
			cluster_arch     = local.clusters.config.cluster_arch
			fault_tolerance_state {
			  domain_awareness_level          = local.clusters.config.fault_tolerance_state.domain_awareness_level
			}
		  }
		  network {
			external_address {
			  ipv4 {
				value = local.clusters.network.virtual_ip
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[0]
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[1]
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[2]
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[3]
			  }
			}
		  }

		  lifecycle {
			ignore_changes = [network.0.smtp_server.0.server.0.password,  links, categories, config.0.cluster_function]
		  }

		  provisioner "local-exec" {
			command = "ssh-keygen -f '~/.ssh/known_hosts' -R '${local.clusters.nodes[0].cvm_ip}'; sshpass -p '${local.clusters.pe_password}' ssh -o StrictHostKeyChecking=no ${local.clusters.pe_username}@${local.clusters.nodes[0].cvm_ip} '/home/nutanix/prism/cli/ncli user reset-password user-name=${local.clusters.nodes[0].username} password=${local.clusters.nodes[0].password}' "

			on_failure = continue
		  }
		  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node]
		}

		# register the cluster to pc
		resource "nutanix_pc_registration_v2" "node-registration" {
		  pc_ext_id = local.cluster_ext_id
		  remote_cluster {
			aos_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = local.clusters.nodes[0].cvm_ip
				  }
				}
				credentials {
				  authentication {
					username = local.clusters.nodes[0].username
					password = local.clusters.nodes[0].password
				  }
				}
			  }
			}
		  }
		  depends_on = [nutanix_cluster_v2.test]
		}
	`, clusterConfig, name)
}

func testAccClusterResourceUpdateConfig(updatedName, pulseStatus string) string {
	return fmt.Sprintf(`
		# cluster config
		%[1]s

		# check if the nodes is un configured or not
		resource "nutanix_clusters_discover_unconfigured_nodes_v2" "test-discover-cluster-node" {
		  ext_id       = local.cluster_ext_id
		  address_type = "IPV4"
		  ip_filter_list {
			ipv4 {
			  value = local.clusters.nodes[0].cvm_ip
			}
		  }
		  depends_on = [data.nutanix_clusters_v2.clusters]

		  ## check if the node is  un configured or not
		  lifecycle {
			postcondition {
			  condition     = length(self.unconfigured_nodes) == 1
			  error_message = "The node ${local.clusters.nodes[0].cvm_ip} are not unconfigured"
			}
		  }
		}


		resource "nutanix_cluster_v2" "test" {
		  name   = "%[2]s"
		  dryrun = false
		  nodes {
			node_list {
			  controller_vm_ip {
				ipv4 {
				  value = local.clusters.nodes[0].cvm_ip
				}
			  }
			}
		  }
		  config {
			cluster_function = local.clusters.config.cluster_functions
			cluster_arch     = local.clusters.config.cluster_arch
			fault_tolerance_state {
			  domain_awareness_level          = local.clusters.config.fault_tolerance_state.domain_awareness_level
			}
		    pulse_status {
		      is_enabled = %[3]s
		      pii_scrubbing_level = "DEFAULT"
		    }
		  }
		  # update the network config, external_address, external data services ip, smtp server
		  network {
			external_address {
			  ipv4 {
				value = local.clusters.network.virtual_ip
			  }
			}
			external_data_services_ip {
			  ipv4 {
				value = local.clusters.network.iscsi_ip
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[0]
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[1]
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[2]
			  }
			}
			ntp_server_ip_list {
			  fqdn {
				value = local.clusters.network.ntp_servers[3]
			  }
			}
			smtp_server {
			  email_address = local.clusters.network.smtp_server.email_address
			  server {
				ip_address {
				  ipv4 {
					value = local.clusters.network.smtp_server.ip
				  }
				}
				port     = local.clusters.network.smtp_server.port
				username = local.clusters.network.smtp_server.username
				password = local.clusters.network.smtp_server.password
			  }
			  type = local.clusters.network.smtp_server.type
			}
		  }

		  lifecycle {
			ignore_changes = [network.0.smtp_server.0.server.0.password,  links, categories, config.0.cluster_function]
		  }

		  provisioner "local-exec" {

			command = "ssh-keygen -f '~/.ssh/known_hosts' -R '${local.clusters.nodes[0].cvm_ip}';  sshpass -p '${local.clusters.pe_password}' ssh -o StrictHostKeyChecking=no ${local.clusters.pe_username}@${local.clusters.nodes[0].cvm_ip} '/home/nutanix/prism/cli/ncli user reset-password user-name=${local.clusters.nodes[0].username} password=${local.clusters.nodes[0].password}' "
			on_failure = continue
		  }
		  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node]
		}

		# register the cluster to pc
		resource "nutanix_pc_registration_v2" "node-registration" {
		  pc_ext_id = local.cluster_ext_id
		  remote_cluster {
			aos_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = local.clusters.nodes[0].cvm_ip
				  }
				}
				credentials {
				  authentication {
					username = local.clusters.nodes[0].username
					password = local.clusters.nodes[0].password
				  }
				}
			  }
			}
		  }
		  depends_on = [nutanix_cluster_v2.test]
		}

`, clusterConfig, updatedName, pulseStatus)
}
