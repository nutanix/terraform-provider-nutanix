package lcmv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	preCheckConfig "github.com/nutanix/ntnx-api-golang-clients/lcm-go-client/v4/models/lcm/v4/common"
	lcmPreCheckResp "github.com/nutanix/ntnx-api-golang-clients/lcm-go-client/v4/models/lcm/v4/operations"
	import1 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourcePreChecksV2() *schema.Resource {
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
			"management_server": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hypervisor_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"entity_update_specs": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"to_version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"skipped_precheck_flags": {
				Type:     schema.TypeList,
				Optional: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceLcmPreChecksV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterId := d.Get("x_cluster_id").(string)
	ntnxRequestId, ok := d.Get("ntnx_request_id").(string)
	if !ok || ntnxRequestId == "" {
		return diag.Errorf("ntnx_request_id is required and cannot be null or empty")
	}

	args := make(map[string]interface{})
	args["X-Cluster-Id"] = clusterId
	args["NTNX-Request-Id"] = ntnxRequestId
	body := preCheckConfig.PrecheckSpec{}

	resp, err := conn.LcmPreChecksAPIInstance.Precheck(&body, args)
	if err != nil {
		return diag.Errorf("error while performing the prechecs: %v", err)
	}
	getResp := resp.Data.GetValue().(lcmPreCheckResp.PrecheckApiResponse)
	TaskRef := getResp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the inventorty to be successful
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFun(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Prechecks task failed: %s", errWaitTask)
	}
	return nil

}

func taskStateRefreshPrismTaskGroupFun(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)
		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(import1.Task)

		if getTaskStat(v.Status) == "CANCELED" || getTaskStat(v.Status) == "FAILED" {
			return v, getTaskStat(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStat(v.Status), nil
	}
}

func getTaskStat(pr *import1.TaskStatus) string {
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
