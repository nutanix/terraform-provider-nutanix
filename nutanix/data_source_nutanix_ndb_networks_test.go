package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraNetworksDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraNetworksDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_networks.test", "networks.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_networks.test", "networks.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_networks.test", "networks.0.managed"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_networks.test", "networks.0.properties.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_networks.test", "networks.0.type"),
				),
			},
		},
	})
}

func testAccEraNetworksDataSourceConfig() string {
	return `
	data "nutanix_ndb_networks" "test" { }
	`
}
