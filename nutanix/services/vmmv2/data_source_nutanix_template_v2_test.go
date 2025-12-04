package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameTemplate = "data.nutanix_template_v2.test"

func TestAccV2NutanixTemplateDatasource_Basic(t *testing.T) {
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
				Config: testTemplateDatasourceConfig(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameTemplate, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameTemplate, "template_name", templateName),
					resource.TestCheckResourceAttr(datasourceNameTemplate, "template_description", templateDesc),
					resource.TestCheckResourceAttrSet(datasourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttr(datasourceNameTemplate, "template_version_spec.0.version_name", "Initial Version"),
					resource.TestCheckResourceAttr(datasourceNameTemplate, "template_version_spec.0.version_description", "Created from VM: "+name),
					resource.TestCheckResourceAttrSet(datasourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(datasourceNameTemplate, "updated_by.#"),
				),
			},
		},
	})
}

func testTemplateDatasourceConfig(name, desc, tempName, tempDesc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
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
				version_description = "Created from VM: %[1]s"
				version_source{
					template_vm_reference{
						ext_id = nutanix_virtual_machine_v2.test.id
					}
				}
			}
			depends_on = [nutanix_virtual_machine_v2.test]
		}

		data "nutanix_template_v2" "test" {
			ext_id = nutanix_template_v2.test.id
		}

`, name, desc, tempName, tempDesc)
}
