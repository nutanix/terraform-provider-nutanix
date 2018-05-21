package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleClusterDataSource_basic(t *testing.T) {
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

// Lookup based on InstanceID
const testAccClusterDataSourceConfig = `
provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}

data "nutanix_clusters" "clusters" {
	metadata = {
		length = 2
	}
}


data "nutanix_cluster" "cluster" {
	cluster_id = "${data.nutanix_clusters.clusters.entities.1.metadata.uuid}"
}`
