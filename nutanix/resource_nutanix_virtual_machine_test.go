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
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "staging"),
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
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "production"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_list"},
			},
		},
	})
}

func TestAccNutanixVirtualMachine_WithDisk(t *testing.T) {
	r := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfigWithDisk(r),
			},
			{
				Config:             testAccNutanixVMConfigWithDisk(r),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccNutanixVMConfigWithDiskUpdate(r),
			},
			{
				Config:             testAccNutanixVMConfigWithDiskUpdate(r),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:            "nutanix_virtual_machine.vm1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_list"},
			},
		}})
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

	categories {
		environment-terraform = "staging"
	}
}
`, r)
}

func testAccNutanixVMConfigWithDisk(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_image" "cirros-034-disk" {
    name        = "test-image-dou-%[1]d"
    source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
    description = "heres a tiny linux image, not an iso, but a real disk!"
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-vm-%[1]d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"

	disk_list = [{
		# data_source_reference in the Nutanix API refers to where the source for
		# the disk device will come from. Could be a clone of a different VM or a
		# image like we're doing here
		data_source_reference = [{
			kind = "image"
			uuid = "${nutanix_image.cirros-034-disk.id}"
		}]
		disk_size_mib = 44
	},
	{
		disk_size_mib = 100
	},
	{
		disk_size_mib = 200
	},
	{
		disk_size_mib = 300
	}]
}
`, r)
}

func testAccNutanixVMConfigWithDiskUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_image" "cirros-034-disk" {
    name        = "test-image-dou-%[1]d"
    source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
    description = "heres a tiny linux image, not an iso, but a real disk!"
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-vm-%[1]d"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"

	disk_list = [{
		# data_source_reference in the Nutanix API refers to where the source for
		# the disk device will come from. Could be a clone of a different VM or a
		# image like we're doing here
		data_source_reference = [{
			kind = "image"
			uuid = "${nutanix_image.cirros-034-disk.id}"
		}]
		disk_size_mib = 44
	},
	{
		disk_size_mib = 100
	},
	{
		disk_size_mib = 400
	}]
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

	categories {
		environment-terraform = "production"
	}
}
`, r)
}
