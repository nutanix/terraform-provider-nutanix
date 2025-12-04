package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVMClone = "nutanix_vm_clone_v2.test"

func TestAccV2NutanixVmsCloneResource_Basic(t *testing.T) {
	// t.Skip("Skipping test as it requires Clone")
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	// stateOn := "power_on"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCloneV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMClone, "name", fmt.Sprintf(`%[1]s-clone`, name)),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_threads_per_core", "2"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsCloneResource_WithGuestCustomization(t *testing.T) {
	// t.Skip("Skipping test as it requires Clone")
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"
	// stateOn := "power_on"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCloneV2WithGuestCustomizationConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMClone, "name", fmt.Sprintf(`%[1]s-clone`, name)),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_threads_per_core", "2"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "memory_size_bytes", "8589934592"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsCloneResource_WithUefiBootConfig(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-vm-%d", r)
	desc := "test vm description"

	datasourceNameVMCloned := "data.nutanix_virtual_machine_v2.test"
	// stateOn := "power_on"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsCloneV2WithUefiBootConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMClone, "name", fmt.Sprintf(`%[1]s-clone`, name)),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "num_threads_per_core", "1"),
					resource.TestCheckResourceAttr(resourceNameVMClone, "memory_size_bytes", "8589934592"),
					// Check on the cloned VM details to verify the boot config
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "memory_size_bytes", "8589934592"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "num_sockets", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "num_threads_per_core", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "power_state", "OFF"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "boot_config.0.uefi_boot.0.boot_order.#", "3"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "boot_config.0.uefi_boot.0.boot_order.0", "CDROM"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "boot_config.0.uefi_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "boot_config.0.uefi_boot.0.boot_order.2", "NETWORK"),
					resource.TestCheckResourceAttr(datasourceNameVMCloned, "boot_config.0.uefi_boot.0.is_secure_boot_enabled", "false"),
				),
			},
		},
	})
}

func testVmsCloneV2Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = (jsondecode(file("%[3]s")))
			vmm    = local.config.vmm
			gs = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    ssh-authorized-keys:\n      - ssh-rsa DUMMYSSH mypass\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
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
			nics {
				network_info {
				  nic_type = "NORMAL_NIC"
				  subnet {
					ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
				  }
				  vlan_mode = "ACCESS"
				}
		    }
			// guest_customization{
			// 	config{
			// 		cloud_init{
			// 			cloud_init_script{
			// 				user_data{
			// 					value="${local.gs}"
			// 				}
			// 			}
			// 		}
			// 	}
			// }
		    boot_config {
			  legacy_boot {
			    boot_order = ["CDROM", "DISK", "NETWORK"]
			  }
		    }
		    power_state = "OFF"

			lifecycle{
				ignore_changes = [
					guest_customization,
					guest_tools
				]
			}
		}

		resource "nutanix_vm_clone_v2" "test" {
			vm_ext_id               = resource.nutanix_virtual_machine_v2.rtest.id
			name                 = "%[1]s-clone"
			num_sockets          = 2
			num_cores_per_socket = 2
			num_threads_per_core = 2
			// guest_customization{
			// 	config{
			// 		cloud_init{
			// 			cloud_init_script{
			// 				user_data{
			// 					value="${local.gs}"
			// 				}
			// 			}
			// 		}
			// 	}
			// }
			// boot_config {
			//   legacy_boot {
			// 	boot_device {
			// 	  boot_device_disk {
			// 		disk_address {
			// 		  bus_type = "IDE"
			// 		  index    = 0
			// 		}
			// 	  }
			// 	//   boot_device_nic {
			// 	// 	mac_address = ""
			// 	//   }
			// 	}
			// 	boot_order = ["CDROM", "DISK", "NETWORK"]
			//   }
			// }
		  }

`, name, desc, filepath)
}

func testVmsCloneV2WithGuestCustomizationConfig(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			clusterUUID = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = (jsondecode(file("%[3]s")))
			vmm    = local.config.vmm
			gs = base64encode("#cloud-config\nusers:\n  - name: ubuntu\n    ssh-authorized-keys:\n      - ssh-rsa DUMMYSSH mypass\n    sudo: ['ALL=(ALL) NOPASSWD:ALL']")
		}

		data "nutanix_subnets_v2" "subnet" {
		  filter = "name eq '${local.vmm.subnet_name}'"
		}
		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.clusterUUID}'"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "rtest"{
		  name				   = "%[1]s"
		  description		   = "%[2]s"
		  num_cores_per_socket = 1
		  num_sockets          = 1
		  memory_size_bytes    = 4 * 1024 * 1024 * 1024
		  cluster {
			ext_id = local.clusterUUID
		  }
		  disks {
			disk_address {
			  bus_type = "SCSI"
			  index    = 0
			}
			backing_info {
			  vm_disk {
				disk_size_bytes = "1073741824"
				storage_container {
				  ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
				}
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
		  power_state = "OFF"
		  lifecycle {
			ignore_changes = [
			  guest_customization,
			  guest_tools
			]
		  }
		}

		resource "nutanix_vm_clone_v2" "test" {
		  vm_ext_id            = resource.nutanix_virtual_machine_v2.rtest.id
		  name                 = "%[1]s-clone"
		  num_sockets          = 2
		  num_cores_per_socket = 2
		  num_threads_per_core = 2
		  memory_size_bytes    = 8 * 1024 * 1024 * 1024

		  guest_customization {
			config {
			  cloud_init {
				cloud_init_script {
				  user_data {
					value = local.gs
				  }
				}
			  }
			}
		  }
		  boot_config {
			legacy_boot {
			  boot_device {
				boot_device_disk {
				  disk_address {
					bus_type = "SCSI"
					index    = 0
				  }
				}
			  }
			  boot_order = ["NETWORK", "DISK", "CDROM"]
			}
		  }
		  lifecycle {
			ignore_changes = [
			  guest_customization,
			  guest_tools
			]
		  }
		}

`, name, desc, filepath)
}

func testVmsCloneV2WithUefiBootConfig(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster0 = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
    ][
    0
  ]
}


resource "nutanix_virtual_machine_v2" "test" {
  name				   = "%[1]s"
  description		   = "%[2]s"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_bytes    = 8 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster0
  }

  boot_config {
    uefi_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }

  power_state = "OFF"
}


resource "nutanix_vm_clone_v2" "test"{
  vm_ext_id = nutanix_virtual_machine_v2.test.id
  name = "%[1]s-clone"
}

// data source to get the cloned vm details
data "nutanix_virtual_machine_v2" "test" {
  ext_id = nutanix_vm_clone_v2.test.id
}

`, name, desc)
}
