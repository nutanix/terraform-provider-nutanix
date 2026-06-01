package microsegv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/directoryserverconfigs"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/prism/v4/config"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixDirectoryServerConfigV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixDirectoryServerConfigV2Create,
		ReadContext:   resourceNutanixDirectoryServerConfigV2Read,
		UpdateContext: resourceNutanixDirectoryServerConfigV2Update,
		DeleteContext: resourceNutanixDirectoryServerConfigV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_service_reference": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_controllers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schemaForIPAddressOrFQDN(),
			},
			"is_default_category_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"matching_criterias": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schemaForMatchingCriteria(),
			},
			"should_keep_default_category_on_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixDirectoryServerConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	bodySpec := import2.NewDirectoryServerConfig()

	bodySpec.DirectoryServiceReference = utils.StringPtr(d.Get("directory_service_reference").(string))

	if v, ok := d.GetOk("domain_controllers"); ok {
		bodySpec.DomainControllers = expandDomainControllers(v.([]interface{}))
	}
	if v, ok := d.GetOk("is_default_category_enabled"); ok {
		bodySpec.IsDefaultCategoryEnabled = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("matching_criterias"); ok {
		bodySpec.MatchingCriterias = expandMatchingCriterias(v.([]interface{}))
	}
	if v, ok := d.GetOk("should_keep_default_category_on_login"); ok {
		bodySpec.ShouldKeepDefaultCategoryOnLogin = utils.BoolPtr(v.(bool))
	}

	aJSON, _ := json.MarshalIndent(bodySpec, "", "  ")
	log.Printf("[DEBUG] Create Directory Server Config Body Spec: %s", string(aJSON))

	req := import3.CreateDirectoryServerConfigRequest{
		Body: bodySpec,
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.CreateDirectoryServerConfig(ctx, &req)
	if err != nil {
		return diag.Errorf("error while creating Directory Server Config: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Directory Server Config (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIdRequest := import5.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching Directory Server Config Task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Directory Server Config Task Response Details: %s", string(aJSON))

	uuid := taskDetails.CompletionDetails[0].Value.GetValue().(string)

	d.SetId(uuid)
	return resourceNutanixDirectoryServerConfigV2Read(ctx, d, meta)
}

func resourceNutanixDirectoryServerConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Id()
	req := import3.GetDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.GetDirectoryServerConfigById(ctx, &req)
	if err != nil {
		return diag.Errorf("error while fetching Directory Server Config: %s", err)
	}

	getResp := resp.Data.GetValue().(import2.DirectoryServerConfig)

	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] Read Directory Server Config Response Details: %s", string(aJSON))

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksDSC(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("directory_service_reference", getResp.DirectoryServiceReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_controllers", flattenDomainControllers(getResp.DomainControllers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_default_category_enabled", getResp.IsDefaultCategoryEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("matching_criterias", flattenMatchingCriterias(getResp.MatchingCriterias)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("should_keep_default_category_on_login", getResp.ShouldKeepDefaultCategoryOnLogin); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", getResp.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixDirectoryServerConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	getReq := import3.GetDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	readResp, err := conn.DirectoryServerConfigsAPIInstance.GetDirectoryServerConfigById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching Directory Server Config: %v", err)
	}

	args := make(map[string]interface{})
	etag := conn.DirectoryServerConfigsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etag)

	updateSpec := import2.NewDirectoryServerConfig()

	if v, ok := d.GetOk("directory_service_reference"); ok {
		updateSpec.DirectoryServiceReference = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("domain_controllers"); ok {
		updateSpec.DomainControllers = expandDomainControllers(v.([]interface{}))
	}
	if v, ok := d.GetOk("is_default_category_enabled"); ok {
		updateSpec.IsDefaultCategoryEnabled = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("matching_criterias"); ok {
		updateSpec.MatchingCriterias = expandMatchingCriterias(v.([]interface{}))
	}
	if v, ok := d.GetOk("should_keep_default_category_on_login"); ok {
		updateSpec.ShouldKeepDefaultCategoryOnLogin = utils.BoolPtr(v.(bool))
	}

	updateReq := import3.UpdateDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  updateSpec,
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.UpdateDirectoryServerConfigById(ctx, &updateReq, args)
	if err != nil {
		return diag.Errorf("error while updating Directory Server Config: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Directory Server Config (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIdRequest := import5.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching Directory Server Config Task: %v", err)
	}

	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update Directory Server Config Task Response Details: %s", string(aJSON))

	return resourceNutanixDirectoryServerConfigV2Read(ctx, d, meta)
}

func resourceNutanixDirectoryServerConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	deleteReq := import3.DeleteDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.DeleteDirectoryServerConfigById(ctx, &deleteReq)
	if err != nil {
		return diag.Errorf("error while deleting Directory Server Config: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Directory Server Config (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIdRequest := import5.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while deleting Directory Server Config Task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Directory Server Config Task Response Details: %s", string(aJSON))

	return nil
}
