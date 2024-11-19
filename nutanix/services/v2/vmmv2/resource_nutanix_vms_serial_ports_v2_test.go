package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmSerialPorts = "nutanix_vm_serial_ports_v4.test"

func TestAccNutanixVmsSerialPortsV4_Basic(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsSerialPortsV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmSerialPorts, "links.#"),
					resource.TestCheckResourceAttr(resourceNameVmSerialPorts, "is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVmSerialPorts, "index", "2"),
				),
			},
		},
	})
}

func TestAccNutanixVmsSerialPortsV4_WithUpdate(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsSerialPortsV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmSerialPorts, "links.#"),
					resource.TestCheckResourceAttr(resourceNameVmSerialPorts, "is_connected", "false"),
					resource.TestCheckResourceAttr(resourceNameVmSerialPorts, "index", "2"),
				),
			},
			{
				Config: testVmsSerialPortsV4ConfigUpdate(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmSerialPorts, "links.#"),
					resource.TestCheckResourceAttr(resourceNameVmSerialPorts, "is_connected", "true"),
					resource.TestCheckResourceAttr(resourceNameVmSerialPorts, "index", "3"),
				),
			},
		},
	})
}

func testVmsSerialPortsV4Config(name, desc string) string {
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
		resource "nutanix_vm_serial_ports_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			is_connected = false
			index = 2
		  }
`, name, desc)
}

func testVmsSerialPortsV4ConfigUpdate(name, desc string) string {
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
		resource "nutanix_vm_serial_ports_v4" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			is_connected = true
			index = 3
		  }
`, name, desc)
}
