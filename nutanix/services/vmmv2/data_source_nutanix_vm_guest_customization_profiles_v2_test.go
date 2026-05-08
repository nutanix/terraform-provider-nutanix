package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmGcProfiles = "data.nutanix_vm_guest_customization_profiles_v2.test"

func TestAccV2NutanixVmGuestCustomizationProfilesDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-list-%d", r)
	desc := "test guest customization profiles list datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfilesDatasourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfiles, "vm_guest_customization_profiles.#"),
				),
			},
		},
	})
}

func testVmGcProfilesDatasourceConfig(name, desc string) string {
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

		data "nutanix_vm_guest_customization_profiles_v2" "test" {
			depends_on = [nutanix_vm_guest_customization_profile_v2.test]
		}
`, name, desc)
}
