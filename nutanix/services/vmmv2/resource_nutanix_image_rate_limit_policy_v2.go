package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import7 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/images/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/imageratelimitpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixImageRateLimitPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixImageRateLimitPolicyV2Create,
		ReadContext:   resourceNutanixImageRateLimitPolicyV2Read,
		UpdateContext: resourceNutanixImageRateLimitPolicyV2Update,
		DeleteContext: resourceNutanixImageRateLimitPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the image rate limit policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Image rate limit policy specification.",
			},
			"rate_limit_kbps": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Network bandwidth in KBps that the rate limited image operation can utilize.",
			},
			"cluster_entity_filter": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"}, false),
						},
						"category_ext_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"matching_cluster_ext_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "External identifier of the Prism Elements where a rate limit is the effective rate limit policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External identifier of the owner of the rate limit policy.",
			},
			"owner_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the owner of the rate limit policy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image rate limit policy creation time.",
			},
			"last_update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time of an image rate limit policy.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
		},
	}
}

func resourceNutanixImageRateLimitPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	body := &import7.RateLimitPolicy{}

	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("rate_limit_kbps"); ok {
		body.RateLimitKbps = utils.Int64Ptr(int64(v.(int)))
	}
	if v, ok := d.GetOk("cluster_entity_filter"); ok {
		body.ClusterEntityFilter = expandRateLimitClusterEntityFilter(v)
	}

	req := import3.CreateRateLimitPolicyRequest{
		Body: body,
	}

	resp, err := conn.ImageRateLimitPoliciesAPIInstance.CreateRateLimitPolicy(ctx, &req)
	if err != nil {
		return diag.Errorf("error creating image rate limit policy: %v", err)
	}

	taskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image rate limit policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	getTaskByIdRequest := import4.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching image rate limit policy create task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Image Rate Limit Policy Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeImageRateLimitPolicy, "Image rate limit policy")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return resourceNutanixImageRateLimitPolicyV2Read(ctx, d, meta)
}

func resourceNutanixImageRateLimitPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	req := import3.GetRateLimitPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.ImageRateLimitPoliciesAPIInstance.GetRateLimitPolicyById(ctx, &req)
	if err != nil {
		return diag.Errorf("error reading image rate limit policy: %v", err)
	}

	policy := resp.Data.GetValue().(import7.RateLimitPolicy)

	if err := d.Set("ext_id", policy.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", policy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", policy.Description); err != nil {
		return diag.FromErr(err)
	}
	if policy.RateLimitKbps != nil {
		if err := d.Set("rate_limit_kbps", int(*policy.RateLimitKbps)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("cluster_entity_filter", flattenRateLimitClusterEntityFilter(policy.ClusterEntityFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("matching_cluster_ext_ids", utils.StringSlice(policy.MatchingClusterExtIds)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", policy.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_name", policy.OwnerName); err != nil {
		return diag.FromErr(err)
	}
	if policy.CreateTime != nil {
		if err := d.Set("create_time", policy.CreateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if policy.LastUpdateTime != nil {
		if err := d.Set("last_update_time", policy.LastUpdateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("tenant_id", policy.TenantId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixImageRateLimitPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	getReq := import3.GetRateLimitPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	getResp, err := conn.ImageRateLimitPoliciesAPIInstance.GetRateLimitPolicyById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error reading image rate limit policy for update: %v", err)
	}

	policy := getResp.Data.GetValue().(import7.RateLimitPolicy)

	if d.HasChange("name") {
		policy.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		policy.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("rate_limit_kbps") {
		policy.RateLimitKbps = utils.Int64Ptr(int64(d.Get("rate_limit_kbps").(int)))
	}
	if d.HasChange("cluster_entity_filter") {
		policy.ClusterEntityFilter = expandRateLimitClusterEntityFilter(d.Get("cluster_entity_filter"))
	}

	updateReq := import3.UpdateRateLimitPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
		Body:  &policy,
	}

	etagValue := conn.ImageRateLimitPoliciesAPIInstance.ApiClient.GetEtag(getResp)
	args := make(map[string]interface{})
	args["If-Match"] = etagValue

	resp, err := conn.ImageRateLimitPoliciesAPIInstance.UpdateRateLimitPolicyById(ctx, &updateReq, args)
	if err != nil {
		return diag.Errorf("error updating image rate limit policy: %v", err)
	}

	taskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image rate limit policy (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return resourceNutanixImageRateLimitPolicyV2Read(ctx, d, meta)
}

func resourceNutanixImageRateLimitPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	req := import3.DeleteRateLimitPolicyByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	resp, err := conn.ImageRateLimitPoliciesAPIInstance.DeleteRateLimitPolicyById(ctx, &req)
	if err != nil {
		return diag.Errorf("error deleting image rate limit policy: %v", err)
	}

	taskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image rate limit policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func expandRateLimitClusterEntityFilter(pr interface{}) *import7.Filter {
	if pr == nil {
		return nil
	}
	prI := pr.([]interface{})
	if len(prI) == 0 || prI[0] == nil {
		return nil
	}
	val := prI[0].(map[string]interface{})

	filter := &import7.Filter{}

	if t, ok := val["type"]; ok {
		filter.Type = common.ExpandEnum[import7.FilterMatchType](t.(string))
	}
	if categoryExtIds, ok := val["category_ext_ids"]; ok {
		categoryExtIdsList := common.InterfaceToSlice(categoryExtIds)
		filter.CategoryExtIds = common.ExpandListOfString(categoryExtIdsList)
	}

	return filter
}

func flattenRateLimitClusterEntityFilter(pr *import7.Filter) []map[string]interface{} {
	if pr == nil {
		return nil
	}
	filter := make(map[string]interface{})
	filter["type"] = common.FlattenPtrEnum(pr.Type)
	filter["category_ext_ids"] = utils.StringSlice(pr.CategoryExtIds)
	return []map[string]interface{}{filter}
}
