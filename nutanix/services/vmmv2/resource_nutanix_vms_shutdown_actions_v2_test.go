package vmmv2_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vm"
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

	// Best-practice acceptance test for the VM power-action resource:
	//   * A single resource.Test exercises the full apply/destroy lifecycle of the
	//     resource across every action (shutdown, reboot, guest_reboot, guest_shutdown)
	//     as ordered steps.
	//   * shutdown/reboot and their guest_* variants are all handled asynchronously
	//     (ACPI signals for the guest ones), so we never assert power state from a racy
	//     immediate read. Instead the action resource blocks until the VM reaches its
	//     terminal power state, and each step asserts the live power state via the
	//     nutanix_virtual_machine_v2 data source (which depends_on the action) rather
	//     than relying on time.Sleep.
	//   * The shutdown/guest_shutdown steps ignore_changes on the VM's power_state so
	//     the async transition to OFF does not race the SDK idempotency plan; the OFF
	//     result is proven by the data source instead of via plan drift.
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineV2Destroy,
		Steps: []resource.TestStep{
			// 1. create a vm
			{
				Config: testVMV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_machine_v2.rtest", "id"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine_v2.rtest", "power_state", "ON"),
				),
			},
			// 2. install ngt on the vm (needs the guest booted, hence the settle window)
			{
				PreConfig: func() {
					time.Sleep(1 * time.Minute)
				},
				Config: testVMV2Config(name, desc) + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_virtual_machine_v2.rtest", "power_state", "ON"),
				),
			},
			// 3. shutdown: action blocks until the VM reaches OFF. power_state is ignored on
			// the VM resource so the async OFF does not race the idempotency plan; the OFF
			// state is verified via the data source.
			{
				Config: testVMV2ConfigIgnorePower(name, desc) + testNGTConfig() + testVmsShutdownV2Config("shutdown") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "OFF"),
				),
			},
			// 4. power the vm back on
			{
				Config: testVMV2Config(name, desc) + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_virtual_machine_v2.rtest", "power_state", "ON"),
				),
			},
			// 5. reboot: VM ends ON, so no drift and no ExpectNonEmptyPlan
			{
				Config: testVMV2Config(name, desc) + testNGTConfig() + testVmsShutdownV2Config("reboot") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "ON"),
				),
			},
			// 6. guest_reboot: VM ends ON
			{
				Config: testVMV2Config(name, desc) + testNGTConfig() + testVmsShutdownV2Config("guest_reboot") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "ON"),
				),
			},
			// 7. guest_shutdown: action blocks until the VM reaches OFF. power_state is ignored
			// on the VM resource so the async OFF does not race the idempotency plan; the OFF
			// state is verified via the data source.
			{
				Config: testVMV2ConfigIgnorePower(name, desc) + testNGTConfig() + testVmsShutdownV2Config("guest_shutdown") + vmDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMShutdown, "power_state", "OFF"),
				),
			},
			// 8. power the vm back on so ngt uninstall/vm delete on destroy succeeds
			{
				Config: testVMV2Config(name, desc) + testNGTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_virtual_machine_v2.rtest", "power_state", "ON"),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineV2Destroy(s *terraform.State) error {
	fmt.Println("Destroying VMs")
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine_v2" {
			continue
		}
		getVmByIdRequest := import1.GetVmByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.VmmAPI.VMAPIInstance.GetVmById(ctx, &getVmByIdRequest)
		if err == nil {
			// delete the vm
			fmt.Printf("Deleting VM with ID: %s\n", rs.Primary.ID)
			deleteVmByIdRequest := import1.DeleteVmByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			_, errVM := conn.VmmAPI.VMAPIInstance.DeleteVmById(ctx, &deleteVmByIdRequest)
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

func testVMV2Config(name, desc string, powerState ...string) string {
	ps := "ON"
	if len(powerState) > 0 {
		ps = powerState[0]
	}
	return testVMV2ConfigWithLifecycle(name, desc, ps, "[guest_tools]")
}

// testVMV2ConfigIgnorePower builds the VM config with power_state added to
// ignore_changes. It is used by the shutdown/guest_shutdown steps: the action
// (not the VM resource) drives the VM to OFF asynchronously, so ignoring
// power_state keeps the post-apply idempotency plan deterministically empty
// while the actual OFF state is still asserted via the nutanix_virtual_machine_v2
// data source. The power-on steps use testVMV2Config so power_state stays active.
func testVMV2ConfigIgnorePower(name, desc string) string {
	return testVMV2ConfigWithLifecycle(name, desc, "ON", "[guest_tools, power_state]")
}

func testVMV2ConfigWithLifecycle(name, desc, ps, ignoreChanges string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = jsondecode(file("%[4]s"))
			images = local.config.images
			vmm    = local.config.vmm
		}

		data "nutanix_images_v2" "ngt-image" {
		  filter = "name eq '${local.images.ngt_image}'"
		}

		data "nutanix_subnets_v2" "subnet" {
		  filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info {
				  virtual_ethernet_nic_network_info {
					nic_type = "NORMAL_NIC"
					subnet {
					  ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
					}
					vlan_mode = "ACCESS"
				  }
				}
			}
			
			boot_config {
				legacy_boot {
				  boot_order = ["CDROM", "DISK", "NETWORK"]
				}
			}
			
			power_state = "%[3]s"
			
			lifecycle {
				ignore_changes = %[5]s
			}
			
			depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_images_v2.ngt-image]			
		}
		
`, name, desc, ps, filepath, ignoreChanges)
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
			vmm    = local.config.vmm
		}

		data "nutanix_subnets_v2" "subnets" {
		  filter = "name eq '${local.vmm.subnet_name}'"
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
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
				nic_network_info{
					virtual_ethernet_nic_network_info{
						nic_type = "NORMAL_NIC"
						subnet{
							ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
						}	
						vlan_mode = "ACCESS"
					}
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
