package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixVirtualMachine_basic(t *testing.T) {
	r := acctest.RandInt()
	resourceName := "nutanix_virtual_machine.vm1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
				),
			},
			{
				Config: testAccNutanixVMConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr(resourceName, "power_state", "ON"),
					resource.TestCheckResourceAttr(resourceName, "memory_size_mib", "186"),
					resource.TestCheckResourceAttr(resourceName, "num_sockets", "2"),
					resource.TestCheckResourceAttr(resourceName, "num_vcpus_per_socket", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine" {
			continue
		}
		for {
			_, err := conn.API.V3.GetVM(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}

	}

	return nil

}

func testAccNutanixVMConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}
output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-%d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"
}
`, r)
}

func testAccNutanixVMConfigUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}


output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-%d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 2
  memory_size_mib      = 186
  power_state          = "ON"
}
`, r)
}
