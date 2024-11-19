package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmCdrom = "data.nutanix_vm_cdrom_v4.test"

func TestAccNutanixVmsCdromDataSourceV4_Basic(t *testing.T) {
	t.Skip("Skipping test as it merged in the virtual_machine_v2 resource")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceCdromV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmCdrom, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdrom, "iso_type"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdrom, "disk_address.#"),
					resource.TestCheckResourceAttr(datasourceNameVmCdrom, "disk_address.0.bus_type", "IDE"),
					resource.TestCheckResourceAttr(datasourceNameVmCdrom, "disk_address.0.index", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdrom, "backing_info.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdrom, "backing_info.0.storage_container.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdrom, "backing_info.0.disk_size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameVmCdrom, "backing_info.0.is_migration_in_progress"),
				),
			},
		},
	})
}

func testVmsDatasourceCdromV4Config(name, desc string) string {
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

		data "nutanix_vm_cdrom_v4" test{
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			ext_id = resource.nutanix_vm_cdroms_v4.rtest.id
		}
`, name, desc)
}
