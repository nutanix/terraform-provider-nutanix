package dataprotectionv2

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
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v16/models/dataprotection/v4/common"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v16/models/dataprotection/v4/config"
	dataprtotectionPrismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v16/models/prism/v4/config"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v16/models/prism/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRecoveryPointsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRecoveryPointsV2Create,
		ReadContext:   ResourceNutanixRecoveryPointsV2Read,
		UpdateContext: ResourceNutanixRecoveryPointsV2Update,
		DeleteContext: ResourceNutanixRecoveryPointsV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": SchemaForLinks(),
			"location_agnostic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE"}, false),
			},
			"recovery_point_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CRASH_CONSISTENT", "APPLICATION_CONSISTENT"}, false),
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_references": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vm_recovery_points": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": SchemaForLinks(),
						"consistency_group_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_agnostic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_recovery_points": SchemaForDiskRecoveryPoints(),
						"vm_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vm_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"application_consistent_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backup_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"FULL_BACKUP", "COPY_BACKUP"}, false),
									},
									"should_include_writers": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"writers": {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"should_store_vss_metadata": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"object_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"dataprotection.v4.common.VssProperties", "dataprotection.v4.r0.b1.common.VssProperties"}, false),
									},
								},
							},
						},
					},
				},
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					log.Printf("[DEBUG] DiffSuppressFunc vm_recovery_points \n")

					// Check if the list has changed
					if d.HasChange("vm_recovery_points") {
						log.Printf("[DEBUG] vm_recovery_points has changed \n")
						oldRaw, newRaw := d.GetChange("vm_recovery_points")
						// Convert to lists of interfaces
						oldList := oldRaw.([]interface{})
						newList := newRaw.([]interface{})
						// Sort lists based on a unique field (e.g., "vm_ext_id") for comparison
						sort.SliceStable(oldList, func(i, j int) bool {
							return oldList[i].(map[string]interface{})["vm_ext_id"].(string) < oldList[j].(map[string]interface{})["vm_ext_id"].(string)
						})
						sort.SliceStable(newList, func(i, j int) bool {
							return newList[i].(map[string]interface{})["vm_ext_id"].(string) < newList[j].(map[string]interface{})["vm_ext_id"].(string)
						})
						aJSON, _ := json.Marshal(oldList)
						log.Printf("[DEBUG] oldList: %v", string(aJSON))
						aJSON, _ = json.Marshal(newList)
						log.Printf("[DEBUG] newList: %v", string(aJSON))
						// Check if lists are equal when vm_ext_id is the same
						if isListEqual(oldList, newList, "vm_ext_id") {
							log.Printf("[DEBUG] vm_recovery_points are equal \n")
							return true
						}
						log.Printf("[DEBUG] vm_recovery_points are not equal \n")
						return false
					}
					log.Printf("[DEBUG] vm_recovery_points has not changed \n")
					return false
				},
			},
			"volume_group_recovery_points": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": SchemaForLinks(),
						"consistency_group_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_agnostic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_group_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"volume_group_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"disk_recovery_points": SchemaForDiskRecoveryPoints(),
					},
				},
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					log.Printf("[DEBUG] DiffSuppressFunc volume_group_recovery_points \n")

					// Check if the list has changed
					if d.HasChange("volume_group_recovery_points") {
						log.Printf("[DEBUG] volume_group_recovery_points has changed \n")
						oldRaw, newRaw := d.GetChange("volume_group_recovery_points")
						// Convert to lists of interfaces
						oldList := oldRaw.([]interface{})
						newList := newRaw.([]interface{})
						// Sort lists based on a unique field (e.g., "volume_group_ext_id") for comparison
						sort.SliceStable(oldList, func(i, j int) bool {
							return oldList[i].(map[string]interface{})["volume_group_ext_id"].(string) < oldList[j].(map[string]interface{})["volume_group_ext_id"].(string)
						})
						sort.SliceStable(newList, func(i, j int) bool {
							return newList[i].(map[string]interface{})["volume_group_ext_id"].(string) < newList[j].(map[string]interface{})["volume_group_ext_id"].(string)
						})
						aJSON, _ := json.Marshal(oldList)
						log.Printf("[DEBUG] oldList: %v", string(aJSON))
						aJSON, _ = json.Marshal(newList)
						log.Printf("[DEBUG] newList: %v", string(aJSON))
						// Check if lists are equal when volume_group_ext_id is the same
						if isListEqual(oldList, newList, "volume_group_ext_id") {
							log.Printf("[DEBUG] volume_group_recovery_points are equal \n")
							return true
						}
						log.Printf("[DEBUG] volume_group_recovery_points are not equal \n")
						return false
					}
					log.Printf("[DEBUG] volume_group_recovery_points has not changed \n")
					return false
				},
			},
		},
	}
}

// Helper function to compare two lists of maps for equality
func isListEqual(oldList, newList []interface{}, key string) bool {
	if len(oldList) != len(newList) {
		return false
	}

	for i := range oldList {
		oldItem := oldList[i].(map[string]interface{})
		newItem := newList[i].(map[string]interface{})

		// Compare all fields of the items
		if oldItem[key] != newItem[key] {
			log.Printf("[DEBUG] both lists are not equal for key: %v", key)
			return false
		}
	}
	log.Printf("[DEBUG] both lists are equal for key: %v", key)
	return true
}

func ResourceNutanixRecoveryPointsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DatasourceNutanixRecoveryPointV2Create \n")

	conn := meta.(*conns.Client).DataProtectionAPI

	body := config.RecoveryPoint{}

	if d.Get("vm_recovery_points") == nil && d.Get("volume_group_recovery_points") == nil {
		return diag.Errorf("Input is invalid because At least one vm_recovery_points or volume_group_recovery_points need to be specified.")
	}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if expirationTime, ok := d.GetOk("expiration_time"); ok {
		expTime, err := time.Parse(time.RFC3339, expirationTime.(string))
		if err != nil {
			return diag.Errorf("error while parsing expiration Time : %v", err)
		}
		body.ExpirationTime = &expTime
	}
	if status, ok := d.GetOk("status"); ok {
		statusMap := map[string]interface{}{
			"COMPLETE": 2,
		}
		pVal := statusMap[status.(string)]
		p := common.RecoveryPointStatus(pVal.(int))
		body.Status = &p
	}
	if recoveryPointType, ok := d.GetOk("recovery_point_type"); ok {
		recoveryPointTypeMap := map[string]interface{}{
			"CRASH_CONSISTENT":       2,
			"APPLICATION_CONSISTENT": 3,
		}
		pVal := recoveryPointTypeMap[recoveryPointType.(string)]
		p := common.RecoveryPointType(pVal.(int))
		body.RecoveryPointType = &p
	}
	if vmRecoveryPoints, ok := d.GetOk("vm_recovery_points"); ok {
		vmRecoveryPointsList, err := expandVMRecoveryPoints(vmRecoveryPoints.([]interface{}))
		if err != nil {
			return diag.Errorf("error while expanding vm recovery points: %v", err)
		}
		aJSON, _ := json.Marshal(vmRecoveryPointsList)
		log.Printf("[DEBUG] VM RecoveryPoint Body: %v", string(aJSON))
		body.VmRecoveryPoints = vmRecoveryPointsList
	}
	if volumeGroupRecoveryPoints, ok := d.GetOk("volume_group_recovery_points"); ok {
		body.VolumeGroupRecoveryPoints = expandVolumeGroupRecoveryPoints(volumeGroupRecoveryPoints.([]interface{}))
	}

	aJSON, _ := json.Marshal(body)
	log.Printf("[DEBUG] RecoveryPoint Body: %v", string(aJSON))

	resp, err := conn.RecoveryPoint.CreateRecoveryPoint(&body)

	if err != nil {
		return diag.Errorf("error while creating recovery point: %v", err)
	}

	TaskRef := resp.Data.GetValue().(dataprtotectionPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for recovery point: (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching recovery point UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.Marshal(rUUID)
	log.Printf("[DEBUG] Create Recovery Point Task Details: %v", string(aJSON))

	uuid := rUUID.CompletionDetails[0].Value

	d.SetId(uuid.GetValue().(string))

	return ResourceNutanixRecoveryPointsV2Read(ctx, d, meta)
}

func ResourceNutanixRecoveryPointsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DatasourceNutanixRecoveryPointV2Read \n")

	conn := meta.(*conns.Client).DataProtectionAPI

	resp, err := conn.RecoveryPoint.GetRecoveryPointById(utils.StringPtr(d.Id()))

	if err != nil {
		return diag.Errorf("error while fetching recovery point: %v", err)
	}

	getResp := resp.Data.GetValue().(config.RecoveryPoint)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location_agnostic_id", getResp.LocationAgnosticId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_time", flattenTime(getResp.CreationTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("expiration_time", flattenTime(getResp.ExpirationTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", flattenStatus(getResp.Status)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("recovery_point_type", flattenRecoveryPointType(getResp.RecoveryPointType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location_references", flattenLocationReferences(getResp.LocationReferences)); err != nil {
		return diag.FromErr(err)
	}

	// Get Vm Recovery Points from the resource
	resourceVmRecoveryPoints := d.Get("vm_recovery_points").([]interface{})
	// Get Vm Recovery Points from the response
	respRecoveryPoints := getResp.VmRecoveryPoints

	// Remove the VM Recovery Points that are present in the resource and in the response
	for _, vmRecoveryPoint := range resourceVmRecoveryPoints {
		for _, respRecoveryPoint := range getResp.VmRecoveryPoints {
			resVmRpExtId := vmRecoveryPoint.(map[string]interface{})["ext_id"]
			respVmRpExtId := utils.StringValue(respRecoveryPoint.ExtId)
			if resVmRpExtId == respVmRpExtId {
				log.Printf("[DEBUG] Removing VM Recovery Point with Ext Id: %v", respVmRpExtId)
				respRecoveryPoints = removeVmRecoveryPointByExtId(respRecoveryPoints, respRecoveryPoint)
			}
		}
	}

	// If there are any VM Recovery Points left in the response, update the resource
	if len(respRecoveryPoints) > 0 {
		if err := d.Set("vm_recovery_points", flattenVMRecoveryPoints(getResp.VmRecoveryPoints)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Get Volume Group Recovery Points from the resource
	resourceVolumeGroupRecoveryPoints := d.Get("volume_group_recovery_points").([]interface{})
	// Get Volume Group Recovery Points from the response
	respVolumeGroupRecoveryPoints := getResp.VolumeGroupRecoveryPoints

	// Remove the Volume Group Recovery Points that are present in the resource and in the response
	for _, volumeGroupRecoveryPoint := range resourceVolumeGroupRecoveryPoints {
		for _, respVolumeGroupRecoveryPoint := range getResp.VolumeGroupRecoveryPoints {
			resVolumeGroupRpExtId := volumeGroupRecoveryPoint.(map[string]interface{})["ext_id"]
			respVolumeGroupRpExtId := utils.StringValue(respVolumeGroupRecoveryPoint.ExtId)
			if resVolumeGroupRpExtId == respVolumeGroupRpExtId {
				log.Printf("[DEBUG] Removing Volume Group Recovery Point with Ext Id: %v", respVolumeGroupRpExtId)
				respVolumeGroupRecoveryPoints = removeVolumeGroupRecoveryPointByExtId(respVolumeGroupRecoveryPoints, respVolumeGroupRecoveryPoint)
			}
		}
	}

	// If there are any Volume Group Recovery Points left in the response, update the resource
	if len(respVolumeGroupRecoveryPoints) > 0 {
		if err := d.Set("volume_group_recovery_points", flattenVolumeGroupRecoveryPoints(getResp.VolumeGroupRecoveryPoints)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func ResourceNutanixRecoveryPointsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// update is supported for expiration_time only
	log.Printf("[DEBUG] DatasourceNutanixRecoveryPointV2Update \n")

	conn := meta.(*conns.Client).DataProtectionAPI

	readResp, err := conn.RecoveryPoint.GetRecoveryPointById(utils.StringPtr(d.Id()))

	if err != nil {
		return diag.Errorf("error while fetching recovery point: %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.RecoveryPoint.ApiClient.GetEtag(readResp)

	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := config.ExpirationTimeSpec{}

	if d.HasChange("expiration_time") {
		expirationTime, ok := d.GetOk("expiration_time")
		if ok {
			expTime, err := time.Parse(time.RFC3339, expirationTime.(string))
			if err != nil {
				return diag.Errorf("error while parsing expiration Time : %v", err)
			}
			body.ExpirationTime = &expTime
		}
	} else {
		return diag.Errorf("expiration_time is the only field that can be updated")
	}

	aJSON, _ := json.Marshal(body)
	log.Printf("[DEBUG] RecoveryPoint Body: %v", string(aJSON))

	resp, err := conn.RecoveryPoint.SetRecoveryPointExpirationTime(utils.StringPtr(d.Id()), &body, args)

	if err != nil {
		return diag.Errorf("error while updating recovery point: %v", err)
	}

	TaskRef := resp.Data.GetValue().(dataprtotectionPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for recovery point (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching recovery point UUID : %v", err)
	}

	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.Marshal(rUUID)
	log.Printf("[DEBUG] Update Recovery Point Task Details: %v", string(aJSON))

	return ResourceNutanixRecoveryPointsV2Read(ctx, d, meta)
}

func ResourceNutanixRecoveryPointsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataProtectionAPI

	resp, err := conn.RecoveryPoint.DeleteRecoveryPointById(utils.StringPtr(d.Id()))

	if err != nil {
		return diag.Errorf("error while deleting recovery point: %v", err)
	}

	TaskRef := resp.Data.GetValue().(dataprtotectionPrismConfig.TaskReference)

	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for recovery point (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func expandVolumeGroupRecoveryPoints(volumeGroupRecoveryPoints []interface{}) []config.VolumeGroupRecoveryPoint {
	if len(volumeGroupRecoveryPoints) == 0 {
		log.Printf("[DEBUG] volume group recovery points is Empty")
		return nil
	}
	var volumeGroupRecoveryPointsList []config.VolumeGroupRecoveryPoint
	for _, volumeGroupRecoveryPoint := range volumeGroupRecoveryPoints {
		volumeGroupRecoveryPointMap := volumeGroupRecoveryPoint.(map[string]interface{})
		volumeGroupRecoveryPointObj := config.VolumeGroupRecoveryPoint{}
		if volumeGroupExtId, ok := volumeGroupRecoveryPointMap["volume_group_ext_id"]; ok {
			volumeGroupRecoveryPointObj.VolumeGroupExtId = utils.StringPtr(volumeGroupExtId.(string))
		}
		volumeGroupRecoveryPointsList = append(volumeGroupRecoveryPointsList, volumeGroupRecoveryPointObj)
	}
	log.Printf("[DEBUG] volumeGroupRecoveryPointsList: %v", volumeGroupRecoveryPointsList)
	return volumeGroupRecoveryPointsList
}

func expandVMRecoveryPoints(vmRecoveryPoints []interface{}) ([]config.VmRecoveryPoint, error) {
	if len(vmRecoveryPoints) == 0 {
		log.Printf("[DEBUG] vm recovery points is Empty")
		return nil, nil
	}
	var vmRecoveryPointsList []config.VmRecoveryPoint
	for _, vmRecoveryPoint := range vmRecoveryPoints {
		vmRecoveryPointMap := vmRecoveryPoint.(map[string]interface{})
		vmRecoveryPointObj := config.VmRecoveryPoint{}
		if vmExtId, ok := vmRecoveryPointMap["vm_ext_id"]; ok {
			vmRecoveryPointObj.VmExtId = utils.StringPtr(vmExtId.(string))
		}
		if applicationConsistentProperties, ok := vmRecoveryPointMap["application_consistent_properties"]; ok {
			appConsistentPropList := applicationConsistentProperties.([]interface{})
			log.Printf("[DEBUG] appConsistentPropList: %v", appConsistentPropList)
			if len(appConsistentPropList) > 0 {
				appConsistentPropMap := appConsistentPropList[0].(map[string]interface{})
				log.Printf("[DEBUG] appConsistentPropMap: %v", appConsistentPropMap)
				if objectType, ok := appConsistentPropMap["object_type"]; ok {
					if objectType == ApplicationConsistentPropertiesVss1 ||
						objectType == ApplicationConsistentPropertiesVss2 {
						appConsistentPropObj, err := expandApplicationConsistentProperties(applicationConsistentProperties)
						if err != nil {
							log.Printf("[ERROR] error while expanding application consistent properties: %v", err)
							return nil, err
						}
						vmRecoveryPointObj.ApplicationConsistentProperties = appConsistentPropObj
					}
				}
			}
		}
		vmRecoveryPointsList = append(vmRecoveryPointsList, vmRecoveryPointObj)
	}
	log.Printf("[DEBUG] vmRecoveryPointsList: %v", vmRecoveryPointsList)
	return vmRecoveryPointsList, nil
}

func expandApplicationConsistentProperties(appConsistentProp interface{}) (*config.OneOfVmRecoveryPointApplicationConsistentProperties, error) {
	if appConsistentProp == nil {
		log.Printf("[DEBUG] application consistent properties is Empty")
		return nil, nil
	}
	log.Printf("[DEBUG] application consistent properties: %v", appConsistentProp)
	appConsistentPropList := appConsistentProp.([]interface{})
	appConsistentPropVal := appConsistentPropList[0].(map[string]interface{})
	oneOfVmRecoveryPointApplicationConsistentPropertiesObj := config.OneOfVmRecoveryPointApplicationConsistentProperties{}
	appConsistentPropObj := common.NewVssProperties()
	if backupType, ok := appConsistentPropVal["backup_type"]; ok {
		backupTypeMap := map[string]interface{}{
			"FULL_BACKUP": 2,
			"COPY_BACKUP": 3,
		}
		pVal := backupTypeMap[backupType.(string)]
		p := common.BackupType(pVal.(int))
		appConsistentPropObj.BackupType = &p
	}
	if shouldIncludeWriters, ok := appConsistentPropVal["should_include_writers"]; ok {
		appConsistentPropObj.ShouldIncludeWriters = utils.BoolPtr(shouldIncludeWriters.(bool))
	}
	if writers, ok := appConsistentPropVal["writers"]; ok {
		appConsistentPropObj.Writers = expandWritersList(writers.([]interface{}))
	}
	if shouldStoreVssMetadata, ok := appConsistentPropVal["should_store_vss_metadata"]; ok {
		appConsistentPropObj.ShouldStoreVssMetadata = utils.BoolPtr(shouldStoreVssMetadata.(bool))
	}
	if objectType, ok := appConsistentPropVal["object_type"]; ok {
		appConsistentPropObj.ObjectType_ = utils.StringPtr(objectType.(string))
	}
	err := oneOfVmRecoveryPointApplicationConsistentPropertiesObj.SetValue(*appConsistentPropObj)
	if err != nil {
		log.Printf("[ERROR] error while setting value for OneOfVmRecoveryPointApplicationConsistentProperties: %v", err)
		return nil, err
	}
	return &oneOfVmRecoveryPointApplicationConsistentPropertiesObj, nil
}

func expandWritersList(writers []interface{}) []string {
	if len(writers) > 0 {
		writersList := make([]string, len(writers))

		for k, v := range writers {
			writersList[k] = v.(string)
		}
		return writersList
	}
	return nil
}

func taskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)

		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(prismConfig.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(pr *prismConfig.TaskStatus) string {
	if pr != nil {
		if *pr == prismConfig.TaskStatus(6) {
			return "FAILED"
		}
		if *pr == prismConfig.TaskStatus(7) {
			return "CANCELED"
		}
		if *pr == prismConfig.TaskStatus(2) {
			return "QUEUED"
		}
		if *pr == prismConfig.TaskStatus(3) {
			return "RUNNING"
		}
		if *pr == prismConfig.TaskStatus(5) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}

// Function to remove a Vm recovery Point with a specific Ext Id from the slice
func removeVmRecoveryPointByExtId(recoveryPoints []config.VmRecoveryPoint, recoveryPoint config.VmRecoveryPoint) []config.VmRecoveryPoint {
	var result []config.VmRecoveryPoint // Create a new slice to hold the result

	for _, rp := range recoveryPoints {
		if utils.StringValue(rp.ExtId) != utils.StringValue(recoveryPoint.ExtId) {
			result = append(result, rp) // Add recovery point to result if the ext id doesn't match
		}
	}
	return result
}

// Function to remove a Volume Group recovery Point with a specific Ext Id from the slice
func removeVolumeGroupRecoveryPointByExtId(recoveryPoints []config.VolumeGroupRecoveryPoint, recoveryPoint config.VolumeGroupRecoveryPoint) []config.VolumeGroupRecoveryPoint {
	var result []config.VolumeGroupRecoveryPoint // Create a new slice to hold the result

	for _, rp := range recoveryPoints {
		if utils.StringValue(rp.ExtId) != utils.StringValue(recoveryPoint.ExtId) {
			result = append(result, rp) // Add recovery point to result if the ext id doesn't match
		}
	}
	return result
}
