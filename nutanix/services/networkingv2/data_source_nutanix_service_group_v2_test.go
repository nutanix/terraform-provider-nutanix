package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameServiceGroup = "data.nutanix_service_group_v2.test"

func TestAccV2NutanixServiceGroupDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-service-%d", r)
	desc := "test service description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceGrpDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameServiceGroup, "name", name),
					resource.TestCheckResourceAttr(datasourceNameServiceGroup, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroup, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroup, "tcp_services.#"),
					resource.TestCheckResourceAttr(datasourceNameServiceGroup, "tcp_services.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroup, "udp_services.#"),
					resource.TestCheckResourceAttr(datasourceNameServiceGroup, "udp_services.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroup, "ext_id"),
				),
			},
		},
	})
}

func testAccServiceGrpDataSourceConfig(name, desc string) string {
	return fmt.Sprintf(`

		resource "nutanix_service_groups_v2" "test" {
			name  = "%[1]s"
			description = "%[2]s"  
			tcp_services {
				start_port = "232"
				end_port = "232"
			}
			udp_services {
				start_port = "232"
				end_port = "232"
			}
		}

		data "nutanix_service_group_v2" "test" {
			ext_id = nutanix_service_groups_v2.test.ext_id
			depends_on = [
				resource.nutanix_service_groups_v2.test
			]
		}
	`, name, desc)
}
