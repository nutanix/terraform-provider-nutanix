package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixServiceGroupDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("nutanix_service_gr")
	description := "this is nutanix service group"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceGroupDataSourceConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_service_group.service_group", "name", name),
					resource.TestCheckResourceAttr("data.nutanix_service_group.service_group", "description", description),
					resource.TestCheckResourceAttr("data.nutanix_service_group.service_group", "service_list.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_service_group.service_group", "service_list.0.protocol", "TCP"),
				),
			},
		},
	})
}

func testAccServiceGroupDataSourceConfig(name, description string) string {
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
		data "nutanix_service_group" "service_group" {
			uuid = "${nutanix_service_group.test.id}"
		}
	
	`, name, description)
}
