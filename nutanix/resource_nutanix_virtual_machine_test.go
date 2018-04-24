package nutanix

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixVirtualMachine_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNutanixVMConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.vm1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "spec.#", "1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "spec.0.resources.#", "1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "spec.0.resources.0.power_state", "NutanixPowerState"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "spec.0.resources.0.memory_size_mb", "memorySize"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "spec.0.resources.0.num_sockets", "NutanixNumSockets"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "spec.0.resources.0.num_vcpus_per_socket", "NutanixNumVCPUs"),
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
	for i := range s.RootModule().Resources {
		if s.RootModule().Resources[i].Type != "nutanix_virtual_machine" {
			continue
		}
		id := string(s.RootModule().Resources[i].Primary.ID)
		if id == "" {
			err := errors.New("ID is already set to the null string")
			return err
		}
		return nil
	}
	return nil
}

const testAccNutanixVMConfig = `
	resource "nutanix_virtual_machine" "vm1" {
		
	}
`
