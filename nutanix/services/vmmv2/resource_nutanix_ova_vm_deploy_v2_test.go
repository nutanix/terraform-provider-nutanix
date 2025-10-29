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
	// ovaNameUpdated := fmt.Sprintf("tf-test-ova-updated-%d", r)

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

func TestAccV2NutanixOvaVmDeployResource_DeployVmFromOvaDoseNotExists(t *testing.T) {
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

func testOvaVMDeployResourceConfigDeployVMFromOva(vmName, vmDescription, ovaName string) string {
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

resource "nutanix_ova_vm_deploy_v2" "test" {
  ext_id = nutanix_ova_v2.test.id
  override_vm_config {
    name              = "${nutanix_virtual_machine_v2.ova-vm.name}-from-ova"
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
        vlan_mode     = "TRUNK"
        trunked_vlans = ["1"]
      }
    }
  }
  cluster_location_ext_id = local.cluster_ext_id
}

data "nutanix_virtual_machines_v2" "vm-from-ova"{
  filter = "name eq '${nutanix_ova_vm_deploy_v2.test.override_vm_config.0.name}' and cluster/extId eq '${local.cluster_ext_id}' and memorySizeBytes eq ${nutanix_ova_vm_deploy_v2.test.override_vm_config.0.memory_size_bytes}"
}


`, filepath, vmName, vmDescription, ovaName)
}

func testOvaVMDeployResourceConfigDeployVMFromOvaDoseNotExists(ovaName string) string {
	return fmt.Sprintf(`

locals {
  config = jsondecode(file("%[1]s"))
  vmm = local.config.vmm
}


data "nutanix_subnets_v2" "subnets" {
  filter = "name eq '${local.vmm.subnet_name}'"
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
        vlan_mode     = "TRUNK"
        trunked_vlans = ["1"]
      }
    }
  }
  cluster_location_ext_id = "1f8b1211-5adf-4da3-8b52-19c743b15aa1"
}
`, filepath)
}

func testOvaVMDeployResourceConfigCreateOvaFromValidURL(ovaName string) string {
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
