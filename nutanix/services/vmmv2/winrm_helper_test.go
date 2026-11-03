package vmmv2_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	pathutil "path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/masterzen/winrm"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	prismTasks "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
	vmmPrismCfg "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	winrmPort                   = 5985
	winrmConnectTimeout         = 30 * time.Minute
	winrmRetryInterval          = 30 * time.Second
	powerOnTimeout              = 5 * time.Minute
	guestCustomizationTimeout   = 15 * time.Minute
	guestCustomizationPollDelay = 10 * time.Second
	domainJoinTimeout           = 10 * time.Minute
	domainJoinPollInterval      = 30 * time.Second
)

// getVMIPFromAPI fetches the VM by ext_id from the Nutanix API and returns
// the first available IPv4 address (learned or static).
func getVMIPFromAPI(vmExtID string) (string, error) {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	req := import1.GetVmByIdRequest{
		ExtId: utils.StringPtr(vmExtID),
	}
	resp, err := conn.VmmAPI.VMAPIInstance.GetVmById(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("failed to get VM %s: %w", vmExtID, err)
	}

	vm := resp.Data.GetValue().(config.Vm)
	for _, nic := range vm.Nics {
		if nic.NetworkInfo == nil {
			continue
		}
		if nic.NetworkInfo.Ipv4Info != nil {
			for _, ip := range nic.NetworkInfo.Ipv4Info.LearnedIpAddresses {
				if ip.Value != nil && *ip.Value != "" {
					return *ip.Value, nil
				}
			}
		}
		if nic.NetworkInfo.Ipv4Config != nil &&
			nic.NetworkInfo.Ipv4Config.IpAddress != nil &&
			nic.NetworkInfo.Ipv4Config.IpAddress.Value != nil {
			return *nic.NetworkInfo.Ipv4Config.IpAddress.Value, nil
		}
	}

	return "", fmt.Errorf("no IP address found for VM %s", vmExtID)
}

// getVMIPFromState extracts the VM ext_id from the terraform state for the
// given resource, then calls the API to get the IP.
func getVMIPFromState(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("resource %s not found in state", resourceName)
	}

	vmExtID := rs.Primary.ID
	if vmExtID == "" {
		return "", fmt.Errorf("resource %s has empty ID", resourceName)
	}

	return getVMIPFromAPI(vmExtID)
}

// waitForWinRM retries connecting to the VM via WinRM until the VM is
// reachable (sysprep complete, WinRM service up) or the timeout expires.
func waitForWinRM(host, user, password string, timeout time.Duration) (*winrm.Client, error) {
	endpoint := winrm.NewEndpoint(host, winrmPort, false, true, nil, nil, nil, 0)

	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		client, err := winrm.NewClient(endpoint, user, password)
		if err != nil {
			lastErr = err
			time.Sleep(winrmRetryInterval)
			continue
		}

		stdout, _, _, err := client.RunPSWithContext(context.Background(), "echo ok")
		if err != nil {
			lastErr = err
			time.Sleep(winrmRetryInterval)
			continue
		}

		if strings.TrimSpace(stdout) == "ok" {
			return client, nil
		}
		lastErr = fmt.Errorf("unexpected WinRM response: %s", stdout)
		time.Sleep(winrmRetryInterval)
	}

	return nil, fmt.Errorf("WinRM connection to %s timed out after %v: %w", host, timeout, lastErr)
}

// runPowerShell executes a PowerShell command via WinRM and returns stdout.
func runPowerShell(client *winrm.Client, script string) (string, error) {
	stdout, stderr, _, err := client.RunPSWithContext(context.Background(), script)
	if err != nil {
		return "", fmt.Errorf("PowerShell execution failed: %w (stderr: %s)", err, stderr)
	}
	return strings.TrimSpace(stdout), nil
}

// waitForDomainJoin polls the VM via WinRM until it reports being joined to
// the expected domain (i.e. not "WORKGROUP"). Handles transient WinRM
// disconnects during domain-join reboots by re-establishing the connection.
func waitForDomainJoin(client **winrm.Client, host, user, password, expectedDomain string) error {
	deadline := time.Now().Add(domainJoinTimeout)
	for time.Now().Before(deadline) {
		actual, err := runPowerShell(*client, `(Get-WmiObject Win32_ComputerSystem).Domain`)
		if err != nil {
			log.Printf("[DEBUG] WinRM command failed during domain-join wait (VM may be rebooting): %v", err)
			time.Sleep(domainJoinPollInterval)
			newClient, reconnErr := waitForWinRM(host, user, password, 5*time.Minute)
			if reconnErr != nil {
				log.Printf("[DEBUG] WinRM reconnect failed: %v", reconnErr)
				continue
			}
			*client = newClient
			continue
		}
		actual = strings.TrimSpace(actual)
		log.Printf("[DEBUG] Domain join check: current=%q, expected=%q", actual, expectedDomain)
		if strings.EqualFold(actual, expectedDomain) {
			return nil
		}
		time.Sleep(domainJoinPollInterval)
	}
	return fmt.Errorf("domain join did not complete within %v (expected %q)", domainJoinTimeout, expectedDomain)
}

// psValidationCommands returns the map of PowerShell commands used to validate
// OS-level settings on a Windows VM via WinRM.
func psValidationCommands() map[string]string {
	return map[string]string{
		"timezone":           `(Get-TimeZone).Id`,
		"computer_name":      `$env:COMPUTERNAME`,
		"domain":             `(Get-WmiObject Win32_ComputerSystem).Domain`,
		"workgroup":          `(Get-WmiObject Win32_ComputerSystem).Workgroup`,
		"system_locale":      `(Get-WinSystemLocale).Name`,
		"ui_language":        `try { $l = Get-WinUILanguageOverride; if ($l) { $l.LanguageTag } else { (Get-UICulture).Name } } catch { (Get-UICulture).Name }`,
		"registered_org":     `(Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion').RegisteredOrganization`,
		"registered_owner":   `(Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion').RegisteredOwner`,
		"dns_server":         `(Get-DnsClientServerAddress -AddressFamily IPv4 | Where-Object { $_.ServerAddresses.Count -gt 0 } | Select-Object -First 1).ServerAddresses[0]`,
		"ip_address":         `(Get-NetIPAddress -AddressFamily IPv4 | Where-Object { $_.InterfaceAlias -notlike 'Loopback*' } | Select-Object -First 1).IPAddress`,
		"first_logon_marker": `if (Test-Path C:\tf_marker*.txt) { (Get-ChildItem C:\tf_marker*.txt | Select-Object -First 1).FullName } else { "NOT_FOUND" }`,
		"auto_logon_count":   `try { $v = (Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon' -ErrorAction Stop).AutoLogonCount; if ($v -ne $null) { $v } else { "0" } } catch { "0" }`,
	}
}

// ensureVMPoweredOn checks the VM power state and powers it on if it's OFF.
// It waits for the power-on task to complete before returning.
func ensureVMPoweredOn(vmExtID string) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	req := import1.GetVmByIdRequest{ExtId: utils.StringPtr(vmExtID)}
	resp, err := conn.VmmAPI.VMAPIInstance.GetVmById(ctx, &req)
	if err != nil {
		return fmt.Errorf("failed to get VM %s: %w", vmExtID, err)
	}

	vm := resp.Data.GetValue().(config.Vm)
	if vm.PowerState != nil && *vm.PowerState == config.PowerState(2) {
		log.Printf("[DEBUG] VM %s is already powered ON", vmExtID)
		return nil
	}

	log.Printf("[DEBUG] VM %s is not ON (state=%v), powering on...", vmExtID, vm.PowerState)
	powerReq := import1.PowerOnVmRequest{ExtId: utils.StringPtr(vmExtID)}
	powerResp, err := conn.VmmAPI.VMAPIInstance.PowerOnVm(ctx, &powerReq)
	if err != nil {
		return fmt.Errorf("failed to power on VM %s: %w", vmExtID, err)
	}

	taskRef := powerResp.Data.GetValue().(vmmPrismCfg.TaskReference)
	taskUUID := utils.StringValue(taskRef.ExtId)
	log.Printf("[DEBUG] Power-on task started: %s", taskUUID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, conn.PrismAPI, taskUUID),
		Timeout: powerOnTimeout,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("power-on task failed for VM %s: %w", vmExtID, err)
	}

	log.Printf("[DEBUG] VM %s powered on successfully", vmExtID)
	return nil
}

// waitForGuestCustomization polls the Prism tasks API to find a
// kApplyVmGuestCustomization task whose affected entity matches the given
// VM UUID. It waits until the task reaches SUCCEEDED status.
func waitForGuestCustomization(vmExtID string) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	log.Printf("[DEBUG] Waiting for kApplyVmGuestCustomization task on VM %s", vmExtID)

	filter := "operation eq 'kApplyVmGuestCustomization'"
	deadline := time.Now().Add(guestCustomizationTimeout)

	for time.Now().Before(deadline) {
		listReq := prismTasks.ListTasksRequest{
			Filter_: utils.StringPtr(filter),
		}
		listResp, err := conn.PrismAPI.TaskRefAPI.ListTasks(ctx, &listReq)
		if err != nil {
			log.Printf("[DEBUG] Error listing tasks: %v, retrying...", err)
			time.Sleep(guestCustomizationPollDelay)
			continue
		}

		tasks, ok := listResp.Data.GetValue().([]prismConfig.Task)
		if !ok || len(tasks) == 0 {
			log.Printf("[DEBUG] No kApplyVmGuestCustomization tasks found yet, retrying...")
			time.Sleep(guestCustomizationPollDelay)
			continue
		}

		for _, task := range tasks {
			if !taskAffectsVM(task, vmExtID) {
				continue
			}

			status := "UNKNOWN"
			if task.Status != nil {
				status = task.Status.GetName()
			}
			log.Printf("[DEBUG] Found GC task %s for VM %s, status: %s",
				utils.StringValue(task.ExtId), vmExtID, status)

			switch status {
			case "SUCCEEDED":
				log.Printf("[DEBUG] Guest customization completed for VM %s", vmExtID)
				return nil
			case "FAILED", "CANCELED":
				return fmt.Errorf("guest customization task %s for VM %s ended with status: %s",
					utils.StringValue(task.ExtId), vmExtID, status)
			}
		}

		log.Printf("[DEBUG] GC task for VM %s not yet succeeded, polling...", vmExtID)
		time.Sleep(guestCustomizationPollDelay)
	}

	return fmt.Errorf("timed out waiting for kApplyVmGuestCustomization task on VM %s after %v",
		vmExtID, guestCustomizationTimeout)
}

// taskAffectsVM returns true if any of the task's EntitiesAffected has an
// ExtId matching the given VM UUID.
func taskAffectsVM(task prismConfig.Task, vmExtID string) bool {
	for _, entity := range task.EntitiesAffected {
		if entity.ExtId != nil && *entity.ExtId == vmExtID {
			return true
		}
	}
	return false
}

// ensureVMReadyForValidation powers on the VM if needed, then waits for
// guest customization to complete (VM is ON with an IP address).
func ensureVMReadyForValidation(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("resource %s not found in state", resourceName)
	}

	vmExtID := rs.Primary.ID
	if vmExtID == "" {
		return "", fmt.Errorf("resource %s has empty ID", resourceName)
	}

	if err := ensureVMPoweredOn(vmExtID); err != nil {
		return "", err
	}

	if err := waitForGuestCustomization(vmExtID); err != nil {
		return "", err
	}

	return getVMIPFromAPI(vmExtID)
}

// testCheckVMSettings is a TestCheckFunc that connects to a deployed/cloned VM
// via WinRM and validates OS-level settings using PowerShell commands.
//
// Supported expected map keys:
//   - "timezone", "computer_name", "domain", "workgroup",
//     "system_locale", "ui_language", "registered_org",
//     "registered_owner", "dns_server", "ip_address",
//     "first_logon_marker", "auto_logon_count"
//
// Empty values are skipped (useful for DHCP where ip_address is unknown).
func testCheckVMSettings(resourceName, winrmUser, winrmPassword string, expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ip, err := ensureVMReadyForValidation(s, resourceName)
		if err != nil {
			return fmt.Errorf("VM readiness check failed for %s: %w", resourceName, err)
		}

		log.Printf("[DEBUG] VM %s ready at %s, connecting via WinRM...", resourceName, ip)
		client, err := waitForWinRM(ip, winrmUser, winrmPassword, winrmConnectTimeout)
		if err != nil {
			return fmt.Errorf("WinRM connection failed for %s (%s): %w", resourceName, ip, err)
		}

		if expectedDomain, ok := expected["domain"]; ok && expectedDomain != "" {
			log.Printf("[DEBUG] VM %s: waiting for domain join to complete (expected %q)...", resourceName, expectedDomain)
			if err := waitForDomainJoin(&client, ip, winrmUser, winrmPassword, expectedDomain); err != nil {
				return fmt.Errorf("VM %s (%s): %w", resourceName, ip, err)
			}
		}

		return validateVMSettingsViaWinRM(client, resourceName, ip, expected)
	}
}

// testCheckVMSettingsWithIP is like testCheckVMSettings but uses a known static IP
// for WinRM connection instead of discovering it from the API. Use this when
// a static IP was explicitly configured in the deployment's ipv4_config.
func testCheckVMSettingsWithIP(resourceName, winrmUser, winrmPassword, connectIP string, expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		vmExtID := rs.Primary.ID
		if vmExtID == "" {
			return fmt.Errorf("resource %s has empty ID", resourceName)
		}

		if err := ensureVMPoweredOn(vmExtID); err != nil {
			return fmt.Errorf("VM readiness check failed for %s: %w", resourceName, err)
		}
		if err := waitForGuestCustomization(vmExtID); err != nil {
			return fmt.Errorf("VM readiness check failed for %s: %w", resourceName, err)
		}

		log.Printf("[DEBUG] VM %s connecting via WinRM to configured IP %s...", resourceName, connectIP)
		client, err := waitForWinRM(connectIP, winrmUser, winrmPassword, winrmConnectTimeout)
		if err != nil {
			return fmt.Errorf("WinRM connection failed for %s (%s): %w", resourceName, connectIP, err)
		}

		if expectedDomain, ok := expected["domain"]; ok && expectedDomain != "" {
			log.Printf("[DEBUG] VM %s: waiting for domain join to complete (expected %q)...", resourceName, expectedDomain)
			if err := waitForDomainJoin(&client, connectIP, winrmUser, winrmPassword, expectedDomain); err != nil {
				return fmt.Errorf("VM %s (%s): %w", resourceName, connectIP, err)
			}
		}

		return validateVMSettingsViaWinRM(client, resourceName, connectIP, expected)
	}
}

// validateVMSettingsViaWinRM runs PowerShell validation commands via the
// provided WinRM client and checks actual values against expected.
func validateVMSettingsViaWinRM(client *winrm.Client, resourceName, host string, expected map[string]string) error {
	psCommands := psValidationCommands()

	for key, expectedVal := range expected {
		if expectedVal == "" {
			continue
		}
		script, ok := psCommands[key]
		if !ok {
			return fmt.Errorf("unknown check key %q", key)
		}

		if key == "first_logon_marker" {
			script = fmt.Sprintf(`if (Test-Path '%s') { '%s' } else { "NOT_FOUND" }`, expectedVal, expectedVal)
		}

		actual, err := runPowerShell(client, script)
		if err != nil {
			return fmt.Errorf("failed to check %s on %s (%s): %w", key, resourceName, host, err)
		}
		actualTrimmed := strings.TrimSpace(actual)
		log.Printf("[DEBUG] VM %s: %s = %q (expected %q)", resourceName, key, actualTrimmed, expectedVal)

		if key == "auto_logon_count" {
			expN, _ := strconv.Atoi(strings.TrimSpace(expectedVal))
			actN, _ := strconv.Atoi(actualTrimmed)
			if expN > 0 && actN <= 0 {
				return fmt.Errorf("VM %s (%s): %s expected positive value (configured %d), got %q",
					resourceName, host, key, expN, actualTrimmed)
			}
			continue
		}

		if !strings.EqualFold(actualTrimmed, strings.TrimSpace(expectedVal)) {
			return fmt.Errorf("VM %s (%s): %s mismatch: expected %q, got %q",
				resourceName, host, key, expectedVal, actualTrimmed)
		}
	}
	return nil
}

// testCheckVMSettingsWithPasswordChange connects to the VM, changes the Administrator
// password (needed when the image has "must change password at next logon" set),
// then runs the same validations as testCheckVMSettings using the new password.
func testCheckVMSettingsWithPasswordChange(resourceName, winrmUser, winrmPassword, newPassword string, expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ip, err := ensureVMReadyForValidation(s, resourceName)
		if err != nil {
			return fmt.Errorf("VM readiness check failed for %s: %w", resourceName, err)
		}

		log.Printf("[DEBUG] VM %s ready at %s, connecting via WinRM to change password...", resourceName, ip)
		client, err := waitForWinRM(ip, winrmUser, winrmPassword, winrmConnectTimeout)
		if err != nil {
			return fmt.Errorf("WinRM connection failed for %s (%s): %w", resourceName, ip, err)
		}

		changePwdCmd := fmt.Sprintf(`net user %s "%s"`, winrmUser, newPassword)
		if _, err := runPowerShell(client, changePwdCmd); err != nil {
			return fmt.Errorf("failed to change password on %s (%s): %w", resourceName, ip, err)
		}
		log.Printf("[DEBUG] VM %s: password changed successfully", resourceName)

		client, err = waitForWinRM(ip, winrmUser, newPassword, winrmConnectTimeout)
		if err != nil {
			return fmt.Errorf("WinRM reconnect with new password failed for %s (%s): %w", resourceName, ip, err)
		}

		if expectedDomain, ok := expected["domain"]; ok && expectedDomain != "" {
			log.Printf("[DEBUG] VM %s: waiting for domain join to complete (expected %q)...", resourceName, expectedDomain)
			if err := waitForDomainJoin(&client, ip, winrmUser, newPassword, expectedDomain); err != nil {
				return fmt.Errorf("VM %s (%s): %w", resourceName, ip, err)
			}
		}

		return validateVMSettingsViaWinRM(client, resourceName, ip, expected)
	}
}

// findFreeIPs runs the nmap-based script to discover free IPs in the given range.
// It returns `count` IPs that are not currently responding to ping in [startIP, endIP].
func findFreeIPs(startIP, endIP string, count int) ([]string, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	scriptPath := pathutil.Join(pathutil.Dir(thisFile), "..", "..", "..", "scripts", "find_free_ips.sh")

	cmd := exec.Command("bash", scriptPath, startIP, endIP, strconv.Itoa(count))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Printf("[DEBUG] Finding %d free IPs in range %s - %s", count, startIP, endIP)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("find_free_ips.sh failed: %v\nstderr: %s", err, stderr.String())
	}

	var ips []string
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip != "" {
			ips = append(ips, ip)
		}
	}

	if len(ips) < count {
		return nil, fmt.Errorf("expected %d free IPs, got %d", count, len(ips))
	}

	log.Printf("[DEBUG] Free IPs found: %v", ips)
	return ips, nil
}
