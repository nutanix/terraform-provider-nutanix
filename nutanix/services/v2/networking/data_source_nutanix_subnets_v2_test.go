package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSubnets = "data.nutanix_subnets_v4.test"

func TestAccNutanixSubnetsDataSourceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSubnets, "subnets.#"),
					resource.TestCheckResourceAttr(datasourceNameSubnets, "subnets.0.is_external", "true"),
					resource.TestCheckResourceAttr(datasourceNameSubnets, "subnets.0.subnet_type", "VLAN"),
					resource.TestCheckResourceAttrSet(datasourceNameSubnets, "subnets.0.cluster_reference"),
					resource.TestCheckResourceAttrSet(datasourceNameSubnets, "subnets.0.links.#"),
				),
			},
		},
	})
}

func testAccSubnetsDataSourceConfig() string {
	return (`
		data "nutanix_subnets_v4" "test" {}
`)
}
