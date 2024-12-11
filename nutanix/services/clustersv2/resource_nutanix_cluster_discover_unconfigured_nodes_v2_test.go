package clustersv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameDiscoverUnconfiguredNodes = "nutanix_clusters_discover_unconfigured_nodes_v2.test"

func TestAccV2NutanixClusterDiscoverUnconfiguredNodesResource_basic(t *testing.T) {
	if testVars.Clusters.Nodes[0].CvmIP == "" {
		t.Skip("Skipping test as No available node to be used for testing")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClustersDataSourceDiscoverUnconfiguredNodesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDiscoverUnconfiguredNodes, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameDiscoverUnconfiguredNodes, "unconfigured_nodes.#"),
					resource.TestCheckResourceAttr(resourceNameDiscoverUnconfiguredNodes, "unconfigured_nodes.0.cvm_ip.0.ipv4.0.value", testVars.Clusters.Nodes[0].CvmIP),
				),
			},
		},
	})
}

func testAccClustersDataSourceDiscoverUnconfiguredNodesConfig() string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}
	locals {
	  cluster_ext_id = [
		for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
	  ][0]
	  config   = (jsondecode(file("%[1]s")))
	  clusters = local.config.clusters
	}

	resource "nutanix_clusters_discover_unconfigured_nodes_v2" "test" {
	  ext_id       = local.cluster_ext_id
	  address_type = "IPV4"
	  ip_filter_list {
		ipv4 {
		  value = local.clusters.nodes[0].cvm_ip
		}
	  }
      depends_on = [data.nutanix_clusters_v2.clusters]
	}
`, filepath)
}
