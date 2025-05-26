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
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vm: (%s) to revert: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching revert vm UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ := json.Marshal(rUUID)
	log.Printf("[DEBUG] revert vm task Details: %v", string(aJSON))

	if err = d.Set("status", getTaskStatus(rUUID.Status)); err != nil {
		return diag.FromErr(err)
	}

	uuid := rUUID.CompletionDetails[0].Value
	d.SetId(uuid.GetValue().(string))

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

// func taskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
// 	return func() (interface{}, string, error) {
// 		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
// 		// encodeUUID := data + ":" + taskUUID
// 		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)

// 		if err != nil {
// 			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
// 		}

// 		// get the group results

// 		v := vresp.Data.GetValue().(prismConfig.Task)

// 		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
// 			return v, getTaskStatus(v.Status),
// 				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
// 		}
// 		return v, getTaskStatus(v.Status), nil
// 	}
// }

// func getTaskStatus(pr *prismConfig.TaskStatus) string {
// 	if pr != nil {
// 		if *pr == prismConfig.TaskStatus(6) {
// 			return "FAILED"
// 		}
// 		if *pr == prismConfig.TaskStatus(7) {
// 			return "CANCELED"
// 		}
// 		if *pr == prismConfig.TaskStatus(2) {
// 			return "QUEUED"
// 		}
// 		if *pr == prismConfig.TaskStatus(3) {
// 			return "RUNNING"
// 		}
// 		if *pr == prismConfig.TaskStatus(5) {
// 			return "SUCCEEDED"
// 		}
// 	}
// 	return "UNKNOWN"
// }
