package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmCdroms = "data.nutanix_vm_cdroms_v4.test"

func TestAccNutanixVmsCdromDataSourceV4_List(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceCdromsV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.iso_type"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.disk_address.#"),
					resource.TestCheckResourceAttr(datasourceNameVmCdroms, "cdroms.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(datasourceNameVmCdroms, "cdroms.0.disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.storage_container.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.disk_size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.is_migration_in_progress"),
				),
			},
		},
	})
}
func TestAccNutanixVmsCdromDataSourceV4_ListWithFilters(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceCdromsV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.iso_type"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.disk_address.#"),
					resource.TestCheckResourceAttr(datasourceNameVmCdroms, "cdroms.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(datasourceNameVmCdroms, "cdroms.0.disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.storage_container.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.disk_size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.is_migration_in_progress"),
				),
			},
			{
				Config: testVmsDatasourceCdromsV4ConfigWithFilters(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVmCdroms, "cdroms.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.iso_type"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.disk_address.#"),
					resource.TestCheckResourceAttr(datasourceNameVmCdroms, "cdroms.0.disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(datasourceNameVmCdroms, "cdroms.0.disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.storage_container.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.disk_size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdroms, "cdroms.0.backing_info.0.is_migration_in_progress"),
				),
			},
		},
	})
}
func TestAccNutanixVmsCdromDataSourceV4_ListNegative(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testVmsDatasourceCdromsV4ConfigWithInvalidFilters(name, desc),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config:      testVmsDatasourceCdromsV4ConfigWithInvalidPage(),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config:      testVmsDatasourceCdromsV4ConfigWithInvalidLimit(),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
		},
	})
}

func testVmsDatasourceCdromsV4Config(name, desc string) string {
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

		resource "nutanix_vm_cdroms_v4" "rtest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 1
			}
			backing_info{
				disk_size_bytes = 1073741824
				storage_container {
				ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
				}
			}
		}

		data "nutanix_vm_cdroms_v4" test{
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_cdroms_v4.rtest
			]
		}
`, name, desc)
}

func testVmsDatasourceCdromsV4ConfigWithFilters(name, desc string) string {
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

		resource "nutanix_vm_cdroms_v4" "rtest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 1
			}
			backing_info{
				disk_size_bytes = 1073741824
				storage_container {
				ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
				}
			}
		}

		data "nutanix_vm_cdroms_v4" test{
			page=0
			limit=1
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_cdroms_v4.rtest
			]
		}
`, name, desc)
}

func testVmsDatasourceCdromsV4ConfigWithInvalidFilters(name, desc string) string {
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

		resource "nutanix_vm_cdroms_v4" "rtest" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			disk_address{
			  bus_type = "IDE"
			  index= 1
			}
			backing_info{
				disk_size_bytes = 1073741824
				storage_container {
				ext_id = "10eb150f-e8b8-4d69-a828-6f23771d3723"
				}
			}
		}

		data "nutanix_vm_cdroms_v4" test{
			page=0
			limit=1
			filter= "name eq 'invalid'"
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_cdroms_v4.rtest
			]
		}
`, name, desc)
}

func testVmsDatasourceCdromsV4ConfigWithInvalidLimit() string {
	return (`

		data "nutanix_virtual_machines_v2" "test" {
		}
		data "nutanix_vm_cdroms_v4" test{
			page=-1
			filter= "name eq 'invalid'"
			vm_ext_id = data.nutanix_virtual_machines_v2.test.vms.0.ext_id
			depends_on = [
				data.nutanix_virtual_machines_v2.test
			]
		}
`)
}

func testVmsDatasourceCdromsV4ConfigWithInvalidPage() string {
	return (`

		data "nutanix_virtual_machines_v2" "test" {
		}
		data "nutanix_vm_cdroms_v4" test{
			limit=0
			filter= "name eq 'invalid'"
			vm_ext_id = data.nutanix_virtual_machines_v2.test.vms.0.ext_id
			depends_on = [
				data.nutanix_virtual_machines_v2.test
			]
		}
`)
}
