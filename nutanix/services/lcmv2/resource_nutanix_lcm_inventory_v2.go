package lcmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	lcmOps "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/operations"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/request/inventory"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
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
			"inventory_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"FULL", "SOFTWARE", "NODE", "RESCAN"}, false),
			},
			"node_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	_, hasInvType := d.GetOk("inventory_type")
	_, hasNodeList := d.GetOk("node_list")

	var body *lcmOps.InventorySpec
	if hasInvType || hasNodeList {
		body = lcmOps.NewInventorySpec()
		if v, ok := d.GetOk("inventory_type"); ok {
			body.InventoryType = common.ExpandEnum[lcmOps.InventoryType](v)
		}
		if v, ok := d.GetOk("node_list"); ok {
			nodeListRaw := v.([]interface{})
			nodeList := make([]string, 0, len(nodeListRaw))
			for _, n := range nodeListRaw {
				nodeList = append(nodeList, n.(string))
			}
			body.NodeList = nodeList
		}
	}

	performInventoryRequest := import1.PerformInventoryRequest{
		XClusterId: clusterID,
		Body:       body,
		Dryrun_:    nil,
	}
	resp, err := conn.LcmInventoryAPIInstance.PerformInventory(ctx, &performInventoryRequest)
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
	getTaskByIdRequest := import4.GetTaskByIdRequest{
		ExtId: taskUUID,
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
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
