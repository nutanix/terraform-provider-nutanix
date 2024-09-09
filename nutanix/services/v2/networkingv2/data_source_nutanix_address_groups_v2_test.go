package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameAddGrps = "data.nutanix_address_groups_v2.test"

func TestAccNutanixAddressGroupsDataSourceV2_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-service-%d", r)
	desc := "test service description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAddGrpsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameAddGrps, "address_groups.#"),
					resource.TestCheckResourceAttr(datasourceNameAddGrps, "address_groups.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameAddGrps, "address_groups.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameAddGrps, "address_groups.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameAddGrps, "address_groups.0.ipv4_addresses.#"),
					resource.TestCheckResourceAttrSet(datasourceNameAddGrps, "address_groups.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameAddGrps, "address_groups.0.created_by"),
				),
			},
		},
	})
}

func testAccAddGrpsDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

		resource "nutanix_address_groups_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			ipv4_addresses{
			value = "10.0.0.0"
			prefix_length = 24
			}
		}

		data "nutanix_address_groups_v2" "test" {
			filter = "name eq '%[1]s'"
			depends_on = [
				resource.nutanix_address_groups_v2.test
			]
		}
	`, name, desc)
}
