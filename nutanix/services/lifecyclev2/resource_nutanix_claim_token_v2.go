package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixClaimTokenV2 manages a Foundation Central claim token. A claim
// token is used to register nodes with Foundation Central; it carries a name, an
// expiry time, and a maximum usage count. Create/Update/Delete are asynchronous
// task-based operations, matching the canonical v2 resources.
func ResourceNutanixClaimTokenV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClaimTokenV2Create,
		ReadContext:   ResourceNutanixClaimTokenV2Read,
		UpdateContext: ResourceNutanixClaimTokenV2Update,
		DeleteContext: ResourceNutanixClaimTokenV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the claim token",
			},
			"expiry_time": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsRFC3339Time,
				Description:  "Expiry time of the claim token",
			},
			"max_usage_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum number of times the claim token can be used for registering nodes",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the claim token was created",
			},
			"current_usage_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of times the claim token has been used for registering nodes",
			},
			"owner_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External ID of the owner",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A HATEOAS style link for the response.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixClaimTokenV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	body := import1.ClaimToken{}
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("expiry_time"); ok {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("error parsing expiry_time: %v", err)
		}
		body.ExpiryTime = &t
	}
	body.MaxUsageCount = utils.IntPtr(d.Get("max_usage_count").(int))

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Claim Token create payload : %s", string(aJSON))

	resp, err := conn.ClaimTokensAPIInstance.CreateClaimToken(&body)
	if err != nil {
		return diag.Errorf("error while creating claim token : %v", err)
	}

	taskReference := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for claim token (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get the claim token UUID from the completed task.
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching claim token task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Claim Token Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeClaimToken, "ClaimToken")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixClaimTokenV2Read(ctx, d, meta)
}

func ResourceNutanixClaimTokenV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.ClaimTokensAPIInstance.GetClaimTokenById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching claim token : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.ClaimToken)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if getResp.ExpiryTime != nil {
		if err := d.Set("expiry_time", getResp.ExpiryTime.Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("max_usage_count", getResp.MaxUsageCount); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("current_usage_count", getResp.CurrentUsageCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLifecycleLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixClaimTokenV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.ClaimTokensAPIInstance.GetClaimTokenById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching claim token : %v", err)
	}

	updateSpec := resp.Data.GetValue().(import1.ClaimToken)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("expiry_time") {
		t, errParse := time.Parse(time.RFC3339, d.Get("expiry_time").(string))
		if errParse != nil {
			return diag.Errorf("error parsing expiry_time: %v", errParse)
		}
		updateSpec.ExpiryTime = &t
	}
	if d.HasChange("max_usage_count") {
		updateSpec.MaxUsageCount = utils.IntPtr(d.Get("max_usage_count").(int))
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", " ")
	log.Printf("[DEBUG] Claim Token update payload : %s", string(aJSON))

	updateResp, err := conn.ClaimTokensAPIInstance.UpdateClaimTokenById(utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating claim token : %v", err)
	}

	taskReference := updateResp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for claim token (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixClaimTokenV2Read(ctx, d, meta)
}

func ResourceNutanixClaimTokenV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.ClaimTokensAPIInstance.DeleteClaimTokenById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting claim token : %v", err)
	}

	taskReference := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := taskReference.ExtId

	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for claim token (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
