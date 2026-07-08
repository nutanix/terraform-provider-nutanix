package volumesv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	volumesPrism "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/prism/v4/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
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
		const two, three = 2, 3
		sharingStatusMap := map[string]interface{}{
			"SHARED":     two,
			"NOT_SHARED": three,
		}
		pVal := sharingStatusMap[sharingStatus.(string)]
		p := volumesClient.SharingStatus(pVal.(int))
		body.SharingStatus = &p
	}
	if targetPrefix, ok := d.GetOk("target_prefix"); ok {
		body.TargetPrefix = utils.StringPtr(targetPrefix.(string))
	}
	if targetName, ok := d.GetOk("target_name"); ok {
		body.TargetName = utils.StringPtr(targetName.(string))
	}
	if enabledAuthentications, ok := d.GetOk("enabled_authentications"); ok {
		const CHAP, NONE = 2, 3
		enabledAuthenticationsMap := map[string]interface{}{
			"CHAP": CHAP,
			"NONE": NONE,
		}
		pVal := enabledAuthenticationsMap[enabledAuthentications.(string)]
		if pVal == nil {
			body.EnabledAuthentications = nil
		} else {
			p := volumesClient.AuthenticationType(pVal.(int))
			body.EnabledAuthentications = &p
		}
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
		const two, three, four, five = 2, 3, 4, 5
		usageTypeMap := map[string]interface{}{
			"USER":          two,
			"INTERNAL":      three,
			"TEMPORARY":     four,
			"BACKUP_TARGET": five,
		}
		pInt := usageTypeMap[usageType.(string)]
		p := volumesClient.UsageType(pInt.(int))
		body.UsageType = &p
	}
	if attachmentType, ok := d.GetOk("attachment_type"); ok {
		const NONE, DIRECT, EXTERNAL = 2, 3, 4
		attachmentTypeMap := map[string]interface{}{
			"NONE":     NONE,
			"DIRECT":   DIRECT,
			"EXTERNAL": EXTERNAL,
		}
		pInt := attachmentTypeMap[attachmentType.(string)]
		if pInt == nil {
			body.AttachmentType = nil
		} else {
			p := volumesClient.AttachmentType(pInt.(int))
			body.AttachmentType = &p
		}
	}
	if protocol, ok := d.GetOk("protocol"); ok {
		const NotAssigned, ISCSI, NVMF = 2, 3, 4
		protocolMap := map[string]interface{}{
			"NotAssigned": NotAssigned,
			"ISCSI":       ISCSI,
			"NVMF":        NVMF,
		}
		pInt := protocolMap[protocol.(string)]
		if pInt == nil {
			body.Protocol = nil
		} else {
			p := volumesClient.Protocol(pInt.(int))
			body.Protocol = &p
		}
	}
	if isHidden, ok := d.GetOk("is_hidden"); ok {
		body.IsHidden = utils.BoolPtr(isHidden.(bool))
	}
	if disks, ok := d.GetOk("disks"); ok {
		body.Disks = expandDisks(disks.([]interface{}))
	}
	resp, err := conn.VolumeAPIInstance.CreateVolumeGroup(&body)
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
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching volume group task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
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

	resp, err := conn.VolumeAPIInstance.GetVolumeGroupById(utils.StringPtr(d.Id()), nil)
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

	// Disks are not part of the Volume Group GET payload (unless expanded) and are
	// managed as a sub-collection. Refresh the size of the disks already tracked in
	// state so that a resize applied here (or done out-of-band) is detected on the
	// next plan. This is best-effort: a failure to list disks must not fail the read.
	if err := refreshVolumeGroupDisksInState(d, meta); err != nil {
		log.Printf("[WARN] unable to refresh volume group disks for %s: %v", d.Id(), err)
	}

	return nil
}

// refreshVolumeGroupDisksInState overlays the API-authoritative disk fields
// (disk_size_bytes, index, description) onto the disks currently tracked in
// state. The config-only disk_data_source_reference and disk_storage_features
// blocks are left untouched to avoid a perpetual diff, since the Volume Group
// disk GET does not echo the data source reference back.
func refreshVolumeGroupDisksInState(d *schema.ResourceData, meta interface{}) error {
	stateDisksRaw, ok := d.GetOk("disks")
	if !ok {
		return nil
	}
	stateDisks := stateDisksRaw.([]interface{})
	if len(stateDisks) == 0 {
		return nil
	}

	apiDisks, err := listVolumeGroupDisksSortedByIndex(meta, d.Id())
	if err != nil {
		return err
	}

	overlap := len(stateDisks)
	if len(apiDisks) < overlap {
		overlap = len(apiDisks)
	}
	for i := 0; i < overlap; i++ {
		diskMap, ok := stateDisks[i].(map[string]interface{})
		if !ok {
			continue
		}
		if apiDisks[i].DiskSizeBytes != nil {
			diskMap["disk_size_bytes"] = int(utils.Int64Value(apiDisks[i].DiskSizeBytes))
		}
		if apiDisks[i].Index != nil {
			diskMap["index"] = utils.IntValue(apiDisks[i].Index)
		}
		if apiDisks[i].Description != nil {
			diskMap["description"] = utils.StringValue(apiDisks[i].Description)
		}
		stateDisks[i] = diskMap
	}

	return d.Set("disks", stateDisks)
}

func ResourceNutanixVolumeGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI
	vgExtID := d.Id()

	// Top-level Volume Group properties (name, description, sharing, iSCSI, load
	// balancing) are updated through UpdateVolumeGroupById. Disks are NOT part of
	// this payload in the v4 API and must be managed through the /disks
	// sub-collection endpoints (see updateVolumeGroupDisks).
	topLevelFields := []string{
		"name",
		"description",
		"should_load_balance_vm_attachments",
		"sharing_status",
		"iscsi_features",
	}
	topLevelChanged := false
	for _, f := range topLevelFields {
		if d.HasChange(f) {
			topLevelChanged = true
			break
		}
	}

	if topLevelChanged {
		readResp, err := conn.VolumeAPIInstance.GetVolumeGroupById(utils.StringPtr(vgExtID), nil)
		if err != nil {
			return diag.Errorf("error while fetching Volume Group for update : %v", err)
		}
		updateSpec := readResp.Data.GetValue().(volumesClient.VolumeGroup)

		// Extract E-Tag header for the If-Match precondition.
		etagValue := conn.VolumeAPIInstance.ApiClient.GetEtag(readResp)
		args := make(map[string]interface{})
		args["If-Match"] = utils.StringPtr(etagValue)

		if d.HasChange("name") {
			updateSpec.Name = utils.StringPtr(d.Get("name").(string))
		}
		if d.HasChange("description") {
			updateSpec.Description = utils.StringPtr(d.Get("description").(string))
		}
		if d.HasChange("should_load_balance_vm_attachments") {
			updateSpec.ShouldLoadBalanceVmAttachments = utils.BoolPtr(d.Get("should_load_balance_vm_attachments").(bool))
		}
		if d.HasChange("sharing_status") {
			updateSpec.SharingStatus = common.ExpandEnum[volumesClient.SharingStatus](d.Get("sharing_status"))
		}
		if d.HasChange("iscsi_features") {
			updateSpec.IscsiFeatures = expandIscsiFeatures(d.Get("iscsi_features").([]interface{}))
		}

		updateResp, err := conn.VolumeAPIInstance.UpdateVolumeGroupById(utils.StringPtr(vgExtID), &updateSpec, args)
		if err != nil {
			return diag.Errorf("error while updating Volume Group : %v", err)
		}

		taskRef := updateResp.Data.GetValue().(volumesPrism.TaskReference)
		if err := waitForVolumeGroupTask(ctx, meta, taskRef.ExtId, d.Timeout(schema.TimeoutUpdate), "update"); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("disks") {
		if err := updateVolumeGroupDisks(ctx, d, meta, vgExtID); err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceNutanixVolumeGroupV2Read(ctx, d, meta)
}

// listVolumeGroupDisksSortedByIndex returns all disks attached to the Volume
// Group sorted by their index. The v4 list endpoint only supports ordering by
// diskSizeBytes, so ordering by index is done client-side to keep the positional
// matching used elsewhere deterministic.
func listVolumeGroupDisksSortedByIndex(meta interface{}, vgExtID string) ([]volumesClient.VolumeDisk, error) {
	conn := meta.(*conns.Client).VolumeAPI

	listResp, err := conn.VolumeAPIInstance.ListVolumeDisksByVolumeGroupId(utils.StringPtr(vgExtID), nil, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	disks := make([]volumesClient.VolumeDisk, 0)
	if listResp != nil && listResp.Data != nil {
		disks = listResp.Data.GetValue().([]volumesClient.VolumeDisk)
	}

	sort.SliceStable(disks, func(i, j int) bool {
		return utils.IntValue(disks[i].Index) < utils.IntValue(disks[j].Index)
	})
	return disks, nil
}

// updateVolumeGroupDisks reconciles the configured disks list against the disks
// currently attached to the Volume Group. Disks are matched positionally (the
// same order in which they were created), which mirrors how the create path
// sends the list. It resizes existing disks, creates newly added disks and
// deletes removed ones. Disk size can only be increased; a shrink is rejected
// before any API call is made, matching the constraint enforced by the fabric.
func updateVolumeGroupDisks(ctx context.Context, d *schema.ResourceData, meta interface{}, vgExtID string) error {
	conn := meta.(*conns.Client).VolumeAPI

	existing, err := listVolumeGroupDisksSortedByIndex(meta, vgExtID)
	if err != nil {
		return fmt.Errorf("error while listing Volume Group disks for update: %v", err)
	}

	desired := expandDisks(d.Get("disks").([]interface{}))

	// Resize existing disks that overlap with the desired list (positional match).
	overlap := len(existing)
	if len(desired) < overlap {
		overlap = len(desired)
	}
	for i := 0; i < overlap; i++ {
		cur := existing[i]
		want := desired[i]
		if want.DiskSizeBytes == nil || cur.DiskSizeBytes == nil {
			continue
		}
		if *want.DiskSizeBytes == *cur.DiskSizeBytes {
			continue
		}
		if *want.DiskSizeBytes < *cur.DiskSizeBytes {
			return fmt.Errorf(
				"disk size can only be increased: disk at index %d has current size %d bytes, requested %d bytes",
				i, *cur.DiskSizeBytes, *want.DiskSizeBytes,
			)
		}
		if cur.ExtId == nil {
			return fmt.Errorf("cannot resize disk at index %d: missing disk external id in API response", i)
		}

		diskExtID := utils.StringValue(cur.ExtId)
		diskReadResp, err := conn.VolumeAPIInstance.GetVolumeDiskById(utils.StringPtr(vgExtID), utils.StringPtr(diskExtID))
		if err != nil {
			return fmt.Errorf("error while fetching Volume Group disk %s for resize: %v", diskExtID, err)
		}
		diskSpec := diskReadResp.Data.GetValue().(volumesClient.VolumeDisk)

		etagValue := conn.VolumeAPIInstance.ApiClient.GetEtag(diskReadResp)
		diskArgs := make(map[string]interface{})
		diskArgs["If-Match"] = utils.StringPtr(etagValue)

		diskSpec.DiskSizeBytes = want.DiskSizeBytes
		// Index is immutable and DiskDataSourceReference must not be sent on an
		// update; leave them out of the payload.
		diskSpec.Index = nil
		diskSpec.DiskDataSourceReference = nil

		updateResp, err := conn.VolumeAPIInstance.UpdateVolumeDiskById(utils.StringPtr(vgExtID), utils.StringPtr(diskExtID), &diskSpec, diskArgs)
		if err != nil {
			return fmt.Errorf("error while resizing Volume Group disk %s: %v", diskExtID, err)
		}
		taskRef := updateResp.Data.GetValue().(volumesPrism.TaskReference)
		if err := waitForVolumeGroupTask(ctx, meta, taskRef.ExtId, d.Timeout(schema.TimeoutUpdate), "disk resize"); err != nil {
			return err
		}
	}

	// Create disks that were added to the configuration.
	for i := len(existing); i < len(desired); i++ {
		newDisk := desired[i]
		createResp, err := conn.VolumeAPIInstance.CreateVolumeDisk(utils.StringPtr(vgExtID), &newDisk)
		if err != nil {
			return fmt.Errorf("error while adding Volume Group disk at index %d: %v", i, err)
		}
		taskRef := createResp.Data.GetValue().(volumesPrism.TaskReference)
		if err := waitForVolumeGroupTask(ctx, meta, taskRef.ExtId, d.Timeout(schema.TimeoutUpdate), "disk add"); err != nil {
			return err
		}
	}

	// Delete disks that were removed from the configuration.
	for i := len(desired); i < len(existing); i++ {
		cur := existing[i]
		if cur.ExtId == nil {
			continue
		}
		diskExtID := utils.StringValue(cur.ExtId)
		deleteResp, err := conn.VolumeAPIInstance.DeleteVolumeDiskById(utils.StringPtr(vgExtID), utils.StringPtr(diskExtID))
		if err != nil {
			return fmt.Errorf("error while removing Volume Group disk %s: %v", diskExtID, err)
		}
		taskRef := deleteResp.Data.GetValue().(volumesPrism.TaskReference)
		if err := waitForVolumeGroupTask(ctx, meta, taskRef.ExtId, d.Timeout(schema.TimeoutUpdate), "disk remove"); err != nil {
			return err
		}
	}

	return nil
}

// waitForVolumeGroupTask blocks until the given Prism task reaches a terminal
// state and logs the resulting task details.
func waitForVolumeGroupTask(ctx context.Context, meta interface{}, taskUUID *string, timeout time.Duration, action string) error {
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: timeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return fmt.Errorf("error waiting for volume group %s task (%s): %s", action, utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return fmt.Errorf("error while fetching volume group %s task (%s): %v", action, utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Volume Group %s Task Details: %s", action, string(aJSON))
	return nil
}

func ResourceNutanixVolumeGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	resp, err := conn.VolumeAPIInstance.DeleteVolumeGroupById(utils.StringPtr(d.Id()))
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
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
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
			const two, three = 2, 3
			enabledAuthenticationsMap := map[string]interface{}{
				"CHAP": two,
				"NONE": three,
			}
			pVal := enabledAuthenticationsMap[enabledAuthentications.(string)]
			p := volumesClient.AuthenticationType(pVal.(int))
			iscsiFeature.EnabledAuthentications = &p
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
