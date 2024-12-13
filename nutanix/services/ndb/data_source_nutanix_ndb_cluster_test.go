package ndb_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraClusterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraClusterDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "status", "UP"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "cloud_type", "NTNX"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "hypervisor_type", "AHV"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "properties.#"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "healthy", "true"),
				),
			},
		},
	})
}

func TestAccEraClusterDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraClusterDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "status", "UP"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "cloud_type", "NTNX"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "hypervisor_type", "AHV"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "properties.#"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_cluster.test", "healthy", "true"),
				),
			},
		},
	})
}

func TestAccEraClusterDataSource_WithWrongID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEraClusterDataSourceConfigWithWrongID(),
				ExpectError: regexp.MustCompile("exit status 1"),
			},
		},
	})
}

func testAccEraClusterDataSourceConfig() string {
	return `
		data "nutanix_ndb_clusters" "test1" {}

		data "nutanix_ndb_cluster" "test" {
			depends_on = [data.nutanix_ndb_clusters.test1]
			cluster_id = data.nutanix_ndb_clusters.test1.clusters[0].id
		}	
	`
}

func testAccEraClusterDataSourceConfigByName() string {
	return `
		data "nutanix_ndb_clusters" "test1" {}

		data "nutanix_ndb_cluster" "test" {
			depends_on = [data.nutanix_ndb_clusters.test1]
			cluster_name = data.nutanix_ndb_clusters.test1.clusters[0].name
		}	
	`
}

func testAccEraClusterDataSourceConfigWithWrongID() string {
	return `
		data "nutanix_ndb_cluster" "test" {
			cluster_id = "0000000-0000-0000-0000-00000000000"
		}	
	`
}
