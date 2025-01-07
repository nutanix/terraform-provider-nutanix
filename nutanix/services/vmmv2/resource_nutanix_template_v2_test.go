package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameTemplate = "nutanix_template_v2.test"

func TestAccV2NutanixTemplateResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testTemplateV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplateV2Config(name, desc, templateName, templateDesc),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_name", templateName),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_description", templateDesc),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "updated_by.#"),
				),
			},
			//update the template name, description and vm config USING template_version_reference
			{
				Config: testTemplateV2UpdateWithTempVersionRefConfig(name, desc, templateName+"-updated", templateDesc+"-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_name", templateName+"-updated"),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_description", templateDesc+"-updated"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "updated_by.#"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_name", "2.0.0"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_description", "updating version from initial to 2.0.0"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_version_reference.0.override_vm_config.0.name", "tf-test-vm-2.0.0"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_version_reference.0.override_vm_config.0.memory_size_bytes", "4294967296"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_version_reference.0.override_vm_config.0.num_cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_version_reference.0.override_vm_config.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_version_reference.0.override_vm_config.0.num_threads_per_core", "2"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_version_reference.0.override_vm_config.0.guest_customization.0.config.0.cloud_init.0.cloud_init_script.0.custom_key_values.0.key_value_pairs.0.name", "locale"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_version_reference.0.override_vm_config.0.guest_customization.0.config.0.cloud_init.0.cloud_init_script.0.custom_key_values.0.key_value_pairs.0.value.0.string", "en-US"),
				),
			},
			//update the template name, description and vm config USING template_vm_reference
			{
				Config: testTemplateV2UpdateWithTempVMRefConfig(name, desc, templateName+"-updated-2", templateDesc+"-updated-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_name", templateName+"-updated-2"),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_description", templateDesc+"-updated-2"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_version_spec.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "updated_by.#"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_name", "3.0.0"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_description", "updating version from initial to 3.0.0"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_vm_reference.0.ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixTemplateResource_RequiredVersionNameOnUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testTemplateV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplateV2Config(name, desc, templateName, templateDesc),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_name", templateName),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_description", templateDesc),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "updated_by.#"),
				),
			},
			// version name is required for update
			{
				Config:      testTemplateV2UpdateWithoutVersionNameConfig(name, desc, templateName, templateDesc),
				ExpectError: regexp.MustCompile("version_name is required for update operation"),
			},
		},
	})
}

func TestAccV2NutanixTemplateResource_RequiredVersionDescriptionOnUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testTemplateV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplateV2Config(name, desc, templateName, templateDesc),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_name", templateName),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_description", templateDesc),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "updated_by.#"),
				),
			},
			// version Description is required for update
			{
				Config:      testTemplateV2UpdateWithoutVersionDescriptionConfig(name, desc, templateName, templateDesc),
				ExpectError: regexp.MustCompile("version_description is required for update operation"),
			},
		},
	})
}

func TestAccV2NutanixTemplateResource_GuestCustomizationSysprep(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testTemplateV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplateV2GuestCustomSysprepConfig(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_name", templateName),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_description", templateDesc),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "updated_by.#"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_vm_reference.0.guest_customization.0.config.0.sysprep.0.sysprep_script.0.custom_key_values.0.key_value_pairs.0.name", "locale"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_vm_reference.0.guest_customization.0.config.0.sysprep.0.sysprep_script.0.custom_key_values.0.key_value_pairs.0.name", "locale"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_vm_reference.0.guest_customization.0.config.0.sysprep.0.sysprep_script.0.custom_key_values.0.key_value_pairs.0.value.0.string", "en-PS"),
				),
			},
		},
	})
}

func TestAccV2NutanixTemplateResource_GuestCustomizationCloudInit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	templateName := fmt.Sprintf("test-temp-%d", r)
	templateDesc := "test temp description"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testTemplateV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTemplateV2GuestCustomCloudInitConfig(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_name", templateName),
					resource.TestCheckResourceAttr(resourceNameTemplate, "template_description", templateDesc),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "template_version_spec.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "update_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplate, "updated_by.#"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_vm_reference.0.guest_customization.0.config.0.cloud_init.0.cloud_init_script.0.custom_key_values.0.key_value_pairs.0.name", "locale"),
					resource.TestCheckResourceAttr(resourceNameTemplate,
						"template_version_spec.0.version_source.0.template_vm_reference.0.guest_customization.0.config.0.cloud_init.0.cloud_init_script.0.custom_key_values.0.key_value_pairs.0.value.0.string", "en-PS"),
				),
			},
		},
	})
}

func testTemplateV2CheckDestroy(state *terraform.State) error {
	fmt.Println("testTemplateV2CheckDestroy")
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	client := conn.VmmAPI.TemplatesAPIInstance
	for _, rs := range state.RootModule().Resources {
		if rs.Type != resourceNameTemplate {
			continue
		}
		_, err := client.GetTemplateById(utils.StringPtr(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("template still exists")
		}
	}
	return nil
}

func testTemplateV2Config(name, desc, tempName, tempDesc string) string {
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
			  version_source{
				template_vm_reference{
				  ext_id = nutanix_virtual_machine_v2.test.id
				}
			  }
			} 
			depends_on = [nutanix_virtual_machine_v2.test]
		}
`, name, desc, tempName, tempDesc)
}

func testTemplateV2UpdateWithTempVersionRefConfig(name, desc, tempName, tempDesc string) string {
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
		  template_name        = "%[3]s"
		  template_description = "%[4]s"
		  template_version_spec {
			version_name        = "2.0.0"
		    version_description = "updating version from initial to 2.0.0"
		  	is_active_version   = true
			version_source {
			  template_version_reference {
				override_vm_config {
				  name                 = "tf-test-vm-2.0.0"
				  memory_size_bytes    = 4 * 1024 * 1024 * 1024 # 4 GB
				  num_cores_per_socket = 2
				  num_sockets          = 2
				  num_threads_per_core = 2
				  guest_customization {
					config {
					  cloud_init {
						cloud_init_script {
						  user_data {
							value = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    ssh-authorized-keys:\n      - ssh-rsa DUMMYSSH mypass\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
						  }
						  custom_key_values {
							key_value_pairs {
							  name = "locale"
							  value {
								string = "en-US"
							  }
							}
						  }
						}
					  }
					}
				  }
				}
			  }
			}
		  }
		  depends_on = [nutanix_virtual_machine_v2.test]
		}		
`, name, desc, tempName, tempDesc)
}

func testTemplateV2UpdateWithTempVMRefConfig(name, desc, tempName, tempDesc string) string {
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

		resource "nutanix_virtual_machine_v2" "test-2"{
			name= "%[1]s-2"
			description =  "%[2]s 2"
			memory_size_bytes    = 4 * 1024 * 1024 * 1024 # 4 GB
		    num_cores_per_socket = 2
		    num_sockets          = 2
		    num_threads_per_core = 2
			cluster {
				ext_id = local.cluster0
			}
		}

		resource "nutanix_template_v2" "test" {
		  template_name        = "%[3]s"
		  template_description = "%[4]s"
		  template_version_spec {
			version_name        = "3.0.0"
		    version_description = "updating version from initial to 3.0.0"
		  	is_active_version   = true
			version_source {
				template_vm_reference{
				  ext_id = nutanix_virtual_machine_v2.test-2.id
				}
			}
		  }
		  depends_on = [nutanix_virtual_machine_v2.test, nutanix_virtual_machine_v2.test-2]
		}		
`, name, desc, tempName, tempDesc)
}

func testTemplateV2UpdateWithoutVersionNameConfig(name, desc, tempName, tempDesc string) string {
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
		      version_description = "updating version from initial to 2.0.0"
			  version_source{
				template_vm_reference{
				  ext_id = "00000000-0000-0000-0000-000000000000"
				}
			  }
			} 
			depends_on = [nutanix_virtual_machine_v2.test]
		}	
`, name, desc, tempName, tempDesc)
}

func testTemplateV2UpdateWithoutVersionDescriptionConfig(name, desc, tempName, tempDesc string) string {
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
              version_name        = "2.0.0"
			  version_source{
				template_vm_reference{
				  ext_id = "00000000-0000-0000-0000-000000000000"
				}
			  }
			} 
			depends_on = [nutanix_virtual_machine_v2.test]
		}	
`, name, desc, tempName, tempDesc)
}

func testTemplateV2GuestCustomSysprepConfig(name, desc, tempName, tempDesc string) string {
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
			  version_source{
				template_vm_reference{
			  		ext_id = nutanix_virtual_machine_v2.test.id
					guest_customization {
					  config {
						sysprep {
						  sysprep_script {
							custom_key_values {
							  key_value_pairs {
								name = "locale"
								value {
								  string = "en-PS"
								}
							  }
							}
						  }
						}
					  }
					}
				}
			  }
			} 
			depends_on = [nutanix_virtual_machine_v2.test]
		}
`, name, desc, tempName, tempDesc)
}

func testTemplateV2GuestCustomCloudInitConfig(name, desc, tempName, tempDesc string) string {
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
			  version_source{
				template_vm_reference{
			  		ext_id = nutanix_virtual_machine_v2.test.id
					guest_customization {
					  config {
						cloud_init {
						  cloud_init_script {
							custom_key_values {
							  key_value_pairs {
								name = "locale"
								value {
								  string = "en-PS"
								}
							  }
							}
						  }
						}
					  }
					}
				}
			  }
			} 
			depends_on = [nutanix_virtual_machine_v2.test]
		}
`, name, desc, tempName, tempDesc)
}
