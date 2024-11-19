package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameTemplateActions = "nutanix_template_guest_os_actions_v2.test"

func TestAccNutanixTemplateActionsV2Resource_Basic(t *testing.T) {
	t.Skip("Skipping test as it is not dependent on template")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTemplateActionsV2Config(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplateActions, "action", "cancel"),
				),
				Destroy: false,
			},
		},
	})
}

func testTemplateActionsV2Config(name, desc, tempName, tempDesc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		data "nutanix_subnets_v2" "subnets" { }

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}
		resource "nutanix_template_v2" "test" {
			template_name = "%[3]s"
			template_description = "%[4]s"
			template_version_spec{
			  version_source{
				template_vm_reference{
				  ext_id = nutanix_virtual_machine_v2.test.id
				}
			  }
			}
			lifecycle {
			  ignore_changes = [
				template_version_spec.0.version_name,
				template_version_spec.0.version_description,
				template_version_spec.0.version_source
			  ]
			} 
			depends_on = [nutanix_virtual_machine_v2.test]
		}

		resource "nutanix_template_guest_os_actions_v2" "test1" {
			ext_id = resource.nutanix_template_v2.test.id
			action = "initiate"
			version_id = resource.nutanix_template_v2.test.template_version_spec.0.ext_id
		}

		resource "nutanix_template_guest_os_actions_v2" "test" {
			ext_id = resource.nutanix_template_v2.test.id
			action = "cancel"
		}
		
`, name, desc, tempName, tempDesc)
}
