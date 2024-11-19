package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVmGpus = "data.nutanix_vm_gpus_v2.test"

func TestAccNutanixVmsGpusDataSourceV4_List(t *testing.T) {
	t.Skip("Skipping test as it requires a VM with GPU attached")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceGpusV4Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.pci_address.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.mode"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.vendor"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.guest_driver_version"),
				),
			},
		},
	})
}

func TestAccNutanixVmsGpusDataSourceV4_WithFilters(t *testing.T) {
	t.Skip("Skipping test as it requires a VM with GPU attached")
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-%d", r)
	desc := "test vm description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmsDatasourceGpusV4ConfigWithFilters(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVmGpus, "page", "0"),
					resource.TestCheckResourceAttr(datasourceNameVmGpus, "limit", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.pci_address.#"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.mode"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.vendor"),
					resource.TestCheckResourceAttrSet(datasourceNameVmGpus, "gpus.0.guest_driver_version"),
				),
			},
		},
	})
}

func testVmsDatasourceGpusV4Config(name, desc string) string {
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

		data "nutanix_vm_gpus_v2" "test"{
			vm_ext_id = resource.nutanix_virtual_machine_v2.test.id
			depends_on = [
				resource.nutanix_vm_disks_v4.resTest
			]
		}
`, name, desc)
}

func testVmsDatasourceGpusV4ConfigWithFilters(name, desc string) string {
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
