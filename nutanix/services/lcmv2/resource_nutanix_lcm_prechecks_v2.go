package lcmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	preCheckConfig "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/common"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixPreChecksV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixLcmPreChecksV2Create,
		ReadContext:   ResourceNutanixLcmPreChecksV2Read,
		UpdateContext: ResourceNutanixLcmPreChecksV2Update,
		DeleteContext: ResourceNutanixLcmPreChecksV2Delete,
		Schema: map[string]*schema.Schema{
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"management_server": {
				Type:     schema.TypeList,
				MaxItems: 1,
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
				Type:     schema.TypeList,
				Required: true,
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
			},
			"skipped_precheck_flags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixLcmPreChecksV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterExtID := d.Get("x_cluster_id").(string)
	body := preCheckConfig.NewPrechecksSpec()

	if managementServer, ok := d.GetOk("management_server"); ok {
		body.ManagementServer = expandManagementServer(managementServer.([]interface{}))
	}
	if entityUpdateSpecs, ok := d.GetOk("entity_update_specs"); ok {
		body.EntityUpdateSpecs = expandEntityUpdateSpecs(entityUpdateSpecs.([]interface{}))
	}
	if skippedPrecheckFlags, ok := d.GetOk("skipped_precheck_flags"); ok {
		body.SkippedPrecheckFlags = expandSystemAutoMgmtFlag(skippedPrecheckFlags.([]interface{}))
	}

	// pass nil for the new dyRun flag
	resp, err := conn.LcmPreChecksAPIInstance.PerformPrechecks(body, utils.StringPtr(clusterExtID), nil)
	if err != nil {
		return diag.Errorf("error while performing the prechecks: %v", err)
	}
	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the LCM prechecks to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for LCM prechecks (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching LCM prechecks task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] LCM Prechecks Task Details: %s", string(aJSON))

	// This is an action resource that does not maintain state.
	// The resource ID is set to the task ExtId for traceability.
	d.SetId(utils.StringValue(taskDetails.ExtId))
	return nil
}

func ResourceNutanixLcmPreChecksV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixLcmPreChecksV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixLcmPreChecksV2Create(ctx, d, meta)
}

func ResourceNutanixLcmPreChecksV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
