package microsegv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/directoryserverconfigs"
	prismMicroseg "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/prism/v4/config"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	prismTasks "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixDirectoryServerConfigV2 configures various aspects of
// identity categorization for Flow Network Security ID-based security.
func ResourceNutanixDirectoryServerConfigV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixDirectoryServerConfigV2Create,
		ReadContext:   ResourceNutanixDirectoryServerConfigV2Read,
		UpdateContext: ResourceNutanixDirectoryServerConfigV2Update,
		DeleteContext: ResourceNutanixDirectoryServerConfigV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:        resourceDirectoryServerConfigSchema(),
		CustomizeDiff: customizeDiffDirectoryServerConfig,
	}
}

func customizeDiffDirectoryServerConfig(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	criterias, ok := d.GetOk("matching_criterias")
	if !ok {
		return nil
	}
	critList, ok := criterias.([]interface{})
	if !ok || len(critList) == 0 {
		return nil
	}

	// Use the raw user config to determine whether the user explicitly set
	// "criteria". During refresh or when transitioning from CONTAINS to ALL,
	// the state may still carry the old criteria value from the API even
	// though the new HCL config omits it.
	rawConfig := d.GetRawConfig()
	var rawCriterias []cty.Value
	if !rawConfig.IsNull() {
		mcAttr := rawConfig.GetAttr("matching_criterias")
		if mcAttr.IsKnown() && !mcAttr.IsNull() {
			rawCriterias = mcAttr.AsValueSlice()
		}
	}

	hasAllMatchType := false
	for i, raw := range critList {
		mc, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		matchType, _ := mc["match_type"].(string)

		if matchType == "ALL" {
			hasAllMatchType = true
			// Only reject if the user explicitly wrote criteria in their HCL
			if i < len(rawCriterias) && !rawCriterias[i].IsNull() {
				rawCrit := rawCriterias[i].GetAttr("criteria")
				if rawCrit.IsKnown() && !rawCrit.IsNull() && rawCrit.AsString() != "" {
					return fmt.Errorf(
						"matching_criterias[%d]: 'criteria' must not be set when match_type is \"ALL\". "+
							"The 'criteria' field is only allowed when match_type is \"CONTAINS\"", i,
					)
				}
			}
		}
	}

	isDefaultCategoryEnabled, _ := d.GetOk("is_default_category_enabled")
	shouldKeep, _ := d.GetOk("should_keep_default_category_on_login")
	defaultEnabled, _ := isDefaultCategoryEnabled.(bool)
	keepOnLogin, _ := shouldKeep.(bool)

	if hasAllMatchType {
		if defaultEnabled {
			return fmt.Errorf(
				"'is_default_category_enabled' must be false when match_type is \"ALL\". " +
					"It can only be set to true when match_type is \"CONTAINS\"",
			)
		}
		if keepOnLogin {
			return fmt.Errorf(
				"'should_keep_default_category_on_login' must be false when match_type is \"ALL\". " +
					"It can only be set to true when match_type is \"CONTAINS\"",
			)
		}
	}

	if keepOnLogin && !defaultEnabled {
		return fmt.Errorf(
			"'should_keep_default_category_on_login' can only be true when " +
				"'is_default_category_enabled' is also true",
		)
	}
	return nil
}

func ResourceNutanixDirectoryServerConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	body := expandDirectoryServerConfigBody(d)

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Directory Server Config Body: %s", string(aJSON))

	createRequest := import3.CreateDirectoryServerConfigRequest{
		Body: body,
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.CreateDirectoryServerConfig(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error creating Directory Server Config: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in CreateDirectoryServerConfig response")
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
		return diag.Errorf("error waiting for Directory Server Config create: %s", errWait)
	}

	getTaskRequest := prismTasks.GetTaskByIdRequest{
		ExtId: taskUUID,
	}
	taskResp, err := taskConn.TaskRefAPI.GetTaskById(ctx, &getTaskRequest)
	if err != nil {
		return diag.Errorf("error fetching Directory Server Config create task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeDirectoryServerConfig, "Directory Server Config")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixDirectoryServerConfigV2Read(ctx, d, meta)
}

func ResourceNutanixDirectoryServerConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	getRequest := import3.GetDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.GetDirectoryServerConfigById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error reading Directory Server Config: %v", err)
	}

	getResp, ok := resp.Data.GetValue().(import2.DirectoryServerConfig)
	if !ok {
		return diag.Errorf("invalid DirectoryServerConfig in response")
	}

	if err := d.Set("ext_id", utils.StringValue(getResp.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("directory_service_reference", utils.StringValue(getResp.DirectoryServiceReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_controllers", flattenIPAddressOrFQDNList(getResp.DomainControllers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("matching_criterias", flattenMatchingCriterias(getResp.MatchingCriterias)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.IsDefaultCategoryEnabled != nil {
		if err := d.Set("is_default_category_enabled", utils.BoolValue(getResp.IsDefaultCategoryEnabled)); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.ShouldKeepDefaultCategoryOnLogin != nil {
		if err := d.Set("should_keep_default_category_on_login", utils.BoolValue(getResp.ShouldKeepDefaultCategoryOnLogin)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("project_ext_id", utils.StringValue(getResp.ProjectExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixDirectoryServerConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	getRequest := import3.GetDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	readResp, err := conn.DirectoryServerConfigsAPIInstance.GetDirectoryServerConfigById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error reading Directory Server Config for update: %v", err)
	}

	body := readResp.Data.GetValue().(import2.DirectoryServerConfig)

	if d.HasChange("directory_service_reference") {
		body.DirectoryServiceReference = utils.StringPtr(d.Get("directory_service_reference").(string))
	}
	if d.HasChange("domain_controllers") {
		if v, ok := d.GetOk("domain_controllers"); ok {
			body.DomainControllers = expandIPAddressOrFQDNList(v.([]interface{}))
		}
	}
	if d.HasChange("is_default_category_enabled") {
		body.IsDefaultCategoryEnabled = utils.BoolPtr(d.Get("is_default_category_enabled").(bool))
	}
	if d.HasChange("matching_criterias") {
		if v, ok := d.GetOk("matching_criterias"); ok {
			body.MatchingCriterias = expandDirectoryServerConfigMatchingCriterias(v.([]interface{}))
		}
		for i := range body.MatchingCriterias {
			if body.MatchingCriterias[i].MatchType != nil && body.MatchingCriterias[i].MatchType.GetName() == "ALL" {
				body.MatchingCriterias[i].Criteria = nil
			}
		}
	}
	if d.HasChange("should_keep_default_category_on_login") {
		body.ShouldKeepDefaultCategoryOnLogin = utils.BoolPtr(d.Get("should_keep_default_category_on_login").(bool))
	}
	if d.HasChange("project_ext_id") {
		return diag.Errorf("error updating project_ext_id: Update of project_ext_id is not supported")
	}

	args := make(map[string]interface{})
	etag := conn.DirectoryServerConfigsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etag)

	updateRequest := import3.UpdateDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  &body,
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.UpdateDirectoryServerConfigById(ctx, &updateRequest, args)
	if err != nil {
		return diag.Errorf("error updating Directory Server Config: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in UpdateDirectoryServerConfigById response")
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
		return diag.Errorf("error waiting for Directory Server Config update: %s", errWait)
	}

	return ResourceNutanixDirectoryServerConfigV2Read(ctx, d, meta)
}

func ResourceNutanixDirectoryServerConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	deleteRequest := import3.DeleteDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.DeleteDirectoryServerConfigById(ctx, &deleteRequest)
	if err != nil {
		return diag.Errorf("error deleting Directory Server Config: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in DeleteDirectoryServerConfigById response")
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
		return diag.Errorf("error waiting for Directory Server Config delete: %s", errWait)
	}

	return nil
}

func expandDirectoryServerConfigBody(d *schema.ResourceData) *import2.DirectoryServerConfig {
	body := import2.NewDirectoryServerConfig()
	if v, ok := d.GetOk("directory_service_reference"); ok {
		body.DirectoryServiceReference = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("domain_controllers"); ok {
		body.DomainControllers = expandIPAddressOrFQDNList(v.([]interface{}))
	}
	if v, ok := d.GetOkExists("is_default_category_enabled"); ok {
		body.IsDefaultCategoryEnabled = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("matching_criterias"); ok {
		body.MatchingCriterias = expandDirectoryServerConfigMatchingCriterias(v.([]interface{}))
	}
	if v, ok := d.GetOkExists("should_keep_default_category_on_login"); ok {
		body.ShouldKeepDefaultCategoryOnLogin = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("project_ext_id"); ok {
		body.ProjectExtId = utils.StringPtr(v.(string))
	}
	return body
}
