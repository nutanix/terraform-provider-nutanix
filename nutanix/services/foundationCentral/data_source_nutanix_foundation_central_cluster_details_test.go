package fc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCClusterDetailsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCClusterDetailsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_clusters_list.cls", "imaged_clusters.#"),
				),
			},
		},
	})
}

func TestAccFCClusterDetailsDataSource_ClusterUUID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCClusterDetailsDataSourceConfigWithUUID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_clusters_list.cls", "imaged_clusters.#"),
					resource.TestCheckResourceAttr("data.nutanix_foundation_central_cluster_details.k1", "storage_node_count", "0"),
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_cluster_details.k1", "imaged_cluster_uuid"),
				),
			},
		},
	})
}

func testAccFCClusterDetailsDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_imaged_clusters_list" "cls" {}
	`
}

func testAccFCClusterDetailsDataSourceConfigWithUUID() string {
	return `
	data "nutanix_foundation_central_imaged_clusters_list" "cls" {}
	
	data "nutanix_foundation_central_cluster_details" "k1"{
		imaged_cluster_uuid = "${data.nutanix_foundation_central_imaged_clusters_list.cls.imaged_clusters[0].imaged_cluster_uuid}"
	}
	`
}
