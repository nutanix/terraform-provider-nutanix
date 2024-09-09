package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameServiceGrps = "data.nutanix_service_groups_v2.test"

func TestAccNutanixServiceGroupsDataSourceV2_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-service-%d", r)
	desc := "test service description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceGrpsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameServiceGrps, "service_groups.#"),
					resource.TestCheckResourceAttr(datasourceNameServiceGrps, "service_groups.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameServiceGrps, "service_groups.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameServiceGrps, "service_groups.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGrps, "service_groups.0.tcp_services.#"),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGrps, "service_groups.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGrps, "service_groups.0.udp_services.#"),
				),
			},
		},
	})
}

func testAccServiceGrpsDataSourceConfig(name, desc string) string {
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

		data "nutanix_service_groups_v2" "test" {
			filter = "name eq '%[1]s'"
			depends_on = [
				resource.nutanix_service_groups_v2.test
			]
		}
	`, name, desc)
}
