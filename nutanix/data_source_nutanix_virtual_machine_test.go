package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixVirtualMachineDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_sockets", "1"),
				),
			},
		},
	})
}

func testAccVMDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
		cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
		? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-%d"
  cluster_uuid = "${local.cluster1}"
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
}

data "nutanix_virtual_machine" "nutanix_virtual_machine" {
	vm_id = "${nutanix_virtual_machine.vm1.id}"
}
`, r)
}
