package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNamevpc = "data.nutanix_vpc_v4.test"

func TestAccNutanixVpcDataSourceV2_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vpc-%d", r)
	desc := "test vpc description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNamevpc, "name", name),
					resource.TestCheckResourceAttr(datasourceNamevpc, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNamevpc, "metadata.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpc, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpc, "snat_ips.#"),
					resource.TestCheckResourceAttrSet(datasourceNamevpc, "external_subnets.#"),
				),
			},
		},
	})
}

func testAccVpcDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

		resource "nutanix_vpc_v4" "test" {
			name =  "%[1]s"
			description = "%[2]s"
			external_subnets{
				subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
			}
		}

		data "nutanix_vpc_v4" "test" {
			ext_id = nutanix_vpc_v4.test.ext_id
			depends_on = [
				resource.nutanix_vpc_v4.test
			]
		}
	`, name, desc)
}
