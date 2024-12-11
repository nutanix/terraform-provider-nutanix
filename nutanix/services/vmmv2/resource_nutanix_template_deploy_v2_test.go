package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameTemplateDeploy = "nutanix_deploy_templates_v2.test"

func TestAccV2NutanixTemplateDeployResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTemplateDeployV2Config(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "number_of_vms", "1"),
					resource.TestCheckResourceAttrSet(resourceNameTemplateDeploy, "override_vm_config_map.#"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.name", "test-tf-template-deploy"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.memory_size_bytes", "4294967296"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameTemplateDeploy, "override_vm_config_map.0.num_cores_per_socket", "1"),
				),
			},
		},
	})
}

func testTemplateDeployV2Config(name, desc, tempName, tempDesc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		    ][0]
			config = jsondecode(file("%[5]s"))
			vmm    = local.config.vmm
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

		resource "nutanix_deploy_templates_v2" "test" {
			ext_id = resource.nutanix_template_v2.test.id
			number_of_vms = 1
			cluster_reference = local.cluster0
			override_vm_config_map{
			  name= "test-tf-template-deploy"
			  memory_size_bytes = 4294967296
			  num_sockets=2
			  num_cores_per_socket=1
			  num_threads_per_core=1
			  nics{
				backing_info{
				  is_connected = true
				  model = "VIRTIO"
				}
				network_info {
				  nic_type = "NORMAL_NIC"
				  subnet {
					ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
				  }
				  vlan_mode = "ACCESS"
				  should_allow_unknown_macs = false
				}
			  }
			}
			depends_on = [
				resource.nutanix_template_v2.test
			]
		}	

`, name, desc, tempName, tempDesc, filepath)
}
