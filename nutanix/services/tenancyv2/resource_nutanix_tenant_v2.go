package tenancyv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/prism/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixTenantV2 manages the lifecycle of a tenant in the tenancy namespace.
func ResourceNutanixTenantV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTenantV2Create,
		ReadContext:   ResourceNutanixTenantV2Read,
		UpdateContext: ResourceNutanixTenantV2Update,
		DeleteContext: ResourceNutanixTenantV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"added_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"admin_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"aggregated_quota": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"vcpu": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retention_expiry_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

// ResourceNutanixTenantV2Create creates a new tenant and polls the create task.
func ResourceNutanixTenantV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	body := import1.Tenant{}
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("url"); ok {
		body.Url = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("added_by"); ok {
		body.AddedBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("domains"); ok {
		body.Domains = common.ExpandListOfString(v.([]interface{}))
	}
	if v, ok := d.GetOk("admin_info"); ok {
		body.AdminInfo = expandLocalUserInfo(v.([]interface{}))
	}
	if v, ok := d.GetOk("aggregated_quota"); ok {
		body.AggregatedQuota = expandQuota(v.([]interface{}))
	}

	createTenantRequest := import4.CreateTenantRequest{
		Body: &body,
	}
	aJSON, _ := json.MarshalIndent(createTenantRequest, "", "  ")
	log.Printf("[DEBUG] Tenant create payload: %s", string(aJSON))

	resp, err := conn.TenantManagementAPIInstance.CreateTenant(ctx, &createTenantRequest)
	if err != nil {
		return diag.Errorf("error while creating tenant: %v", err)
	}

	taskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for tenant (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIDRequest := import5.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIDRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant create task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeTenant, "Tenant")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixTenantV2Read(ctx, d, meta)
}

// ResourceNutanixTenantV2Read reads a tenant by its external identifier.
func ResourceNutanixTenantV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	getTenantRequest := import4.GetTenantByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Tenant)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", getResp.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("added_by", getResp.AddedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domains", getResp.Domains); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("aggregated_quota", flattenQuota(getResp.AggregatedQuota)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.RetentionExpiryTime != nil {
		if err := d.Set("retention_expiry_time", getResp.RetentionExpiryTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("state_info", flattenTenantStateInfo(getResp.StateInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	// admin_info: the API never returns the password, so preserve the configured
	// value while updating the other fields from the response.
	adminInfo := flattenLocalUserInfo(getResp.AdminInfo)
	if adminInfo != nil {
		if v, ok := d.GetOk("admin_info"); ok {
			configured := v.([]interface{})
			if len(configured) > 0 && configured[0] != nil {
				if pwd, ok := configured[0].(map[string]interface{})["password"]; ok {
					adminInfo[0]["password"] = pwd
				}
			}
		}
	}
	if err := d.Set("admin_info", adminInfo); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// ResourceNutanixTenantV2Update updates a tenant using a read-modify-write flow.
func ResourceNutanixTenantV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	getTenantRequest := import4.GetTenantByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant: %v", err)
	}

	updateSpec := resp.Data.GetValue().(import1.Tenant)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("url") {
		updateSpec.Url = utils.StringPtr(d.Get("url").(string))
	}
	if d.HasChange("added_by") {
		updateSpec.AddedBy = utils.StringPtr(d.Get("added_by").(string))
	}
	if d.HasChange("domains") {
		updateSpec.Domains = common.ExpandListOfString(d.Get("domains").([]interface{}))
	}
	if d.HasChange("admin_info") {
		updateSpec.AdminInfo = expandLocalUserInfo(d.Get("admin_info").([]interface{}))
	}
	if d.HasChange("aggregated_quota") {
		updateSpec.AggregatedQuota = expandQuota(d.Get("aggregated_quota").([]interface{}))
	}

	// Extract the E-Tag from the read response for optimistic concurrency.
	etagValue := conn.TenantManagementAPIInstance.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updateTenantRequest := import4.UpdateTenantByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  &updateSpec,
	}
	aJSON, _ := json.MarshalIndent(updateTenantRequest, "", "  ")
	log.Printf("[DEBUG] Tenant update payload: %s", string(aJSON))

	updateResp, err := conn.TenantManagementAPIInstance.UpdateTenantById(ctx, &updateTenantRequest, args)
	if err != nil {
		return diag.Errorf("error while updating tenant: %v", err)
	}

	taskRef := updateResp.Data.GetValue().(import3.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for tenant (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return ResourceNutanixTenantV2Read(ctx, d, meta)
}

// ResourceNutanixTenantV2Delete deletes a tenant and polls the delete task.
func ResourceNutanixTenantV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	deleteTenantRequest := import4.DeleteTenantByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantManagementAPIInstance.DeleteTenantById(ctx, &deleteTenantRequest)
	if err != nil {
		return diag.Errorf("error while deleting tenant: %v", err)
	}

	taskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for tenant (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}
