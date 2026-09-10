package microsegv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/request/directoryserverconfigs"
	prismMicroseg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	prismTasks "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixCategoryMappingV2 manages the mapping between a group in
// Active Directory and a Prism Central category.
func ResourceNutanixCategoryMappingV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixCategoryMappingV2Create,
		ReadContext:   ResourceNutanixCategoryMappingV2Read,
		UpdateContext: ResourceNutanixCategoryMappingV2Update,
		DeleteContext: ResourceNutanixCategoryMappingV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceCategoryMappingSchema(),
	}
}

func ResourceNutanixCategoryMappingV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	body := expandCategoryMappingBody(d)

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Category Mapping Body: %s", string(aJSON))

	createRequest := import3.CreateCategoryMappingRequest{
		Body: body,
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.CreateCategoryMapping(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error creating Category Mapping: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in CreateCategoryMapping response")
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
		return diag.Errorf("error waiting for Category Mapping create: %s", errWait)
	}

	getTaskRequest := prismTasks.GetTaskByIdRequest{
		ExtId: taskUUID,
	}
	taskResp, err := taskConn.TaskRefAPI.GetTaskById(ctx, &getTaskRequest)
	if err != nil {
		return diag.Errorf("error fetching Category Mapping create task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeCategoryMapping, "Category Mapping")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixCategoryMappingV2Read(ctx, d, meta)
}

func ResourceNutanixCategoryMappingV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	getRequest := import3.GetDsCategoryMappingByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.GetDsCategoryMappingById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error reading Category Mapping: %v", err)
	}

	getResp, ok := resp.Data.GetValue().(import2.CategoryMapping)
	if !ok {
		return diag.Errorf("invalid CategoryMapping in response")
	}

	if err := d.Set("ext_id", utils.StringValue(getResp.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(getResp.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_name", utils.StringValue(getResp.CategoryName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_value", utils.StringValue(getResp.CategoryValue)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ad_info", flattenAdInfo(getResp.AdInfo)); err != nil {
		return diag.FromErr(err)
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

func ResourceNutanixCategoryMappingV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	getRequest := import3.GetDsCategoryMappingByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	readResp, err := conn.DirectoryServerConfigsAPIInstance.GetDsCategoryMappingById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error reading Category Mapping for update: %v", err)
	}

	body := readResp.Data.GetValue().(import2.CategoryMapping)

	if d.HasChange("name") {
		body.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("category_name") {
		body.CategoryName = utils.StringPtr(d.Get("category_name").(string))
	}
	if d.HasChange("category_value") {
		body.CategoryValue = utils.StringPtr(d.Get("category_value").(string))
	}
	if d.HasChange("ad_info") {
		if v, ok := d.GetOk("ad_info"); ok {
			body.AdInfo = expandAdInfo(v.([]interface{}))
		}
	}
	if d.HasChange("project_ext_id") {
		return diag.Errorf("error updating project_ext_id: Update of project_ext_id is not supported")
	}
	args := make(map[string]interface{})
	etag := conn.DirectoryServerConfigsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etag)

	updateRequest := import3.UpdateDsCategoryMappingByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  &body,
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.UpdateDsCategoryMappingById(ctx, &updateRequest, args)
	if err != nil {
		return diag.Errorf("error updating Category Mapping: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in UpdateDsCategoryMappingById response")
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
		return diag.Errorf("error waiting for Category Mapping update: %s", errWait)
	}

	return ResourceNutanixCategoryMappingV2Read(ctx, d, meta)
}

func ResourceNutanixCategoryMappingV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	deleteRequest := import3.DeleteDsCategoryMappingByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.DeleteDsCategoryMappingById(ctx, &deleteRequest)
	if err != nil {
		return diag.Errorf("error deleting Category Mapping: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return diag.Errorf("invalid TaskReference in DeleteDsCategoryMappingById response")
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
		return diag.Errorf("error waiting for Category Mapping delete: %s", errWait)
	}

	return nil
}

func expandCategoryMappingBody(d *schema.ResourceData) *import2.CategoryMapping {
	body := import2.NewCategoryMapping()
	body.Name = utils.StringPtr(d.Get("name").(string))
	body.CategoryName = utils.StringPtr(d.Get("category_name").(string))
	body.CategoryValue = utils.StringPtr(d.Get("category_value").(string))
	if v, ok := d.GetOk("ad_info"); ok {
		body.AdInfo = expandAdInfo(v.([]interface{}))
	}
	if v, ok := d.GetOk("project_ext_id"); ok {
		body.ProjectExtId = utils.StringPtr(v.(string))
	}
	return body
}
