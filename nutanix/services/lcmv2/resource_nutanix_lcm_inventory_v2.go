package lcmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixLcmPerformInventoryV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixLcmPerformInventoryV2Create,
		ReadContext:   ResourceNutanixLcmPerformInventoryV2Read,
		UpdateContext: ResourceNutanixLcmPerformInventoryV2Update,
		DeleteContext: ResourceNutanixLcmPerformInventoryV2Delete,
		Schema: map[string]*schema.Schema{
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ResourceNutanixLcmPerformInventoryV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterExtID := d.Get("x_cluster_id").(string)
	var clusterID *string
	if clusterExtID != "" {
		clusterID = utils.StringPtr(clusterExtID)
	} else {
		clusterID = nil
	}
	// pass nil for the body as it is not required and its implemented in hercules Sdk
	// it will be implemented in the future releases of terraform
	resp, err := conn.LcmInventoryAPIInstance.PerformInventory(nil, clusterID, nil)
	if err != nil {
		return diag.Errorf("error while performing the inventory: %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the LCM inventory to be performed
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for LCM inventory (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching LCM inventory task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Perform LCM Inventory Task Details: %s", string(aJSON))

	// This is an action resource that does not maintain state.
	// The resource ID is set to the task ExtId for traceability.
	d.SetId(utils.StringValue(taskDetails.ExtId))
	return nil
}

func ResourceNutanixLcmPerformInventoryV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixLcmPerformInventoryV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixLcmPerformInventoryV2Create(ctx, d, meta)
}

func ResourceNutanixLcmPerformInventoryV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
