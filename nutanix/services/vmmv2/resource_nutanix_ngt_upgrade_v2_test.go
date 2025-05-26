package vmmv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameNGTUpgrade = "nutanix_ngt_upgrade_v2.test"

func TestAccV2NutanixNGTUpgradeResource_UpgradeNGTWithRebootPreferenceSetToIMMEDIATE(t *testing.T) {
	t.Skip("This test case is skip since NGT upgrade is failing from v4 api: https://jira.nutanix.com/browse/ENG-665842")
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
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
					time.Sleep(timeSleep)
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
					time.Sleep(timeSleep)
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTUpgradeResourceConfigRebootIMMEDIATE(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTUpgrade, "guest_os_version"),
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

func TestAccV2NutanixNGTUpgradeResource_UpgradeNGTWithRebootPreferenceSetToLATER(t *testing.T) {
	t.Skip("This test case is skip since NGT upgrade is failing from v4 api: https://jira.nutanix.com/browse/ENG-665842")
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
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
					time.Sleep(timeSleep)
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
					time.Sleep(timeSleep)
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTUpgradeResourceConfigRebootLATER(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTUpgrade, "guest_os_version"),
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

func TestAccV2NutanixNGTUpgradeResource_UpgradeNGTWithRebootPreferenceSetToSKIP(t *testing.T) {
	t.Skip("This test case is skip since NGT upgrade is failing from v4 api: https://jira.nutanix.com/browse/ENG-665842")
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
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
					time.Sleep(timeSleep)
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
					time.Sleep(timeSleep)
				},
				Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTUpgradeResourceConfigRebootSKIP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNGTUpgrade, "guest_os_version"),
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

func testNGTUpgradeResourceConfigRebootLATER() string {
	return `
	resource "nutanix_ngt_upgrade_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id

		reboot_preference {
			schedule_type = "LATER"
			schedule {
				start_time = "2026-08-01T00:00:00Z"
			}
		}
	}`
}

func testNGTUpgradeResourceConfigRebootSKIP() string {
	return `
	resource "nutanix_ngt_upgrade_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id

		reboot_preference {
			schedule_type = "SKIP"
		}
	}`
}

func testNGTUpgradeResourceConfigRebootIMMEDIATE() string {
	return `
	resource "nutanix_ngt_upgrade_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id

		reboot_preference {
			schedule_type = "IMMEDIATE"
		}
	}`
}
