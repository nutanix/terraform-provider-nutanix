package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmSerialPorts = "data.nutanix_vm_serial_ports_v4.test"

func TestAccNutanixVmsSerialPortsDataSourceV4_List(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceVmSerialPortsConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmSerialPorts, "serial_ports.#"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.#", "2"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.0.is_connected", "false"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.0.index", "2"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.1.is_connected", "true"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.1.index", "3"),
				),
			},
		},
	})
}

func TestAccNutanixVmsSerialPortsDataSourceV4_ListWithFilters(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceVmSerialPortsConfigWithFilters(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmSerialPorts, "serial_ports.#"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.0.is_connected", "false"),
					resource.TestCheckResourceAttr(datasourceNameVmSerialPorts, "serial_ports.0.index", "2"),
				),
			},
		},
	})
}

func testVmsDatasourceVmSerialPortsConfig(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}
		resource "nutanix_vm_serial_ports_v4" "rtest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			is_connected = false
			index = 2
		}

		resource "nutanix_vm_serial_ports_v4" "rtest2" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			is_connected = true
			index = 3
			depends_on = [
				resource.nutanix_vm_serial_ports_v4.rtest
			]
		}

		data "nutanix_vm_serial_ports_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_serial_ports_v4.rtest,
				resource.nutanix_vm_serial_ports_v4.rtest2
			]
		}
`, name, desc)
}

func testVmsDatasourceVmSerialPortsConfigWithFilters(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
		cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
		}
	
		resource "nutanix_virtual_machine_v2" "test"{
			name= "%[1]s"
			description =  "%[2]s"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
		}
		resource "nutanix_vm_serial_ports_v4" "rtest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			is_connected = false
			index = 2
		}

		resource "nutanix_vm_serial_ports_v4" "rtest2" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			is_connected = true
			index = 3
			depends_on = [
				resource.nutanix_vm_serial_ports_v4.rtest
			]
		}

		data "nutanix_vm_serial_ports_v4" "test" {
			page=0
			limit=1
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_serial_ports_v4.rtest,
				resource.nutanix_vm_serial_ports_v4.rtest2
			]
		}
`, name, desc)
}
