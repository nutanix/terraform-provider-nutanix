package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmNetworkDevices = "data.nutanix_vm_network_devices_v2.test"

func TestAccNutanixVmNetworkDeviceDataSourceV2_List(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceNetworkDevicesV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmNetworkDevices, "network_devices.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmNetworkDevices, "network_devices.0.network_info.#"),
				),
			},
		},
	})
}
func TestAccNutanixVmNetworkDeviceDataSourceV2_ListWithFilters(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceNetworkDevicesV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmNetworkDevices, "network_devices.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmNetworkDevices, "network_devices.0.network_info.#"),
				),
			},
			{
				Config: testVmsDatasourceNetworkDevicesV4ConfigWithFilters(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVmNetworkDevices, "network_devices.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmNetworkDevices, "network_devices.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmNetworkDevices, "network_devices.0.network_info.#"),
				),
			},
		},
	})
}

func testVmsDatasourceNetworkDevicesV4Config(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		config = jsondecode(file("%[3]s"))
		vmm = local.config.vmm
		}
	
		resource "nutanix_virtual_machine_v2" "rtest"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			lifecycle {
				ignore_changes = [
					nics,
				]
			}
		}

		resource "nutanix_vm_network_device_v2" "ntest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			network_info {
			  nic_type = "DIRECT_NIC"
			  subnet {
				ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
			  }
			}
		}

		data "nutanix_vm_network_devices_v2" "test"{
			vm_ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			depends_on = [
				resource.nutanix_vm_network_device_v2.ntest
			]
		}
`, name, desc, filepath)
}

func testVmsDatasourceNetworkDevicesV4ConfigWithFilters(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		data "nutanix_subnets_v2" "subnets" {
			filter = "name eq '${local.vmm.subnet_name}'"
		}

		locals {
			cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		resource "nutanix_virtual_machine_v2" "rtest"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			lifecycle {
				ignore_changes = [
					nics,
				]
			}
		}

		resource "nutanix_vm_network_device_v2" "ntest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			network_info {
			nic_type = "DIRECT_NIC"
			subnet {
				ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
			}
			}
		}

		data "nutanix_vm_network_devices_v2" test{
			page=0
			limit=1
			vm_ext_id = resource.nutanix_virtual_machine_v2.rtest.id
			// depends_on = [
			// 	resource.nutanix_vm_network_device_v2.ntest
			// ]
		}
`, name, desc, filepath)
}
