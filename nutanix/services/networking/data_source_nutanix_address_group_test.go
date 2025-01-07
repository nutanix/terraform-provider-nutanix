package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixAddressGroupDataSource_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressGroupDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_address_group.addr_group", "ip_address_block_list.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_address_group.addr_group", "ip_address_block_list.0.prefix_length", "24"),
					resource.TestCheckResourceAttr("data.nutanix_address_group.addr_group", "description", "test address group resource"),
				),
			},
		},
	})
}

func testAccAddressGroupDataSourceConfig(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_address_group" "test_address" {
  			name = "test-%[1]d"
  			description = "test address group resource"

  			ip_address_block_list {
    			ip = "10.0.0.0"
    			prefix_length = 24
  			}
		}

		data "nutanix_address_group" "addr_group" {
			uuid = "${nutanix_address_group.test_address.id}"
		}
	`, r)
}
