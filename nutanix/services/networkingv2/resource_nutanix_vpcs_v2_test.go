package networkingv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVpc = "nutanix_vpc_v2.test"
const datasourceNameVpcProject = "data.nutanix_vpc_v2.project_test"
const datasourceNameVpcProjectList = "data.nutanix_vpcs_v2.list_test"

func TestAccV2NutanixVpcResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vlanID := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	desc := "test vpc description"
	updatedName := fmt.Sprintf("updated-vpc-%d", r)
	updatedDesc := "updated vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfig(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.0.category_ids.#"),
					testCheckMetadataCategoryIDsContain(resourceNameVpc, "nutanix_category_v2.test"),
				),
			},
			{
				Config: testVpcConfig(updatedName, updatedDesc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.0.category_ids.#"),
					testCheckMetadataCategoryIDsContain(resourceNameVpc, "nutanix_category_v2.test"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_WithExternallyRoutablePrefixes(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	vlanID := acctest.RandIntRange(1, 999)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithExtRoutablePrefix(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.0.category_ids.#"),
					testCheckMetadataCategoryIDsContain(resourceNameVpc, "nutanix_category_v2.test"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
				),
			},
		},
	})
}

// TestAccV2NutanixVpcResource_WithMultipleExternallyRoutablePrefixes tests the fix for issue
// where multiple ipv4 blocks within a single externally_routable_prefixes block were not
// being properly sent to the API - only the first one was processed.
func TestAccV2NutanixVpcResource_WithMultipleExternallyRoutablePrefixes(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vpc-multi-prefix-%d", r)
	vlanID := acctest.RandIntRange(1, 999)
	desc := "test vpc with multiple externally routable prefixes"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithMultipleExtRoutablePrefixes(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					// Verify that all 3 ipv4 prefixes are configured
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.#", "3"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.0.ip.0.value", "172.16.3.0"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.1.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.1.ip.0.value", "172.16.4.0"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.2.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.2.ip.0.value", "192.168.120.0"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
				),
			},
			{
				// Remove one ipv4 block: update from 3 to 2 prefixes
				Config: testVpcConfigWithTwoExtRoutablePrefixes(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.0.ip.0.value", "172.16.3.0"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.1.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefixes.0.ipv4.1.ip.0.value", "172.16.4.0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_WithDHCP(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	vlanID := acctest.RandIntRange(1, 999)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithDHCP(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "common_dhcp_options.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_WithTransitType(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vpc-%d", r)
	vlanID := acctest.RandIntRange(1, 999)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcConfigWithTransitType(name, desc, vlanID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "TRANSIT"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "links.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "snat_ips.#"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "common_dhcp_options.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_DefaultProjectAndSharing(t *testing.T) {
	r := acctest.RandInt()
	vlanID := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-test-vpc-proj-%d", r)
	desc := "vpc for project sharing test"
	updatedDesc := "vpc updated with sharing"
	unsharedDesc := "vpc after unsharing"
	projectName := fmt.Sprintf("tf-vpc-share-proj-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcProjectConfig(name, desc, vlanID, projectName, "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVpc, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(datasourceNameVpcProject, "project_ext_id", defaultProjectUUID),
				),
			},
			{
				Config: testVpcProjectConfig(name, updatedDesc, vlanID, projectName, "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVpc, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameVpc, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVpcProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameVpcProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVpcProjectList, "vpcs.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameVpcProjectList, "vpcs.0.ext_id", resourceNameVpc, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVpcProjectList, "vpcs.0.shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameVpcProjectList, "vpcs.0.shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
				),
			},
			{
				Config: testVpcProjectConfig(name, unsharedDesc, vlanID, projectName, "empty"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "description", unsharedDesc),
					resource.TestCheckResourceAttr(resourceNameVpc, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameVpcProject, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameVpcProjectList, "vpcs.#", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcResource_NonDefaultProjectSharingFails(t *testing.T) {
	r := acctest.RandInt()
	vlanID := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-test-vpc-nondfl-%d", r)
	desc := "vpc in non-default project"
	project1Name := fmt.Sprintf("tf-vpc-proj1-%d", r)
	project2Name := fmt.Sprintf("tf-vpc-proj2-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVpcNonDefaultProjectConfig(name, desc, vlanID, project1Name, project2Name, "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVpc, "ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameVpc, "project_ext_id", "nutanix_project_v2.project1", "ext_id"),
				),
			},
			{
				Config:      testVpcNonDefaultProjectConfig(name, desc, vlanID, project1Name, project2Name, "share"),
				ExpectError: regexp.MustCompile("Can only share VPCs owned by the default project"),
			},
		},
	})
}

// TestAccV2NutanixVpcResource_ScopeAdvertiseAndExtSubnetType covers the VPC
// attributes added in "Networking: VPC: Add support for VPC association with
// kubernetes clusters": scope and should_advertise_connected_subnets. Both are
// plain VPC-level fields, so a REGULAR VPC with a single external subnet
// exercises them without extra infrastructure. The test creates the VPC
// (scope=VMS / advertise=true) and then updates the advertise flag (-> false)
// to cover the create-expand and update paths.
//
// Notes:
//   - scope stays VMS in both steps. VMS_AND_CONTAINERS and CONTAINERS require a
//     non-empty kubernetes_clusters list (server rejects otherwise with
//     NETWORKING-10069 "Kubernetes cluster list cannot be empty if VPC scope is
//     VMS_AND_CONTAINERS"), and those cluster ext_ids cannot be provisioned in a
//     generic acceptance environment. So the update path is driven by
//     should_advertise_connected_subnets instead.
//   - should_advertise_connected_subnets is created with true (not false)
//     on purpose: the resource expands it via d.GetOk, which treats a bool
//     false as "unset" and never sends it, so a false value at create time is
//     not actually exercised. Creating true -> updating false covers both the
//     create (GetOk) and update (d.HasChange) paths.
//   - supported_multiple_external_subnet_type is left unset in config and only
//     asserted as "set": it is Optional+Computed and the server computes/
//     overrides it (e.g. a v4.4 PC returns ONLY_NONAT regardless of the
//     requested value), so pinning a literal would be brittle and cause a
//     perpetual diff.
//   - The kubernetes_clusters block is intentionally not covered here for the
//     same NKE-provisioning reason noted above.
func TestAccV2NutanixVpcResource_ScopeAdvertiseAndExtSubnetType(t *testing.T) {
	r := acctest.RandInt()
	vlanID := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-test-vpc-attrs-%d", r)
	desc := "vpc scope/advertise test"
	updatedDesc := desc + " updated"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create the VPC with the new attributes set to their initial values.
			{
				Config: testVpcConfigWithScopeAdvertise(name, desc, vlanID, "VMS", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttr(resourceNameVpc, "scope", "VMS"),
					resource.TestCheckResourceAttr(resourceNameVpc, "should_advertise_connected_subnets", "true"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "supported_multiple_external_subnet_type"),
				),
			},
			// Update: flip the advertise flag (scope stays VMS - see notes above).
			{
				Config: testVpcConfigWithScopeAdvertise(name, updatedDesc, vlanID, "VMS", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttr(resourceNameVpc, "scope", "VMS"),
					resource.TestCheckResourceAttr(resourceNameVpc, "should_advertise_connected_subnets", "false"),
					resource.TestCheckResourceAttrSet(resourceNameVpc, "supported_multiple_external_subnet_type"),
				),
			},
		},
	})
}

// TestAccV2NutanixVpcResource_WithKubernetesClusters covers the
// kubernetes_clusters association block on the VPC v2 resource (ext_id,
// gateway_nodes_selector, namespace_selector and pod_network). Associating a
// Kubernetes cluster requires an NKP/Karbon cluster that is already registered
// with Prism Central; its extId is read from
// test_config_v2.json -> networking.kubernetes_cluster_ext_id (populate it e.g.
// via testenv/e2e_nkp_deploy.py). The test is skipped when that fixture is empty
// so the suite still runs in environments without a registered cluster.
//
// Server-side scope validations (confirmed against the platform) drive the
// config choices here:
//   - A non-empty kubernetes_clusters list requires scope != VMS.
//   - When BOTH namespace_selector and gateway_nodes_selector are set, scope
//     MUST be VMS_AND_CONTAINERS (a gateway_nodes_selector also forbids
//     scope=CONTAINERS). So this test uses VMS_AND_CONTAINERS with both
//     selectors to exercise the richest expand/flatten path.
//
// Step 1 creates the association; step 2 updates the namespace label value to
// exercise the update (d.HasChange) path. NOTE: the associated cluster's pod
// CIDR (including host_slice) is immutable server-side, so host_slice is kept
// constant across steps -- changing it fails with "pod CIDR is immutable".
func TestAccV2NutanixVpcResource_WithKubernetesClusters(t *testing.T) {
	k8sExtID := testVars.Networking.KubernetesClusterExtID
	if k8sExtID == "" {
		t.Skip("networking.kubernetes_cluster_ext_id is not set in test_config_v2.json; " +
			"skipping VPC kubernetes_clusters association test")
	}

	r := acctest.RandInt()
	vlanID := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-test-vpc-k8s-%d", r)
	desc := "vpc with kubernetes cluster association"
	updatedDesc := desc + " updated"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create the VPC associated with the Kubernetes cluster.
			{
				Config: testVpcConfigWithKubernetesClusters(name, desc, vlanID, k8sExtID, "ns1", 24),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "name", name),
					resource.TestCheckResourceAttr(resourceNameVpc, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVpc, "vpc_type", "REGULAR"),
					resource.TestCheckResourceAttr(resourceNameVpc, "scope", "VMS_AND_CONTAINERS"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.ext_id", k8sExtID),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.gateway_nodes_selector.0.match_labels.0.name", "env"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.gateway_nodes_selector.0.match_labels.0.value", "prod"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.namespace_selector.0.match_labels.0.name", "nutanix.com/vpc-namespace"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.namespace_selector.0.match_labels.0.value", "ns1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.pod_network.0.cidr.0.ipv4.0.ip.0.value", "172.20.0.0"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.pod_network.0.cidr.0.ipv4.0.prefix_length", "16"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.pod_network.0.host_slice", "24"),
				),
			},
			// Update: change the namespace label value and the pod host_slice.
			{
				Config: testVpcConfigWithKubernetesClusters(name, updatedDesc, vlanID, k8sExtID, "ns2", 24),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpc, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVpc, "scope", "VMS_AND_CONTAINERS"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.ext_id", k8sExtID),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.namespace_selector.0.match_labels.0.value", "ns2"),
					resource.TestCheckResourceAttr(resourceNameVpc, "kubernetes_clusters.0.pod_network.0.host_slice", "24"),
				),
			},
		},
	})
}

func testVpcConfig(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name = "terraform-test-subnet-vpc"
		description = "test subnet description"
		cluster_reference = local.cluster0
		subnet_type = "VLAN"
		network_id = %[3]d
		is_external = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list{
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_category_v2" "test" {
		key = "tf-test-category-key-%[3]d"
		value = "tf-test-category-value-%[3]d"
		description = "test category for vpc"
	}

	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = nutanix_subnet_v2.test.id
		}
		metadata {
			category_ids = [nutanix_category_v2.test.id]
		}
		vpc_type = "REGULAR"
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID)
}

func testVpcConfigWithExtRoutablePrefix(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name = "terraform-test-subnet-vpc"
		description = "test subnet description"
		cluster_reference = local.cluster0
		subnet_type = "VLAN"
		network_id = %[3]d
		is_external = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list{
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_category_v2" "test" {
		key = "tf-test-category-key-%[1]s"
		value = "tf-test-category-value-%[1]s"
		description = "test category for vpc"
	}

	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
			external_ips {
			  ipv4 {
				value         = "192.168.0.24"
				prefix_length = 32
			  }
			}
			external_ips {
			  ipv4 {
				value         = "192.168.0.25"
				prefix_length = 32
			  }
			}
	   	}
		externally_routable_prefixes{
		  ipv4{
			ip{
			  value = "172.30.0.0"
			  prefix_length = 32
			}
			prefix_length = 16
		  }
		}
		metadata {
			category_ids = [nutanix_category_v2.test.id]
		}
		vpc_type = "REGULAR"
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID)
}

// testVpcConfigWithMultipleExtRoutablePrefixes tests the fix for the bug where
// multiple ipv4 blocks within a single externally_routable_prefixes were not all
// being sent to the API - only the first one was processed.
func testVpcConfigWithMultipleExtRoutablePrefixes(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name              = "terraform-test-subnet-vpc-multi"
		description       = "test subnet description"
		cluster_reference = local.cluster0
		subnet_type       = "VLAN"
		network_id        = %[3]d
		is_external       = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list {
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_vpc_v2" "test" {
		name        = "%[1]s"
		description = "%[2]s"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		# Test case: multiple ipv4 blocks within a single externally_routable_prefixes block
		# All 3 prefixes should be configured on the VPC
		externally_routable_prefixes {
			ipv4 {
				ip {
					value         = "172.16.3.0"
					prefix_length = 32
				}
				prefix_length = 24
			}
			ipv4 {
				ip {
					value         = "172.16.4.0"
					prefix_length = 32
				}
				prefix_length = 24
			}
			ipv4 {
				ip {
					value         = "192.168.120.0"
					prefix_length = 32
				}
				prefix_length = 24
			}
		}
		vpc_type = "REGULAR"
		depends_on = [nutanix_subnet_v2.test]
	}

`, name, desc, vlanID)
}

// testVpcConfigWithTwoExtRoutablePrefixes is the same as testVpcConfigWithMultipleExtRoutablePrefixes
// but with only 2 ipv4 blocks, used to test removal of one block (update from 3 to 2).
func testVpcConfigWithTwoExtRoutablePrefixes(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name              = "terraform-test-subnet-vpc-multi"
		description       = "test subnet description"
		cluster_reference = local.cluster0
		subnet_type       = "VLAN"
		network_id        = %[3]d
		is_external       = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list {
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_vpc_v2" "test" {
		name        = "%[1]s"
		description = "%[2]s"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		externally_routable_prefixes {
			ipv4 {
				ip {
					value         = "172.16.3.0"
					prefix_length = 32
				}
				prefix_length = 24
			}
			ipv4 {
				ip {
					value         = "172.16.4.0"
					prefix_length = 32
				}
				prefix_length = 24
			}
		}
		vpc_type = "REGULAR"
		depends_on = [nutanix_subnet_v2.test]
	}

`, name, desc, vlanID)
}

func testVpcConfigWithDHCP(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
	 	name              = "terraform-test-subnet-vpc"
	  	description       = "test subnet description"
		  cluster_reference = local.cluster0
		  subnet_type       = "VLAN"
		  network_id        = %[3]d
		  is_external       = true
		  ip_config {
			ipv4 {
			  ip_subnet {
				ip {
				  value = "192.168.0.0"
				}
				prefix_length = 24
			  }
			  default_gateway_ip {
				value = "192.168.0.1"
			  }
			  pool_list {
				start_ip {
				  value = "192.168.0.20"
				}
				end_ip {
				  value = "192.168.0.30"
				}
			  }
			}
		  }
		  depends_on = [data.nutanix_clusters_v2.clusters]
		}
		resource "nutanix_vpc_v2" "test" {
		  name        = "%[1]s"
		  description = "%[2]s"
		  external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
			external_ips {
			  ipv4 {
				value         = "192.168.0.24"
				prefix_length = 32
			  }
			}
			external_ips {
			  ipv4 {
				value         = "192.168.0.25"
				prefix_length = 32
			  }
			}
		  }

		  externally_routable_prefixes {
			ipv4 {
			  ip {
				value         = "172.30.0.0"
				prefix_length = 32
			  }
			  prefix_length = 16
			}
		  }
		  depends_on = [nutanix_subnet_v2.test]
		}

`, name, desc, vlanID)
}

func testVpcConfigWithTransitType(name, desc string, vlanID int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name = "terraform-test-subnet-vpc"
		description = "test subnet description"
		cluster_reference = local.cluster0
		subnet_type = "VLAN"
		network_id = %[3]d
		is_external = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list{
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}
	resource "nutanix_vpc_v2" "test" {
		name =  "%[1]s"
		description = "%[2]s"
		external_subnets{
		  subnet_reference = nutanix_subnet_v2.test.id
		}
		vpc_type = "TRANSIT"
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID)
}

func testVpcProjectConfig(name, desc string, vlanID int, projectName, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `shared_with_projects = [nutanix_project_v2.share_project.ext_id]`
	case "empty":
		shareBlock = `shared_with_projects = []`
	}
	listDataSource := ""
	if shareState != "none" {
		listDataSource = `
	data "nutanix_vpcs_v2" "list_test" {
		filter     = "sharedWithProjects/any(p:p eq '${nutanix_project_v2.share_project.id}')"
		depends_on = [nutanix_vpc_v2.test]
	}`
	}
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name              = "terraform-test-subnet-vpc-proj"
		description       = "test subnet for vpc project test"
		cluster_reference = local.cluster0
		subnet_type       = "VLAN"
		network_id        = %[3]d
		is_external       = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list {
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_project_v2" "share_project" {
		name       = "%[4]s"
		project_id = "%[4]s"
		description = "project for vpc sharing test"
	}

	resource "nutanix_vpc_v2" "test" {
		name        = "%[1]s"
		description = "%[2]s"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		vpc_type = "REGULAR"
		%[5]s
		depends_on = [nutanix_subnet_v2.test, nutanix_project_v2.share_project]
	}

	data "nutanix_vpc_v2" "project_test" {
		ext_id = nutanix_vpc_v2.test.id
	}
	%[6]s
`, name, desc, vlanID, projectName, shareBlock, listDataSource)
}

func testVpcNonDefaultProjectConfig(name, desc string, vlanID int, proj1, proj2, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `shared_with_projects = [nutanix_project_v2.project2.ext_id]`
	case "empty":
		shareBlock = `shared_with_projects = []`
	}
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_project_v2" "project1" {
		name       = "%[4]s"
		project_id = "%[4]s"
		description = "first project for vpc test"
	}

	resource "nutanix_project_v2" "project2" {
		name       = "%[5]s"
		project_id = "%[5]s"
		description = "second project for vpc test"
	}
	
	resource "nutanix_subnet_v2" "test" {
		name              = "terraform-test-subnet-vpc-nondfl"
		description       = "test subnet for vpc non-default project test"
		cluster_reference = local.cluster0
		subnet_type       = "VLAN"
		network_id        = %[3]d
		is_external       = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list {
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		shared_with_projects = [nutanix_project_v2.project1.ext_id]
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_vpc_v2" "test" {
		name           = "%[1]s"
		description    = "%[2]s"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		vpc_type       = "REGULAR"
		project_ext_id = nutanix_project_v2.project1.ext_id
		%[6]s
		depends_on = [nutanix_subnet_v2.test, nutanix_project_v2.project1, nutanix_project_v2.project2]
	}
`, name, desc, vlanID, proj1, proj2, shareBlock)
}

func testVpcConfigWithScopeAdvertise(name, desc string, vlanID int, scope string, advertise bool) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name              = "tf-test-subnet-vpc-attrs-%[3]d"
		description       = "external subnet for vpc attrs test"
		cluster_reference = local.cluster0
		subnet_type       = "VLAN"
		network_id        = %[3]d
		is_external       = true
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "192.168.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "192.168.0.1"
				}
				pool_list {
					start_ip {
						value = "192.168.0.20"
					}
					end_ip {
						value = "192.168.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_vpc_v2" "test" {
		name        = "%[1]s"
		description = "%[2]s"
		vpc_type    = "REGULAR"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		scope                              = "%[4]s"
		should_advertise_connected_subnets = %[5]t
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID, scope, advertise)
}

func testVpcConfigWithKubernetesClusters(name, desc string, vlanID int, k8sExtID, nsValue string, hostSlice int) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet_v2" "test" {
		name              = "tf-test-subnet-vpc-k8s-%[3]d"
		description       = "external subnet for vpc kubernetes cluster test"
		cluster_reference = local.cluster0
		subnet_type       = "VLAN"
		network_id        = %[3]d
		is_external       = true
		# External subnet must NOT overlap the associated cluster's pod CIDR
		# (default 192.168.0.0/16) or the VPC pod_network (172.20.0.0/16), or the
		# platform rejects VPC create with "Pod CIDR ... overlaps with external subnet".
		ip_config {
			ipv4 {
				ip_subnet {
					ip {
						value = "172.30.0.0"
					}
					prefix_length = 24
				}
				default_gateway_ip {
					value = "172.30.0.1"
				}
				pool_list {
					start_ip {
						value = "172.30.0.20"
					}
					end_ip {
						value = "172.30.0.30"
					}
				}
			}
		}
		depends_on = [data.nutanix_clusters_v2.clusters]
	}

	resource "nutanix_vpc_v2" "test" {
		name        = "%[1]s"
		description = "%[2]s"
		vpc_type    = "REGULAR"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		scope = "VMS_AND_CONTAINERS"
		kubernetes_clusters {
			ext_id = "%[4]s"
			gateway_nodes_selector {
				match_labels {
					name  = "env"
					value = "prod"
				}
			}
			namespace_selector {
				match_labels {
					name  = "nutanix.com/vpc-namespace"
					value = "%[5]s"
				}
			}
			pod_network {
				cidr {
					ipv4 {
						ip {
							value = "172.20.0.0"
						}
						prefix_length = 16
					}
				}
				host_slice = %[6]d
			}
		}
		depends_on = [nutanix_subnet_v2.test]
	}
`, name, desc, vlanID, k8sExtID, nsValue, hostSlice)
}
