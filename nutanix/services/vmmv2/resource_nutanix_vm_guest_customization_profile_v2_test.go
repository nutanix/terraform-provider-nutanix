package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameVmGcProfile = "nutanix_vm_guest_customization_profile_v2.test"

func TestAccV2NutanixVmGuestCustomizationProfileResource_BasicSysprepParams(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-%d", r)
	desc := "test gc profile description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmGcProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileSysprepParamsConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "update_time"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.general_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.general_settings.0.computer_name.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.general_settings.0.computer_name.0.use_vm_name", "true"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.general_settings.0.timezone", "Pacific Standard Time"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.locale_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.locale_settings.0.ui_language", "en-US"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.network_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.network_settings.0.nic_config_list.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.network_settings.0.nic_config_list.0.ipv4_config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.network_settings.0.nic_config_list.0.ipv4_config.0.use_dhcp", "true"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.workgroup_or_domain_info.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.workgroup_or_domain_info.0.workgroup.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.sysprep_params.0.workgroup_or_domain_info.0.workgroup.0.name", "WORKGROUP"),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "links.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmGuestCustomizationProfileResource_WithUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-%d", r)
	updatedName := fmt.Sprintf("test-gc-profile-%d-updated", r)
	desc := "test gc profile description"
	updatedDesc := "test gc profile description updated"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmGcProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileSysprepParamsConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", desc),
				),
			},
			{
				Config: testVmGcProfileSysprepParamsConfig(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", updatedDesc),
				),
			},
		},
	})
}

func TestAccV2NutanixVmGuestCustomizationProfileResource_AnswerFile(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-af-%d", r)
	desc := "answer file profile"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVmGcProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileAnswerFileConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.answer_file.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "config.0.sysprep_config.0.customization.0.answer_file.0.unattend_xml"),
				),
			},
		},
	})
}

func testAccCheckNutanixVmGcProfileDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	client := conn.VmmAPI.VmGuestCustomizationProfilesAPIInstance

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_vm_guest_customization_profile_v2" {
			continue
		}
		_, err := client.GetVmGuestCustomizationProfileById(utils.StringPtr(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("VM Guest Customization Profile still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testVmGcProfileSysprepParamsConfig(name, desc string) string {
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
`, name, desc)
}

func testVmGcProfileAnswerFileConfig(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "test" {
			name        = "%s"
			description = "%s"
			config {
				sysprep_config {
					customization {
						answer_file {
							unattend_xml = "<unattend xmlns='urn:schemas-microsoft-com:unattend'></unattend>"
						}
					}
				}
			}
		}
`, name, desc)
}
