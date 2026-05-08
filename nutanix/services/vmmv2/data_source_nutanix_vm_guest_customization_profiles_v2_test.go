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
	desc := "test gc profiles list datasource"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfilesDatasourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmGcProfiles, "profiles.#"),
				),
			},
		},
	})
}

func testVmGcProfilesDatasourceConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "test" {
			name        = "%s"
			description = "%s"
			config {
				sysprep_config {
					customization {
						sysprep_params {
							general_settings {
								computer_name {
									use_vm_name = true
								}
								timezone = "Pacific Standard Time"
							}
							locale_settings {
								ui_language   = "en-US"
								system_locale = "en-US"
								user_locale   = "en-US"
							}
							network_settings {
								nic_config_list {
									ipv4_config {
										use_dhcp = true
									}
								}
							}
							workgroup_or_domain_info {
								workgroup {
									name = "WORKGROUP"
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
