package volumesv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/config"
	volumesPrism "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/prism/v4/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Creates a new Volume Disk.
func ResourceNutanixVolumeGroupDiskV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Creates a new Volume Disk.",
		CreateContext: ResourceNutanixVolumeGroupDiskV2Create,
		ReadContext:   ResourceNutanixVolumeGroupDiskV2Read,
		UpdateContext: ResourceNutanixVolumeGroupDiskV2Update,
		DeleteContext: ResourceNutanixVolumeGroupDiskV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				const expectedPartsCount = 2
				parts := strings.Split(d.Id(), "/")
				if len(parts) != expectedPartsCount {
					return nil, fmt.Errorf("invalid import uuid (%q), expected volume_group_ext_id/disk_ext_id", d.Id())
				}
				d.Set("volume_group_ext_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"volume_group_ext_id": {
				Description: "The external identifier of the volume group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ext_id": {
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"index": {
				Description: "Index of the disk in a Volume Group. This field is optional and immutable.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"disk_size_bytes": {
				Description: "Size of the disk in bytes. This field is mandatory during Volume Group creation if a new disk is being created on the storage container.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"description": {
				Description: "Volume Disk description. This is an optional field.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"disk_data_source_reference": {
				Description: "Disk Data Source Reference.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Description: "The external identifier of the Data Source Reference.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "The name of the Data Source Reference.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"uris": {
							Description: "The uri list of the Data Source Reference.",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"entity_type": {
							Description:  "The Entity Type of the Data Source Reference.",
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"STORAGE_CONTAINER", "VM_DISK", "VOLUME_DISK", "DISK_RECOVERY_POINT"}, false),
						},
					},
				},
			},
			"disk_storage_features": {
				Description: "Storage optimization features which must be enabled on the Volume Disks. This is an optional field. If omitted, the disks will honor the Volume Group specific storage features setting.",
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
										Description: "The flash mode is enabled or not.",
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
		},
	}
}

func ResourceNutanixVolumeGroupDiskV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id")

	body := volumesClient.VolumeDisk{}

	if diskSizeBytes, ok := d.GetOk("disk_size_bytes"); ok {
		diskSize := int64(diskSizeBytes.(int))
		body.DiskSizeBytes = utils.Int64Ptr(diskSize)
	}
	if index, ok := d.GetOk("index"); ok {
		body.Index = utils.IntPtr(index.(int))
	}
	if description, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(description.(string))
	}
	if diskDataSourceReference, ok := d.GetOk("disk_data_source_reference"); ok {
		body.DiskDataSourceReference = expandDiskDataSourceReference(diskDataSourceReference)
	}
	if diskStorageFeatures, ok := d.GetOk("disk_storage_features"); ok {
		body.DiskStorageFeatures = expandDiskStorageFeatures(diskStorageFeatures.([]interface{}))
	}

	log.Printf("[DEBUG] Volume Disk Body body.DiskDataSourceReference.Uris : %v", body.DiskDataSourceReference.Uris)
	resp, err := conn.VolumeAPIInstance.CreateVolumeDisk(utils.StringPtr(volumeGroupExtID.(string)), &body)
	if err != nil {
		return diag.Errorf("error while creating Volume Disk : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the volume disk to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for volume disk (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching volume disk task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Volume Disk Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVolumeGroupDisk, "Volume disk")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixVolumeGroupDiskV2Read(ctx, d, meta)
}

func ResourceNutanixVolumeGroupDiskV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id")

	volumeDiskExtID := d.Id() // d.Id gives volume_group_ext_id not volume_disk_ext_id

	resp, err := conn.VolumeAPIInstance.GetVolumeDiskById(utils.StringPtr(volumeGroupExtID.(string)), utils.StringPtr(volumeDiskExtID))
	if err != nil {
		return diag.Errorf("error while fetching volume Disk : %v", err)
	}
	getResp := resp.Data.GetValue().(volumesClient.VolumeDisk)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("index", getResp.Index); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_size_bytes", getResp.DiskSizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_data_source_reference", flattenDiskDataSourceReference(getResp.DiskDataSourceReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_storage_features", flattenDiskStorageFeatures(getResp.DiskStorageFeatures)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixVolumeGroupDiskV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id")
	volumeDiskExtID := d.Id()

	resp, err := conn.VolumeAPIInstance.GetVolumeDiskById(utils.StringPtr(volumeGroupExtID.(string)), utils.StringPtr(volumeDiskExtID))
	if err != nil {
		return diag.Errorf("error while updating Volume Disk : %v", err)
	}
	updateSpec := resp.Data.GetValue().(volumesClient.VolumeDisk)

	if d.HasChange("index") {
		index := d.Get("index").(int)
		updateSpec.Index = &index
	} else {
		updateSpec.Index = nil
	}
	if d.HasChange("disk_size_bytes") {
		diskSizeBytes := int64(d.Get("disk_size_bytes").(int))
		updateSpec.DiskSizeBytes = &diskSizeBytes
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateSpec.Description = &description
	}
	if d.HasChange("disk_storage_features") {
		diskStorageFeatures := d.Get("disk_storage_features").([]interface{})
		updateSpec.DiskStorageFeatures = expandDiskStorageFeatures(diskStorageFeatures)
	}
	if d.HasChange("disk_data_source_reference") {
		diskDataSourceReference := d.Get("disk_data_source_reference").([]interface{})
		updateSpec.DiskDataSourceReference = expandDiskDataSourceReference(diskDataSourceReference)
	} else {
		updateSpec.DiskDataSourceReference = nil
	}

	updateResp, err := conn.VolumeAPIInstance.UpdateVolumeDiskById(utils.StringPtr(volumeGroupExtID.(string)), utils.StringPtr(volumeDiskExtID), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating Volume Disk : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the volume disk to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for volume disk (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching volume disk update task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update Volume Disk Task Details: %s", string(aJSON))

	return ResourceNutanixVolumeGroupDiskV2Read(ctx, d, meta)
}

func ResourceNutanixVolumeGroupDiskV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	volumeGroupExtID := d.Get("volume_group_ext_id")
	volumeDiskExtID := d.Get("ext_id")

	resp, err := conn.VolumeAPIInstance.DeleteVolumeDiskById(utils.StringPtr(volumeGroupExtID.(string)), utils.StringPtr(volumeDiskExtID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching volume Disk : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the volume disk to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for volume disk (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details for logging
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching volume disk delete task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Volume Disk Task Details: %s", string(aJSON))

	return nil
}

func expandDiskStorageFeatures(diskStorageFeatures []interface{}) *volumesClient.DiskStorageFeatures {
	if len(diskStorageFeatures) > 0 {
		diskStorageFeature := volumesClient.DiskStorageFeatures{}

		val := diskStorageFeatures[0].(map[string]interface{})

		if flashMode, ok := val["flash_mode"]; ok {
			diskStorageFeature.FlashMode = expandFlashMode(flashMode.([]interface{}))
		}
		return &diskStorageFeature
	}
	return nil
}

func expandDiskDataSourceReference(entityReference interface{}) *config.EntityReference {
	if entityReference != nil {
		entityReferenceI := entityReference.([]interface{})
		val := entityReferenceI[0].(map[string]interface{})

		diskDataSourceReference := config.EntityReference{}

		if extID, ok := val["ext_id"]; ok {
			diskDataSourceReference.ExtId = utils.StringPtr(extID.(string))
		}
		if name, ok := val["name"]; ok {
			diskDataSourceReference.Name = utils.StringPtr(name.(string))
		}
		if uris, ok := val["uris"]; ok {
			uriSlice := make([]*string, len(uris.([]interface{})))
			for i, uri := range uris.([]interface{}) {
				uriSlice[i] = utils.StringPtr(uri.(string))
			}
			diskDataSourceReference.Uris = utils.StringValueSlice(uriSlice)
		}
		if entityType, ok := val["entity_type"]; ok {
			const zero, one, four, twenty, twentyone, twentytwo = 0, 1, 4, 20, 21, 22
			subMap := map[string]interface{}{
				"UNKNOWN":             zero,
				"REDACTED":            one,
				"STORAGE_CONTAINER":   four,
				"VM_DISK":             twenty,
				"VOLUME_DISK":         twentyone,
				"DISK_RECOVERY_POINT": twentytwo,
			}

			pInt := subMap[entityType.(string)]
			p := config.EntityType(pInt.(int))

			diskDataSourceReference.EntityType = &p
		}
		log.Printf("[DEBUG] Disk Data Source Reference : %v", diskDataSourceReference)
		return &diskDataSourceReference
	}
	return nil
}
