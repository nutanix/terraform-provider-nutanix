package clustersv2_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	clusterResourceName                    = "nutanix_cluster_v2.test"
	resourceNameDiscoverUnConfigNode       = "nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node"
	clusterResourceNameRegistration        = "nutanix_pc_registration_v2.node-registration"
	dataSourceNameClusterData              = "data.nutanix_cluster_v2.cluster"
	dataSourceNameGetClusterCategoriesData = "data.nutanix_clusters_v2.get-cluster-categories"
)

func TestAccV2NutanixClusterResource_CreateClusterWithMinimumConfig(t *testing.T) {
	if testVars.Clusters.Nodes[0].CvmIP == "" {
		t.Skip("Skipping test as No available node to be used for testing")
	}
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-cluster-%d", r)

	clusterProfileResourceName := "nutanix_cluster_profile_v2.test"
	clusterProfileDataSourceName := "data.nutanix_cluster_profile_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckNutanixClusterDestroy,
			testAccCheckClusterProfileDestroy,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Step 1: Plan the cluster with minimum config")
				},
				Config:   testAccClusterResourceMinimumConfig(name, "", ""),
				PlanOnly: false,
			},
			{
				PreConfig: func() {
					fmt.Println("Step 2: Create the cluster with minimum config")
				},
				Config: testAccClusterResourceMinimumConfig(name, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceName, "name", name),
					resource.TestCheckResourceAttr(clusterResourceName, "dryrun", "false"),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
				),
			},
			// Create cluster profile and associate with cluster
			{
				PreConfig: func() {
					fmt.Println("Step 3: Associating the cluster profile with the cluster")
				},
				Config: testAccClusterResourceMinimumConfig(name, "cluster_profile_ext_id = nutanix_cluster_profile_v2.test.id", "nutanix_category_v2.cat-1.id, nutanix_category_v2.cat-2.id, nutanix_category_v2.cat-3.id") +
					testAccClusterProfileResourceConfig("tf-first-cluster-profile") +
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
					resource.TestCheckResourceAttr(clusterResourceName, "categories.#", "3"),
					// Check categories on the resource (order-independent check)
					checkCategories(clusterResourceName, "categories", []string{
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

					resource.TestCheckResourceAttr(clusterProfileResourceName, "name", "tf-first-cluster-profile"),
					resource.TestCheckResourceAttrSet(clusterProfileResourceName, "ext_id"),
					resource.TestCheckResourceAttrPair(clusterProfileResourceName, "id", clusterProfileDataSourceName, "ext_id"),
					resource.TestCheckResourceAttrPair(clusterProfileResourceName, "name", clusterProfileDataSourceName, "name"),
				),
			},
			// Get Cluster Profile Details
			{
				PreConfig: func() {
					fmt.Println("Step 5: Getting the cluster profile details")
				},
				Config: testAccClusterResourceMinimumConfig(name, "cluster_profile_ext_id = nutanix_cluster_profile_v2.test.id", "") + testAccClusterProfileResourceConfig("tf-first-cluster-profile"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterProfileResourceName, "name", "tf-first-cluster-profile"),
					resource.TestCheckResourceAttr(clusterProfileDataSourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttrPair(clusterProfileDataSourceName, "clusters.0.ext_id", clusterResourceName, "id"),
				),
			},
			// de-associate the cluster profile from the cluster
			{
				PreConfig: func() {
					fmt.Println("Step 4: De-associating the cluster profile from the cluster")
				},
				Config: testAccClusterResourceMinimumConfig(name, "", "") + testAccClusterProfileResourceConfig("tf-first-cluster-profile") +
					`				
					# get the cluster data
					data "nutanix_cluster_v2" "cluster" {
						ext_id = nutanix_cluster_v2.test.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceName, "cluster_profile_ext_id", ""),
					resource.TestCheckResourceAttr(clusterResourceName, "categories.#", "0"),
					resource.TestCheckResourceAttr(dataSourceNameClusterData, "categories.#", "0"),
					resource.TestCheckResourceAttr(dataSourceNameGetClusterCategoriesData, "cluster_entities.0.categories.#", "0"),
				),
			},
			// Get Cluster Profile Details after de-association
			{
				PreConfig: func() {
					fmt.Println("Step 6: Getting the cluster profile details after de-association")
				},
				Config: testAccClusterResourceMinimumConfig(name, "", "") + testAccClusterProfileResourceConfig("tf-first-cluster-profile") +
					`				
					# get the cluster data
					data "nutanix_cluster_v2" "cluster" {
						ext_id = nutanix_cluster_v2.test.id
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterProfileDataSourceName, "clusters.#", "0"),
					resource.TestCheckResourceAttr(clusterResourceName, "cluster_profile_ext_id", ""),
					resource.TestCheckResourceAttr(clusterResourceName, "categories.#", "0"),
					resource.TestCheckResourceAttr(dataSourceNameClusterData, "categories.#", "0"),
					resource.TestCheckResourceAttr(dataSourceNameGetClusterCategoriesData, "cluster_entities.0.categories.#", "0"),
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
					resource.TestCheckResourceAttr(clusterResourceName, "name", name),
					resource.TestCheckResourceAttr(clusterResourceName, "dryrun", "false"),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.external_address.0.ipv4.0.value", testVars.Clusters.Network.VirtualIP),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.0.fqdn.0.value", testVars.Clusters.Network.NTPServers[0]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.1.fqdn.0.value", testVars.Clusters.Network.NTPServers[1]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.2.fqdn.0.value", testVars.Clusters.Network.NTPServers[2]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.3.fqdn.0.value", testVars.Clusters.Network.NTPServers[3]),

					resource.TestCheckResourceAttrSet(clusterResourceNameRegistration, "pc_ext_id"),
					resource.TestCheckResourceAttr(clusterResourceNameRegistration, "remote_cluster.0.aos_remote_cluster_spec.0.remote_cluster.0.address.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(clusterResourceNameRegistration, "remote_cluster.0.aos_remote_cluster_spec.0.remote_cluster.0.credentials.0.authentication.0.username", testVars.Clusters.Nodes[0].Username),
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
					checkCategories(clusterResourceName, "categories", []string{
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
				Taint:  []string{clusterResourceName},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceName, "categories.#", "0"),
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
					resource.TestCheckResourceAttr(clusterResourceName, "name", name+"-updated"),
					resource.TestCheckResourceAttr(clusterResourceName, "dryrun", "false"),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.pulse_status.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.pulse_status.0.pii_scrubbing_level", "DEFAULT"),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.external_address.0.ipv4.0.value", testVars.Clusters.Network.VirtualIP),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.external_data_services_ip.0.ipv4.0.value", testVars.Clusters.Network.IscsiIP),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.0.fqdn.0.value", testVars.Clusters.Network.NTPServers[0]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.1.fqdn.0.value", testVars.Clusters.Network.NTPServers[1]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.2.fqdn.0.value", testVars.Clusters.Network.NTPServers[2]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.3.fqdn.0.value", testVars.Clusters.Network.NTPServers[3]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.email_address", testVars.Clusters.Network.SMTPServer.EmailAddress),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.server.0.ip_address.0.ipv4.0.value", testVars.Clusters.Network.SMTPServer.IP),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.server.0.port", strconv.Itoa(testVars.Clusters.Network.SMTPServer.Port)),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.server.0.username", testVars.Clusters.Network.SMTPServer.Username),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.type", testVars.Clusters.Network.SMTPServer.Type),

					// check on list cluster data source for categories
					resource.TestCheckResourceAttr(dataSourceNameClusters, "cluster_entities.0.categories.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceNameClusters, "cluster_entities.0.categories.0", "nutanix_category_v2.test", "id"),
				),
			},
			// Step 9: Disable the cluster pulse status and check on cluster resource for config
			{
				PreConfig: func() {
					time.Sleep(10 * time.Second) // 10-second delay
				},
				Config: testAccClusterResourceUpdateConfig(name+"-updated", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterResourceName, "name", name+"-updated"),
					resource.TestCheckResourceAttr(clusterResourceName, "dryrun", "false"),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.node_list.0.controller_vm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
					resource.TestCheckResourceAttr(clusterResourceName, "nodes.0.number_of_nodes", "1"),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.cluster_arch", testVars.Clusters.Config.ClusterArch),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.fault_tolerance_state.0.domain_awareness_level", testVars.Clusters.Config.FaultToleranceState.DomainAwarenessLevel),
					resource.TestCheckResourceAttr(clusterResourceName, "config.0.pulse_status.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.external_address.0.ipv4.0.value", testVars.Clusters.Network.VirtualIP),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.external_data_services_ip.0.ipv4.0.value", testVars.Clusters.Network.IscsiIP),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.0.fqdn.0.value", testVars.Clusters.Network.NTPServers[0]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.1.fqdn.0.value", testVars.Clusters.Network.NTPServers[1]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.2.fqdn.0.value", testVars.Clusters.Network.NTPServers[2]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.ntp_server_ip_list.3.fqdn.0.value", testVars.Clusters.Network.NTPServers[3]),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.email_address", testVars.Clusters.Network.SMTPServer.EmailAddress),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.server.0.ip_address.0.ipv4.0.value", testVars.Clusters.Network.SMTPServer.IP),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.server.0.port", strconv.Itoa(testVars.Clusters.Network.SMTPServer.Port)),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.server.0.username", testVars.Clusters.Network.SMTPServer.Username),
					resource.TestCheckResourceAttr(clusterResourceName, "network.0.smtp_server.0.type", testVars.Clusters.Network.SMTPServer.Type),
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
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixClusterDestroy,
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

func testAccClusterResourceMinimumConfig(name, clusterProfileExtID string, categories string) string {
	// Build the cluster_profile_ext_id line - always include it to ensure Terraform detects changes
	var clusterProfileExtIDLine string
	if clusterProfileExtID == "" {
		// Explicitly set to empty string to ensure Terraform detects the change from non-empty to empty
		clusterProfileExtIDLine = "cluster_profile_ext_id = \"\""
	} else if strings.HasPrefix(clusterProfileExtID, "cluster_profile_ext_id =") {
		// Already a full line, use as-is
		clusterProfileExtIDLine = clusterProfileExtID
	} else {
		// Just the value part (e.g., "nutanix_cluster_profile_v2.test.id")
		clusterProfileExtIDLine = fmt.Sprintf("cluster_profile_ext_id = %s", clusterProfileExtID)
	}

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

		  %[4]s

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
`, clusterConfig, name, categories, clusterProfileExtIDLine)
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

func testAccClusterProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_cluster_profile_v2" "test" {
  name = "%s"
  description = "Example First Cluster Profile created via Terraform"
  allowed_overrides = ["NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"]

  name_server_ip_list {
    ipv4 { value = "240.29.254.180" }
    ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
  }

  ntp_server_ip_list {
    ipv4 { value = "240.29.254.180" }
    ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
    fqdn { value = "ntp.example.com" }
  }

  smtp_server {
    email_address = "email@example.com"
    type = "SSL"
    server {
      ip_address {
        ipv4 { value = "240.29.254.180" }
        ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
        fqdn { value = "smtp.example.com" }
      }
      port     = 587
      username = "example_user"
      password = "example_password"
    }
  }

  nfs_subnet_white_list = ["10.110.106.45/255.255.255.255"]

  snmp_config {
    is_enabled = true
    users {
      username  = "snmpuser1"
      auth_type = "MD5"
      auth_key  = "Test_SNMP_user_authentication_key"
      priv_type = "DES"
      priv_key  = "Test_SNMP_user_encryption_key"
    }
    transports {
      protocol = "UDP"
      port     = 21
    }
    traps {
      address {
        ipv4 {
					value         = "240.29.254.180"
					prefix_length = 24
				}
        ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
      }
      username         = "trapuser"
      protocol         = "UDP"
      port             = 59
      should_inform    = false
      engine_id        = "0x1234567890abcdef12"
      version          = "V2"
      receiver_name    = "trap-receiver"
      community_string = "snmp-server community public RO 192.168.1.0 255.255.255.0"
    }
  }

  rsyslog_server_list {
    server_name      = "testServer1"
    port             = 29
    network_protocol = "UDP"
    ip_address {
      ipv4 { value = "240.29.254.180" }
      ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
    }
    modules {
      name                     = "CASSANDRA"
      log_severity_level       = "EMERGENCY"
      should_log_monitor_files = true
    }
    modules {
      name                     = "CURATOR"
      log_severity_level       = "ERROR"
      should_log_monitor_files = false
    }
  }

  pulse_status {
    is_enabled          = false
    pii_scrubbing_level = "DEFAULT"
  }

  lifecycle {
    ignore_changes = [
      smtp_server.0.server.0.password,
      snmp_config.0.users.0.auth_key,
      snmp_config.0.users.0.priv_key
    ]
  }
}

data "nutanix_cluster_profile_v2" "test" {
  ext_id = nutanix_cluster_profile_v2.test.id
}
`, name)
}
