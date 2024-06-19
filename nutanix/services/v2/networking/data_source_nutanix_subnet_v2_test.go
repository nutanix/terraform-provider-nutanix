package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSubnet = "data.nutanix_subnet_v4.test"

func TestAccNutanixSubnetDataSourceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameSubnet, "is_external", "true"),
					resource.TestCheckResourceAttr(datasourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttrSet(datasourceNameSubnet, "cluster_reference"),
					resource.TestCheckResourceAttrSet(datasourceNameSubnet, "links.#"),
				),
			},
		},
	})
}

func testAccSubnetDataSourceConfig() string {
	return (`

		data "nutanix_subnets_v4" "test" {}

		data "nutanix_subnet_v4" "test" {
			ext_id = data.nutanix_subnets_v4.test.subnets.0.ext_id
		}
`)
}
