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
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import7 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/images/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixImagePlacementV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixImagePlacementV2Create,
		ReadContext:   ResourceNutanixImagePlacementV2Read,
		UpdateContext: ResourceNutanixImagePlacementV2Update,
		DeleteContext: ResourceNutanixImagePlacementV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HARD", "SOFT"}, false),
			},
			"image_entity_filter": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"}, false),
						},
						"category_ext_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"cluster_entity_filter": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"}, false),
						},
						"category_ext_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"enforcement_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "SUSPENDED"}, false),
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SUSPEND", "RESUME"}, false),
			},
			"should_cancel_running_tasks": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixImagePlacementV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	body := &import7.PlacementPolicy{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if placementType, ok := d.GetOk("placement_type"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"SOFT": two,
			"HARD": three,
		}
		pVal := subMap[placementType.(string)]
		p := import7.PlacementType(pVal.(int))
		body.PlacementType = &p
	}
	if imageEntityFilter, ok := d.GetOk("image_entity_filter"); ok {
		body.ImageEntityFilter = expandEntityFilter(imageEntityFilter)
	}
	if clusterEntityFilter, ok := d.GetOk("cluster_entity_filter"); ok {
		body.ClusterEntityFilter = expandEntityFilter(clusterEntityFilter)
	}
	if enforcementState, ok := d.GetOk("enforcement_state"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"ACTIVE":    two,
			"SUSPENDED": three,
		}
		pVal := subMap[enforcementState.(string)]
		p := import7.EnforcementState(pVal.(int))
		body.EnforcementState = &p
	}
	resp, err := conn.ImagesPlacementAPIInstance.CreatePlacementPolicy(body)
	if err != nil {
		return diag.Errorf("error while creating Image placement policy : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the image placement policy to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image placement policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching image placement policy create task (%s): %v", utils.StringValue(taskUUID), errorMessage["message"])
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Image Placement Policy Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeImagePlacementPolicy, "Image placement policy")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixImagePlacementV2Read(ctx, d, meta)
}

func ResourceNutanixImagePlacementV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.ImagesPlacementAPIInstance.GetPlacementPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching image placement policy : %v", err)
	}

	getResp := resp.Data.GetValue().(import7.PlacementPolicy)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("placement_type", flattenPlacementType(getResp.PlacementType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("image_entity_filter", flattenEntityFilter(getResp.ImageEntityFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_entity_filter", flattenEntityFilter(getResp.ClusterEntityFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enforcement_state", flattenEnforcementState(getResp.EnforcementState)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdateTime != nil {
		t := getResp.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixImagePlacementV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.ImagesPlacementAPIInstance.GetPlacementPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching image placement policy : %v", err)
	}

	respImagePlacements := resp.Data.GetValue().(import7.PlacementPolicy)
	updateSpec := respImagePlacements

	changed := false

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
		changed = true
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("placement_type") {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"SOFT": two,
			"HARD": three,
		}
		pVal := subMap[d.Get("placement_type").(string)]
		p := import7.PlacementType(pVal.(int))
		updateSpec.PlacementType = &p
		changed = true
	}
	if d.HasChange("image_entity_filter") {
		updateSpec.ImageEntityFilter = expandEntityFilter(d.Get("image_entity_filter"))
		changed = true
	}
	if d.HasChange("cluster_entity_filter") {
		updateSpec.ClusterEntityFilter = expandEntityFilter(d.Get("cluster_entity_filter"))
		changed = true
	}
	if d.HasChange("enforcement_state") {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"ACTIVE":    two,
			"SUSPENDED": three,
		}
		pVal := subMap[d.Get("enforcement_state").(string)]
		p := import7.EnforcementState(pVal.(int))
		updateSpec.EnforcementState = &p
		changed = true
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		if action == "SUSPEND" {
			suspendAction(ctx, conn, d, meta)
		} else if action == "RESUME" {
			resumeAction(ctx, conn, d, meta)
		}
	}

	if changed {
		updateResp, er := conn.ImagesPlacementAPIInstance.UpdatePlacementPolicyById(utils.StringPtr(d.Id()), &updateSpec)
		if er != nil {
			return diag.Errorf("error while updating image placement policy : %v", err)
		}
		TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
		taskUUID := TaskRef.ExtId

		taskconn := meta.(*conns.Client).PrismAPI
		// Wait for the image placement policy to be updated
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for image placement policy (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
		}
	}

	return ResourceNutanixImagePlacementV2Read(ctx, d, meta)
}

func suspendAction(ctx context.Context, conn *vmm.Client, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	extID := d.Id()

	readResp, err := conn.ImagesPlacementAPIInstance.GetPlacementPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while reading placement policy : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	body := &import7.SuspendPlacementPolicyConfig{}

	if shouldCancelRunningTasks, ok := d.GetOk("should_cancel_running_tasks"); ok {
		body.ShouldCancelRunningTasks = utils.BoolPtr(shouldCancelRunningTasks.(bool))
	}

	resp, err := conn.ImagesPlacementAPIInstance.SuspendPlacementPolicy(utils.StringPtr(extID), body, args)
	if err != nil {
		return diag.Errorf("error while suspend Image placement policy : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the image placement policy to be suspended
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image placement policy (%s) to suspend: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func resumeAction(ctx context.Context, conn *vmm.Client, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	extID := d.Id()
	readResp, err := conn.ImagesPlacementAPIInstance.GetPlacementPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while reading placement policy : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.ImagesPlacementAPIInstance.ResumePlacementPolicy(utils.StringPtr(extID), args)
	if err != nil {
		return diag.Errorf("error while resume Image placement policy : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the image placement policy to be resumed
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image placement policy (%s) to resume: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func ResourceNutanixImagePlacementV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.ImagesPlacementAPIInstance.DeletePlacementPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while deleting image placement policy : %v", errorMessage["message"])
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the image placement policy to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for image placement policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandEntityFilter(pr interface{}) *import7.Filter {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		entityFilter := &import7.Filter{}

		// entity_filter.ObjectType_ = utils.StringPtr("vmm.v4.r0.b1.images.config.Filter")

		if ftype, ok := val["type"]; ok {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"CATEGORIES_MATCH_ALL": two,
				"CATEGORIES_MATCH_ANY": three,
			}
			pVal := subMap[ftype.(string)]
			p := import7.FilterMatchType(pVal.(int))
			entityFilter.Type = &p
		}
		if categoryExtIds, ok := val["category_ext_ids"]; ok {
			categoriesList := categoryExtIds.([]interface{})
			categories := make([]string, len(categoriesList))

			for k, v := range categoriesList {
				categories[k] = v.(string)
			}
			entityFilter.CategoryExtIds = categories
		}

		return entityFilter
	}
	return nil
}
