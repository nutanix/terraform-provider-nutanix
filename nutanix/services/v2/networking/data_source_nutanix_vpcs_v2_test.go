package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNamevpcs = "data.nutanix_vpcs_v2.test"

func TestAccNutanixVpcsDataSourceV2_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vpc-%d", r)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.#"),
					resource.TestCheckResourceAttr(datasourceNamevpcs, "vpcs.0.name", name),
					resource.TestCheckResourceAttr(datasourceNamevpcs, "vpcs.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.metadata.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.snat_ips.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpcs, "vpcs.0.external_subnets.#"),
				),
			},
		},
	})
}

func testAccVpcsDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

		resource "nutanix_vpc_v2" "rtest" {
			name =  "%[1]s"
			description = "%[2]s"
			external_subnets{
				subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
			}
		}

		data "nutanix_vpcs_v2" "test" {
			filter = "name eq '%[1]s'"
			depends_on = [
				resource.nutanix_vpc_v2.rtest
			]
		}
	`, name, desc)
}
