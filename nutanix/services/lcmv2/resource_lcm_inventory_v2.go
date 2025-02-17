package lcmv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	lcmInventoryResp "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/operations"
	import1 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceLcmInventoryV2() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ntnx_request_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceLcmPerformInventoryV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterId := d.Get("x_cluster_id").(string)
	ntnxRequestId, ok := d.Get("ntnx_request_id").(string)
	if !ok || ntnxRequestId == "" {
		return diag.Errorf("ntnx_request_id is required and cannot be null or empty")
	}

	args := make(map[string]interface{})
	args["X-Cluster-Id"] = clusterId
	args["NTNX-Request-Id"] = ntnxRequestId

	resp, err := conn.LcmInventoryAPIInstance.PerformInventory(&clusterId, args)
	if err != nil {
		return diag.Errorf("error while performing the inventory: %v", err)
	}
	getResp := resp.Data.GetValue().(lcmInventoryResp.InventoryApiResponse)
	TaskRef := getResp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the inventorty to be successful
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroup(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Perform inventory task failed: %s", errWaitTask)
	}
	return nil
}

func taskStateRefreshPrismTaskGroup(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)
		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(import1.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(pr *import1.TaskStatus) string {
	const two, three, five, six, seven = 2, 3, 5, 6, 7
	if pr != nil {
		if *pr == import1.TaskStatus(six) {
			return "FAILED"
		}
		if *pr == import1.TaskStatus(seven) {
			return "CANCELED"
		}
		if *pr == import1.TaskStatus(two) {
			return "QUEUED"
		}
		if *pr == import1.TaskStatus(three) {
			return "RUNNING"
		}
		if *pr == import1.TaskStatus(five) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}
