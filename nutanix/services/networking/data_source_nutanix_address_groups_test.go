package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixAddressGroupsDataSource_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressGroupsDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_address_group.addr_group", "ip_address_block_list.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_address_group.addr_group", "description", "test address groups resource"),
					resource.TestCheckResourceAttr("data.nutanix_address_groups.addr_groups", "entities.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_address_groups.addr_groups", "entities.0.address_group.#", "1"),
				),
			},
		},
	})
}

func testAccAddressGroupsDataSourceConfig(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_address_group" "test_address" {
  			name = "test-%[1]d"
  			description = "test address groups resource"

  			ip_address_block_list {
    			ip = "10.0.0.0"
    			prefix_length = 24
  			}
		}

		data "nutanix_address_group" "addr_group" {
			uuid = nutanix_address_group.test_address.id
		}

		data "nutanix_address_groups" "addr_groups" {
			depends_on = [nutanix_address_group.test_address]
		}
	`, r)
}
