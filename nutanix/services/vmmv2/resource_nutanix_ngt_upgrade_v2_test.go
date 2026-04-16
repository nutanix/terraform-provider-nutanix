package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameNGTUpgrade = "nutanix_ngt_upgrade_v2.test"

func TestAccV2NutanixNGTUpgradeResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// t.Log("Sleeping for 2 Minute waiting vm to reboot")
					// time.Sleep(timeSleep)
					t.Log("Upgrading NGT")
				},
				Config: testNGTUpgradeResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTUpgrade, "guest_os_version"),
					resource.TestCheckResourceAttrSet(datasourceNameNGTConfiguration, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_installed", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_reachable", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_enabled", "true"),
				),
			},
			{
				PreConfig: func() {
					t.Log("Upgrading NGT again")
					time.Sleep(30 * time.Second)
				},
				Config:      testNGTUpgradeResourceConfigDoublicated(),
				ExpectError: regexp.MustCompile(`Failed to perform the operation on VM with UUID '.+' as Nutanix Guest Tools is already upgraded\.`),
			},
		},
	})
}

func testNGTUpgradeResourceConfig() string {
	return fmt.Sprintf(`
	locals {
		config = (jsondecode(file("%[1]s")))
		vmm    = local.config.vmm
	}
	data "nutanix_virtual_machines_v2" "ngt-vm-upgrade" {
		filter = "name eq '${local.vmm.ngt.ngt_upgrade_vm_name}'"
	}
	resource "nutanix_ngt_upgrade_v2" "test" {
		ext_id = data.nutanix_virtual_machines_v2.ngt-vm-upgrade.vms.0.ext_id
		reboot_preference {
			schedule_type = "IMMEDIATE"
		}
	}
	data "nutanix_ngt_configuration_v2" "test" {
		ext_id = data.nutanix_virtual_machines_v2.ngt-vm-upgrade.vms.0.ext_id
	}
		
	`, filepath)
}

func testNGTUpgradeResourceConfigDoublicated() string {
	return fmt.Sprintf(`
	locals {
		config = (jsondecode(file("%[1]s")))
		vmm    = local.config.vmm
	}
	data "nutanix_virtual_machines_v2" "ngt-vm-upgrade-1" {
		filter = "name eq '${local.vmm.ngt.ngt_upgrade_vm_name}'"
	}
	resource "nutanix_ngt_upgrade_v2" "test-1" {
		ext_id = data.nutanix_virtual_machines_v2.ngt-vm-upgrade-1.vms.0.ext_id
		reboot_preference {
			schedule_type = "IMMEDIATE"
		}
	}

		
	`, filepath)
}
