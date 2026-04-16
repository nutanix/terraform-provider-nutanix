package clustersv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameCluster = "data.nutanix_cluster_v2.test"

func TestAccV2NutanixClusterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameCluster, "config.0.cluster_function.0", "PRISM_CENTRAL"),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfig() string {
	return `
data "nutanix_clusters_v2" "test" {}

locals {
	      clusterId = [for cluster in data.nutanix_clusters_v2.test.cluster_entities: 
						cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"][0]
}
					  
data "nutanix_cluster_v2" "test" {
	ext_id = local.clusterId
}`
}
