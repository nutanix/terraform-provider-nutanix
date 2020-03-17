package nutanix

import (
	"regexp"
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

func TestAccNutanixClusterByNameDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterByNameDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "id"),
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "name"),
				),
			},
		},
	})
}

func TestAccNutanixClusterByNameNotExistingDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterByNameNotExistingDataSourceConfig,
				ExpectError: regexp.MustCompile("Did not find cluster with name *"),
			},
		},
	})
}

const testAccClusterDataSourceConfig = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	cluster_id = data.nutanix_clusters.clusters.entities.0.metadata.uuid
}`

const testAccClusterByNameDataSourceConfig = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	name = data.nutanix_clusters.clusters.entities.0.name
}`

const testAccClusterByNameNotExistingDataSourceConfig = `
data "nutanix_clusters" "clusters" {}


data "nutanix_cluster" "cluster" {
	name = "ThisDoesNotExist"
}`
