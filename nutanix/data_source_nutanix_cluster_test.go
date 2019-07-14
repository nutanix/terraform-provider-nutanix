package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixClusterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "id"),
				),
			},
		},
	})
}

const testAccClusterDataSourceConfig = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	cluster_id = data.nutanix_clusters.clusters.entities.1.metadata.uuid
}`
