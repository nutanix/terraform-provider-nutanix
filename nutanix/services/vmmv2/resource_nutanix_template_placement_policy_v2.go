package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	catalogConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/catalogCommon/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/prism/v4/config"
	vmmConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/templateplacementpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixTemplatePlacementPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixTemplatePlacementPolicyV2Create,
		ReadContext:   resourceNutanixTemplatePlacementPolicyV2Read,
		UpdateContext: resourceNutanixTemplatePlacementPolicyV2Update,
		DeleteContext: resourceNutanixTemplatePlacementPolicyV2Delete,
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
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"placement_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SOFT",
				ValidateFunc: validation.StringInSlice([]string{"SOFT"}, false),
			},
			"cluster_filter": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_ext_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"}, false),
						},
					},
				},
			},
			"content_filter": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_ext_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"}, false),
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixTemplatePlacementPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	body := &vmmConfig.TemplatePlacementPolicy{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if pt, ok := d.GetOk("placement_type"); ok {
		body.PlacementType = common.ExpandEnum[catalogConfig.ContentPlacementType](pt.(string))
	}
	if cf, ok := d.GetOk("cluster_filter"); ok {
		body.ClusterFilter = expandCategoriesFilter(cf)
	}
	if cf, ok := d.GetOk("content_filter"); ok {
		body.ContentFilter = expandCategoriesFilter(cf)
	}

	req := import3.CreateTemplatePlacementPolicyRequest{
		Body: body,
	}
	resp, err := conn.TemplatePlacementPoliciesAPIInstance.CreateTemplatePlacementPolicy(ctx, &req)
	if err != nil {
		return diag.Errorf("error creating template placement policy: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template placement policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIdRequest := import4.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching template placement policy create task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Template Placement Policy Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeTemplatePlacementPolicy, "Template placement policy")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return resourceNutanixTemplatePlacementPolicyV2Read(ctx, d, meta)
}

func resourceNutanixTemplatePlacementPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	req := import3.GetTemplatePlacementPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TemplatePlacementPoliciesAPIInstance.GetTemplatePlacementPolicyById(ctx, &req)
	if err != nil {
		return diag.Errorf("error reading template placement policy: %v", err)
	}

	policy := resp.Data.GetValue().(vmmConfig.TemplatePlacementPolicy)

	if err := d.Set("ext_id", policy.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", policy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", policy.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("placement_type", common.FlattenPtrEnum(policy.PlacementType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_filter", flattenCategoriesFilter(policy.ClusterFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("content_filter", flattenCategoriesFilter(policy.ContentFilter)); err != nil {
		return diag.FromErr(err)
	}
	if policy.CreateTime != nil {
		if err := d.Set("create_time", policy.CreateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_by", policy.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if policy.UpdateTime != nil {
		if err := d.Set("update_time", policy.UpdateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("updated_by", policy.UpdatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", policy.TenantId); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNutanixTemplatePlacementPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	getReq := import3.GetTemplatePlacementPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	readResp, err := conn.TemplatePlacementPoliciesAPIInstance.GetTemplatePlacementPolicyById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error reading template placement policy for update: %v", err)
	}

	updateSpec := readResp.Data.GetValue().(vmmConfig.TemplatePlacementPolicy)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("placement_type") {
		updateSpec.PlacementType = common.ExpandEnum[catalogConfig.ContentPlacementType](d.Get("placement_type").(string))
	}
	if d.HasChange("cluster_filter") {
		updateSpec.ClusterFilter = expandCategoriesFilter(d.Get("cluster_filter"))
	}
	if d.HasChange("content_filter") {
		updateSpec.ContentFilter = expandCategoriesFilter(d.Get("content_filter"))
	}

	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	updateReq := import3.UpdateTemplatePlacementPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  &updateSpec,
	}
	updateResp, err := conn.TemplatePlacementPoliciesAPIInstance.UpdateTemplatePlacementPolicyById(ctx, &updateReq, args)
	if err != nil {
		return diag.Errorf("error updating template placement policy: %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template placement policy (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return resourceNutanixTemplatePlacementPolicyV2Read(ctx, d, meta)
}

func resourceNutanixTemplatePlacementPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	req := import3.DeleteTemplatePlacementPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TemplatePlacementPoliciesAPIInstance.DeleteTemplatePlacementPolicyById(ctx, &req)
	if err != nil {
		return diag.Errorf("error deleting template placement policy: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template placement policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandCategoriesFilter(pr interface{}) *catalogConfig.CategoriesFilter {
	if pr == nil {
		return nil
	}
	prI := pr.([]interface{})
	if len(prI) == 0 || prI[0] == nil {
		return nil
	}
	val := prI[0].(map[string]interface{})

	filter := &catalogConfig.CategoriesFilter{}

	if t, ok := val["type"]; ok && t.(string) != "" {
		filter.Type = common.ExpandEnum[catalogConfig.CategoriesMatchType](t.(string))
	}
	if catExtIds, ok := val["category_ext_ids"]; ok {
		categoryExtIdsList := common.InterfaceToSlice(catExtIds)
		filter.CategoryExtIds = common.ExpandListOfString(categoryExtIdsList)
	}

	return filter
}

func flattenCategoriesFilter(pr *catalogConfig.CategoriesFilter) []map[string]interface{} {
	if pr == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0)
	filter := make(map[string]interface{})

	filter["type"] = common.FlattenPtrEnum(pr.Type)
	filter["category_ext_ids"] = utils.StringSlice(pr.CategoryExtIds)

	result = append(result, filter)
	return result
}
