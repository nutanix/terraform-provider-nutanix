package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmShutdown = "data.nutanix_virtual_machine_v2.test"

func TestAccNutanixVmsShutdownV4_Basic(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	stateOn := "power_on"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsShutdownV4Config(name, desc, stateOn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmShutdown, "power_state", "OFF"),
				),
			},
		},
	})
}

func TestAccNutanixVmsShutdownV4_WithError(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	stateOn := "power_on"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testVmsShutdownV4ConfigWithError(name, desc, stateOn),
				ExpectError: regexp.MustCompile("guest_power_state_transition_config  attribute is not optional"),
			},
		},
	})
}

func testVmsShutdownV4Config(name, desc, state string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}

		data "nutanix_subnets_v2" "subnets" { }
	
		resource "nutanix_virtual_machine_v2" "rtest"{
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
						ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
					}	
					vlan_mode = "ACCESS"
				}
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					disk_size_bytes = 21474836480
					data_source {
						reference{
							image_reference{
								image_ext_id = "5867f64e-7d0a-4b04-a72e-e26a4dbbaea2"
							}
						}
					}
				}
			}
		}
		resource "nutanix_vm_power_action" "test" {
			ext_id= resource.nutanix_virtual_machine_v2.rtest.id
			action = "%[3]s"
			depends_on = [
				resource.nutanix_virtual_machine_v2.rtest
			]
		}

		resource "nutanix_vm_shutdown_action_v4" "vmShuts" {
			ext_id= resource.nutanix_virtual_machine_v2.rtest.id
			action = "shutdown"
			depends_on =[
				resource.nutanix_vm_power_action.test
			]
		}

		data "nutanix_virtual_machine_v2" "test"{
			ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			depends_on = [
				resource.nutanix_vm_shutdown_action_v4.vmShuts
			]
		}
`, name, desc, state)
}

func testVmsShutdownV4ConfigWithError(name, desc, state string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}

		data "nutanix_subnets_v2" "subnets" { }
	
		resource "nutanix_virtual_machine_v2" "rtest"{
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
						ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
					}	
					vlan_mode = "ACCESS"
				}
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
		}
		resource "nutanix_vm_power_action" "test" {
			ext_id= resource.nutanix_virtual_machine_v2.rtest.id
			action = "%[3]s"
			depends_on = [
				resource.nutanix_virtual_machine_v2.rtest
			]
		}

		resource "nutanix_vm_shutdown_action_v4" "vmShuts" {
			ext_id= resource.nutanix_virtual_machine_v2.rtest.id
			action = "shutdown"
			guest_power_state_transition_config{
				should_fail_on_script_failure = false
			  }
			depends_on =[
				resource.nutanix_vm_power_action.test
			]
		}

`, name, desc, state)
}
