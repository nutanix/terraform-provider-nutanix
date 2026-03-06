package multidomainv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/resourcegroups"
	multidomainPrism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/prism/v4/config"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixResourceGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixResourceGroupV2Create,
		ReadContext:   ResourceNutanixResourceGroupV2Read,
		UpdateContext: ResourceNutanixResourceGroupV2Update,
		DeleteContext: ResourceNutanixResourceGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"placement_targets": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schemaResourceGroupPlacementTargets(),
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
		},
	}
}

func ResourceNutanixResourceGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	body := config.NewResourceGroup()
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("project_ext_id"); ok {
		body.ProjectExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		body.TenantId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("placement_targets"); ok {
		body.PlacementTargets = expandResourceGroupPlacementTargets(v.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create ResourceGroup Body: %s", string(aJSON))

	createReq := import1.CreateResourceGroupRequest{
		Body: body,
	}
	resp, err := conn.ResourceGroups.CreateResourceGroup(ctx, &createReq)
	if err != nil {
		return diag.Errorf("error creating ResourceGroup: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("create resource group response did not contain task reference")
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for resource group create task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	getTaskReq := import3.GetTaskByIdRequest{ExtId: taskUUID}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskReq)
	if err != nil {
		return diag.Errorf("error fetching resource group create task: %v", err)
	}
	taskDetails, ok := taskResp.Data.GetValue().(prismConfig.Task)
	if !ok {
		return diag.Errorf("error parsing task response")
	}

	values := commonUtils.ExtractCompletionDetailsFromTask(taskDetails, utils.CompletionDetailsNameResourceGroup)
	if len(values) == 0 {
		return diag.Errorf("resource group ext_id not found in task completion details")
	}
	d.SetId(values[0])

	return ResourceNutanixResourceGroupV2Read(ctx, d, meta)
}

func ResourceNutanixResourceGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	extID := d.Id()
	getReq := import1.GetResourceGroupByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.ResourceGroups.GetResourceGroupById(ctx, &getReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Data == nil {
		d.SetId("")
		return nil
	}

	var rg config.ResourceGroup
	switch v := resp.Data.GetValue().(type) {
	case config.ResourceGroup:
		rg = v
	case *config.ResourceGroup:
		if v == nil {
			d.SetId("")
			return nil
		}
		rg = *v
	default:
		d.SetId("")
		return nil
	}

	_ = d.Set("name", utils.StringValue(rg.Name))
	_ = d.Set("project_ext_id", utils.StringValue(rg.ProjectExtId))
	_ = d.Set("tenant_id", utils.StringValue(rg.TenantId))
	_ = d.Set("ext_id", utils.StringValue(rg.ExtId))
	_ = d.Set("created_by", utils.StringValue(rg.CreatedBy))
	_ = d.Set("last_updated_by", utils.StringValue(rg.LastUpdatedBy))
	if rg.CreateTime != nil {
		_ = d.Set("create_time", rg.CreateTime.Format("2006-01-02T15:04:05Z07:00"))
	}
	if rg.LastUpdateTime != nil {
		_ = d.Set("last_update_time", rg.LastUpdateTime.Format("2006-01-02T15:04:05Z07:00"))
	}
	_ = d.Set("placement_targets", flattenResourceGroupPlacementTargets(rg.PlacementTargets))
	_ = d.Set("links", flattenLinks(rg.Links))

	return nil
}

func ResourceNutanixResourceGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	body := config.NewResourceGroup()
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("project_ext_id"); ok {
		body.ProjectExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		body.TenantId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("placement_targets"); ok {
		body.PlacementTargets = expandResourceGroupPlacementTargets(v.([]interface{}))
	}

	updateReq := import1.UpdateResourceGroupByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  body,
	}
	resp, err := conn.ResourceGroups.UpdateResourceGroupById(ctx, &updateReq)
	if err != nil {
		return diag.Errorf("error updating ResourceGroup: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("update resource group response did not contain task reference")
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for resource group update task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	return ResourceNutanixResourceGroupV2Read(ctx, d, meta)
}

func ResourceNutanixResourceGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	deleteReq := import1.DeleteResourceGroupByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.ResourceGroups.DeleteResourceGroupById(ctx, &deleteReq)
	if err != nil {
		return diag.Errorf("error deleting ResourceGroup: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("delete resource group response did not contain task reference")
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for resource group delete task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	d.SetId("")
	return nil
}

func expandResourceGroupPlacementTargets(in []interface{}) []config.TargetDetails {
	if len(in) == 0 {
		return nil
	}
	out := make([]config.TargetDetails, 0, len(in))
	for _, raw := range in {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		t := config.TargetDetails{}
		if v, ok := m["cluster_ext_id"].(string); ok && v != "" {
			t.ClusterExtId = utils.StringPtr(v)
		}
		if v, ok := m["storage_containers"].([]interface{}); ok && len(v) > 0 {
			t.StorageContainers = expandResourceGroupStorageContainers(v)
		}
		out = append(out, t)
	}
	return out
}

func expandResourceGroupStorageContainers(in []interface{}) []config.StorageContainerDetails {
	if len(in) == 0 {
		return nil
	}
	out := make([]config.StorageContainerDetails, 0, len(in))
	for _, raw := range in {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if v, ok := m["ext_id"].(string); ok && v != "" {
			out = append(out, config.StorageContainerDetails{ExtId: utils.StringPtr(v)})
		}
	}
	return out
}
