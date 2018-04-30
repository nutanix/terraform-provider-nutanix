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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNutanixVMConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.vm1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "power_state", "ON"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "memory_size_mib", "2048"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "num_sockets", "1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "num_vcpus_per_socket", "1"),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*NutanixClient)

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
provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}

variable clusterid {
  default = "000567f3-1921-c722-471d-0cc47ac31055"
}

resource "nutanix_virtual_machine" "vm1" {
  metadata {
    kind = "vm"
    name = "metadata-name-test-dou-%d"
  }

  name = "test-dou-%d"

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.clusterid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  power_state          = "ON"

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.test.id}"
    }

    ip_endpoint_list = {
			ip = "192.168.0.10"
			type = "ASSIGNED"
    }
  }]
}

resource "nutanix_subnet" "test" {
  metadata = {
    kind = "subnet"
  }

  name        = "dou_vlan0_test"
  description = "Dou Vlan 0"

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.clusterid}"
  }

  vlan_id     = 201
  subnet_type = "VLAN"

  prefix_length      = 24
  default_gateway_ip = "192.168.0.1"
  subnet_ip          = "192.168.0.0"

  dhcp_options {
    boot_file_name   = "bootfile"
    tftp_server_name = "192.168.0.252"
    domain_name      = "nutanix"
  }

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "calm.io"]
}

`, r, r+1)
}
