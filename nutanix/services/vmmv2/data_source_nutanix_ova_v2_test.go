package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameOva = "data.nutanix_ova_v2.test"

func TestAccV2NutanixOvaDatasource_GetOvaDetails(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-%d", r)
	vmDescription := "VM for OVA terraform testing"
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)

	config := testOvaResourceConfigCreateOvaFromVM(vmName, vmDescription, ovaName)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List all Ovas
			{
				Config: config + testOvasDatasourceConfigGetOvaDetails(),
				Check: resource.ComposeTestCheckFunc(
					// ova checks
					resource.TestCheckResourceAttrPair(resourceNameOva, "ext_id", datasourceNameOva, "ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "cluster_location_ext_ids.0", datasourceNameOva, "cluster_location_ext_ids.0"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "size_bytes", datasourceNameOva, "size_bytes"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "create_time", datasourceNameOva, "create_time"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "name", datasourceNameOva, "name"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "parent_vm", datasourceNameOva, "parent_vm"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "disk_format", datasourceNameOva, "disk_format"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.description", datasourceNameOva, "vm_config.0.description"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.memory_size_bytes", datasourceNameOva, "vm_config.0.memory_size_bytes"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.name", datasourceNameOva, "vm_config.0.name"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.num_sockets", datasourceNameOva, "vm_config.0.num_sockets"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.num_cores_per_socket", datasourceNameOva, "vm_config.0.num_cores_per_socket"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.num_threads_per_core", datasourceNameOva, "vm_config.0.num_threads_per_core"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.machine_type", datasourceNameOva, "vm_config.0.machine_type"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.0", datasourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.0"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.1", datasourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.1"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.2", datasourceNameOva, "vm_config.0.boot_config.0.legacy_boot.0.boot_order.2"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.disks.#", datasourceNameOva, "vm_config.0.disks.#"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.disks.0.backing_info.0.vm_disk.0.disk_size_bytes", datasourceNameOva, "vm_config.0.disks.0.backing_info.0.vm_disk.0.disk_size_bytes"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.disks.0.disk_address.0.bus_type", datasourceNameOva, "vm_config.0.disks.0.disk_address.0.bus_type"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.nics.#", datasourceNameOva, "vm_config.0.nics.#"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "vm_config.0.nics.0.network_info.0.nic_type", datasourceNameOva, "vm_config.0.nics.0.network_info.0.nic_type"),
				),
			},
		},
	})
}

func testOvasDatasourceConfigGetOvaDetails() string {
	return `
data "nutanix_ova_v2" "test" {
   ext_id = nutanix_ova_v2.test.id
}
`
}
