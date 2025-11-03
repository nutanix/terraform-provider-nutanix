package vmmv2_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameOva = "nutanix_ova_v2.test"
const resourceNameVMForOva = "nutanix_virtual_machine_v2.ova-vm"

func TestAccV2NutanixOvaResource_CreateOvaFromVM(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-%d", r)
	vmDescription := "VM for OVA terraform testing"
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)
	ovaNameUpdated := fmt.Sprintf("tf-test-ova-updated-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOvaResourceConfigCreateOvaFromVM(vmName, vmDescription, ovaName),
				Check: resource.ComposeTestCheckFunc(
					// vm checks
					resource.TestCheckResourceAttrSet(resourceNameVMForOva, "id"),
					resource.TestCheckResourceAttrSet(resourceNameVMForOva, "cluster.0.ext_id"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "name", vmName),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "description", vmDescription),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "num_threads_per_core", "2"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "num_cores_per_socket", "4"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "memory_size_bytes", strconv.Itoa(4*1024*1024*1024)), // 4 GiB
					resource.TestCheckResourceAttr(resourceNameVMForOva, "boot_config.0.legacy_boot.0.boot_order.0", "CDROM"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "boot_config.0.legacy_boot.0.boot_order.1", "DISK"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "boot_config.0.legacy_boot.0.boot_order.2", "NETWORK"),
					resource.TestCheckResourceAttrSet(resourceNameVMForOva, "disks.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameVMForOva, "disks.0.backing_info.0.vm_disk.0.disk_ext_id"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "disks.0.backing_info.0.vm_disk.0.disk_size_bytes", strconv.Itoa(10*1024*1024*1024)), // 10 GiB
					resource.TestCheckResourceAttrSet(resourceNameVMForOva, "disks.0.backing_info.0.vm_disk.0.storage_container.0.ext_id"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttrSet(resourceNameVMForOva, "nics.0.ext_id"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttrSet(resourceNameVMForOva, "nics.0.network_info.0.subnet.0.ext_id"),
					resource.TestCheckResourceAttr(resourceNameVMForOva, "nics.0.network_info.0.vlan_mode", "TRUNK"),

					// ova checks
					resource.TestCheckResourceAttrSet(resourceNameOva, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "cluster_location_ext_ids.0"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "create_time"),
					resource.TestCheckResourceAttr(resourceNameOva, "name", ovaName),
					resource.TestCheckResourceAttrPair(resourceNameOva, "parent_vm", resourceNameVMForOva, "name"),
					resource.TestCheckResourceAttr(resourceNameOva, "disk_format", "QCOW2"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "source.0.ova_vm_source.0.vm_ext_id", resourceNameVMForOva, "id"),
					resource.TestCheckResourceAttr(resourceNameOva, "source.0.ova_vm_source.0.disk_file_format", "QCOW2"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.description", resourceNameVMForOva, "description"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.memory_size_bytes", resourceNameVMForOva, "memory_size_bytes"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.name", resourceNameVMForOva, "name"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.num_sockets", resourceNameVMForOva, "num_sockets"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.num_cores_per_socket", resourceNameVMForOva, "num_cores_per_socket"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.num_threads_per_core", resourceNameVMForOva, "num_threads_per_core"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.machine_type", resourceNameVMForOva, "machine_type"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.0", resourceNameVMForOva, "boot_config.0.legacy_boot.0.boot_order.0"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.1", resourceNameVMForOva, "boot_config.0.legacy_boot.0.boot_order.1"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.2", resourceNameVMForOva, "boot_config.0.legacy_boot.0.boot_order.2"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.disks.#", resourceNameVMForOva, "disks.#"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.disks.0.backing_info.0.vm_disk.0.disk_size_bytes", resourceNameVMForOva, "disks.0.backing_info.0.vm_disk.0.disk_size_bytes"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.disks.0.disk_address.0.bus_type", resourceNameVMForOva, "disks.0.disk_address.0.bus_type"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.nics.#", resourceNameVMForOva, "nics.#"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.nics.0.network_info.0.nic_type", resourceNameVMForOva, "nics.0.network_info.0.nic_type"),
				),
			},
			// update ova
			{
				Config: testOvaResourceConfigCreateOvaFromVM(vmName, vmDescription, ovaNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameOva, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "cluster_location_ext_ids.0"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "size_bytes"),
					resource.TestCheckResourceAttr(resourceNameOva, "name", ovaNameUpdated),
				),
			},
		},
	})
}

func TestAccV2NutanixOvaResource_CreateOvaFromVMDoseNotExists(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testOvaResourceConfigCreateOvaFromVMDoseNotExists(ovaName),
				ExpectError: regexp.MustCompile("Failed to perform the operation as the VM entity is not present in the backend database."),
			},
		},
	})
}

func TestAccV2NutanixOvaResource_CreateOvaFromValidUrl(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)
	ovaNameUpdated := fmt.Sprintf("tf-test-ova-updated-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOvaResourceConfigCreateOvaFromValidURL(ovaName),
				Check: resource.ComposeTestCheckFunc(
					// ova checks
					resource.TestCheckResourceAttrSet(resourceNameOva, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "cluster_location_ext_ids.0"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "create_time"),
					resource.TestCheckResourceAttr(resourceNameOva, "name", ovaName),
					resource.TestCheckResourceAttr(resourceNameOva, "source.0.ova_url_source.0.url", testVars.VMM.OvaURL),
					resource.TestCheckResourceAttr(resourceNameOva, "source.0.ova_url_source.0.should_allow_insecure_url", "true"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "vm_config.#"),
				),
			},
			// update ova
			{
				Config: testOvaResourceConfigCreateOvaFromValidURL(ovaNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameOva, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "cluster_location_ext_ids.0"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "size_bytes"),
					resource.TestCheckResourceAttr(resourceNameOva, "name", ovaNameUpdated),
					resource.TestCheckResourceAttr(resourceNameOva, "source.0.ova_url_source.0.url", testVars.VMM.OvaURL),
					resource.TestCheckResourceAttr(resourceNameOva, "source.0.ova_url_source.0.should_allow_insecure_url", "true"),
					resource.TestCheckResourceAttrSet(resourceNameOva, "vm_config.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixOvaResource_CreateOvaFromInvalidUrl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "nutanix_clusters_v2" "clusters" {
						filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
						limit  = 1
					}
					locals {
						cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
					}
					resource "nutanix_ova_v2" "test" {
						name = "tf-test-ova-invalid-url"
						source {
							ova_url_source {
								url = "https://invalid-url.com/path/to/ova/file.ova"
								should_allow_insecure_url = true
							}
						}
						cluster_location_ext_ids = [local.cluster_ext_id]
						disk_format = "QCOW2"
					}
					`,
				ExpectError: regexp.MustCompile("Failed to upload the OVA as the dependent backend service has encountered an error."),
			},
		},
	})
}

func testOvaResourceConfigCreateOvaFromVM(vmName, vmDescription, ovaName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
  config = jsondecode(file("%[1]s"))
  vmm = local.config.vmm
}


data "nutanix_subnets_v2" "subnets" {
  filter = "name eq '${local.vmm.subnet_name}'"
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}


resource "nutanix_virtual_machine_v2" "ova-vm" {
  name        = "%[2]s"
  description = "%[3]s"
  num_sockets = 2
  num_threads_per_core = 2
  num_cores_per_socket = 4
  cluster {
    ext_id = local.cluster_ext_id
  }
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        disk_size_bytes = 10 * 1024 * 1024 * 1024 # 10 GiB
        storage_container {
          ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
        }
      }
    }
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
      }
      vlan_mode     = "TRUNK"
      trunked_vlans = ["1"]
    }
  }
  memory_size_bytes = 4 * 1024 * 1024 * 1024 # 4 GiB
  power_state = "OFF"
}


resource "nutanix_ova_v2" "test" {
  name = "%[4]s"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.ova-vm.id
      disk_file_format = "QCOW2"
    }
  }
}

`, filepath, vmName, vmDescription, ovaName)
}

func testOvaResourceConfigCreateOvaFromVMDoseNotExists(ovaName string) string {
	return fmt.Sprintf(`
resource "nutanix_ova_v2" "test" {
  name = "%[1]s"
  source {
    ova_vm_source {
      vm_ext_id        = "123e4567-e89b-12d3-a456-426614174000"
      disk_file_format = "QCOW2"
    }
  }
}
`, ovaName)
}

func testOvaResourceConfigCreateOvaFromValidURL(ovaName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
  config = jsondecode(file("%[1]s"))
  vmm = local.config.vmm
}

resource "nutanix_ova_v2" "test" {
  name = "%[2]s"
  source {
    ova_url_source {
      url = local.vmm.ova_url
      should_allow_insecure_url = true
    }
  }
  cluster_location_ext_ids = [local.cluster_ext_id]
}
`, filepath, ovaName)
}
