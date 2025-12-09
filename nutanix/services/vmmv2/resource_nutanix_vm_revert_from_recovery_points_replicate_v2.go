package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmPrismConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRevertVMRecoveryPointV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRevertVMRecoveryPointV2Create,
		ReadContext:   ResourceNutanixRevertVMRecoveryPointV2Read,
		UpdateContext: ResourceNutanixRevertVMRecoveryPointV2Update,
		DeleteContext: ResourceNutanixRevertVMRecoveryPointV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_recovery_point_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// ResourceNutanixRevertVMRecoveryPointV2Create to Restore Recovery Point
func ResourceNutanixRevertVMRecoveryPointV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixRevertVMRecoveryPointV2Create \n")

	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching Vm : %v", err)
	}
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	body := config.RevertParams{}
	rpExtID := d.Get("ext_id").(string)

	if v, ok := d.GetOk("vm_recovery_point_ext_id"); ok {
		body.VmRecoveryPointExtId = utils.StringPtr(v.(string))
	}

	resp, err := conn.VMAPIInstance.RevertVm(utils.StringPtr(rpExtID), &body, args)
	if err != nil {
		return diag.Errorf("error while reverting vm : %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be reverted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM revert (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching VM revert task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Revert VM Task Details: %s", string(aJSON))

	if err = d.Set("status", common.FlattenPtrEnum(taskDetails.Status)); err != nil {
		return diag.FromErr(err)
	}

	uuid, err := common.ExtractCompletionDetailFromTask(taskDetails, utils.CompletionDetailsNameVMExtIDs, "VM")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid)

	return ResourceNutanixRevertVMRecoveryPointV2Read(ctx, d, meta)
}

func ResourceNutanixRevertVMRecoveryPointV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRevertVMRecoveryPointV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixRevertVMRecoveryPointV2Read(ctx, d, meta)
}

func ResourceNutanixRevertVMRecoveryPointV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
