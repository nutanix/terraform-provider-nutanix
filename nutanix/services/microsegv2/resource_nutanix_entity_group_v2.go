package microsegv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	prismMicroseg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixEntityGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixEntityGroupV2Create,
		ReadContext:   ResourceNutanixEntityGroupV2Read,
		UpdateContext: ResourceNutanixEntityGroupV2Update,
		DeleteContext: ResourceNutanixEntityGroupV2Delete,
		// CustomizeDiff: entityGroupCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceEntityGroupSchema(),
	}
}

func ResourceNutanixEntityGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	body := import2.NewEntityGroup()
	body.Name = utils.StringPtr(d.Get("name").(string))
	if v, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("allowed_config"); ok {
		body.AllowedConfig = expandAllowedConfig(v.([]interface{}))
	}
	if v, ok := d.GetOk("except_config"); ok {
		body.ExceptConfig = expandExceptConfig(v.([]interface{}))
	}
	if v, ok := d.GetOk("policy_ext_ids"); ok {
		body.PolicyExtIds = common.ExpandListOfString(v.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Entity Group Body: %s", string(aJSON))

	resp, err := conn.EntityGroupsAPIInstance.CreateEntityGroup(body)
	if err != nil {
		return diag.Errorf("error creating Entity Group: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in CreateEntityGroup response")
	}
	taskUUID := taskRef.ExtId

	taskConn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskConn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for Entity Group create: %s", errWait)
	}

	taskResp, err := taskConn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error fetching Entity Group create task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeEntityGroup, "Entity group")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixEntityGroupV2Read(ctx, d, meta)
}

func ResourceNutanixEntityGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.EntityGroupsAPIInstance.GetEntityGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error reading Entity Group: %v", err)
	}

	getResp, ok := resp.Data.GetValue().(import2.EntityGroup)
	if !ok {
		return diag.Errorf("invalid EntityGroup in response")
	}

	_ = d.Set("ext_id", utils.StringValue(getResp.ExtId))
	_ = d.Set("name", utils.StringValue(getResp.Name))
	_ = d.Set("description", utils.StringValue(getResp.Description))
	_ = d.Set("allowed_config", flattenAllowedConfig(getResp.AllowedConfig))
	_ = d.Set("except_config", flattenExceptConfig(getResp.ExceptConfig))
	_ = d.Set("policy_ext_ids", getResp.PolicyExtIds)
	_ = d.Set("owner_ext_id", utils.StringValue(getResp.OwnerExtId))
	_ = d.Set("tenant_id", utils.StringValue(getResp.TenantId))
	_ = d.Set("links", flattenLinksEntityGroup(getResp.Links))
	_ = d.Set("creation_time", utils.TimeStringValue(getResp.CreationTime))
	_ = d.Set("last_update_time", utils.TimeStringValue(getResp.LastUpdateTime))

	return nil
}

func ResourceNutanixEntityGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	body := import2.NewEntityGroup()
	body.Name = utils.StringPtr(d.Get("name").(string))
	if v, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("allowed_config"); ok {
		body.AllowedConfig = expandAllowedConfig(v.([]interface{}))
	}
	if v, ok := d.GetOk("except_config"); ok {
		body.ExceptConfig = expandExceptConfig(v.([]interface{}))
	}
	if v, ok := d.GetOk("policy_ext_ids"); ok {
		body.PolicyExtIds = common.ExpandListOfString(v.([]interface{}))
	}

	args := make(map[string]interface{})
	readResp, err := conn.EntityGroupsAPIInstance.GetEntityGroupById(utils.StringPtr(d.Id()))
	if err == nil {
		etag := conn.EntityGroupsAPIInstance.ApiClient.GetEtag(readResp)
		args["If-Match"] = utils.StringPtr(etag)
	}

	resp, err := conn.EntityGroupsAPIInstance.UpdateEntityGroupById(utils.StringPtr(d.Id()), body, args)
	if err != nil {
		return diag.Errorf("error updating Entity Group: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in UpdateEntityGroupById response")
	}
	taskUUID := taskRef.ExtId

	taskConn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskConn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for Entity Group update: %s", errWait)
	}

	return ResourceNutanixEntityGroupV2Read(ctx, d, meta)
}

func ResourceNutanixEntityGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.EntityGroupsAPIInstance.DeleteEntityGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error deleting Entity Group: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in DeleteEntityGroupById response")
	}
	taskUUID := taskRef.ExtId

	taskConn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskConn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for Entity Group delete: %s", errWait)
	}

	return nil
}
