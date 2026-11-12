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
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "cloud_type"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "hypervisor_type"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "healthy"),
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
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "status"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "cloud_type"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "hypervisor_type"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_cluster.test", "healthy"),
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

		locals {
			// Use NDBCluster if it exists, otherwise pick the first available cluster
			target_cluster = contains([for c in data.nutanix_ndb_clusters.test1.clusters : c.name], "NDBCluster") ? [for c in data.nutanix_ndb_clusters.test1.clusters : c if c.name == "NDBCluster"][0] : data.nutanix_ndb_clusters.test1.clusters[0]
		}

		data "nutanix_ndb_cluster" "test" {
			cluster_id = local.target_cluster.id
		}	
	`
}

func testAccEraClusterDataSourceConfigByName() string {
	return `
		data "nutanix_ndb_clusters" "test1" {}

		locals {
			// Use NDBCluster if it exists, otherwise pick the first available cluster
			target_cluster = contains([for c in data.nutanix_ndb_clusters.test1.clusters : c.name], "NDBCluster") ? [for c in data.nutanix_ndb_clusters.test1.clusters : c if c.name == "NDBCluster"][0] : data.nutanix_ndb_clusters.test1.clusters[0]
		}

		data "nutanix_ndb_cluster" "test" {
			cluster_name = local.target_cluster.name
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
