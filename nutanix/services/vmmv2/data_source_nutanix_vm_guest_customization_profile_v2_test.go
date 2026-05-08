package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmGcProfile = "data.nutanix_vm_guest_customization_profile_v2.test"

func TestAccV2NutanixVmGuestCustomizationProfileDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-ds-%d", r)
	desc := "test guest customization profile datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileDatasourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "name", name),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfile, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfile, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "config.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmGcProfile, "config.0.sysprep_config.#", "1"),
				),
			},
		},
	})
}

func testVmGcProfileDatasourceConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "test" {
			name        = "%[1]s"
			description = "%[2]s"
			config {
				sysprep_config {
					customization {
						sysprep_params {
							general_settings {
								computer_name {
									use_vm_name = true
								}
							}
						}
					}
				}
			}
		}

		data "nutanix_vm_guest_customization_profile_v2" "test" {
			ext_id = resource.nutanix_vm_guest_customization_profile_v2.test.id
		}
`, name, desc)
}
