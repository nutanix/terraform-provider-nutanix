package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixVirtualMachineClone_basic(t *testing.T) {
	r := acctest.RandInt()
	vmName := acctest.RandomWithPrefix("test-clone-vm")
	resourceName := "nutanix_virtual_machine_clone.vm2"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMCloneConfig(r, vmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceName, "name", vmName),
				),
			},
		},
	})
}

func TestAccNutanixVirtualMachineClone_WithBootDeviceOrderChange(t *testing.T) {
	r := acctest.RandInt()
	vmName := acctest.RandomWithPrefix("test-clone-vm")
	resourceName := "nutanix_virtual_machine_clone.vm2"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMCloneConfigWithBootOrder(r, vmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceName, "name", vmName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "boot_device_order_list.0", "DISK"),
					resource.TestCheckResourceAttr(resourceName, "boot_device_order_list.1", "NETWORK"),
					resource.TestCheckResourceAttr(resourceName, "boot_device_order_list.2", "CDROM"),
				),
			},
		},
	})
}

func testAccNutanixVMCloneConfig(r int, name string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_virtual_machine" "vm1" {
			name         = "test-dou-%d"
			cluster_uuid = "${local.cluster1}"

			boot_device_order_list = ["DISK", "CDROM"]
			boot_type            = "LEGACY"
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186
		}

		resource "nutanix_virtual_machine_clone" "vm2"{
			vm_uuid = nutanix_virtual_machine.vm1.id
			name = "%s"
			num_vcpus_per_socket = 2
		}
	`, r, name)
}

func testAccNutanixVMCloneConfigWithBootOrder(r int, name string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

			locals {
				cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
				? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
			}

			resource "nutanix_virtual_machine" "vm1" {
				name         = "test-dou-%d"
				cluster_uuid = "${local.cluster1}"

				boot_device_order_list = ["DISK", "CDROM", "NETWORK"]
				boot_type            = "LEGACY"
				num_vcpus_per_socket = 1
				num_sockets          = 1
				memory_size_mib      = 186
			}
			resource "nutanix_virtual_machine_clone" "vm2"{
				vm_uuid = nutanix_virtual_machine.vm1.id
				name = "%s"
				num_vcpus_per_socket = 2
				num_sockets          = 2
				boot_device_order_list = ["DISK","NETWORK","CDROM"]
			}
	`, r, name)
}
