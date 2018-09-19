package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixVirtualMachinesDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMSSDataSourceConfig(acctest.RandIntRange(0, 100), acctest.RandIntRange(200, 300)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_virtual_machines.unittest", "entities.#"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccVMSSDataSourceConfig(rNumVM1 int, rNumVM2 int) string {
	return fmt.Sprintf(`
#data "nutanix_clusters" "clusters" {}

#output "cluster" {
#  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
#}

#resource "nutanix_virtual_machine" "vm1" {
#   name = "unittest-dou-vm%d"

#   cluster_reference = {
# 	  kind = "cluster"
# 	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
#   }

#   num_vcpus_per_socket = 1
#   num_sockets          = 1
#   memory_size_mib      = 186

#}

#resource "nutanix_virtual_machine" "vm2" {
#   name = "unittest-dou-vm%d"

#   cluster_reference = {
# 	  kind = "cluster"
# 	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
#   }

#   num_vcpus_per_socket = 1
#   num_sockets          = 1
#   memory_size_mib      = 186

#}

data "nutanix_virtual_machines" "unittest" {
#	depends_on = ["nutanix_virtual_machine.vm1", "nutanix_virtual_machine.vm2"]
}`, rNumVM1, rNumVM2)
}
