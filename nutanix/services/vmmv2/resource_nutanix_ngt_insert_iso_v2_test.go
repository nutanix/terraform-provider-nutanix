package vmmv2_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	pathfilepath "path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameNGTInsertISO = "nutanix_ngt_insert_iso_v2.test"
	ressourceVMNGT           = "nutanix_virtual_machine_v2.ngt-vm"
	datasourceVMNGT          = "data.nutanix_virtual_machine_v2.ngt-vm-refresh"
	timeSleep                = 2 * time.Minute

	// Retry only when the failure is the known transient UVM secure connection issue.
	ngtInsertIsoRetryMaxAttempts = 10
	ngtInsertIsoRetrySubstring   = "failed to establish secure connection with the UVM."
	ngtInsertIsoRetryChildEnvVar = "NUTANIX_NGT_INSERT_ISO_V2_TEST_CHILD"
)

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmHaveNGTTest(t *testing.T) {
	runNGTInsertIsoAccTestWithRetry(t, "TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmHaveNGTTest", func(t *testing.T) {
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
						t.Log("Inserting NGT Iso")
					},
					Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTInsertIsoConfig("true", "insert"),
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
	})
}

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmHaveNGTIsConfigFalse(t *testing.T) {
	runNGTInsertIsoAccTestWithRetry(t, "TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmHaveNGTIsConfigFalse", func(t *testing.T) {
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
						t.Log("Inserting NGT Iso")
					},
					Config: testPreEnvConfig(vmName, r) + testNGTInstallationResourceConfigIMMEDIATEReboot() + testNGTInsertIsoConfig("false", "insert"),
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
	})
}

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmDoseNotHaveNGT(t *testing.T) {
	runNGTInsertIsoAccTestWithRetry(t, "TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmDoseNotHaveNGT", func(t *testing.T) {
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
					Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("true", "insert"),
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
	})
}

func TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmDoseNotHaveNGTIsConfigFalse(t *testing.T) {
	runNGTInsertIsoAccTestWithRetry(t, "TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmDoseNotHaveNGTIsConfigFalse", func(t *testing.T) {
		r := acctest.RandInt()
		vmName := fmt.Sprintf("tf-test-vm-ngt-%d", r)

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { acc.TestAccPreCheck(t) },
			Providers:    acc.TestAccProviders,
			CheckDestroy: testAccCheckNutanixVirtualMachineV2Destroy,
			Steps: []resource.TestStep{
				{
					PreConfig: func() {
						fmt.Println("Step 1: Creating and Powering on the VM")
					},
					Config: testPreEnvConfig(vmName, r),
				},
				// Step 2: Insert the NGT ISO on vm
				{
					PreConfig: func() {
						fmt.Println("Step 2: Inserting the NGT ISO on vm")
					},
					Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("false", "insert"),
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
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_enabled", datasourceVMNGT, "guest_tools.0.is_enabled"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_installed", datasourceVMNGT, "guest_tools.0.is_installed"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_reachable", datasourceVMNGT, "guest_tools.0.is_reachable"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "version", datasourceVMNGT, "guest_tools.0.version"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "guest_os_version", datasourceVMNGT, "guest_tools.0.guest_os_version"),
						resource.TestCheckResourceAttrPair(datasourceVMNGT, "guest_tools.0.capabilities.#", resourceNameNGTInsertISO, "capablities.#"),
						resource.TestCheckResourceAttrPair(datasourceVMNGT, "guest_tools.0.capabilities.0", resourceNameNGTInsertISO, "capablities.0"),
						resource.TestCheckResourceAttrPair(datasourceVMNGT, "guest_tools.0.capabilities.1", resourceNameNGTInsertISO, "capablities.1"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "cd_roms.0.iso_type", "GUEST_TOOLS"),
					),
				},
				// Step 3: Eject the NGT ISO
				{
					PreConfig: func() {
						fmt.Println("Step 3: Ejecting the NGT ISO")
					},
					Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("false", "eject"),
				},
				// Step 4: check the NGT ISO is ejected
				{
					PreConfig: func() {
						fmt.Println("Step 4: Checking the NGT ISO is ejected")
					},
					Config: testPreEnvConfig(vmName, r) + testNGTInsertIsoConfig("false", "eject"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.is_enabled", "true"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.is_installed", "false"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.is_reachable", "false"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.is_iso_inserted", "true"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.version", ""),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.guest_os_version", ""),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.capabilities.#", "2"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.capabilities.0", "SELF_SERVICE_RESTORE"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "guest_tools.0.capabilities.1", "VSS_SNAPSHOT"),
						resource.TestCheckResourceAttr(datasourceVMNGT, "cd_roms.0.iso_type", "OTHER"),
						resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "available_version"),
						resource.TestCheckResourceAttrSet(resourceNameNGTInsertISO, "ext_id"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_iso_inserted", "true"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_config_only", "false"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_enabled", "true"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_installed", "false"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "is_reachable", "false"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.#", "2"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.0", "SELF_SERVICE_RESTORE"),
						resource.TestCheckResourceAttr(resourceNameNGTInsertISO, "capablities.1", "VSS_SNAPSHOT"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_enabled", ressourceVMNGT, "guest_tools.0.is_enabled"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_installed", ressourceVMNGT, "guest_tools.0.is_installed"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "is_reachable", ressourceVMNGT, "guest_tools.0.is_reachable"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "version", ressourceVMNGT, "guest_tools.0.version"),
						resource.TestCheckResourceAttrPair(resourceNameNGTInsertISO, "guest_os_version", ressourceVMNGT, "guest_tools.0.guest_os_version"),
						resource.TestCheckResourceAttrPair(ressourceVMNGT, "guest_tools.0.capabilities.#", resourceNameNGTInsertISO, "capablities.#"),
						resource.TestCheckResourceAttrPair(ressourceVMNGT, "guest_tools.0.capabilities.0", resourceNameNGTInsertISO, "capablities.0"),
						resource.TestCheckResourceAttrPair(ressourceVMNGT, "guest_tools.0.capabilities.1", resourceNameNGTInsertISO, "capablities.1"),
						resource.TestCheckResourceAttr(ressourceVMNGT, "cd_roms.0.iso_type", "OTHER"),
					),
				},
			},
		})
	})
}

func runNGTInsertIsoAccTestWithRetry(t *testing.T, testName string, fn func(t *testing.T)) {
	t.Helper()

	// Child execution: run the real acceptance test body once (no recursion).
	if os.Getenv(ngtInsertIsoRetryChildEnvVar) == "1" {
		fn(t)
		return
	}

	pkgDir := ngtInsertIsoTestPackageDir(t)
	testRe := fmt.Sprintf("^%s$", testName)

	for attempt := 1; attempt <= ngtInsertIsoRetryMaxAttempts; attempt++ {
		// Use -v so progress from the child test is visible, and set a large timeout
		// because acceptance tests can take a long time.
		cmd := exec.Command("go", "test", "-v", "-run", testRe, "-count=1", "-timeout", "500m")
		cmd.Dir = pkgDir
		cmd.Env = append(os.Environ(), ngtInsertIsoRetryChildEnvVar+"=1")

		// Stream child output live (so outer `go test ... > file` logs show progress),
		// but also buffer it so we can detect the retryable substring.
		var buf bytes.Buffer
		cmd.Stdout = io.MultiWriter(&buf, os.Stdout)
		cmd.Stderr = io.MultiWriter(&buf, os.Stderr)

		err := cmd.Run()
		out := buf.Bytes()
		if err == nil {
			if attempt > 1 {
				t.Logf("passed after %d attempt(s)", attempt)
			}
			return
		}

		// Retry only for the known transient UVM secure connection issue.
		if bytes.Contains(out, []byte(ngtInsertIsoRetrySubstring)) && attempt < ngtInsertIsoRetryMaxAttempts {
			t.Logf("attempt %d/%d failed with retryable error (%q); retrying", attempt, ngtInsertIsoRetryMaxAttempts, ngtInsertIsoRetrySubstring)
			continue
		}

		t.Fatalf("failed (attempt %d/%d): %v\n%s", attempt, ngtInsertIsoRetryMaxAttempts, err, string(out))
	}
}

func ngtInsertIsoTestPackageDir(t *testing.T) string {
	t.Helper()

	// Use the directory of this source file as the package dir for `go test`.
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		// Fallback to current working directory if caller info isn't available.
		if wd, err := os.Getwd(); err == nil {
			return wd
		}
		return "."
	}
	return pathfilepath.Dir(file)
}

func testNGTInsertIsoConfig(configMode, action string) string {
	return fmt.Sprintf(`
	resource "nutanix_ngt_insert_iso_v2" "test" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		capablities = ["SELF_SERVICE_RESTORE","VSS_SNAPSHOT"]
		is_config_only = %s
		action = "%s"
	}
		
	data "nutanix_virtual_machine_v2" "ngt-vm-refresh" {
		ext_id = nutanix_virtual_machine_v2.ngt-vm.id
		depends_on = [nutanix_ngt_insert_iso_v2.test]
	}	
		`, configMode, action)
}
