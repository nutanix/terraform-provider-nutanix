package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameServiceGroups = "data.nutanix_service_groups_v2.test"

func TestAccV2NutanixServiceGroupsDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-service-%d", r)
	desc := "test service description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceGrpsDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroups, "service_groups.#"),
					checkAttributeLength(datasourceNameServiceGroups, "service_groups", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixServiceGroupsDataSource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-service-%d", r)
	desc := "test service description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceGrpsDataSourceWithFilterConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroups, "service_groups.#"),
					resource.TestCheckResourceAttr(datasourceNameServiceGroups, "service_groups.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameServiceGroups, "service_groups.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameServiceGroups, "service_groups.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroups, "service_groups.0.tcp_services.#"),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroups, "service_groups.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroups, "service_groups.0.udp_services.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixServiceGroupsDataSource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceGrpsDataSourceWithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameServiceGroups, "service_groups.#"),
					resource.TestCheckResourceAttr(datasourceNameServiceGroups, "service_groups.#", "0"),
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
			depends_on = [
				resource.nutanix_service_groups_v2.test
			]
		}
	`, name, desc)
}

func testAccServiceGrpsDataSourceWithFilterConfig(name, desc string) string {
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

func testAccServiceGrpsDataSourceWithInvalidFilterConfig() string {
	return `



		data "nutanix_service_groups_v2" "test" {
			filter = "name eq 'invalid_filter'"

		}
	`
}
