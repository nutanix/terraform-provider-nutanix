package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVMShutdown = "data.nutanix_virtual_machine_v2.test"

func TestAccV2NutanixVmsShutdownResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm power action "

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineV2Destroy,
		Steps: []resource.TestStep{
			// 1. create a vm with ngt
			{
				Config: testVMV2Config(name, desc, "ON"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
			// 2. install ngt on the vm
			{
				PreConfig: func() {
					//sleep for 1 minute before installing ngt
					time.Sleep(1 * time.Minute)
				},
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_virtual_machine_v2.rtest", "power_state", "ON"),
				),
			},
			// 3. create a vm shutdown action
			{
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig() + testVmsShutdownV2Config("shutdown"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
			// 4. check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 1 minute to allow the vm to shut down
					time.Sleep(1 * time.Minute)
				},
				Config: testVMV2Config(name, desc, "OFF") + testNGTConfig() + testVmsShutdownV2Config("shutdown") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "OFF"),
				),
			},
			// 5. power on the vm
			{
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
			// 6. check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to power on
					time.Sleep(1 * time.Minute)
				},
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_virtual_machine_v2.rtest", "power_state", "ON"),
				),
			},
			// 7. reboot the vm
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to power on
					time.Sleep(1 * time.Minute)
				},
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig() + testVmsShutdownV2Config("reboot"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
			// 8. check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to reboot
					time.Sleep(1 * time.Minute)
				},
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig() + testVmsShutdownV2Config("reboot") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "ON"),
				),
			},
			// 9. guest_reboot the vm
			{
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig() + testVmsShutdownV2Config("guest_reboot"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
			// 10. check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 2 Minute to allow the vm to reboot
					time.Sleep(timeSleep)
				},
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig() + testVmsShutdownV2Config("guest_reboot") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "ON"),
				),
			},
			// 11. guest_shutdown the vm
			{
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig() + testVmsShutdownV2Config("guest_shutdown"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
			// 12. check the power state of the vm
			{
				PreConfig: func() {
					//sleep for 2 Minute to allow the vm to shut down
					time.Sleep(timeSleep)
				},
				Config: testVMV2Config(name, desc, "OFF") + testNGTConfig() + testVmsShutdownV2Config("guest_shutdown") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "OFF"),
				),
			},
			// 13. power on the vm to uninstall ngt and delete the vm
			{
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
			// 14. check the power state of the vm before uninstalling ngt
			{
				PreConfig: func() {
					//sleep for 1 Minute to allow the vm to power on
					time.Sleep(1 * time.Minute)
				},
				Config: testVMV2Config(name, desc, "ON") + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineV2Destroy(s *terraform.State) error {
	fmt.Println("Destroying VMs")
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine_v2" {
			continue
		}
		_, err := conn.VmmAPI.VMAPIInstance.GetVmById(utils.StringPtr(rs.Primary.ID))
		if err == nil {
			// delete the vm
			fmt.Printf("Deleting VM with ID: %s\n", rs.Primary.ID)
			_, errVM := conn.VmmAPI.VMAPIInstance.DeleteVmById(utils.StringPtr(rs.Primary.ID))
			if errVM != nil {
				return errVM
			}
		}
	}
	return nil
}

func TestAccV2NutanixVmsShutdownResource_WithError(t *testing.T) {
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

func testVMV2Config(name, desc, powerState string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = jsondecode(file("%[4]s"))
			preEnv = local.config.pre_env
			vmm    = local.config.vmm
		}

		data "nutanix_images_v2" "ngt-image" {
		  filter = "name eq '${local.preEnv.ngt_image.name}'"
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
			memory_size_bytes = 4 * 1024 * 1024 * 1024
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
			
			depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_images_v2.ngt-image]			
		}
		
`, name, desc, powerState, filepath)
}

func testNGTConfig() string {
	return `		
		resource "nutanix_ngt_installation_v2" "test" {
			ext_id = nutanix_virtual_machine_v2.rtest.id
			credential {
				username = local.vmm.ngt.credential.username
				password = local.vmm.ngt.credential.password
			}
			reboot_preference {
				schedule_type = "IMMEDIATE"
			}
			capablities = ["VSS_SNAPSHOT"]
			depends_on = [nutanix_virtual_machine_v2.rtest]
		}
	`
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
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = jsondecode(file("%[4]s"))
			preEnv = local.config.pre_env
			vmm    = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
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

`, name, desc, state, filepath)
}

const vmDataSource = ` 
		data "nutanix_virtual_machine_v2" "test"{
			ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			depends_on = [
				resource.nutanix_vm_shutdown_action_v2.vmShuts
			]
		}
`
