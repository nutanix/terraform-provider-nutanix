package multidomainv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/projects"
	multidomainPrism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/prism/v4/config"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixProjectV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixProjectV2Create,
		ReadContext:   ResourceNutanixProjectV2Read,
		UpdateContext: ResourceNutanixProjectV2Update,
		DeleteContext: ResourceNutanixProjectV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"modified_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"links": schemaForLinks(),
		},
	}
}

func ResourceNutanixProjectV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	body := &config.Project{}
	
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		body.Id = utils.StringPtr(v.(string))
	}
	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Project Body: %s", string(aJSON))

	createReq := import1.CreateProjectRequest{
		Body: body,
	}
	resp, err := conn.Projects.CreateProject(ctx, &createReq)
	if err != nil {
		return diag.Errorf("error creating Project: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("create project response did not contain task reference")
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
		return diag.Errorf("error waiting for project create task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	getTaskReq := import3.GetTaskByIdRequest{ExtId: taskUUID}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskReq)
	if err != nil {
		return diag.Errorf("error fetching project create task: %v", err)
	}
	taskDetails, ok := taskResp.Data.GetValue().(prismConfig.Task)
	if !ok {
		return diag.Errorf("error parsing task response")
	}

	uuid, err := commonUtils.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeProject, "Project")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixProjectV2Read(ctx, d, meta)
}

func ResourceNutanixProjectV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	extID := d.Id()
	getReq := import1.GetProjectByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.Projects.GetProjectById(ctx, &getReq)
	if err != nil {
		return diag.FromErr(err)
	}

	project := resp.Data.GetValue().(config.Project)
	
	if err := d.Set("name", utils.StringValue(project.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", utils.StringValue(project.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_id", utils.StringValue(project.Id)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", utils.StringValue(project.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(project.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", project.State.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_default", utils.BoolValue(project.IsDefault)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_system_defined", utils.BoolValue(project.IsSystemDefined)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", utils.StringValue(project.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_by", utils.StringValue(project.UpdatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_timestamp", utils.Int64Value(project.CreatedTimestamp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("modified_timestamp", utils.Int64Value(project.ModifiedTimestamp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(project.Links)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixProjectV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI
  extID := d.Get("ext_id").(string)
	getReq := import1.GetProjectByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.Projects.GetProjectById(ctx, &getReq)
	if err != nil {
		return diag.FromErr(err)
	}
  
	// Extract E-Tag Header
	etagValue := conn.APIClientInstance.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)
	getResp := resp.Data.GetValue().(config.Project)

	if d.HasChange("name") {
		return diag.Errorf("error while updating name: Update of name is not supported")
	}
	if d.HasChange("project_id") {
		return diag.Errorf("error while updating project_id: Update of project_id is not supported")
	}
	updateSpec := &config.Project{
		Name: getResp.Name,
		Id: getResp.Id,
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}

	updateReq := import1.UpdateProjectByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  updateSpec,
	}
	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Update Project Body: %s", string(aJSON))
	updateResp, err := conn.Projects.UpdateProjectById(ctx, &updateReq, args)
	if err != nil {
		return diag.Errorf("error updating Project: %v", err)
	}

	taskRef, ok := updateResp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("update project response did not contain task reference")
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
		return diag.Errorf("error waiting for project update task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	return ResourceNutanixProjectV2Read(ctx, d, meta)
}

func ResourceNutanixProjectV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI
  
	extID := d.Get("ext_id").(string)
	getReq := import1.GetProjectByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.Projects.GetProjectById(ctx, &getReq)
	if err != nil {
		return diag.FromErr(err)
	}
  
	// Extract E-Tag Header
	etagValue := conn.APIClientInstance.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)
	
	deleteReq := import1.DeleteProjectByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	deleteResp, err := conn.Projects.DeleteProjectById(ctx, &deleteReq, args)
	if err != nil {
		return diag.Errorf("error deleting Project: %v", err)
	}

	taskRef, ok := deleteResp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("delete project response did not contain task reference")
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
		return diag.Errorf("error waiting for project delete task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	d.SetId("")
	return nil
}
