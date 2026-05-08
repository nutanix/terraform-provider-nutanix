package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmGcProfile = "nutanix_vm_guest_customization_profile_v2.test"

func TestAccV2NutanixVmGuestCustomizationProfileResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-%d", r)
	desc := "test guest customization profile description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileResourceConfigBasic(name, desc),
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
				),
			},
		},
	})
}

func TestAccV2NutanixVmGuestCustomizationProfileResource_WithUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-%d", r)
	updatedName := fmt.Sprintf("test-gc-profile-%d-updated", r)
	desc := "test guest customization profile description"
	updatedDesc := "test guest customization profile description updated"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileResourceConfigBasic(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", desc),
				),
			},
			{
				Config: testVmGcProfileResourceConfigBasic(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmGcProfile, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVmGcProfile, "description", updatedDesc),
				),
			},
		},
	})
}

func TestAccV2NutanixVmGuestCustomizationProfileResource_WithAnswerFile(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-gc-profile-af-%d", r)
	desc := "test guest customization profile with answer file"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmGcProfileResourceConfigAnswerFile(name, desc),
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

func testVmGcProfileResourceConfigBasic(name, desc string) string {
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
`, name, desc)
}

func testVmGcProfileResourceConfigAnswerFile(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_vm_guest_customization_profile_v2" "test" {
			name        = "%[1]s"
			description = "%[2]s"
			config {
				sysprep_config {
					customization {
						answer_file {
							unattend_xml = "<unattend><settings pass='specialize'></settings></unattend>"
						}
					}
				}
			}
		}
`, name, desc)
}
