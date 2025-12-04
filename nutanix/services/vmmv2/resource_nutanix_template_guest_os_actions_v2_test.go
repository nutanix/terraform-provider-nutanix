package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameTemplateActions = "nutanix_template_guest_os_actions_v2.test"

func TestAccV2NutanixTemplateActionsResource_Basic(t *testing.T) {
	//t.Skip("Skipping test as it is not dependent on template")
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("tf-test-temp-%d", r)
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
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 =  [
				  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
				][0]
			config = jsondecode(file("%[5]s"))
			vmm = local.config.vmm
		}
	
		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
					}	
					vlan_mode = "ACCESS"
				}
			}
			power_state = "ON"
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
			depends_on = [nutanix_template_v2.test]
		}

		resource "nutanix_template_guest_os_actions_v2" "test" {
			ext_id = resource.nutanix_template_v2.test.id
			action = "cancel"
			depends_on = [nutanix_template_guest_os_actions_v2.test1]
		}
		
`, name, desc, tempName, tempDesc, filepath)
}
