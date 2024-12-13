package networking_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourcesAddressGroup = "nutanix_address_group.test"

func TestAccNutanixAddressGroup(t *testing.T) {
	name := acctest.RandomWithPrefix("nutanix_address_gr")
	description := "this is nutanix address group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixAddressGroupConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcesAddressGroup, "name", name),
					resource.TestCheckResourceAttr(resourcesAddressGroup, "description", description),
					resource.TestCheckResourceAttr(resourcesAddressGroup, "ip_address_block_list.#", "1"),
					resource.TestCheckResourceAttr(resourcesAddressGroup, "ip_address_block_list.0.prefix_length", "24"),
				),
			},
			{
				ResourceName:      resourcesAddressGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixAddressGroupDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_address_grp" {
			continue
		}
		for {
			_, err := conn.API.V3.GetVM(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}
	}

	return nil
}

func testAccNutanixAddressGroupConfig(name, description string) string {
	return fmt.Sprintf(`
		resource "nutanix_address_group" "test" {
			name        = "%[1]s"
			description = "%[2]s"
			ip_address_block_list {
    			ip = "10.0.0.0"
   	 			prefix_length = 24
  			}
		}
`, name, description)
}
