package vmmv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameNGTInsertISO = "nutanix_ngt_insert_iso_v2.test"

func TestAccNutanixNGTInsertIsoV2Resource_InsertNGTIsoIntoVmHaveNGT(t *testing.T) {
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
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to power on")
					time.Sleep(2 * time.Minute)
					t.Log("Installing NGT")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.1", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
				),
			},
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(2 * time.Minute)
					t.Log("Inserting NGT Iso")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTInsertIsoConfig("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "guest_os_version"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccNutanixNGTInsertIsoV2Resource_InsertNGTIsoIntoVmHaveNGTIsConfigFalse(t *testing.T) {
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
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to power on")
					time.Sleep(2 * time.Minute)
					t.Log("Installing NGT")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "guest_os_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInstallation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "is_iso_inserted", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "capablities.1", "VSS_SNAPSHOT"),
					resource.TestCheckResourceAttr(resourceNameNGTInstallation, "reboot_preference.0.schedule_type", "IMMEDIATE"),
				),
			},
			{
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(2 * time.Minute)
					t.Log("Inserting NGT Iso")
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTInsertIsoConfig("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "guest_os_version"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccNutanixNGTInsertIsoV2Resource_InsertNGTIsoIntoVmDoseNotHaveNGT(t *testing.T) {
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
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(2 * time.Minute)
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "available_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "guest_os_version", ""),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func TestAccNutanixNGTInsertIsoV2Resource_InsertNGTIsoIntoVmDoseNotHaveNGTIsConfigFalse(t *testing.T) {
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
				PreConfig: func() {
					t.Log("Sleeping for 2 Minute waiting vm to reboot")
					time.Sleep(2 * time.Minute)
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "available_version"),
					resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "guest_os_version", ""),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "false"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
					resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
				),
			},
		},
	})
}

func testNGTInsertIsoConfig(configMode string) string {
	return fmt.Sprintf(`	
	resource "nutanix_ngt_insert_iso_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		capablities = ["SELF_SERVICE_RESTORE","VSS_SNAPSHOT"]
		is_config_only = %s
	}`, configMode)
}
