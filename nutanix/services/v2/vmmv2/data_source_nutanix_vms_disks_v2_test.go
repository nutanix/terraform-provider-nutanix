package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmDisks = "data.nutanix_vm_disks_v4.test"

func TestAccNutanixVmsDisksDataSourceV4_List(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceDisksV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.0.disk_address.#"),
					resource.TestCheckResourceAttr(datasourceNameVmDisks, "disks.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(datasourceNameVmDisks, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.0.backing_info.#"),
				),
			},
		},
	})
}

func TestAccNutanixVmsDisksDataSourceV4_WithFilters(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceDisksV4ConfigWithFilters(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVmDisks, "page", "0"),
					resource.TestCheckResourceAttr(datasourceNameVmDisks, "limit", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.0.disk_address.#"),
					resource.TestCheckResourceAttr(datasourceNameVmDisks, "disks.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(datasourceNameVmDisks, "disks.0.disk_address.0.index", "0"),
					resource.TestCheckResourceAttrSet(datasourceNameVmDisks, "disks.0.backing_info.#"),
				),
			},
		},
	})
}

func testVmsDatasourceDisksV4Config(name, desc string) string {
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
		resource "nutanix_vm_disks_v4" "resTest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 0
			}
			backing_info{
			  vm_disk{
				disk_size_bytes = 1073741824*2
				storage_container {
				  ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
				}
			  }
			}
		}

		data "nutanix_vm_disks_v4" "test"{
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_disks_v4.resTest
			]
		}
`, name, desc)
}

func testVmsDatasourceDisksV4ConfigWithFilters(name, desc string) string {
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
		resource "nutanix_vm_disks_v4" "resTest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 0
			}
			backing_info{
			  vm_disk{
				disk_size_bytes = 1073741824*2
				storage_container {
				  ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
				}
			  }
			}
		}

		data "nutanix_vm_disks_v4" "test"{
			page = 0
			limit = 1
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_disks_v4.resTest
			]
		}
`, name, desc)
}
