package lcmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/common"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourcePreloadArtifactsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourcePreloadArtifactsV2Create,
		ReadContext:   ResourcePreloadArtifactsV2Read,
		UpdateContext: ResourcePreloadArtifactsV2Update,
		DeleteContext: ResourcePreloadArtifactsV2Delete,
		Schema: map[string]*schema.Schema{
			"x_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
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
						"entity_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func ResourcePreloadArtifactsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI

	clusterId := d.Get("x_cluster_id").(string)

	body := common.NewPreloadSpec()

	entityUpdateSpecs := d.Get("entity_update_specs").([]interface{})

	body.EntityUpdateSpecs = expandEntityUpdateSpecs(entityUpdateSpecs)

	resp, err := conn.LcmEntitiesAPIInstance.PreloadArtifacts(body, utils.StringPtr(clusterId))
	if err != nil {
		return diag.Errorf("error while Perform Preload Artifacts: %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the Config Update to be successful
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroup(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("LCM Upgrade task failed: %s", errWaitTask)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching the Preload Artifacts task : %v", err)
	}

	task := resourceUUID.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(task, "", "  ")
	log.Printf("[DEBUG] LCM Preload Artifacts Task Response: %s", string(aJSON))

	// randomly generating the id
	d.SetId(utils.GenUUID())

	return ResourcePreloadArtifactsV2Read(ctx, d, meta)
}

func ResourcePreloadArtifactsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourcePreloadArtifactsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourcePreloadArtifactsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
