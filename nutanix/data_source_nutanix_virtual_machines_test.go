package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixVMSDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMSSDataSourceConfig(acctest.RandIntRange(0, 100), acctest.RandIntRange(200, 300)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machines.unittest", "entities.#", "2"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccVMSSDataSourceConfig(rNumVM1 int, rNumVM2 int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_virtual_machine" "vm1" {
  name = "unittest-dou-vm%d"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"

}

resource "nutanix_virtual_machine" "vm2" {
  name = "unittest-dou-vm%d"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"

}

data "nutanix_virtual_machines" "unittest" {
	metadata = {
		length = 2
	}
}`, rNumVM1, rNumVM2)
}
