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
	import5 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixTemplateActionsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTemplateActionsV2Create,
		ReadContext:   ResourceNutanixTemplateActionsV2Read,
		UpdateContext: ResourceNutanixTemplateActionsV2Update,
		DeleteContext: ResourceNutanixTemplateActionsV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"initiate", "complete", "cancel"}, false),
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_active_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
			},
		},
	}
}

func ResourceNutanixTemplateActionsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	taskconn := meta.(*conns.Client).PrismAPI

	extID := d.Get("ext_id").(string)
	action := d.Get("action").(string)

	var taskUUID *string

	if action == "initiate" {
		versionID := d.Get("version_id").(string)
		spec := import5.InitiateGuestUpdateSpec{}

		spec.VersionId = &versionID
		resp, err := conn.TemplatesAPIInstance.InitiateGuestUpdate(utils.StringPtr(extID), &spec)
		if err != nil {
			return diag.FromErr(err)
		}
		TaskRef := resp.Data.GetValue().(import1.TaskReference)
		taskUUID = TaskRef.ExtId

		// Wait for the guest OS update to be initiated
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for guest OS update initiation (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
		}
	}

	if action == "complete" {
		versionName := d.Get("version_name").(string)
		versionDesc := d.Get("version_description").(string)
		isActiveVersion := d.Get("is_active_version").(bool)

		spec := import5.CompleteGuestUpdateSpec{}
		spec.VersionName = &versionName
		spec.VersionDescription = &versionDesc
		spec.IsActiveVersion = &isActiveVersion

		resp, err := conn.TemplatesAPIInstance.CompleteGuestUpdate(utils.StringPtr(extID), &spec)
		if err != nil {
			return diag.FromErr(err)
		}
		TaskRef := resp.Data.GetValue().(import1.TaskReference)
		taskUUID = TaskRef.ExtId

		// Wait for the guest OS update to complete
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for guest OS update completion (%s) to finish: %s", utils.StringValue(taskUUID), errWaitTask)
		}
	}

	if action == "cancel" {
		resp, err := conn.TemplatesAPIInstance.CancelGuestUpdate(utils.StringPtr(extID))
		if err != nil {
			return diag.FromErr(err)
		}
		TaskRef := resp.Data.GetValue().(import1.TaskReference)
		taskUUID = TaskRef.ExtId

		// Wait for the guest OS update to be cancelled
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for guest OS update cancellation (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
		}
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching template guest OS action task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Template Guest OS Action (%s) Task Details: %s", action, string(aJSON))

	// This is an action resource that does not maintain state.
	// The resource ID is set to the task ExtId for traceability.
	d.SetId(utils.StringValue(taskDetails.ExtId))
	return nil
}

func ResourceNutanixTemplateActionsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixTemplateActionsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixTemplateActionsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
