package vmmv2

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ExpandDiskFunc converts a list of disk config (interface{}) to SDK disk types.
// Used by ApplyDiskDeletions, ApplyDiskUpdates, ApplyDiskAdditions so callers can pass their resource-specific expandDisk.
type ExpandDiskFunc func(disks []interface{}) []config.Disk

func getEtagHeader(resp interface{}, conn *vmm.Client) *string {
	// Extract E-Tag Header
	etagValue := conn.VMAPIInstance.ApiClient.GetEtag(resp)
	return utils.StringPtr(etagValue)
}

func isVmmEtagMismatchErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "If-Match header value passed") ||
		strings.Contains(msg, "VM_ETAG_MISMATCH") ||
		strings.Contains(msg, "VMM-30303")
}

// StripDataSourceFromDiskBackingInfo removes the data_source key from a disk's backing_info.vm_disk map.
// Call this on each updated disk map before expanding for UpdateDiskById so the API accepts the payload.
func StripDataSourceFromDiskBackingInfo(disk interface{}) {
	diskMap, ok := disk.(map[string]interface{})
	if !ok {
		return
	}
	backingInfoRaw, ok := diskMap["backing_info"]
	if !ok {
		return
	}
	backingInfoSlice, ok := backingInfoRaw.([]interface{})
	if !ok || len(backingInfoSlice) == 0 {
		return
	}
	backingInfoMap, ok := backingInfoSlice[0].(map[string]interface{})
	if !ok {
		return
	}
	vmDiskArray, ok := backingInfoMap["vm_disk"].([]interface{})
	if !ok || len(vmDiskArray) == 0 {
		return
	}
	vmDiskMap, ok := vmDiskArray[0].(map[string]interface{})
	if !ok {
		return
	}
	if vmDiskMap["data_source"] != nil {
		delete(vmDiskMap, "data_source")
	}
}

func waitForDiskTask(ctx context.Context, d *schema.ResourceData, meta interface{}, taskUUID *string, timeoutType string, operation string) diag.Diagnostics {
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(timeoutType),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for disk (%s) to %s: %s", utils.StringValue(taskUUID), operation, errWait)
	}
	return nil
}

// ApplyDiskDeletions deletes the given disks from the VM and waits for each task.
func ApplyDiskDeletions(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *vmm.Client, vmID string, deletedDisks []interface{}, expandDisk ExpandDiskFunc) diag.Diagnostics {
	if len(deletedDisks) == 0 {
		return nil
	}
	for _, disk := range deletedDisks {
		diskInputs := expandDisk([]interface{}{disk})
		if len(diskInputs) == 0 {
			continue
		}
		diskInput := diskInputs[0]
		diskExtID := diskInput.ExtId

		readVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmID))
		if err != nil {
			return diag.Errorf("error while fetching vm : %v", err)
		}
		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readVMResp, conn)

		resp, err := conn.VMAPIInstance.DeleteDiskById(utils.StringPtr(vmID), diskExtID, args)
		if err != nil {
			return diag.Errorf("error while deleting Disk : %v", err)
		}
		taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
		if err := waitForDiskTask(ctx, d, meta, taskRef.ExtId, schema.TimeoutDelete, "be deleted"); err != nil {
			return err
		}
	}
	return nil
}

// ApplyDiskUpdates updates the given disks on the VM and waits for each task.
// Strips data_source from each disk's backing_info before sending.
func ApplyDiskUpdates(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *vmm.Client, vmID string, updatedDisks []interface{}, expandDisk ExpandDiskFunc) diag.Diagnostics {
	if len(updatedDisks) == 0 {
		return nil
	}
	for _, disk := range updatedDisks {
		StripDataSourceFromDiskBackingInfo(disk)
		diskInputs := expandDisk([]interface{}{disk})
		if len(diskInputs) == 0 {
			continue
		}
		diskInput := diskInputs[0]
		diskExtID := diskInput.ExtId

		readVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmID))
		if err != nil {
			return diag.Errorf("error while fetching vm : %v", err)
		}
		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readVMResp, conn)

		resp, err := conn.VMAPIInstance.UpdateDiskById(utils.StringPtr(vmID), diskExtID, &diskInput, args)
		if err != nil {
			return diag.Errorf("error while updating Disk : %v", err)
		}
		taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
		if err := waitForDiskTask(ctx, d, meta, taskRef.ExtId, schema.TimeoutUpdate, "be updated"); err != nil {
			return err
		}
	}
	return nil
}

// ApplyDiskAdditions creates the given disks on the VM and waits for each task.
func ApplyDiskAdditions(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *vmm.Client, vmID string, addedDisks []interface{}, expandDisk ExpandDiskFunc) diag.Diagnostics {
	if len(addedDisks) == 0 {
		return nil
	}
	for _, disk := range addedDisks {
		diskInputs := expandDisk([]interface{}{disk})
		if len(diskInputs) == 0 {
			continue
		}
		diskInput := diskInputs[0]

		readVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmID))
		if err != nil {
			return diag.Errorf("error while fetching vm : %v", err)
		}
		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readVMResp, conn)

		resp, err := conn.VMAPIInstance.CreateDisk(utils.StringPtr(vmID), &diskInput, args)
		if err != nil {
			return diag.Errorf("error while creating Disk : %v", err)
		}
		taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
		if err := waitForDiskTask(ctx, d, meta, taskRef.ExtId, schema.TimeoutCreate, "add"); err != nil {
			return err
		}
	}
	return nil
}
