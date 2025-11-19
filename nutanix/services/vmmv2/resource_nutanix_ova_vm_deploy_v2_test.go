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

const resourceNameOvaVMDeploy = "nutanix_ova_vm_deploy_v2.test"
const datasourceVMFromOva = "data.nutanix_virtual_machines_v2.vm-from-ova"

func TestAccV2NutanixOvaVmDeployResource_DeployVMFromOva(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-%d", r)
	vmDescription := "VM for OVA terraform testing"
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOvaVMDeployResourceConfigDeployVMFromOva(vmName, vmDescription, ovaName),
				Check: resource.ComposeTestCheckFunc(
					// ova vm deploy  checks
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.name", fmt.Sprintf("%s-from-ova", vmName)),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.memory_size_bytes", strconv.Itoa(8*1024*1024*1024)), // 8GB
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.nics.#", "1"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.nics.0.network_info.0.nic_type", "NORMAL_NIC"),

					// vm data source checks
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.#", "1"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.name", fmt.Sprintf("%s-from-ova", vmName)),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.memory_size_bytes", strconv.Itoa(8*1024*1024*1024)), // 8GB
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.boot_config.0.legacy_boot.0.boot_order.#", "3"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.description", vmDescription),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.disks.#", "1"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.nics.#", "1"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.nics.0.network_info.0.nic_type", "NORMAL_NIC"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.num_cores_per_socket", "4"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.num_threads_per_core", "2"),
					resource.TestCheckResourceAttr(datasourceVMFromOva, "vms.0.power_state", "OFF"),
				),
			},
		},
	})
}

func TestAccV2NutanixOvaVmDeployResource_DeployVmFromOvaDoesNotExist(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testOvaVMDeployResourceConfigDeployVMFromOvaDoseNotExists(ovaName),
				ExpectError: regexp.MustCompile("Failed to perform the operation as the backend service could not find the entity."),
			},
		},
	})
}

func TestAccV2NutanixOvaVmDeployResource_BasicUpdate(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-basic-%d", r)
	vmNameUpdated := fmt.Sprintf("tf-test-vm-ova-basic-updated-%d", r)
	vmDescription := "VM for basic OVA update testing"
	vmDescriptionUpdated := "VM for basic OVA update testing - updated"
	ovaName := fmt.Sprintf("tf-test-ova-basic-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOvaVMDeployResourceConfigDeployVMFromOva(vmName, vmDescription, ovaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.name", fmt.Sprintf("%s-from-ova", vmName)),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.memory_size_bytes", strconv.Itoa(8*1024*1024*1024)), // 8GB
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_cores_per_socket", "4"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_threads_per_core", "2"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.power_state", "OFF"),
					resource.TestCheckResourceAttrSet(resourceNameOvaVMDeploy, "id"),
				),
			},
			{
				Config: testOvaVMDeployResourceConfigDeployVMFromOvaUpdated(vmNameUpdated, vmDescriptionUpdated, ovaName),
				Check: resource.ComposeTestCheckFunc(
					// Basic update test - just name and power state change
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.name", fmt.Sprintf("%s-from-ova", vmNameUpdated)),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.memory_size_bytes", strconv.Itoa(16*1024*1024*1024)), // 16GB
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_sockets", "4"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_threads_per_core", "1"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.power_state", "ON"),
					resource.TestCheckResourceAttrSet(resourceNameOvaVMDeploy, "id"),
				),
			},
		},
	})
}

func TestAccV2NutanixOvaVmDeployResource_FullUpdate(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-full-%d", r)
	vmNameUpdated := fmt.Sprintf("tf-test-vm-ova-full-updated-%d", r)
	vmDescription := "VM for comprehensive OVA update testing"
	vmDescriptionUpdated := "VM for comprehensive OVA update testing - updated"
	ovaName := fmt.Sprintf("tf-test-ova-full-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOvaVMDeployResourceConfigForFullUpdate(vmName, vmDescription, ovaName, "initial"),
				Check: resource.ComposeTestCheckFunc(
					// Initial configuration checks
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.name", fmt.Sprintf("%s-from-ova", vmName)),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.memory_size_bytes", strconv.Itoa(4*1024*1024*1024)), // 4GB
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_threads_per_core", "1"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.power_state", "OFF"),
					resource.TestCheckResourceAttrSet(resourceNameOvaVMDeploy, "id"),
				),
			},
			{
				Config: testOvaVMDeployResourceConfigForFullUpdate(vmNameUpdated, vmDescriptionUpdated, ovaName, "updated"),
				Check: resource.ComposeTestCheckFunc(
					// Updated configuration checks - all properties changed
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.name", fmt.Sprintf("%s-from-ova", vmNameUpdated)),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.memory_size_bytes", strconv.Itoa(12*1024*1024*1024)), // 12GB
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_sockets", "6"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.num_threads_per_core", "2"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.power_state", "ON"),
					resource.TestCheckResourceAttrSet(resourceNameOvaVMDeploy, "id"),
				),
			},
		},
	})
}

func testOvaVMDeployResourceConfigDeployVMFromOva(vmName, vmDescription, ovaName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

data "nutanix_subnets_v2" "subnets" {
  limit = 1
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}


resource "nutanix_virtual_machine_v2" "ova-vm" {
  name        = "%[1]s"
  description = "%[2]s"
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
  memory_size_bytes = 4 * 1024 * 1024 * 1024 # 4 GiB
  power_state = "OFF"
}


resource "nutanix_ova_v2" "test" {
  name = "%[3]s"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.ova-vm.id
      disk_file_format = "QCOW2"
    }
  }
}

resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = nutanix_ova_v2.test.id
  override_vm_config {
    name                 = "${nutanix_virtual_machine_v2.ova-vm.name}-from-ova"
    memory_size_bytes    = 8 * 1024 * 1024 * 1024 # 8 GiB
    num_sockets          = 2
    num_cores_per_socket = 4
    num_threads_per_core = 2
    power_state          = "OFF"
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
        }
        vlan_mode     = "ACCESS"
      }
    }
  }
  cluster_location_ext_id = local.cluster_ext_id
}

data "nutanix_virtual_machines_v2" "vm-from-ova"{
  filter = "name eq '${nutanix_ova_vm_deploy_v2.test.override_vm_config.0.name}' and cluster/extId eq '${local.cluster_ext_id}' and memorySizeBytes eq ${nutanix_ova_vm_deploy_v2.test.override_vm_config.0.memory_size_bytes}"
}

`, vmName, vmDescription, ovaName)
}

func testOvaVMDeployResourceConfigDeployVMFromOvaDoseNotExists(ovaName string) string {
	return `

data "nutanix_subnets_v2" "subnets" {
  limit = 1
}

resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = "9fdb1211-5adf-4da3-8b52-19c743b15aa1"
  override_vm_config {
    name              = "tf-test-vm-ova-from-ova"
    memory_size_bytes = 8 * 1024 * 1024 * 1024 # 8 GiB
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
        }
        vlan_mode     = "ACCESS"
      }
    }
  }
  cluster_location_ext_id = "1f8b1211-5adf-4da3-8b52-19c743b15aa1"
}
`
}

func testOvaVMDeployResourceConfigDeployVMFromOvaUpdated(vmName, vmDescription, ovaName string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

data "nutanix_subnets_v2" "subnets" {
  limit = 1
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

resource "nutanix_virtual_machine_v2" "ova-vm" {
  name        = "%[1]s"
  description = "%[2]s"
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
  memory_size_bytes = 4 * 1024 * 1024 * 1024 # 4 GiB
  power_state = "OFF"
}

resource "nutanix_ova_v2" "test" {
  name = "%[3]s"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.ova-vm.id
      disk_file_format = "QCOW2"
    }
  }
}

resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = nutanix_ova_v2.test.id
  override_vm_config {
    name                 = "${nutanix_virtual_machine_v2.ova-vm.name}-from-ova"
    memory_size_bytes    = 16 * 1024 * 1024 * 1024 # 16 GiB (updated)
    num_sockets          = 4                       # updated
    num_cores_per_socket = 2                       # updated  
    num_threads_per_core = 1                       # updated
    power_state          = "ON"                    # updated
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
        }
        vlan_mode     = "ACCESS"
      }
    }
  }
  cluster_location_ext_id = local.cluster_ext_id
}

`, vmName, vmDescription, ovaName)
}

func testOvaVMDeployResourceConfigForFullUpdate(vmName, vmDescription, ovaName, stage string) string {
	if stage == "initial" {
		// Initial configuration with smaller resources
		return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

data "nutanix_subnets_v2" "subnets" {
  limit = 1
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

resource "nutanix_virtual_machine_v2" "ova-vm" {
  name        = "%[1]s"
  description = "%[2]s"
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
  memory_size_bytes = 4 * 1024 * 1024 * 1024 # 4 GiB
  power_state = "OFF"
}

resource "nutanix_ova_v2" "test" {
  name = "%[3]s"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.ova-vm.id
      disk_file_format = "QCOW2"
    }
  }
}

resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = nutanix_ova_v2.test.id
  override_vm_config {
    name                 = "${nutanix_virtual_machine_v2.ova-vm.name}-from-ova"
    memory_size_bytes    = 4 * 1024 * 1024 * 1024 # 4 GiB - initial smaller size
    num_sockets          = 2                       # initial config
    num_cores_per_socket = 2                       # initial config
    num_threads_per_core = 1                       # initial config
    power_state          = "OFF"                   # initial state
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
        }
        vlan_mode     = "ACCESS"
      }
    }
  }
  cluster_location_ext_id = local.cluster_ext_id
}

`, vmName, vmDescription, ovaName)
	}
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

data "nutanix_subnets_v2" "subnets" {
  limit = 1
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

resource "nutanix_virtual_machine_v2" "ova-vm" {
  name        = "%[1]s"
  description = "%[2]s"
  num_sockets = 2
  num_threads_per_core = 2
  num_cores_per_socket = 1
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
  memory_size_bytes = 4 * 1024 * 1024 * 1024 # 4 GiB
  power_state = "OFF"
}

resource "nutanix_ova_v2" "test" {
  name = "%[3]s"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.ova-vm.id
      disk_file_format = "QCOW2"
    }
  }
}

resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = nutanix_ova_v2.test.id
  override_vm_config {
    name                 = "${nutanix_virtual_machine_v2.ova-vm.name}-from-ova"
    memory_size_bytes    = 12 * 1024 * 1024 * 1024 # 12 GiB - updated larger size
    num_sockets          = 6                        # updated config
    num_cores_per_socket = 1                        # updated config
    num_threads_per_core = 2                        # updated config
    power_state          = "ON"                     # updated state
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
        }
        vlan_mode     = "ACCESS"
      }
    }
  }
  cluster_location_ext_id = local.cluster_ext_id
}

`, vmName, vmDescription, ovaName)
}

func TestAccV2NutanixOvaVmDeployResource_DiskUpdate(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-disk-%d", r)
	vmDescription := "VM for OVA disk testing"
	ovaName := fmt.Sprintf("tf-test-ova-disk-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOvaVMDeployResourceConfigWithDisk(vmName, vmDescription, ovaName, "20"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.name", fmt.Sprintf("%s-from-ova", vmName)),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.#", "1"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.0.backing_info.0.vm_disk.0.disk_size_bytes", strconv.Itoa(20*1024*1024*1024)), // 15GB
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.0.disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(resourceNameOvaVMDeploy, "id"),
				),
			},
			{
				Config: testOvaVMDeployResourceConfigWithDisk(vmName, vmDescription, ovaName, "25"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.name", fmt.Sprintf("%s-from-ova", vmName)),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.#", "1"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.0.backing_info.0.vm_disk.0.disk_size_bytes", strconv.Itoa(25*1024*1024*1024)), // 25GB
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.0.disk_address.0.bus_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceNameOvaVMDeploy, "override_vm_config.0.disks.0.disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(resourceNameOvaVMDeploy, "id"),
				),
			},
		},
	})
}

func testOvaVMDeployResourceConfigWithDisk(vmName, vmDescription, ovaName, diskSizeGB string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(a:a eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

data "nutanix_subnets_v2" "subnets" {
  limit = 1
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

resource "nutanix_virtual_machine_v2" "ova-vm" {
  name        = "%[1]s"
  description = "%[2]s"
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
  memory_size_bytes = 4 * 1024 * 1024 * 1024 # 4 GiB
  power_state = "OFF"
}

resource "nutanix_ova_v2" "test" {
  name = "%[3]s"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.ova-vm.id
      disk_file_format = "QCOW2"
    }
  }
}

resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = nutanix_ova_v2.test.id
  override_vm_config {
    name                 = "${nutanix_virtual_machine_v2.ova-vm.name}-from-ova"
    memory_size_bytes    = 8 * 1024 * 1024 * 1024 # 8 GiB
    num_sockets          = 2
    num_cores_per_socket = 4
    num_threads_per_core = 2
    power_state          = "OFF"
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
        }
        vlan_mode     = "ACCESS"
      }
    }
    disks {
      disk_address {
        bus_type = "SCSI"
        index    = 1
      }
      backing_info {
        vm_disk {
          disk_size_bytes = %[4]s * 1024 * 1024 * 1024 # %[4]s GB additional disk
          storage_container {
            ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
          }
        }
      }
    }
  }
  cluster_location_ext_id = local.cluster_ext_id
}

`, vmName, vmDescription, ovaName, diskSizeGB)
}
