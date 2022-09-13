package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraClusterDataSource_basic(t *testing.T) {
	// r := randIntBetween(31, 40)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraClusterDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "status", "UP"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "cloud_type", "NTNX"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "hypervisor_type", "AHV"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "hypervisor_version", "master"),
					resource.TestCheckResourceAttrSet("data.nutanix_era_cluster.test", "properties.#"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "healthy", "true"),
				),
			},
		},
	})
}

func TestAccEraClusterDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraClusterDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "status", "UP"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "cloud_type", "NTNX"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "hypervisor_type", "AHV"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "hypervisor_version", "master"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "properties.#", "0"),
					resource.TestCheckResourceAttr("data.nutanix_era_cluster.test", "healthy", "true"),
				),
			},
		},
	})
}

func testAccEraClusterDataSourceConfig() string {
	return `
		data "nutanix_era_clusters" "test1" {}

		data "nutanix_era_cluster" "test" {
			depends_on = [data.nutanix_era_clusters.test1]
			cluster_id = data.nutanix_era_clusters.test1.clusters[0].id
		}	
	`
}

func testAccEraClusterDataSourceConfigByName() string {
	return `
		data "nutanix_era_clusters" "test1" {}

		data "nutanix_era_cluster" "test" {
			depends_on = [data.nutanix_era_clusters.test1]
			cluster_name = data.nutanix_era_clusters.test1.clusters[0].name
		}	
	`
}
