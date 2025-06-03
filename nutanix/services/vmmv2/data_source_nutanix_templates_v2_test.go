package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameTemplates = "data.nutanix_templates_v2.test"

func TestAccV2NutanixTemplateDatasource_ListAllTemplates(t *testing.T) {
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
				Config: testTemplatesDatasourceConfig(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixTemplateDatasource_ListAllTemplatesWithFilter(t *testing.T) {
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
				Config: testTemplatesDatasourceFilterConfig(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.#"),
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameTemplates, "templates.0.template_name", templateName),
					resource.TestCheckResourceAttr(datasourceNameTemplates, "templates.0.template_description", templateDesc),
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.0.template_version_spec.#"),
					resource.TestCheckResourceAttr(datasourceNameTemplates, "templates.0.template_version_spec.0.version_name", "Initial Version"),
					resource.TestCheckResourceAttr(datasourceNameTemplates, "templates.0.template_version_spec.0.version_description", "Created from VM: "+name),
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.0.create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.0.update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.0.created_by.#"),
					resource.TestCheckResourceAttrSet(datasourceNameTemplates, "templates.0.updated_by.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixTemplateDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTemplatesDatasourceInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameTemplates, "filter", "templateName eq 'invalid'"),
				),
			},
		},
	})
}

func testTemplatesDatasourceConfig(name, desc, tempName, tempDesc string) string {
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
		data "nutanix_templates_v2" "test" {
			depends_on = [ nutanix_template_v2.test ]
		}

`, name, desc, tempName, tempDesc)
}

func testTemplatesDatasourceFilterConfig(name, desc, tempName, tempDesc string) string {
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
		data "nutanix_templates_v2" "test" {
			filter = "templateName eq '%[3]s'"
			depends_on = [ nutanix_template_v2.test ]
		}

`, name, desc, tempName, tempDesc)
}

func testTemplatesDatasourceInvalidFilterConfig() string {
	return `
		data "nutanix_templates_v2" "test" {
			filter = "templateName eq 'invalid'"
		}
	`
}
