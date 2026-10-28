package dataprotectionv2_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clustersReq "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/request/clusters"
	import2 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/request/protectionpolicies"
	dpConfig "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	protectedResources "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/request/protectedresources"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	prismManagement "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	import3 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/categories"
	domainmanagerReq "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/domainmanager"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vm"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/request/volumegroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func waitForVMToBeProtected(resourceName, attributeName, desiredValue string, maxRetries int, retryInterval, sleepTime time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var lastValue string
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.VmmAPI.VMAPIInstance
		ctx := context.Background()

		for i := 0; i < maxRetries; i++ {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource not found: %s", resourceName)
			}

			getVmByIdRequest := import1.GetVmByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			vmResp, err := client.GetVmById(ctx, &getVmByIdRequest)
			if err != nil {
				return fmt.Errorf("error getting vm by id: %v", err)
			}

			// read the attribute value from the response
			vm := vmResp.Data.GetValue().(config.Vm)
			lastValue = config.ProtectionType.GetName(*vm.ProtectionType)
			if lastValue == desiredValue {
				log.Printf("[DEBUG] VM is %s\n", lastValue)
				time.Sleep(sleepTime)
				return nil // Desired value reached
			}

			log.Printf("[DEBUG] Waiting for vm to be protected:  attribute %q to be %q. Current value: %q\n", attributeName, desiredValue, lastValue)
			// Wait before retrying
			time.Sleep(retryInterval)
		}

		return fmt.Errorf("VM: failed to reach desired value for attribute %q: expected %q, got %q after %d retries", attributeName, desiredValue, lastValue, maxRetries)
	}
}

// waitForProtectedResourceAvailable polls GET protected-resources/{ext_id} until it succeeds
// (HTTP 200 / no error). After success, sleeps settleTime so replication can stabilize
// before promote/restore (same pattern as Ansible's post-protection wait).
//
// If expectedRPOSeconds is set (e.g. 60 for minutely restore), also wait until a
// replication_states entry reports that recovery_point_objective_seconds. Magneto
// restore rejects entities that are only sync-protected (RPO=0).
func waitForProtectedResourceAvailable(resourceName string, maxRetries int, retryInterval, settleTime time.Duration, expectedRPOSeconds ...int64) resource.TestCheckFunc {
	var wantRPO int64
	checkRPO := false
	if len(expectedRPOSeconds) > 0 {
		wantRPO = expectedRPOSeconds[0]
		checkRPO = true
	}

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		extID := rs.Primary.ID
		if extID == "" {
			return fmt.Errorf("resource %s has empty id", resourceName)
		}

		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.DataProtectionAPI.ProtectedResource
		ctx := context.Background()

		var lastErr error
		var lastRPOs []int64
		for i := 0; i < maxRetries; i++ {
			req := protectedResources.GetProtectedResourceByIdRequest{
				ExtId: utils.StringPtr(extID),
			}
			resp, err := client.GetProtectedResourceById(ctx, &req)
			if err != nil {
				lastErr = err
				log.Printf("[DEBUG] Waiting for protected resource %s (attempt %d/%d): %v", extID, i+1, maxRetries, err)
				time.Sleep(retryInterval)
				continue
			}

			if checkRPO {
				pr := resp.Data.GetValue().(dpConfig.ProtectedResource)
				lastRPOs = lastRPOs[:0]
				for _, st := range pr.ReplicationStates {
					lastRPOs = append(lastRPOs, utils.Int64Value(st.RecoveryPointObjectiveSeconds))
				}
				// Magneto restore checks the protection-policy schedule (minutely / RPO=60),
				// not that replication_states is already populated. UI can show Protected
				// with Sync Status empty while replication_states stays [] for a long time.
				// Only fail-fast if states are present but report a different RPO (e.g. sync 0).
				if len(lastRPOs) > 0 {
					matched := false
					for _, rpo := range lastRPOs {
						if rpo == wantRPO {
							matched = true
							break
						}
					}
					if !matched {
						lastErr = fmt.Errorf("replication_states RPO %v, want %d", lastRPOs, wantRPO)
						log.Printf("[DEBUG] Waiting for protected resource %s minutely RPO (attempt %d/%d): %v", extID, i+1, maxRetries, lastErr)
						time.Sleep(retryInterval)
						continue
					}
				} else {
					log.Printf("[DEBUG] Protected resource %s GET ok with empty replication_states (attempt %d/%d); proceeding (expect policy RPO=%d)",
						extID, i+1, maxRetries, wantRPO)
				}
			}

			log.Printf("[DEBUG] Protected resource %s is available (attempt %d/%d); settling %s", extID, i+1, maxRetries, settleTime)
			if settleTime > 0 {
				time.Sleep(settleTime)
			}
			return nil
		}

		if checkRPO {
			return fmt.Errorf("protected resource %s not ready with RPO=%d after %d retries (~%s): last error: %v",
				extID, wantRPO, maxRetries, time.Duration(maxRetries)*retryInterval, lastErr)
		}
		return fmt.Errorf("protected resource %s not available after %d retries (~%s): last error: %v",
			extID, maxRetries, time.Duration(maxRetries)*retryInterval, lastErr)
	}
}

func testCheckDestroyProtectedResourceAndCleanup(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	vmClient := conn.VmmAPI.VMAPIInstance
	categoryClient := conn.PrismAPI.CategoriesAPIInstance
	ppClient := conn.DataPoliciesAPI.ProtectionPolicies
	ctx := context.Background()

	vmExtID := ""
	ppExtID := ""
	categoryExtID := ""

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_virtual_machine_v2" {
			log.Printf("[DEBUG] Checking if VM still exists\n")
			vmExtID = rs.Primary.ID
		}

		if rs.Type == "nutanix_protection_policy_v2" {
			log.Printf("[DEBUG] Checking if Protection Policy still exists\n")
			ppExtID = rs.Primary.ID
		}

		if rs.Type == "nutanix_category_v2" {
			log.Printf("[DEBUG] Checking if Category still exists\n")
			categoryExtID = rs.Primary.ID
		}
	}

	// delete vm
	if vmExtID != "" {
		getVmByIdRequest := import1.GetVmByIdRequest{
			ExtId: utils.StringPtr(vmExtID),
		}
		readResp, err := vmClient.GetVmById(ctx, &getVmByIdRequest)
		if err == nil {
			args := make(map[string]interface{})
			etag := vmClient.ApiClient.GetEtag(readResp)
			args["If-Match"] = utils.StringPtr(etag)
			deleteVmByIdRequest := import1.DeleteVmByIdRequest{
				ExtId: utils.StringPtr(vmExtID),
			}
			_, err = vmClient.DeleteVmById(ctx, &deleteVmByIdRequest, args)
			if err != nil {
				return fmt.Errorf("error: VM still exists: %v", err)
			}
			log.Printf("[DEBUG] VM deleted\n")
		}
	}

	// delete protection policy
	if ppExtID != "" {
		getProtectionPolicyByIdRequest := import2.GetProtectionPolicyByIdRequest{
			ExtId: utils.StringPtr(ppExtID),
		}
		_, err := ppClient.GetProtectionPolicyById(ctx, &getProtectionPolicyByIdRequest)
		if err == nil {
			log.Printf("[DEBUG] Protection Policy still exists")
			deleteProtectionPolicyByIdRequest := import2.DeleteProtectionPolicyByIdRequest{
				ExtId: utils.StringPtr(ppExtID),
			}
			_, err = ppClient.DeleteProtectionPolicyById(ctx, &deleteProtectionPolicyByIdRequest)
			if err != nil {
				return fmt.Errorf("error: Protection Policy still exists : %v", err)
			}
			log.Printf("[DEBUG] Protection Policy deleted\n")
		}
	}

	// delete category
	if categoryExtID != "" {
		getCategoryByIdRequest := import3.GetCategoryByIdRequest{
			ExtId: utils.StringPtr(categoryExtID),
		}
		_, err := categoryClient.GetCategoryById(ctx, &getCategoryByIdRequest)
		if err == nil {
			log.Printf("[DEBUG] Category still exists")
			deleteCategoryByIdRequest := import3.DeleteCategoryByIdRequest{
				ExtId: utils.StringPtr(categoryExtID),
			}
			_, err = categoryClient.DeleteCategoryById(ctx, &deleteCategoryByIdRequest)
			if err != nil {
				return fmt.Errorf("error: Category still exists : %v", err)
			}
			log.Printf("[DEBUG] Category deleted\n")
		}
	}

	return nil
}

func testCheckDestroyProtectedResourceAndCleanupForPromoteVM(state *terraform.State) error {
	// check on protection policy, category and vm on local site
	testCheckDestroyProtectedResourceAndCleanup(state)

	// Best-effort cleanup of promoted VM on remote (may already be deleted in Check).
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machines_v2" {
			continue
		}
		vmExtID := rs.Primary.Attributes["vms.0.ext_id"]
		if vmExtID == "" {
			continue
		}

		connRemote := acc.TestAccProvider2.Meta().(*conns.Client)
		client := connRemote.VmmAPI.VMAPIInstance
		ctx := context.Background()

		getVmByIdRequest := import1.GetVmByIdRequest{
			ExtId: utils.StringPtr(vmExtID),
		}
		readResp, err := client.GetVmById(ctx, &getVmByIdRequest)
		if err != nil {
			// Already cleaned up in the promote Check — OK.
			log.Printf("[DEBUG] Promoted VM already gone: %v\n", err)
			continue
		}
		args := make(map[string]interface{})
		etag := client.ApiClient.GetEtag(readResp)
		args["If-Match"] = utils.StringPtr(etag)

		deleteVmByIdRequest := import1.DeleteVmByIdRequest{
			ExtId: utils.StringPtr(vmExtID),
		}
		if _, err = client.DeleteVmById(ctx, &deleteVmByIdRequest, args); err != nil {
			return fmt.Errorf("error: Promoted VM still exists: %v", err)
		}
		log.Printf("[DEBUG] Promoted VM deleted\n")
	}

	return nil
}

// dropResourceFromState removes a managed resource from the test state so the
// harness skip provider Destroy for it (used after sync promote removes the source VM).
func dropResourceFromState(address string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		delete(s.RootModule().Resources, address)
		return nil
	}
}

func deleteRestoredVM(vmName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider2.Meta().(*conns.Client)
		client := conn.VmmAPI.VMAPIInstance
		ctx := context.Background()

		filter := fmt.Sprintf("startswith(name, '%s')", vmName)

		listVmsRequest := import1.ListVmsRequest{
			Filter_: utils.StringPtr(filter),
		}
		resp, err := client.ListVms(ctx, &listVmsRequest)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if resp.Data == nil {
			return fmt.Errorf("no data returned from list vms on Remote site")
		}
		vms := resp.Data.GetValue().([]config.Vm)

		vm := vms[0]
		getVmByIdRequest := import1.GetVmByIdRequest{
			ExtId: vm.ExtId,
		}
		readResp, err := client.GetVmById(ctx, &getVmByIdRequest)
		if err == nil {
			args := make(map[string]interface{})
			etag := client.ApiClient.GetEtag(readResp)
			args["If-Match"] = utils.StringPtr(etag)
			deleteVmByIdRequest := import1.DeleteVmByIdRequest{
				ExtId: vm.ExtId,
			}
			_, err = client.DeleteVmById(ctx, &deleteVmByIdRequest, args)
			if err != nil {
				return fmt.Errorf("error: Restored VM still exists: %v", err)
			}
			log.Printf("[DEBUG] Restored VM deleted\n")
			return nil
		}

		return nil
	}
}

func deleteRestoredVg(vgName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider2.Meta().(*conns.Client)
		client := conn.VolumeAPI.VolumeAPIInstance
		ctx := context.Background()
		filter := fmt.Sprintf("startswith(name, '%s')", vgName)

		listVolumeGroupsRequest := import4.ListVolumeGroupsRequest{
			Filter_: utils.StringPtr(filter),
		}
		resp, err := client.ListVolumeGroups(ctx, &listVolumeGroupsRequest)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if resp.Data == nil {
			return fmt.Errorf("no data returned from list Volume Groups on Remote site")
		}
		vgs := resp.Data.GetValue().([]volumesClient.VolumeGroup)

		vg := vgs[0]
		if err == nil {
			deleteVolumeGroupByIdRequest := import4.DeleteVolumeGroupByIdRequest{
				ExtId: vg.ExtId,
			}
			_, err = client.DeleteVolumeGroupById(ctx, &deleteVolumeGroupByIdRequest)
			if err != nil {
				return fmt.Errorf("error: Restored Volume Group still exists: %v", err)
			}
			log.Printf("[DEBUG] Restored Volume Group deleted\n")
			return nil
		}

		return nil
	}
}

// Captured during PromoteVG Check so CheckDestroy can unregister after terraform
// destroy (null_resource may no longer expose triggers in the post-destroy state).
var (
	promoteVGPCExtID             string
	promoteVGRemoteClusterExtID  string
)

// capturePromoteVGUnregisterIDs stores PC + remote PE ext IDs from
// null_resource.wait_stretch_cleared triggers for post-test unregister.
func capturePromoteVGUnregisterIDs() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "null_resource" {
				continue
			}
			if v := rs.Primary.Attributes["triggers.remote_cluster_ext_id"]; v != "" {
				promoteVGRemoteClusterExtID = v
				promoteVGPCExtID = rs.Primary.Attributes["triggers.pc_ext_id"]
				log.Printf("[DEBUG] Captured PromoteVG unregister ids: pc=%s remotePE=%s",
					promoteVGPCExtID, promoteVGRemoteClusterExtID)
				return nil
			}
		}
		return fmt.Errorf("wait_stretch_cleared triggers.remote_cluster_ext_id not found in state")
	}
}

// testCheckDestroyPromoteVGAndUnregister cleans up leftover PP/category resources,
// then unregisters the remote PE that PromoteVG registered for the test.
// pc_registration_v2 Delete is a no-op, so unregister must be explicit. Stretch
// params are cleared first by null_resource.wait_stretch_cleared's destroy
// provisioner (KB 12749).
func testCheckDestroyPromoteVGAndUnregister(state *terraform.State) error {
	if err := testCheckDestroyProtectedResourceAndCleanup(state); err != nil {
		return err
	}
	return unregisterPromoteVGRemotePE()
}

func unregisterPromoteVGRemotePE() error {
	pcExtID := promoteVGPCExtID
	remoteClusterExtID := promoteVGRemoteClusterExtID
	if remoteClusterExtID == "" || pcExtID == "" {
		log.Printf("[DEBUG] No remote PE ext_id captured; skipping unregister")
		return nil
	}

	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	getReq := domainmanagerReq.GetDomainManagerByIdRequest{
		ExtId: utils.StringPtr(pcExtID),
	}
	readResp, err := conn.PrismAPI.DomainManagerAPIInstance.GetDomainManagerById(ctx, &getReq)
	if err != nil {
		return fmt.Errorf("error fetching PC before unregister: %v", err)
	}

	args := map[string]interface{}{
		"If-Match": utils.StringPtr(conn.ClusterAPI.ClusterEntityAPI.ApiClient.GetEtag(readResp)),
	}
	body := prismManagement.NewClusterReference()
	body.ExtId = utils.StringPtr(remoteClusterExtID)
	unregisterReq := domainmanagerReq.UnregisterRequest{
		ExtId: utils.StringPtr(pcExtID),
		Body:  body,
	}
	resp, err := conn.PrismAPI.DomainManagerAPIInstance.Unregister(ctx, &unregisterReq, args)
	if err != nil {
		if isUnregisterPRI55102(err) {
			return handlePromoteVGUnregisterPRI55102(ctx, conn, remoteClusterExtID, err)
		}
		return fmt.Errorf("error unregistering remote PE %s: %v", remoteClusterExtID, err)
	}

	taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := utils.StringValue(taskRef.ExtId)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, conn.PrismAPI, taskUUID),
		Timeout: 30 * time.Minute,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for remote PE unregister (%s): %v", taskUUID, err)
	}
	log.Printf("[DEBUG] Remote PE %s unregistered from PC %s\n", remoteClusterExtID, pcExtID)
	// Unregister API can leave partial PE/PC multicluster state (PRI-55102 on later runs).
	if err := clearPartialUnregisterState(pcExtID, remoteClusterExtID); err != nil {
		return fmt.Errorf("remote PE %s unregistered, but failed clearing partial unregister state: %v", remoteClusterExtID, err)
	}
	clearPromoteVGUnregisterIDs()
	return nil
}

func clearPromoteVGUnregisterIDs() {
	promoteVGPCExtID = ""
	promoteVGRemoteClusterExtID = ""
}

func isUnregisterPRI55102(err error) bool {
	return err != nil && strings.Contains(err.Error(), "PRI-55102")
}

// handlePromoteVGUnregisterPRI55102 treats PRI-55102 as "already unregistered"
// when the PE is not in active registration, then clears leftover ncli cluster
// state so the next run can register again. Payload IDs are not swapped.
func handlePromoteVGUnregisterPRI55102(ctx context.Context, conn *conns.Client, remoteClusterExtID string, unregisterErr error) error {
	active, verifyErr := promoteVGPEActivelyRegistered(ctx, conn, remoteClusterExtID)
	if verifyErr != nil {
		return fmt.Errorf("unregister returned PRI-55102 for PE %s, but could not verify registration state: %v (unregister err: %v)",
			remoteClusterExtID, verifyErr, unregisterErr)
	}
	if active {
		return fmt.Errorf("unregister returned PRI-55102 for PE %s, but PE still appears actively registered: %v",
			remoteClusterExtID, unregisterErr)
	}

	log.Printf("[DEBUG] PRI-55102 for PE %s: not actively registered; clearing partial unregister state", remoteClusterExtID)
	if err := clearPartialUnregisterState(promoteVGPCExtID, remoteClusterExtID); err != nil {
		return fmt.Errorf("PRI-55102 for PE %s (already unregistered), but failed clearing partial unregister state: %v",
			remoteClusterExtID, err)
	}
	clearPromoteVGUnregisterIDs()
	return nil
}

// promoteVGPEActivelyRegistered returns false when the PE is gone from the
// clusters API and/or ncli reports no remote multicluster connection.
func promoteVGPEActivelyRegistered(ctx context.Context, conn *conns.Client, remoteClusterExtID string) (bool, error) {
	getClusterReq := clustersReq.GetClusterByIdRequest{
		ExtId: utils.StringPtr(remoteClusterExtID),
	}
	_, err := conn.ClusterAPI.ClusterEntityAPI.GetClusterById(ctx, &getClusterReq)
	if err != nil {
		log.Printf("[DEBUG] GetClusterById(%s) after PRI-55102: %v (treating as not actively registered)", remoteClusterExtID, err)
		return false, nil
	}

	out, ncliErr := runPCNcli(fmt.Sprintf("multicluster get-cluster-state cluster-id=%s", remoteClusterExtID))
	if ncliErr != nil {
		// State already gone, or PC has no row for this UUID.
		if clusterStateAbsent(out, ncliErr) {
			log.Printf("[DEBUG] ncli get-cluster-state: cluster state absent for %s", remoteClusterExtID)
			return false, nil
		}
		return false, fmt.Errorf("ncli get-cluster-state: %v (output: %s)", ncliErr, out)
	}
	if ncliFieldTrue(out, "Remote Connection Exists") || ncliFieldTrue(out, "Is Multicluster") {
		return true, nil
	}
	log.Printf("[DEBUG] ncli get-cluster-state for %s shows no active remote connection", remoteClusterExtID)
	return false, nil
}

func ncliFieldTrue(out, field string) bool {
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, field) {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		return strings.EqualFold(strings.TrimSpace(parts[1]), "true")
	}
	return false
}

func clusterStateAbsent(out string, err error) bool {
	msg := out
	if err != nil {
		msg += " " + err.Error()
	}
	lower := strings.ToLower(msg)
	return strings.Contains(msg, "cannot be retrieved") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "No cluster") ||
		strings.Contains(lower, "never added to prism central") ||
		strings.Contains(lower, "not registered")
}

func ncliSoftOK(out string, err error) bool {
	if err == nil {
		return true
	}
	return clusterStateAbsent(out, err) ||
		strings.Contains(strings.ToLower(out), "already")
}

// clearPartialUnregisterState cleans both sides after an unregister-only test
// (no blacklist-removal). Sequence:
//  1. PE: remove-from-multicluster ... force=true local_only=true
//  2. PE: remove-from-multicluster ... cleanup_only=true (auto-confirm Y)
//  3. PC: delete-cluster-state ... skip_prechecks=true
//  4. Verify PE no longer lists PC_UUID and PC no longer lists PE_UUID
func clearPartialUnregisterState(pcExtID, remoteClusterExtID string) error {
	peIP := testVars.DataProtection.RemoteClusterPE
	pcIP, err := pcSSHHostForCleanup()
	if err != nil {
		return err
	}
	peUser := testVars.PEUsername
	pePass := testVars.PEPassword
	if peIP == "" || peUser == "" || pePass == "" || pcExtID == "" {
		return fmt.Errorf("missing pe/pc cleanup inputs (remote_cluster_pe/pe_username/pe_password/pc_ext_id)")
	}
	pePassQ := strings.ReplaceAll(pePass, `'`, `'\''`)

	// 1. Clear PE-side registration locally (PC already dropped the PE).
	step1 := fmt.Sprintf(
		"source /etc/profile; /home/nutanix/prism/cli/ncli multicluster remove-from-multicluster external-ip-address-or-svm-ips=%s username=%s password='%s' force=true local_only=true",
		pcIP, peUser, pePassQ,
	)
	log.Printf("[DEBUG] Partial-unregister cleanup step 1 (PE local_only) pe=%s pc=%s", peIP, pcIP)
	out1, err1 := runPESSH(step1)
	log.Printf("[DEBUG] step1 output: %s", out1)
	if !ncliSoftOK(out1, err1) {
		return fmt.Errorf("PE local_only remove-from-multicluster: %v (output: %s)", err1, out1)
	}

	// 2. PE-side cleanup_only. ncli uses System.console().readLine() (needs a
	// PTY); plain stdin pipes fail. script(1) allocates a PTY and we answer Y.
	step2 := fmt.Sprintf(
		"script -q -c \"/home/nutanix/prism/cli/ncli multicluster remove-from-multicluster cluster-id=%s username=%s password='%s' external-ip-address-or-svm-ips=%s cleanup_only=true\" /dev/null <<'EOF'\nY\nEOF",
		pcExtID, peUser, pePassQ, pcIP,
	)
	log.Printf("[DEBUG] Partial-unregister cleanup step 2 (PE cleanup_only) pe=%s pcUUID=%s", peIP, pcExtID)
	out2, err2 := runPESSH(step2)
	log.Printf("[DEBUG] step2 output: %s", out2)
	if !ncliSoftOK(out2, err2) {
		return fmt.Errorf("PE cleanup_only remove-from-multicluster: %v (output: %s)", err2, out2)
	}

	// 3. Remove stale PE state from the PC.
	step3 := fmt.Sprintf(
		"multicluster delete-cluster-state cluster-id=%s skip_prechecks=true", remoteClusterExtID)
	log.Printf("[DEBUG] Partial-unregister cleanup step 3 (PC delete-cluster-state) peUUID=%s", remoteClusterExtID)
	out3, err3 := runPCNcli(step3)
	log.Printf("[DEBUG] step3 output: %s", out3)
	if !ncliSoftOK(out3, err3) {
		return fmt.Errorf("PC delete-cluster-state: %v (output: %s)", err3, out3)
	}

	// 4. Verify both sides.
	peState, peErr := runPESSH("/home/nutanix/prism/cli/ncli multicluster get-cluster-state")
	log.Printf("[DEBUG] step4 PE get-cluster-state: %s", peState)
	if peErr == nil && strings.Contains(peState, pcExtID) {
		return fmt.Errorf("after cleanup, PE still lists PC %s in get-cluster-state", pcExtID)
	}

	pcState, pcErr := runPCNcli(fmt.Sprintf("multicluster get-cluster-state cluster-id=%s", remoteClusterExtID))
	log.Printf("[DEBUG] step4 PC get-cluster-state: %s", pcState)
	if pcErr == nil && (ncliFieldTrue(pcState, "Remote Connection Exists") || ncliFieldTrue(pcState, "Is Multicluster")) {
		return fmt.Errorf("after cleanup, PC still shows PE %s as actively registered: %s", remoteClusterExtID, pcState)
	}
	if pcErr != nil && !clusterStateAbsent(pcState, pcErr) {
		return fmt.Errorf("after cleanup, unexpected PC get-cluster-state error: %v (output: %s)", pcErr, pcState)
	}

	log.Printf("[DEBUG] Partial unregister cleanup complete for PE %s / PC %s", remoteClusterExtID, pcExtID)
	return nil
}

func runPESSH(remoteCmd string) (string, error) {
	peIP := testVars.DataProtection.RemoteClusterPE
	user := testVars.SshPeUsername
	pass := testVars.SshPePassword
	if peIP == "" || user == "" || pass == "" {
		return "", fmt.Errorf("ssh_pe_username/ssh_pe_password/remote_cluster_pe missing from test_config_v2.json")
	}
	cmd := exec.Command("sshpass", "-p", pass,
		"ssh", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("%s@%s", user, peIP),
		remoteCmd,
	)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runPCNcli(ncliArgs string) (string, error) {
	host, err := pcSSHHostForCleanup()
	if err != nil {
		return "", err
	}
	user := testVars.SshPcUsername
	pass := testVars.SshPcPassword
	if user == "" || pass == "" {
		return "", fmt.Errorf("ssh_pc_username/ssh_pc_password missing from test_config_v2.json")
	}
	cmd := exec.Command("sshpass", "-p", pass,
		"ssh", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("%s@%s", user, host),
		"/home/nutanix/prism/cli/ncli "+ncliArgs,
	)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// pcSSHHostForCleanup resolves the PC CVM address for ncli SSH.
// Direct PC uses NUTANIX_ENDPOINT; NC often embeds the IP in the central
// hostname (e.g. nconprem-10-44-76-117...).
func pcSSHHostForCleanup() (string, error) {
	if ep := os.Getenv("NUTANIX_ENDPOINT"); ep != "" {
		return ep, nil
	}
	central := os.Getenv("NUTANIX_CENTRAL_ENDPOINT")
	if central == "" {
		return "", fmt.Errorf("set NUTANIX_ENDPOINT or NUTANIX_CENTRAL_ENDPOINT to resolve PC SSH host")
	}
	re := regexp.MustCompile(`(\d{1,3})-(\d{1,3})-(\d{1,3})-(\d{1,3})`)
	m := re.FindStringSubmatch(central)
	if m == nil {
		return "", fmt.Errorf("cannot extract PC IP from NUTANIX_CENTRAL_ENDPOINT=%q", central)
	}
	return fmt.Sprintf("%s.%s.%s.%s", m[1], m[2], m[3], m[4]), nil
}
