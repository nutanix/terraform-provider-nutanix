package volumesv2

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	volumesPrism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v17/models/prism/v4/config"
	volumesClient "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v17/models/volumes/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v17/models/volumes/v4/request/volumegroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	volumesSDK "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/volumes"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixVolumeGroupV2 CRUD for Volume Group.
func ResourceNutanixVolumeGroupV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Creates a new Volume Group.",
		CreateContext: ResourceNutanixVolumeGroupV2Create,
		ReadContext:   ResourceNutanixVolumeGroupV2Read,
		UpdateContext: ResourceNutanixVolumeGroupV2Update,
		DeleteContext: ResourceNutanixVolumeGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Volume Group name. This is an Required field.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Volume Group description. This is an optional field.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"should_load_balance_vm_attachments": {
				Description: "Indicates whether to enable Volume Group load balancing for VM attachments. This cannot be enabled if there are iSCSI client attachments already associated with the Volume Group, and vice-versa. This is an optional field.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"sharing_status": {
				Description:  "Indicates whether the Volume Group can be shared across multiple iSCSI initiators. The mode cannot be changed from SHARED to NOT_SHARED on a Volume Group with multiple attachments. Similarly, a Volume Group cannot be associated with more than one attachment as long as it is in exclusive mode. This is an optional field.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NOT_SHARED", "SHARED"}, false),
			},
			"target_prefix": {
				Description: "The specifications contain the target prefix for external clients as the value. This is an optional field.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"target_name": {
				Description: "Name of the external client target that will be visible and accessible to the client. This is an optional field.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"enabled_authentications": {
				Description:  "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CHAP", "NONE"}, false),
			},
			"iscsi_features": {
				Description: "iSCSI specific settings for the Volume Group. This is an optional field.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_secret": {
							Description: "Target secret in case of a CHAP authentication. This field must only be provided in case the authentication type is not set to CHAP. This is an optional field and it cannot be retrieved once configured.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"enabled_authentications": {
							Description:  "The authentication type enabled for the Volume Group. This is an optional field. If omitted, authentication is not configured for the Volume Group. If this is set to CHAP, the target/client secret must be provided.",
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"CHAP", "NONE"}, false),
						},
					},
				},
			},
			"created_by": {
				Description: "Service/user who created this Volume Group. This is an optional field.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"cluster_reference": {
				Description: "The UUID of the cluster that will host the Volume Group. This is a mandatory field for creating a Volume Group on Prism Central.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"storage_features": {
				Description: "Storage optimization features which must be enabled on the Volume Group. This is an optional field.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flash_mode": {
							Description: "Once configured, this field will avoid down migration of data from the hot tier unless the overrides field is specified for the virtual disks.",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Description: "Indicates whether the flash mode is enabled for the Volume Group.",
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"usage_type": {
				Description:  "Expected usage type for the Volume Group. This is an indicative hint on how the caller will consume the Volume Group. This is an optional field.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "INTERNAL", "TEMPORARY", "BACKUP_TARGET"}, false),
			},
			"attachment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"EXTERNAL", "NONE", "DIRECT"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NOT_ASSIGNED", "ISCSI", "NVMF"}, false),
			},
			"is_hidden": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"disk_size_bytes": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disk_data_source_reference": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uris": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"entity_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"STORAGE_CONTAINER", "VM_DISK", "VOLUME_DISK", "DISK_RECOVERY_POINT"}, false),
									},
								},
							},
						},
						"disk_storage_features": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flash_mode": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixVolumeGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO_VG] Creating Volume Group")
	conn := meta.(*conns.Client).VolumeAPI

	body := volumesClient.VolumeGroup{}

	// Required field
	if name, nok := d.GetOk("name"); nok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if shouldLoadBalanceVMAttachments, ok := d.GetOk("should_load_balance_vm_attachments"); ok {
		body.ShouldLoadBalanceVmAttachments = utils.BoolPtr(shouldLoadBalanceVMAttachments.(bool))
	}
	if sharingStatus, ok := d.GetOk("sharing_status"); ok {
		body.SharingStatus = common.ExpandEnum[volumesClient.SharingStatus](sharingStatus.(string))
	}
	if targetPrefix, ok := d.GetOk("target_prefix"); ok {
		body.TargetPrefix = utils.StringPtr(targetPrefix.(string))
	}
	if targetName, ok := d.GetOk("target_name"); ok {
		body.TargetName = utils.StringPtr(targetName.(string))
	}
	if enabledAuthentications, ok := d.GetOk("enabled_authentications"); ok {
		body.EnabledAuthentications = common.ExpandEnum[volumesClient.AuthenticationType](enabledAuthentications.(string))
	}
	if iscsiFeatures, ok := d.GetOk("iscsi_features"); ok {
		body.IscsiFeatures = expandIscsiFeatures(iscsiFeatures.([]interface{}))
	}
	if createdBy, ok := d.GetOk("created_by"); ok {
		body.CreatedBy = utils.StringPtr(createdBy.(string))
	}
	// Required field
	if clusterReference, ok := d.GetOk("cluster_reference"); ok {
		body.ClusterReference = utils.StringPtr(clusterReference.(string))
	}
	if storageFeatures, ok := d.GetOk("storage_features"); ok {
		body.StorageFeatures = expandStorageFeatures(storageFeatures.([]interface{}))
	}
	if usageType, ok := d.GetOk("usage_type"); ok {
		body.UsageType = common.ExpandEnum[volumesClient.UsageType](usageType.(string))
	}
	if attachmentType, ok := d.GetOk("attachment_type"); ok {
		const NONE, DIRECT, EXTERNAL = 2, 3, 4
		body.AttachmentType = common.ExpandEnum[volumesClient.AttachmentType](attachmentType.(string))
	}
	if protocol, ok := d.GetOk("protocol"); ok {
		body.Protocol = common.ExpandEnum[volumesClient.Protocol](protocol.(string))
	}
	if isHidden, ok := d.GetOk("is_hidden"); ok {
		body.IsHidden = utils.BoolPtr(isHidden.(bool))
	}
	if disks, ok := d.GetOk("disks"); ok {
		body.Disks = expandDisks(disks.([]interface{}))
	}
	if projectExtID, ok := d.GetOk("project_ext_id"); ok {
		body.ProjectExtId = utils.StringPtr(projectExtID.(string))
	}
	createVolumeGroupRequest := import1.CreateVolumeGroupRequest{
		Body: &body,
	}

	aJSON, _ := json.MarshalIndent(createVolumeGroupRequest, "", "  ")
	log.Printf("[DEBUG] Create Volume Group Request: %s", string(aJSON))

	resp, err := conn.VolumeAPIInstance.CreateVolumeGroup(ctx, &createVolumeGroupRequest)
	if err != nil {
		return diag.Errorf("error while creating Volume Group : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the volume group to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for volume group (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	getTaskByIdRequest := import2.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching volume group task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Volume Group Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVolumeGroup, "Volume group")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(uuid))
	d.Set("ext_id", utils.StringValue(uuid))

	return ResourceNutanixVolumeGroupV2Read(ctx, d, meta)
}

func ResourceNutanixVolumeGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	getVolumeGroupByIdRequest := import1.GetVolumeGroupByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.VolumeAPIInstance.GetVolumeGroupById(ctx, &getVolumeGroupByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group : %v", err)
	}

	getResp := resp.Data.GetValue().(volumesClient.VolumeGroup)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("should_load_balance_vm_attachments", getResp.ShouldLoadBalanceVmAttachments); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sharing_status", flattenSharingStatus(getResp.SharingStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_prefix", getResp.TargetPrefix); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_name", getResp.TargetName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled_authentications", flattenEnabledAuthentications(getResp.EnabledAuthentications)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iscsi_features", flattenIscsiFeatures(getResp.IscsiFeatures)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_reference", getResp.ClusterReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_features", flattenStorageFeatures(getResp.StorageFeatures)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("usage_type", flattenUsageType(getResp.UsageType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_hidden", getResp.IsHidden); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", getResp.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixVolumeGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("project_ext_id") {
		return diag.Errorf("error while updating project_ext_id: Update of project_ext_id is not supported")
	}

	conn := meta.(*conns.Client).VolumeAPI

	// Read-modify-write: fetch the current Volume Group, apply only the changed
	// fields onto it, then send the Update so server-populated fields are preserved.
	getVolumeGroupByIdRequest := import1.GetVolumeGroupByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	readResp, err := conn.VolumeAPIInstance.GetVolumeGroupById(ctx, &getVolumeGroupByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching Volume Group : %v", err)
	}

	body := readResp.Data.GetValue().(volumesClient.VolumeGroup)

	// Clear immutable / server-managed fields that the Update API rejects if
	// echoed back in the request body (e.g. the cluster reference cannot be
	// updated, ext_id/created_time/links/tenant_id are read-only, and
	// attachments/attachment_type/protocol/disks are configured out-of-band).
	body.ClusterReference = nil
	body.ExtId = nil
	body.Links = nil
	body.TenantId = nil
	body.Attachments = nil
	body.AttachmentType = nil
	body.Protocol = nil
	body.Disks = nil
	body.HydrationStatus = nil
	// The server also derives/echoes the external client target name & prefix.
	// Re-sending the server-populated value makes the backend attempt to re-set
	// the same iSCSI target name (which collides with the VG itself), so only
	// send these when the user is explicitly changing them below.
	body.TargetName = nil
	body.TargetPrefix = nil

	changed := false
	if d.HasChange("name") {
		body.Name = utils.StringPtr(d.Get("name").(string))
		changed = true
	}
	if d.HasChange("description") {
		body.Description = utils.StringPtr(d.Get("description").(string))
		changed = true
	}
	if d.HasChange("should_load_balance_vm_attachments") {
		body.ShouldLoadBalanceVmAttachments = utils.BoolPtr(d.Get("should_load_balance_vm_attachments").(bool))
		changed = true
	}
	if d.HasChange("sharing_status") {
		body.SharingStatus = common.ExpandEnum[volumesClient.SharingStatus](d.Get("sharing_status").(string))
		changed = true
	}
	if d.HasChange("target_prefix") {
		body.TargetPrefix = utils.StringPtr(d.Get("target_prefix").(string))
		changed = true
	}
	if d.HasChange("target_name") {
		body.TargetName = utils.StringPtr(d.Get("target_name").(string))
		changed = true
	}
	if d.HasChange("enabled_authentications") {
		body.EnabledAuthentications = common.ExpandEnum[volumesClient.AuthenticationType](d.Get("enabled_authentications").(string))
		changed = true
	}
	if d.HasChange("iscsi_features") {
		body.IscsiFeatures = expandIscsiFeatures(d.Get("iscsi_features").([]interface{}))
		changed = true
	}
	if d.HasChange("created_by") {
		body.CreatedBy = utils.StringPtr(d.Get("created_by").(string))
		changed = true
	}
	if d.HasChange("storage_features") {
		body.StorageFeatures = expandStorageFeatures(d.Get("storage_features").([]interface{}))
		changed = true
	}
	if d.HasChange("usage_type") {
		body.UsageType = common.ExpandEnum[volumesClient.UsageType](d.Get("usage_type").(string))
		changed = true
	}
	if d.HasChange("is_hidden") {
		body.IsHidden = utils.BoolPtr(d.Get("is_hidden").(bool))
		changed = true
	}

	if changed {
		updateVolumeGroupByIdRequest := import1.UpdateVolumeGroupByIdRequest{
			ExtId: utils.StringPtr(d.Id()),
			Body:  &body,
		}

		aJSON, _ := json.MarshalIndent(updateVolumeGroupByIdRequest, "", "  ")
		log.Printf("[DEBUG] Update Volume Group Request: %s", string(aJSON))

		// Read-modify-write: fetch the current Volume Group, apply only the changed
		// fields onto it, then send the Update so server-populated fields are preserved.
		getVolumeGroupByIdRequest := import1.GetVolumeGroupByIdRequest{
			ExtId: utils.StringPtr(d.Id()),
		}
		readResp, err := conn.VolumeAPIInstance.GetVolumeGroupById(ctx, &getVolumeGroupByIdRequest)
		if err != nil {
			return diag.Errorf("error while fetching Volume Group : %v", err)
		}

		// Extract E-Tag Header for the If-Match precondition.
		etagValue := conn.VolumeAPIInstance.ApiClient.GetEtag(readResp)
		args := make(map[string]interface{})
		args["If-Match"] = utils.StringPtr(etagValue)

		resp, err := conn.VolumeAPIInstance.UpdateVolumeGroupById(ctx, &updateVolumeGroupByIdRequest, args)
		if err != nil {
			return diag.Errorf("error while updating Volume Group : %v", err)
		}

		TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
		taskUUID := TaskRef.ExtId

		taskconn := meta.(*conns.Client).PrismAPI
		// Wait for the volume group to be updated
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for volume group (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
		}

		// Get task details for logging
		getTaskByIdRequest := import2.GetTaskByIdRequest{
			ExtId: utils.StringPtr(*taskUUID),
		}
		taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
		if err != nil {
			return diag.Errorf("error while fetching volume group update task (%s): %v", utils.StringValue(taskUUID), err)
		}
		taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
		aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
		log.Printf("[DEBUG] Update Volume Group Task Details: %s", string(aJSON))
	}

	// Handle disks update.
	//
	// Disks are NOT part of the Volume Group update body (see body.Disks = nil
	// above); they are managed out-of-band through the dedicated Volume Disk
	// API. Following the same approach as the VM resource, we diff the old vs
	// new disk config into added / deleted / updated buckets and apply each set
	// through its own create / delete / update call.
	if d.HasChange("disks") {
		oldDisks, newDisks := d.GetChange("disks")
		newAddedDisks, oldDeletedDisks, updatedDisks := diffVolumeGroupDisks(oldDisks.([]interface{}), newDisks.([]interface{}))

		if err := applyVolumeGroupDiskDeletions(ctx, d, meta, conn, d.Id(), oldDeletedDisks); err != nil {
			return err
		}
		if err := applyVolumeGroupDiskUpdates(ctx, d, meta, conn, d.Id(), updatedDisks); err != nil {
			return err
		}
		if err := applyVolumeGroupDiskAdditions(ctx, d, meta, conn, d.Id(), newAddedDisks); err != nil {
			return err
		}
	}

	return ResourceNutanixVolumeGroupV2Read(ctx, d, meta)
}

func ResourceNutanixVolumeGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	deleteVolumeGroupByIdRequest := import1.DeleteVolumeGroupByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.VolumeAPIInstance.DeleteVolumeGroupById(ctx, &deleteVolumeGroupByIdRequest)
	if err != nil {
		return diag.Errorf("error while Deleting Volume group : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the volume group to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for volume group (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details for logging
	getTaskByIdRequest := import2.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching volume group delete task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Volume Group Task Details: %s", string(aJSON))

	return nil
}

func expandIscsiFeatures(iscsiFeaturesList interface{}) *volumesClient.IscsiFeatures {
	if len(iscsiFeaturesList.([]interface{})) > 0 {
		iscsiFeature := &volumesClient.IscsiFeatures{}
		iscsiFeaturesI := iscsiFeaturesList.([]interface{})
		if iscsiFeaturesI[0] == nil {
			return nil
		}
		val := iscsiFeaturesI[0].(map[string]interface{})

		if targetSecret, ok := val["target_secret"]; ok {
			iscsiFeature.TargetSecret = utils.StringPtr(targetSecret.(string))
		}

		if enabledAuthentications, ok := val["enabled_authentications"]; ok {
			iscsiFeature.EnabledAuthentications = common.ExpandEnum[volumesClient.AuthenticationType](enabledAuthentications.(string))
		}
		log.Printf("[INFO_VG] iscsiFeature.EnabledAuthentications: %v", *iscsiFeature.EnabledAuthentications)
		log.Printf("[INFO_VG] iscsiFeature.TargetSecret: %v", *iscsiFeature.TargetSecret)
		return iscsiFeature
	}
	return nil
}

func expandStorageFeatures(storageFeaturesList []interface{}) *volumesClient.StorageFeatures {
	if len(storageFeaturesList) > 0 {
		storageFeature := volumesClient.StorageFeatures{}

		val := storageFeaturesList[0].(map[string]interface{})

		if flashMode, ok := val["flash_mode"]; ok {
			storageFeature.FlashMode = expandFlashMode(flashMode.([]interface{}))
		}
		return &storageFeature
	}
	return nil
}

func expandFlashMode(flashModeList []interface{}) *volumesClient.FlashMode {
	if len(flashModeList) > 0 {
		flashMode := volumesClient.FlashMode{}

		val := flashModeList[0].(map[string]interface{})

		if isEnabled, ok := val["is_enabled"]; ok {
			flashMode.IsEnabled = utils.BoolPtr(isEnabled.(bool))
		}
		return &flashMode
	}
	return nil
}

// diskIndex returns the "index" value of a disk config map, or -1 when it is
// missing or not set. The volume group disk schema keys disks by their
// (immutable) index, so it is used as the identity when diffing.
func diskIndex(disk interface{}) int {
	diskMap, ok := disk.(map[string]interface{})
	if !ok {
		return -1
	}
	index, ok := diskMap["index"].(int)
	if !ok {
		return -1
	}
	return index
}

// diffVolumeGroupDisks splits the old vs new disk config into the disks that
// were added, removed, and updated. Identity is the disk index (disks in the
// volume group schema have no ext_id attribute). A disk that exists in both
// old and new but whose contents changed is reported as updated; a disk with no
// index (server-assigned) is always treated as a new addition. This mirrors the
// VM resource's diffConfig, keyed on index instead of ext_id.
func diffVolumeGroupDisks(oldValue []interface{}, newValue []interface{}) ([]interface{}, []interface{}, []interface{}) {
	newlyAdded := make([]interface{}, 0)
	removed := make([]interface{}, 0)
	updated := make([]interface{}, 0)

	oldByIndex := make(map[int]interface{})
	for _, oldItem := range oldValue {
		if idx := diskIndex(oldItem); idx >= 0 {
			oldByIndex[idx] = oldItem
		}
	}
	newByIndex := make(map[int]interface{})
	for _, newItem := range newValue {
		if idx := diskIndex(newItem); idx >= 0 {
			newByIndex[idx] = newItem
		}
	}

	// Additions and updates.
	for _, newItem := range newValue {
		idx := diskIndex(newItem)
		if idx < 0 {
			// No index: server assigns it, so this is always a new disk.
			newlyAdded = append(newlyAdded, newItem)
			continue
		}
		oldItem, exists := oldByIndex[idx]
		if !exists {
			newlyAdded = append(newlyAdded, newItem)
			continue
		}
		if !reflect.DeepEqual(oldItem, newItem) {
			updated = append(updated, newItem)
		}
	}

	// Removals: indices present in old but not in new.
	for _, oldItem := range oldValue {
		idx := diskIndex(oldItem)
		if idx < 0 {
			continue
		}
		if _, exists := newByIndex[idx]; !exists {
			removed = append(removed, oldItem)
		}
	}

	return newlyAdded, removed, updated
}

// listVolumeGroupDisks fetches the disks currently on the volume group keyed by
// index. The Volume Group GET does not return disks inline, so the ext_id
// needed for update / delete calls has to come from the dedicated disks
// endpoint.
func listVolumeGroupDisksByIndex(ctx context.Context, conn *volumesSDK.Client, vgExtID string) (map[int]volumesClient.VolumeDisk, error) {
	listVolumeDisksRequest := import1.ListVolumeDisksByVolumeGroupIdRequest{
		VolumeGroupExtId: utils.StringPtr(vgExtID),
	}
	listResp, err := conn.VolumeAPIInstance.ListVolumeDisksByVolumeGroupId(ctx, &listVolumeDisksRequest)
	if err != nil {
		return nil, err
	}

	byIndex := make(map[int]volumesClient.VolumeDisk)
	if listResp == nil || listResp.Data == nil {
		return byIndex, nil
	}
	disks, ok := listResp.GetData().([]volumesClient.VolumeDisk)
	if !ok {
		return byIndex, nil
	}
	for _, disk := range disks {
		if disk.Index != nil {
			byIndex[*disk.Index] = disk
		}
	}
	return byIndex, nil
}

func waitForVolumeDiskTask(ctx context.Context, d *schema.ResourceData, meta interface{}, taskUUID *string, timeoutType string, operation string) diag.Diagnostics {
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(timeoutType),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for volume disk (%s) to %s: %s", utils.StringValue(taskUUID), operation, errWait)
	}
	return nil
}

// applyVolumeGroupDiskAdditions creates the given disks on the volume group.
func applyVolumeGroupDiskAdditions(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *volumesSDK.Client, vgExtID string, addedDisks []interface{}) diag.Diagnostics {
	if len(addedDisks) == 0 {
		return nil
	}
	for _, disk := range expandDisks(addedDisks) {
		diskBody := disk
		createVolumeDiskRequest := import1.CreateVolumeDiskRequest{
			VolumeGroupExtId: utils.StringPtr(vgExtID),
			Body:             &diskBody,
		}

		aJSON, _ := json.MarshalIndent(createVolumeDiskRequest, "", "  ")
		log.Printf("[DEBUG] Create Volume Disk Request: %s", string(aJSON))

		resp, err := conn.VolumeAPIInstance.CreateVolumeDisk(ctx, &createVolumeDiskRequest)
		if err != nil {
			return diag.Errorf("error while creating Volume Disk : %v", err)
		}
		taskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
		if err := waitForVolumeDiskTask(ctx, d, meta, taskRef.ExtId, schema.TimeoutCreate, "create"); err != nil {
			return err
		}
	}
	return nil
}

// applyVolumeGroupDiskUpdates updates the given disks on the volume group. The
// index and disk_data_source_reference are immutable and are stripped before
// sending, matching the standalone volume disk resource.
func applyVolumeGroupDiskUpdates(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *volumesSDK.Client, vgExtID string, updatedDisks []interface{}) diag.Diagnostics {
	if len(updatedDisks) == 0 {
		return nil
	}
	currentByIndex, err := listVolumeGroupDisksByIndex(ctx, conn, vgExtID)
	if err != nil {
		return diag.Errorf("error while listing Volume Disks : %v", err)
	}

	for _, disk := range expandDisks(updatedDisks) {
		if disk.Index == nil {
			continue
		}
		currentDisk, exists := currentByIndex[*disk.Index]
		if !exists {
			return diag.Errorf("error while updating Volume Disk : no existing disk found at index %d", *disk.Index)
		}

		getResp, err := conn.VolumeAPIInstance.GetVolumeDiskById(ctx, &import1.GetVolumeDiskByIdRequest{
			VolumeGroupExtId: utils.StringPtr(vgExtID),
			ExtId:            currentDisk.ExtId,
		})
		if err != nil {
			return diag.Errorf("error while getting Volume Disk : %v", err)
		}

		eTag := conn.VolumeAPIInstance.ApiClient.GetEtag(getResp)

		updateBody := disk
		updateBody.Index = nil
		updateBody.DiskDataSourceReference = nil

		updateVolumeDiskByIdRequest := import1.UpdateVolumeDiskByIdRequest{
			VolumeGroupExtId: utils.StringPtr(vgExtID),
			ExtId:            currentDisk.ExtId,
			Body:             &updateBody,
		}

		aJSON, _ := json.MarshalIndent(updateVolumeDiskByIdRequest, "", "  ")
		log.Printf("[DEBUG] Update Volume Disk Request: %s", string(aJSON))

		args := make(map[string]interface{})
		args["If-Match"] = utils.StringPtr(eTag)

		resp, err := conn.VolumeAPIInstance.UpdateVolumeDiskById(ctx, &updateVolumeDiskByIdRequest, args)
		if err != nil {
			return diag.Errorf("error while updating Volume Disk : %v", err)
		}
		taskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
		if err := waitForVolumeDiskTask(ctx, d, meta, taskRef.ExtId, schema.TimeoutUpdate, "update"); err != nil {
			return err
		}
	}
	return nil
}

// applyVolumeGroupDiskDeletions deletes the given disks from the volume group.
func applyVolumeGroupDiskDeletions(ctx context.Context, d *schema.ResourceData, meta interface{}, conn *volumesSDK.Client, vgExtID string, deletedDisks []interface{}) diag.Diagnostics {
	if len(deletedDisks) == 0 {
		return nil
	}
	currentByIndex, err := listVolumeGroupDisksByIndex(ctx, conn, vgExtID)
	if err != nil {
		return diag.Errorf("error while listing Volume Disks : %v", err)
	}

	for _, disk := range expandDisks(deletedDisks) {
		if disk.Index == nil {
			continue
		}
		currentDisk, exists := currentByIndex[*disk.Index]
		if !exists {
			// Already gone; nothing to delete.
			continue
		}

		deleteVolumeDiskByIdRequest := import1.DeleteVolumeDiskByIdRequest{
			VolumeGroupExtId: utils.StringPtr(vgExtID),
			ExtId:            currentDisk.ExtId,
		}
		resp, err := conn.VolumeAPIInstance.DeleteVolumeDiskById(ctx, &deleteVolumeDiskByIdRequest)
		if err != nil {
			return diag.Errorf("error while deleting Volume Disk : %v", err)
		}
		taskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
		if err := waitForVolumeDiskTask(ctx, d, meta, taskRef.ExtId, schema.TimeoutDelete, "delete"); err != nil {
			return err
		}
	}
	return nil
}

func expandDisks(disks []interface{}) []volumesClient.VolumeDisk {
	if len(disks) == 0 {
		return nil
	}

	disksList := make([]volumesClient.VolumeDisk, len(disks))

	for k, v := range disks {
		disk := volumesClient.VolumeDisk{}

		diskI := v.(map[string]interface{})

		if index, ok := diskI["index"]; ok {
			disk.Index = utils.IntPtr(index.(int))
		}
		if diskSizeBytes, ok := diskI["disk_size_bytes"]; ok {
			diskSize := int64(diskSizeBytes.(int))
			disk.DiskSizeBytes = utils.Int64Ptr(diskSize)
		}
		if description, ok := diskI["description"]; ok {
			disk.Description = utils.StringPtr(description.(string))
		}
		if diskDataSourceReference, ok := diskI["disk_data_source_reference"]; ok {
			disk.DiskDataSourceReference = expandDiskDataSourceReference(diskDataSourceReference.([]interface{}))
		}
		if diskStorageFeatures, ok := diskI["disk_storage_features"]; ok {
			disk.DiskStorageFeatures = expandDiskStorageFeatures(diskStorageFeatures.([]interface{}))
		}
		disksList[k] = disk
	}
	return disksList
}
