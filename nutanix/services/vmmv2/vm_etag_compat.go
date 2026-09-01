package vmmv2

import (
	"context"
	"fmt"
	"log"
	"time"

	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vm"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// TEMPORARY COMPAT — REMOVE_AFTER: provider v2.6.0 (two releases after reintroduction).
//
// PC still requires If-Match (eTag) on several VMM power / create / delete ops
// (PLAT-104 / MISSING_ETAG_HEADER). Once the platform no longer requires eTag:
//  1. Delete this file (vm_etag_compat.go).
//  2. Strip vmEtagArgs(...) / powerOnVMWithEtag / powerOffVMWithEtag and ", args"
//     from call sites in resource_nutanix_virtual_machine_v2.go, helper.go, and
//     resource_nutanix_ova_vm_deploy_v2.go.
//  3. Collapse callForPowerOnVM / callForPowerOffVM back to a single-shot API +
//     task wait (no maxPowerRetries loop).
//  4. Changelog: remove temporary VMM If-Match compatibility.

const (
	// maxPowerRetries is used for both API-layer retries (power-on/power-off call)
	// and task-layer retries (wait for task).
	maxPowerRetries     = 10
	powerTaskRetryDelay = 5 * time.Second
)

// vmEtagArgs GETs the VM and returns SDK header args with a fresh If-Match eTag.
// REMOVE_AFTER: v2.6.0 — delete with this file.
func vmEtagArgs(ctx context.Context, conn *vmm.Client, vmID string) (map[string]interface{}, error) {
	getVmByIdRequest := import3.GetVmByIdRequest{
		ExtId: utils.StringPtr(vmID),
	}
	readResp, err := conn.VMAPIInstance.GetVmById(ctx, &getVmByIdRequest)
	if err != nil {
		return nil, fmt.Errorf("error while fetching VM for eTag: %w", err)
	}
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)
	return args, nil
}

// powerOnVMWithEtag performs a single PowerOnVm call with a fresh If-Match eTag.
// REMOVE_AFTER: v2.6.0 — delete with this file; callers can call PowerOnVm without args.
func powerOnVMWithEtag(ctx context.Context, conn *vmm.Client, vmID *string) (import1.TaskReference, error) {
	args, err := vmEtagArgs(ctx, conn, utils.StringValue(vmID))
	if err != nil {
		return import1.TaskReference{}, err
	}
	powerOnVmRequest := import3.PowerOnVmRequest{
		ExtId: vmID,
	}
	resp, err := conn.VMAPIInstance.PowerOnVm(ctx, &powerOnVmRequest, args)
	if err != nil {
		return import1.TaskReference{}, fmt.Errorf("error powering on VM: %v", err)
	}
	taskRef, err := extractTaskReferenceFromResponse(resp)
	if err != nil {
		return import1.TaskReference{}, fmt.Errorf("error extracting task reference from power on response: %v", err)
	}
	log.Printf("[DEBUG] PowerOn Response: TaskReference ExtId: %s", utils.StringValue(taskRef.ExtId))
	return taskRef, nil
}

// powerOffVMWithEtag performs a single PowerOffVm call with a fresh If-Match eTag.
// REMOVE_AFTER: v2.6.0 — delete with this file; callers can call PowerOffVm without args.
func powerOffVMWithEtag(ctx context.Context, conn *vmm.Client, vmID *string) (import1.TaskReference, error) {
	args, err := vmEtagArgs(ctx, conn, utils.StringValue(vmID))
	if err != nil {
		return import1.TaskReference{}, err
	}
	powerOffVmRequest := import3.PowerOffVmRequest{
		ExtId: vmID,
	}
	resp, err := conn.VMAPIInstance.PowerOffVm(ctx, &powerOffVmRequest, args)
	if err != nil {
		return import1.TaskReference{}, fmt.Errorf("error powering off VM: %v", err)
	}
	taskRef, err := extractTaskReferenceFromResponse(resp)
	if err != nil {
		return import1.TaskReference{}, fmt.Errorf("error extracting task reference from power off response: %v", err)
	}
	log.Printf("[DEBUG] PowerOff Response: TaskReference ExtId: %s", utils.StringValue(taskRef.ExtId))
	return taskRef, nil
}
