package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmNetworkDevice = "nutanix_vm_network_device_v2.test"

func TestAccNutanixVmsNetworkDeviceV2Resource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNICVmPreEnvConfig2(vmName) + testVmsNetworkDeviceV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmNetworkDevice, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameVmNetworkDevice, "backing_info.#"),
					resource.TestCheckResourceAttrSet(resourceNameVmNetworkDevice, "network_info.#"),
					resource.TestCheckResourceAttr(resourceNameVmNetworkDevice, "network_info.0.nic_type", "DIRECT_NIC"),
					resource.TestCheckResourceAttr(resourceNameVmNetworkDevice, "network_info.0.vlan_mode", "ACCESS"),
					resource.TestCheckResourceAttrSet(resourceNameVmNetworkDevice, "network_info.0.subnet.0.ext_id"),
				),
			},
		},
	})
}

func testVmsNetworkDeviceV2Config() string {
	return `
		data "nutanix_subnets_v2" "subnet" {
		  filter = "name eq '${local.vmm.subnet_name}'"
		}	

		resource "nutanix_vm_network_device_v2" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.nic-vm.id
			network_info {
			  nic_type = "DIRECT_NIC"
			  subnet {
				ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
			  }
			}
            depends_on = [data.nutanix_subnets_v2.subnet, resource.nutanix_virtual_machine_v2.nic-vm]
		}
`
}

func testNICVmPreEnvConfig2(vmName string) string {

	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}
		
		locals {
		  clusterUUID = [
			for cluster in data.nutanix_clusters.clusters.entities :
			cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		  ][0]
		  config = (jsondecode(file("%[1]s")))
		  vmm    = local.config.vmm
		}

	
		
		data "nutanix_storage_containers_v2" "sc" {
		  filter = "clusterExtId eq '${local.clusterUUID}'"
		  limit = 1
		}
		
				
		
		resource "nutanix_virtual_machine_v2" "nic-vm" {
		  name                 = "%[2]s"
		  description          = "vm to test network device "
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
				  ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
				}
			  }
			}
		  }
		  
		
		  
		  power_state = "OFF"

		  lifecycle {
			ignore_changes = [guest_tools, nics]
		  }
		
		  depends_on = [data.nutanix_clusters.clusters, data.nutanix_storage_containers_v2.sc]
		}
			
`, filepath, vmName)
}
