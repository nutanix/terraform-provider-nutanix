package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNamefips = "data.nutanix_floating_ips_v4.test"

func TestAccNutanixFloatingIPsDataSourceV2_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFipsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNamefips, "floating_ips.#"),
					resource.TestCheckResourceAttr(datasourceNamefips, "floating_ips.#", "1"),
					resource.TestCheckResourceAttr(datasourceNamefips, "floating_ips.0.name", name),
					resource.TestCheckResourceAttr(datasourceNamefips, "floating_ips.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNamefips, "floating_ips.0.metadata.#"),
					resource.TestCheckResourceAttrSet(datasourceNamefips, "floating_ips.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNamefips, "floating_ips.0.association.#"),
					resource.TestCheckResourceAttrSet(datasourceNamefips, "floating_ips.0.external_subnet_reference"),
				),
			},
		},
	})
}

func testAccFipsDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

		resource "nutanix_floating_ip_v4" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
		}

		data "nutanix_floating_ips_v4" "test" {
			filter = "name eq '%[1]s'"
			depends_on = [
				resource.nutanix_floating_ip_v4.test
			]
		}
	`, name, desc)
}
