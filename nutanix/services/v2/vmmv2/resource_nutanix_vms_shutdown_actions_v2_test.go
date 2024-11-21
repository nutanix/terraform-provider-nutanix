package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmShutdown = "data.nutanix_virtual_machine_v2.test"

func TestAccNutanixVmsShutdownV2Resource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm power action "

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// create a vm with ngt
			{
				Config: testVmV2Config(name, desc, "ON"),
			},
			// create a vm shutdown action
			{
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("shutdown"),
			},
			// check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 1 minute to allow the vm to shut down
					time.Sleep(1 * time.Minute)
				},
				Config: testVmV2Config(name, desc, "OFF") + testVmsShutdownV2Config("shutdown") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmShutdown, "power_state", "OFF"),
				),
			},
			// power on the vm
			{
				Config: testVmV2Config(name, desc, "ON"),
			},
			// check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to power on
					time.Sleep(1 * time.Minute)
				},
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("shutdown") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmShutdown, "power_state", "ON"),
				),
			},
			// reboot the vm
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to power on
					time.Sleep(1 * time.Minute)
				},
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("reboot"),
			},
			// check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to reboot
					time.Sleep(1 * time.Minute)
				},
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("reboot") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmShutdown, "power_state", "ON"),
				),
			},
			// guest_reboot the vm
			{
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("guest_reboot"),
			},
			// check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 2 Minute to allow the vm to reboot
					time.Sleep(2 * time.Minute)
				},
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("guest_reboot") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmShutdown, "power_state", "ON"),
				),
			},
			// guest_shutdown the vm
			{
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("guest_shutdown"),
			},
			// check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 2 Minute to allow the vm to shut down
					time.Sleep(2 * time.Minute)
				},
				Config: testVmV2Config(name, desc, "OFF") + testVmsShutdownV2Config("guest_shutdown") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVmShutdown, "power_state", "OFF"),
				),
			},
			// power on the vm to uninstall ngt and delete the vm
			{
				Config: testVmV2Config(name, desc, "ON") + testVmsShutdownV2Config("guest_shutdown") + vmDataSource,
			},
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to power on
					time.Sleep(1 * time.Minute)
				},
			},
		},
	})
}

func TestAccNutanixVmsShutdownV2Resource_WithError(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
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

func testVmV2Config(name, desc, powerState string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = jsondecode(file("%[4]s"))
			vmm    = local.config.vmm
		}

		data "nutanix_images_v2" "ngt-image" {
		  filter = "name eq '${local.vmm.image_name}'"
		}

		data "nutanix_subnets_v2" "subnet" {
		  filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "rtest"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}

			disks {
				disk_address {
				  bus_type = "SCSI"
				  index    = 0
				}
				backing_info {
				  vm_disk {
					data_source {
					  reference {
						image_reference {
						  image_ext_id = data.nutanix_images_v2.ngt-image.images[0].ext_id
						}
					  }
					}
					disk_size_bytes = 20 * 1024 * 1024 * 1024
				  }
				}
			}
			
			cd_roms {
				disk_address {
				  bus_type = "IDE"
				  index    = 0
				}
			}
			
			nics {
				network_info {
				  nic_type = "NORMAL_NIC"
				  subnet {
					ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
				  }
				  vlan_mode = "ACCESS"
				}
			}
			
			boot_config {
				legacy_boot {
				  boot_order = ["CDROM", "DISK", "NETWORK"]
				}
			}
			
			power_state = "%[3]s"
			
			lifecycle {
				ignore_changes = [guest_tools]
			}
			
			depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_images_v2.ngt-image, data.nutanix_storage_containers_v2.ngt-sc]			
		}

		resource "nutanix_ngt_installation_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.rtest.id
			credential {
				username = local.vmm.ngt.credential.username
				password = local.vmm.ngt.credential.password
			}
			reboot_preference {
				schedule_type = "IMMEDIATE"
			}
			capablities = ["SELF_SERVICE_RESTORE", "VSS_SNAPSHOT"]
			depends_on = [nutanix_virtual_machine_v2.rtest]
		}

		
`, name, desc, powerState, filepath)
}

func testVmsShutdownV2Config(action string) string {
	return fmt.Sprintf(`
		resource "nutanix_vm_shutdown_action_v2" "vmShuts" {
			ext_id= resource.nutanix_virtual_machine_v2.rtest.id
			action = "%[1]s"
		}
		`, action)
}

func testVmsShutdownV4ConfigWithError(name, desc, state string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}

		data "nutanix_subnets_v2" "subnets" { }

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}'"
		  limit = 1
		}
	
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
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			power_state = "ON"
		}


		resource "nutanix_vm_shutdown_action_v2" "vmShuts" {
			ext_id= resource.nutanix_virtual_machine_v2.rtest.id
			action = "shutdown"
			guest_power_state_transition_config{
				should_fail_on_script_failure = false
			  }
		}

`, name, desc, state)
}

const vmDataSource = ` 
		data "nutanix_virtual_machine_v2" "test"{
			ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			depends_on = [
				resource.nutanix_vm_shutdown_action_v2.vmShuts
			]
		}
`
