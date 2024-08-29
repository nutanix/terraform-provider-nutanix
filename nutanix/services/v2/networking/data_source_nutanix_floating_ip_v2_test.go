package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNamefip = "data.nutanix_floating_ip_v2.test"

func TestAccNutanixFloatingIPDataSourceV2_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-fip-%d", r)
	desc := "test fip description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFipDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNamefip, "name", name),
					resource.TestCheckResourceAttr(datasourceNamefip, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNamefip, "metadata.#"),
					resource.TestCheckResourceAttrSet(datasourceNamefip, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNamefip, "association.#"),
					resource.TestCheckResourceAttrSet(datasourceNamefip, "external_subnet_reference"),
				),
			},
		},
	})
}

func testAccFipDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

		resource "nutanix_floating_ip_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			external_subnet_reference = "bd319622-1a45-4075-811a-2b0399bf9a49"
		}

		data "nutanix_floating_ip_v2" "test" {
			ext_id = nutanix_floating_ip_v2.test.ext_id
		}
	`, name, desc)
}
