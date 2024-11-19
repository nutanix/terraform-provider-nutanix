package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmNetworkDeviceMigrate = "nutanix_vm_network_device_migrate_v2.test"

func TestAccNutanixVmsNetworkDeviceMigrateV2Resource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("nic-vm-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmPreEnvConfig(r) + testVmWithNicAndDiskConfig(vmName) + testVmsNetworkDeviceMigrateV4AssignConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmNetworkDeviceMigrate, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVmNetworkDeviceMigrate, "ip_address.0.value", testVars.VMM.AssignedIP),
				),
			},
			{
				Config: testVmPreEnvConfig(r) +
					testVmWithNicAndDiskConfig(vmName) + testVmsNetworkDeviceMigrateV4ReleaseConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmNetworkDeviceMigrate, "ext_id"),
				),
			},
		},
	})
}

func testVmsNetworkDeviceMigrateV4ReleaseConfig() string {
	return `
	

		resource "nutanix_vm_network_device_migrate_v2" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test-vm.id
			ext_id    = resource.nutanix_virtual_machine_v2.test-vm.nics.0.ext_id
			subnet {
				ext_id = nutanix_subnet_v2.subnet.ext_id
			}
			migrate_type = "RELEASE_IP"

		}


`
}

func testVmsNetworkDeviceMigrateV4AssignConfig() string {
	return `

		resource "nutanix_vm_network_device_migrate_v2" "test" {
			vm_ext_id = resource.nutanix_virtual_machine_v2.test-vm.id
			ext_id    = resource.nutanix_virtual_machine_v2.test-vm.nics.0.ext_id
			subnet {
				ext_id = nutanix_subnet_v2.subnet.ext_id
			}
			migrate_type = "ASSIGN_IP"
			ip_address {
			  value = local.vmm.assigned_ip
			}
		}

`
}
