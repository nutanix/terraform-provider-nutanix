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
	resourceNameCluster                    = "nutanix_cluster_v2.test"
	resourceNameDiscoverUnConfigNode       = "nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node"
	resourceNameClusterRegistration        = "nutanix_pc_registration_v2.node-registration"
	dataSourceNameClusterData              = "data.nutanix_cluster_v2.cluster"
	dataSourceNameGetClusterCategoriesData = "data.nutanix_clusters_v2.get-cluster-categories"
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
				Config:   testAccClusterResourceMinimumConfig(name, ""),
				PlanOnly: false,
			},
			{
				PreConfig: func() {
					time.Sleep(10 * time.Second) // 10-second delay
				},
				Config: testAccClusterResourceMinimumConfig(name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCluster, "name", name),
					resource.TestCheckResourceAttr(resourceNameCluster, "dryrun", "false"),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(resourceNameCluster, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(resourceNameCluster, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
				),
			},
			// Step 3: Associate categories to the cluster and check on list cluster data source for categories
			{
				Config: testAccClusterResourceMinimumConfig(name, "nutanix_category_v2.cat-1.id, nutanix_category_v2.cat-2.id, nutanix_category_v2.cat-3.id") +
					`				
					# get the cluster data
					data "nutanix_cluster_v2" "cluster" {
						ext_id = nutanix_cluster_v2.test.id
					}

					# get the clusters data from the data source
					data "nutanix_clusters_v2" "get-cluster-categories" {
						filter = "name eq '${nutanix_cluster_v2.test.name}'"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					// Check categories count on the resource itself (TypeSet)
					resource.TestCheckResourceAttr(resourceNameCluster, "categories.#", "3"),
					// Check categories on the resource (order-independent check)
					checkCategories(resourceNameCluster, "categories", []string{
						"nutanix_category_v2.cat-1",
						"nutanix_category_v2.cat-2",
						"nutanix_category_v2.cat-3",
					}),
					// Check categories count on the data source
					resource.TestCheckResourceAttr(dataSourceNameClusterData, "categories.#", "3"),
					// Check categories on the data source (order-independent check)
					checkCategories(dataSourceNameClusterData, "categories", []string{
						"nutanix_category_v2.cat-1",
						"nutanix_category_v2.cat-2",
						"nutanix_category_v2.cat-3",
					}),

					// Check categories count on the data source
					resource.TestCheckResourceAttr(dataSourceNameGetClusterCategoriesData, "cluster_entities.0.categories.#", "3"),
					// Check categories on the data source (order-independent check)
					checkCategories(dataSourceNameGetClusterCategoriesData, "cluster_entities.0.categories", []string{
						"nutanix_category_v2.cat-1",
						"nutanix_category_v2.cat-2",
						"nutanix_category_v2.cat-3",
					}),
				),
			},
		},
	})
}

func TestAccV2NutanixClusterResource_CreateClusterWithAllConfig(t *testing.T) {
	if testVars.Clusters.Nodes[0].CvmIP == "" {
		t.Skip("Skipping test as No available node to be used for testing")
	}
	r := acctest.RandIntRange(1, 10000)
	name := fmt.Sprintf("tf-test-cluster-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckNutanixClusterDestroy,
			testAccCheckNutanixClusterCategoriesDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Plan
			{
				Config:   testAccClusterResourceAllConfig(name),
				PlanOnly: false,
			},
			// Step 2: Apply
			{
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
			// ############################################## Associate categories with cluster ##############################################
			// Step 3: Associate categories to the cluster and check on list cluster data source for categories
			{
				Config: testAccClusterResourceAllConfig(name) + testAccClusterResourceAssociateCategoriesConfig(r),
				Check: resource.ComposeTestCheckFunc(
					// check on list cluster data source for categories (order-independent)
					checkCategories(dataSourceNameClusters, "cluster_entities.0.categories", []string{
						"nutanix_category_v2.cat-1",
						"nutanix_category_v2.cat-2",
						"nutanix_category_v2.cat-3",
					}),

					// check on cluster data source for categories (order-independent)
					checkCategories(dataSourceNameCluster, "categories", []string{
						"nutanix_category_v2.cat-1",
						"nutanix_category_v2.cat-2",
						"nutanix_category_v2.cat-3",
					}),
				),
			},
			// Step 4: Check on cluster resource for categories
			{

				Config: testAccClusterResourceAllConfig(name) + testAccClusterResourceAssociateCategoriesConfig(r),
				Check: resource.ComposeTestCheckFunc(
					// check on cluster resource for categories (order-independent)
					checkCategories(resourceNameCluster, "categories", []string{
						"nutanix_category_v2.cat-1",
						"nutanix_category_v2.cat-2",
						"nutanix_category_v2.cat-3",
					}),
				),
			},
			// Step 5: Disassociate categories from cluster
			{
				Config: testAccClusterResourceAllConfig(name),
			},
			// Step 6: Check if categories are disassociated from cluster, data source check for categories
			{
				// Check if categories are disassociated from cluster
				Config: testAccClusterResourceAllConfig(name) + `
					# List all cluster to tests categories
					data "nutanix_clusters_v2" "list-cluster" {
						filter = "name eq '${nutanix_cluster_v2.test.name}'"
					}

					# get the cluster data source to test categories
					data "nutanix_cluster_v2" "get-cluster" {
						ext_id = nutanix_cluster_v2.test.id
					}

				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_cluster_v2.get-cluster", "categories.#", "0"),
					resource.TestCheckResourceAttr("data.nutanix_clusters_v2.list-cluster", "cluster_entities.0.categories.#", "0"),
				),
			},
			// Step 7: Check if categories are disassociated from cluster, resource check for categories
			{
				Config: testAccClusterResourceAllConfig(name),
				Taint:  []string{resourceNameCluster},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCluster, "categories.#", "0"),
				),
			},
			// ############################################## Update cluster config ##############################################
			// Step 8: Update cluster config and check on cluster resource for config
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
			// Step 9: Disable the cluster pulse status and check on cluster resource for config
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

func TestAccV2NutanixClusterResource_ExpandCluster(t *testing.T) {
	if testVars.Clusters.Nodes[1].CvmIP == "" &&
		testVars.Clusters.Nodes[2].CvmIP == "" &&
		testVars.Clusters.Nodes[3].CvmIP == "" {
		t.Skip("Skipping test as No available nodes to be used for testing")
	}

	r := acctest.RandIntRange(1, 10000)
	clusterName := fmt.Sprintf("tf-3node-cluster-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// step 1: create cluster with 3 nodes
			{
				PreConfig: func() {
					fmt.Println("Step 1: Creating a cluster with 3 nodes")
				},
				Config: testAcc3NodeClustersConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					// add node to cluster check
					// unconfigured Nodes check
					resource.TestCheckResourceAttr(resourceNameDiscoverUnconfiguredClusterNodes, "unconfigured_nodes.#", "3"),
				),
			},

			// step 2: register a cluster to pc
			{
				PreConfig: func() {
					fmt.Println("Step 2: Registering the cluster to prism central")
				},
				Config: testAcc3NodeClustersConfig(clusterName) + clusterRegistrationConfig(),
				Check: resource.ComposeTestCheckFunc(
					//Cluster Check
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "nodes.0.node_list.#", "3"),
					checkNodesIPs([]string{
						testVars.Clusters.Nodes[0].CvmIP,
						testVars.Clusters.Nodes[1].CvmIP,
						testVars.Clusters.Nodes[2].CvmIP,
					}),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_function.0", testVars.Clusters.Config.ClusterFunctions[0]),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
				),
			},

			// step 3: expand cluster by adding 4th node, and updating the name
			{
				PreConfig: func() {
					fmt.Println("Step 3: Expanding the cluster by adding 4th node")
				},
				Config: testAccExpandClustersConfig(clusterName + "_add_node"),
				Check: resource.ComposeTestCheckFunc(
					//Cluster Check
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "name", clusterName+"_add_node"),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "nodes.0.node_list.#", "4"),
					checkNodesIPs([]string{
						testVars.Clusters.Nodes[0].CvmIP,
						testVars.Clusters.Nodes[1].CvmIP,
						testVars.Clusters.Nodes[2].CvmIP,
						testVars.Clusters.Nodes[3].CvmIP,
					}),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_function.0", testVars.Clusters.Config.ClusterFunctions[0]),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
				),
			},

			// step 4: remove node from cluster by reducing to 3 nodes and updating the name
			{
				PreConfig: func() {
					fmt.Println("Step 4: Removing a node from the cluster")
					t.Log("Sleeping for 10 Minute before removing the node")
					time.Sleep(10 * time.Minute)
				},
				Config: testAccRemoveNodeClustersConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					//Cluster Check
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "nodes.0.node_list.#", "3"),
					checkNodesIPs([]string{
						testVars.Clusters.Nodes[0].CvmIP,
						testVars.Clusters.Nodes[2].CvmIP,
						testVars.Clusters.Nodes[3].CvmIP,
					}),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_function.0", testVars.Clusters.Config.ClusterFunctions[0]),
					resource.TestCheckResourceAttr(resourceName3NodesCluster, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
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

func testAccClusterResourceMinimumConfig(name string, categories string) string {
	return fmt.Sprintf(`
		# cluster config
		%[1]s

		# create a new category
		resource "nutanix_category_v2" "cat-1" {
			key         = "test-cat1-key-%[2]s"
			value       = "test-cat1-value-%[2]s"
			description = "first category for cluster"
			# Delay 5 minutes before destroying the resource to make sure that synced data is deleted
			provisioner "local-exec" {
				command    = "sleep 300"
				when       = destroy
				on_failure = continue
			}
		}

		resource "nutanix_category_v2" "cat-2" {
			key         = "test-cat2-key-%[2]s"
			value       = "test-cat2-value-%[2]s"
			description = "second category for cluster"
			# Delay 5 minutes before destroying the resource to make sure that synced data is deleted
			provisioner "local-exec" {
				command    = "sleep 300"
				when       = destroy
				on_failure = continue
			}
		}

		resource "nutanix_category_v2" "cat-3" {
			key         = "test-cat3-key-%[2]s"
			value       = "test-cat3-value-%[2]s"
			description = "third category for cluster"
			# Delay 5 minutes before destroying the resource to make sure that synced data is deleted
			provisioner "local-exec" {
				command    = "sleep 300"
				when       = destroy
				on_failure = continue
			}
		}

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

		  # associate categories to the cluster
		  categories = [%[3]s]
		  provisioner "local-exec" {
			command = "ssh-keygen -f '~/.ssh/known_hosts' -R '${local.clusters.nodes[0].cvm_ip}';  sshpass -p '${local.clusters.pe_password}' ssh -o StrictHostKeyChecking=no ${local.clusters.pe_username}@${local.clusters.nodes[0].cvm_ip} '/home/nutanix/prism/cli/ncli user reset-password user-name=${local.clusters.nodes[0].username} password=${local.clusters.nodes[0].password}' "

			on_failure = continue
		  }
          # Set lifecycle to ignore changes
		  lifecycle {
			ignore_changes = [network.0.smtp_server.0.server.0.password,  links, config.0.cluster_function]
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
`, clusterConfig, name, categories)
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


		# create a new category
		resource "nutanix_category_v2" "test" {
			key         = "%[2]s-key"
			value       = "%[2]s-value"
			description = "test-cat-cluster-description"
			provisioner "local-exec" {
				command = "sleep 120"
				when = destroy
				on_failure = continue
			}
		  lifecycle {
			ignore_changes = [key,  value]
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
		  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node, nutanix_category_v2.test]
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

		# create a new category
		resource "nutanix_category_v2" "test" {
			key         = "%[2]s-key"
			value       = "%[2]s-value"
			description = "test-cat-cluster-description"
			provisioner "local-exec" {
				command = "sleep 120"
				when = destroy
				on_failure = continue
			}
		  lifecycle {
			ignore_changes = [key,  value]
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
		  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node, nutanix_category_v2.test]
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

func testAccClusterResourceAssociateCategoriesConfig(r int) string {
	return fmt.Sprintf(`
		# create a new category
		resource "nutanix_category_v2" "cat-1" {
			key         = "test-cat1-key-%[1]d"
			value       = "test-cat1-value-%[1]d"
			description = "first category for cluster"
		}

		resource "nutanix_category_v2" "cat-2" {
			key         = "test-cat2-key-%[1]d"
			value       = "test-cat2-value-%[1]d"
			description = "second category for cluster"
		}

		resource "nutanix_category_v2" "cat-3" {
			key         = "test-cat3-key-%[1]d"
			value       = "test-cat3-value-%[1]d"
			description = "third category for cluster"
		}

		# associate categories with cluster
		resource "nutanix_cluster_categories_v2" "test" {
			cluster_ext_id = nutanix_cluster_v2.test.id
			categories = [nutanix_category_v2.cat-1.id, nutanix_category_v2.cat-2.id, nutanix_category_v2.cat-3.id]
		}

		# List all cluster to tests categories
		data "nutanix_clusters_v2" "test" {
			filter = "name eq '${nutanix_cluster_v2.test.name}'"
			depends_on = [nutanix_cluster_categories_v2.test]
		}

		# get the cluster data source to test categories
		data "nutanix_cluster_v2" "test" {
			ext_id = nutanix_cluster_v2.test.id
			depends_on = [nutanix_cluster_categories_v2.test]
		}
	`, r)
}

func testAcc3NodeClustersConfig(clusterName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {}

locals {
  pc_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]
  config   = (jsondecode(file("%[2]s")))
  clusters = local.config.clusters
}


############################ cluster with 3 nodes

## check if the nodes is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-nodes" {
  ext_id       = local.pc_ext_id
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
    ignore_changes = [links, categories, config.0.cluster_function]
  }

  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.cluster-nodes]
}
`, clusterName, filepath)
}

func clusterRegistrationConfig() string {
	return `

## we need to register on of 3 nodes cluster to pc
resource "nutanix_pc_registration_v2" "nodes-registration" {
  pc_ext_id = local.pc_ext_id
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
  `
}

func testAccExpandClustersConfig(clusterName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {}

locals {
  pc_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]
  config   = (jsondecode(file("%[2]s")))
  clusters = local.config.clusters
}


############################ cluster with 3 nodes

## check if the nodes is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-nodes" {
  ext_id       = local.pc_ext_id
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
    node_list {
      controller_vm_ip {
        ipv4 {
          value = local.clusters.nodes[3].cvm_ip
        }
      }
      should_skip_host_networking   = false
      should_skip_pre_expand_checks = true
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
    ignore_changes = [links, categories, config.0.cluster_function]
  }

  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.cluster-nodes]
}
`, clusterName, filepath)
}

func testAccRemoveNodeClustersConfig(clusterName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {}

locals {
  pc_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]
  config   = (jsondecode(file("%[2]s")))
  clusters = local.config.clusters
}


############################ cluster with 3 nodes

## check if the nodes is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-nodes" {
  ext_id       = local.pc_ext_id
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
          value = local.clusters.nodes[2].cvm_ip
        }
      }
    }
    node_list {
      controller_vm_ip {
        ipv4 {
          value = local.clusters.nodes[3].cvm_ip
        }
      }
      should_skip_host_networking   = false
      should_skip_pre_expand_checks = true
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
    ignore_changes = [links, categories, config.0.cluster_function]
  }

  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.cluster-nodes]
}
`, clusterName, filepath)
}
