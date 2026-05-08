package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmPrismConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmGuestCustomizationProfileV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmGuestCustomizationProfileV2Create,
		ReadContext:   ResourceNutanixVmGuestCustomizationProfileV2Read,
		UpdateContext: ResourceNutanixVmGuestCustomizationProfileV2Update,
		DeleteContext: ResourceNutanixVmGuestCustomizationProfileV2Delete,
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
			"config": schemaForVmGcProfileConfig(false),
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"updated_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixVmGuestCustomizationProfileV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	body := *config.NewVmGuestCustomizationProfile()

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if cfg, ok := d.GetOk("config"); ok {
		cfgList := cfg.([]interface{})
		if len(cfgList) > 0 {
			body.Config = expandVmGcProfileConfig(cfgList)
		}
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] VM Guest Customization Profile Create Request Body: %s", string(aJSON))

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.CreateVmGuestCustomizationProfile(&body)
	if err != nil {
		return diag.Errorf("error while creating VM Guest Customization Profile: %v", err)
	}

	taskRef := resp.Data.GetValue().(vmmPrismConfig.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Guest Customization Profile (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] VM Guest Customization Profile Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVmGuestCustomizationProfile, "VM Guest Customization Profile")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixVmGuestCustomizationProfileV2Read(ctx, d, meta)
}

func ResourceNutanixVmGuestCustomizationProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile: %v", err)
	}

	profile := resp.Data.GetValue().(config.VmGuestCustomizationProfile)

	if err := d.Set("ext_id", utils.StringValue(profile.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(profile.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenVmGcProfileLinks(profile.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(profile.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", utils.StringValue(profile.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config", flattenVmGcProfileConfig(profile.Config)); err != nil {
		return diag.FromErr(err)
	}
	if profile.CreateTime != nil {
		if err := d.Set("create_time", profile.CreateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if profile.UpdateTime != nil {
		if err := d.Set("update_time", profile.UpdateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_by", flattenVmGcProfileUserReference(profile.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_by", flattenVmGcProfileUserReference(profile.UpdatedBy)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixVmGuestCustomizationProfileV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	readResp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile for update: %v", err)
	}

	updateSpec := readResp.Data.GetValue().(config.VmGuestCustomizationProfile)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("config") {
		if cfg, ok := d.GetOk("config"); ok {
			cfgList := cfg.([]interface{})
			if len(cfgList) > 0 {
				updateSpec.Config = expandVmGcProfileConfig(cfgList)
			}
		}
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] VM Guest Customization Profile Update Request Body: %s", string(aJSON))

	args := make(map[string]interface{})
	etagValue := conn.VmGuestCustomizationProfilesAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etagValue)

	updateResp, err := conn.VmGuestCustomizationProfilesAPIInstance.UpdateVmGuestCustomizationProfileById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating VM Guest Customization Profile: %v", err)
	}

	taskRef := updateResp.Data.GetValue().(vmmPrismConfig.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Guest Customization Profile (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return ResourceNutanixVmGuestCustomizationProfileV2Read(ctx, d, meta)
}

func ResourceNutanixVmGuestCustomizationProfileV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	readResp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile for delete: %v", err)
	}

	args := make(map[string]interface{})
	etagValue := conn.VmGuestCustomizationProfilesAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etagValue)

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.DeleteVmGuestCustomizationProfileById(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting VM Guest Customization Profile: %v", err)
	}

	taskRef := resp.Data.GetValue().(vmmPrismConfig.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Guest Customization Profile (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}
