package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraNetworkDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraNetworkDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "cluster_id"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_network.test", "managed", "false"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "type"),
				),
			},
		},
	})
}

func TestAccEraNetworkDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraNetworkDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "cluster_id"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_network.test", "managed", "false"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network.test", "type"),
				),
			},
		},
	})
}

func testAccEraNetworkDataSourceConfig() string {
	return `
	data "nutanix_ndb_networks" "name" { }
	
	data "nutanix_ndb_network" "test" {
		id = data.nutanix_ndb_networks.name.networks.0.id
	}
	`
}

func testAccEraNetworkDataSourceConfigByName() string {
	return `
	data "nutanix_ndb_networks" "name" { }
	
	data "nutanix_ndb_network" "test" {
		name = data.nutanix_ndb_networks.name.networks.0.name
	}
	`
}
