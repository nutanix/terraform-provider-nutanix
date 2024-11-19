package vmmv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameNGTConfiguration = "data.nutanix_ngt_configuration_v2.test"

func TestAccNutanixNGTConfigurationV2Datasource_GetNGTConfigurationForVM(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// create a vm and power on
			{
				PreConfig: func() {
					t.Log("Creating and Powering on the VM")
				},
				Config: testPreEnvConfig(vmName, r),
			},
			// install NGT on the VM
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to power on")
					time.Sleep(2 * time.Minute)
					t.Log("Installing NGT")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot(),
			},
			// get NGT configuration for the VM
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(2 * time.Minute)
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTConfigurationDatasource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNGTConfiguration, "guest_os_version"),
					resource.TestCheckResourceAttrSet(datasourceNameNGTConfiguration, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_installed", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_reachable", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_enabled", "true"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.#", "2"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccNutanixNGTConfigurationV4Datasource_GetNGTConfigurationForVM_NGTNotInstalled(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Log("Creating and Powering on the VM")
				},
				Config: testPreEnvConfig(vmName, r),
			},
			{
				Config: testPreEnvConfig(vmName, r) + testNGTConfigurationDatasource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_enabled", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_installed", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_reachable", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_vm_mobility_drivers_installed", "false"),
					resource.TestCheckResourceAttr(datasourceNameNGTConfiguration, "is_vss_snapshot_capable", "false"),
				),
			},
		},
	})
}

func testNGTConfigurationDatasource() string {
	return `
	data "nutanix_ngt_configuration_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
	}
  `
}
