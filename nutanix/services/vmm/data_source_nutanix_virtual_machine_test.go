package vmm_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixVirtualMachineDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_sockets", "1"),
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "is_vcpu_hard_pinned", "true"),
				),
			},
		},
	})
}

func TestAccNutanixVirtualMachineDataSource_WithDisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigWithDisk(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_sockets", "1"),
				),
			},
		},
	})
}

func TestAccNutanixVirtualMachineDataSource_withDiskContainer(t *testing.T) {
	datasourceName := "data.nutanix_virtual_machine.nutanix_virtual_machine"
	vmName := acctest.RandomWithPrefix("test-dou-vm")
	containerUUID := os.Getenv("NUTANIX_STORAGE_CONTAINER")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceWithDiskContainer(vmName, containerUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "vm_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.0.disk_size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.0.disk_size_mib"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.0.storage_config.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.0.storage_config.0.storage_container_reference.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.0.storage_config.0.storage_container_reference.0.kind"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.0.storage_config.0.storage_container_reference.0.uuid"),
					resource.TestCheckResourceAttrSet(datasourceName, "disk_list.0.storage_config.0.storage_container_reference.0.name"),
				),
			},
		},
	})
}

func testAccVMDataSourceWithDiskContainer(vmName, containerUUID string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_virtual_machine" "vm-disk" {
			name                 = "%s"
			cluster_uuid         = local.cluster1
			num_vcpus_per_socket = 1
			num_sockets          = 1
			memory_size_mib      = 186

			disk_list {
				# disk_size_mib = 300
				disk_size_bytes = 68157440
				disk_size_mib   = 65

				storage_config {
					storage_container_reference {
						kind = "storage_container"
						uuid = "%s"
					}
				}
			}
		}

		data "nutanix_virtual_machine" "nutanix_virtual_machine" {
			vm_id = nutanix_virtual_machine.vm-disk.id
		}
	`, vmName, containerUUID)
}

func testAccVMDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-%d"
  cluster_uuid = local.cluster1
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  is_vcpu_hard_pinned  = true
}

data "nutanix_virtual_machine" "nutanix_virtual_machine" {
	vm_id = nutanix_virtual_machine.vm1.id
}
`, r)
}

func testAccVMDataSourceConfigWithDisk(r int) string {
	return fmt.Sprintf(`
	data "nutanix_clusters" "clusters" {}

	locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
	}

	resource "nutanix_image" "cirros-034-disk" {
		name        = "test-image-dou-vm-create-%[1]d"
		source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
		description = "heres a tiny linux image, not an iso, but a real disk!"
	}

	resource "nutanix_virtual_machine" "vm1" {
	  name = "test-dou-vm-%[1]d"
	  cluster_uuid = "${local.cluster1}"
	  num_vcpus_per_socket = 1
	  num_sockets          = 1
	  memory_size_mib      = 186

		disk_list {
			data_source_reference = {
				kind = "image"
				uuid = "${nutanix_image.cirros-034-disk.id}"
			}

			device_properties {
				disk_address = {
					device_index = 0,
					adapter_type = "SCSI"
				}
				device_type = "DISK"
			}
		}

		disk_list {
			disk_size_mib = 100
		}

		disk_list {
			disk_size_mib = 200
		}

		disk_list {
			disk_size_mib = 300
		}
	}

	data "nutanix_virtual_machine" "nutanix_virtual_machine" {
		vm_id = "${nutanix_virtual_machine.vm1.id}"
	}

	`, r)
}
