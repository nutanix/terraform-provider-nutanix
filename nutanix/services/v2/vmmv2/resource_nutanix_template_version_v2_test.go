package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameTemplateVersion = "nutanix_template_version_v4.test"

func TestAccNutanixTemplateVersionV4_Basic(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
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
				Config: testTemplateVersionV4Config(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameTemplateVersion, "template_version_spec.0.create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplateVersion, "template_version_spec.0.created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplateVersion, "template_version_spec.0.vm_spec.#"),
					resource.TestCheckResourceAttr(resourceNameTemplateVersion, "template_version_spec.0.version_name", "Second temp"),
					resource.TestCheckResourceAttr(resourceNameTemplateVersion, "template_version_spec.0.version_description", "second desc"),
					resource.TestCheckResourceAttr(resourceNameTemplateVersion, "template_version_spec.0.is_active_version", "true"),
				),
			},
		},
	})
}

func TestAccNutanixTemplateVersionV4_WithDisk(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
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
				Config: testTemplateVersionV4ConfigWithDisk(name, desc, templateName, templateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameTemplateVersion, "template_version_spec.0.create_time"),
					resource.TestCheckResourceAttrSet(resourceNameTemplateVersion, "template_version_spec.0.created_by.#"),
					resource.TestCheckResourceAttrSet(resourceNameTemplateVersion, "template_version_spec.0.vm_spec.#"),
					resource.TestCheckResourceAttr(resourceNameTemplateVersion, "template_version_spec.0.version_name", "Second temp"),
					resource.TestCheckResourceAttr(resourceNameTemplateVersion, "template_version_spec.0.version_description", "second desc"),
					resource.TestCheckResourceAttr(resourceNameTemplateVersion, "template_version_spec.0.is_active_version", "true"),
				),
			},
		},
	})
}

func testTemplateVersionV4Config(name, desc, tempName, tempDesc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
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

		resource "nutanix_template_version_v4" "test" {
			template_ext_id = resource.nutanix_template_v2.test.id
			template_version_spec{
			  version_name = "Second temp"
			  version_description = "second desc"
			  version_source{
				template_version_reference{
				  version_id= nutanix_template_v2.test.template_version_spec.0.ext_id
				  override_vm_config{
					num_sockets=1
					num_threads_per_core=2
					memory_size_bytes= 1073741824
					num_cores_per_socket = 1    
				  }
				}
			  }
			  is_active_version = true
			}
			 lifecycle {
			  ignore_changes = [
				template_version_spec.0.version_source
			  ]
			}
		  }
`, name, desc, tempName, tempDesc)
}

func testTemplateVersionV4ConfigWithDisk(name, desc, tempName, tempDesc string) string {
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
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = "1073741824"
						storage_container{
							ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
						}
					}
				}
			}
			nics{
				network_info{
					nic_type = "NORMAL_NIC"
					subnet{
						ext_id = data.nutanix_subnets_v2.subnets.subnets.1.ext_id
					}	
					vlan_mode = "ACCESS"
				}
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

		resource "nutanix_template_version_v4" "test" {
			template_ext_id = resource.nutanix_template_v2.test.id
			template_version_spec{
			  version_name = "Second temp"
			  version_description = "second desc"
			  version_source{
				template_version_reference{
				  version_id= nutanix_template_v2.test.template_version_spec.0.ext_id
				  override_vm_config{
					num_sockets=1
					num_threads_per_core=2
					memory_size_bytes= 1073741824
					num_cores_per_socket = 1 
					nics{
						network_info{
							nic_type = "NORMAL_NIC"
							subnet{
								ext_id = data.nutanix_subnets_v2.subnets.subnets.1.ext_id
							}	
							vlan_mode = "ACCESS"
						}
					}   
				  }
				}
			  }
			  is_active_version = true
			}
			 lifecycle {
			  ignore_changes = [
				template_version_spec.0.version_source
			  ]
			}
		  }
`, name, desc, tempName, tempDesc)
}
