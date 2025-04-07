package dataprotectionv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	dataprotectionPrismConfig "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRestoreProtectedResourceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRestoreProtectedResourceV2Create,
		ReadContext:   ResourceNutanixRestoreProtectedResourceV2Read,
		UpdateContext: ResourceNutanixRestoreProtectedResourceV2Update,
		DeleteContext: ResourceNutanixRestoreProtectedResourceV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restore_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// ResourceNutanixRestoreProtectedResourceV2Create to Restore Protected Resource
// This Resource is action resource and does not have any state
// resource id is set to random UUID
func ResourceNutanixRestoreProtectedResourceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataProtectionAPI

	extID := d.Get("ext_id").(string)

	bodySpec := config.NewProtectedResourceRestoreSpec()

	if clusterExtID, ok := d.GetOk("cluster_ext_id"); ok {
		bodySpec.ClusterExtId = utils.StringPtr(clusterExtID.(string))
	}
	if restoreTime, ok := d.GetOk("restore_time"); ok {
		bodySpec.RestoreTime = utils.Time(restoreTime.(time.Time))
	}

	resp, err := conn.ProtectedResource.RestoreProtectedResource(utils.StringPtr(extID), bodySpec)
	if err != nil {
		return diag.Errorf("Error while restoring protected resource: %s", err)
	}

	TaskRef := resp.Data.GetValue().(dataprotectionPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for task to complete: %s", err)
	}

	// Get UUID from TASK API

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("Error while getting task by ID: %s", err)
	}

	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Restore Protected Resource Task Details: %s", aJSON)

	d.SetId(utils.GenUUID())

	return nil
}

// ResourceNutanixRestoreProtectedResourceV2Read to Restore Protected Resource
func ResourceNutanixRestoreProtectedResourceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixRestoreProtectedResourceV2Update to Restore Protected Resource
func ResourceNutanixRestoreProtectedResourceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixRestoreProtectedResourceV2Delete to Restore Protected Resource
func ResourceNutanixRestoreProtectedResourceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
