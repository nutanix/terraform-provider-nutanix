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

const resourceServiceGroup = "nutanix_service_group.test"

func TestAccNutanixServiceGroup(t *testing.T) {
	name := acctest.RandomWithPrefix("nutanix_service_gr")
	description := "this is nutanix service group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixServiceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixServiceGroupConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceServiceGroup, "name", name),
					resource.TestCheckResourceAttr(resourceServiceGroup, "description", description),
					resource.TestCheckResourceAttr(resourceServiceGroup, "service_list.#", "1"),
					resource.TestCheckResourceAttr(resourceServiceGroup, "service_list.0.protocol", "TCP"),
				),
			},
			{
				ResourceName:      resourceServiceGroup,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixServiceGroupDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_service_grp" {
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

func testAccNutanixServiceGroupConfig(name, description string) string {
	return fmt.Sprintf(`
		resource "nutanix_service_group" "test" {
  			name = "%[1]s"
  			description = "%[2]s"

  			service_list {
      			protocol = "TCP"
      			tcp_port_range_list {
        		start_port = 22
        		end_port = 22
      		}
      		tcp_port_range_list {
       			start_port = 2222
        		end_port = 2222
      		}
		}
	}
`, name, description)
}
