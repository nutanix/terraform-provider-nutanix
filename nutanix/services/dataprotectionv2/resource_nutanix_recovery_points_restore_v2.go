package dataprotectionv2

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	dataprtotectionPrismConfig "github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRecoveryPointRestoreV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRecoveryPointRestoreV2Create,
		ReadContext:   ResourceNutanixRecoveryPointRestoreV2Read,
		UpdateContext: ResourceNutanixRecoveryPointRestoreV2Update,
		DeleteContext: ResourceNutanixRecoveryPointRestoreV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_recovery_point_restore_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_recovery_point_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"volume_group_recovery_point_restore_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_group_recovery_point_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"volume_group_override_spec": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"vm_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"volume_group_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// ResourceNutanixRecoveryPointRestoreV2Create to Restore Recovery Point
func ResourceNutanixRecoveryPointRestoreV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixRecoveryPointRestoreV2Create \n")

	conn := meta.(*conns.Client).DataProtectionAPI

	body := config.RecoveryPointRestorationSpec{}
	rpExtID := d.Get("ext_id").(string)

	if clusterExtID, ok := d.GetOk("cluster_ext_id"); ok {
		body.ClusterExtId = utils.StringPtr(clusterExtID.(string))
	}
	if vmRecoveryPointRestoreOverrides, ok := d.GetOk("vm_recovery_point_restore_overrides"); ok {
		body.VmRecoveryPointRestoreOverrides = expandVMRecoveryPointRestoreOverrides(vmRecoveryPointRestoreOverrides)
	}
	if volumeGroupRecoveryPointRestoreOverrides, ok := d.GetOk("volume_group_recovery_point_restore_overrides"); ok {
		body.VolumeGroupRecoveryPointRestoreOverrides = expandVolumeGroupRecoveryPointRestoreOverrides(volumeGroupRecoveryPointRestoreOverrides)
	}

	resp, err := conn.RecoveryPoint.RestoreRecoveryPoint(utils.StringPtr(rpExtID), &body)
	if err != nil {
		return diag.Errorf("error while replicating recovery point: %v", err)
	}

	taskRef := resp.Data.GetValue().(dataprtotectionPrismConfig.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the recovery point restore operation to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for recovery point restore (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching restore recovery point task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Restore Recovery Point Task Details: %s", string(aJSON))

	// Extract VM and Volume Group ExtIds from completion details
	vmExtIds := make([]string, 0)
	vgExtIds := make([]string, 0)
	vmExtIdsStrs := common.ExtractCompletionDetailsFromTask(taskDetails, utils.CompletionDetailsNameVMExtIDs)
	for _, vmExtIdsStr := range vmExtIdsStrs {
		vmExtIds = append(vmExtIds, strings.Split(vmExtIdsStr, ",")...)
	}
	vgExtIdsStrs := common.ExtractCompletionDetailsFromTask(taskDetails, utils.CompletionDetailsNameVGExtIDs)
	for _, vgExtIdsStr := range vgExtIdsStrs {
		vgExtIds = append(vgExtIds, strings.Split(vgExtIdsStr, ",")...)
	}
	d.Set("vm_ext_ids", vmExtIds)
	d.Set("volume_group_ext_ids", vgExtIds)
	d.SetId(utils.StringValue(taskDetails.ExtId))

	return nil
}

func expandVolumeGroupRecoveryPointRestoreOverrides(vgRecoveryPoints interface{}) []config.VolumeGroupRecoveryPointRestoreOverride {
	volumeGroupRecoveryPointRestoreOverrides := make([]config.VolumeGroupRecoveryPointRestoreOverride, 0)
	for _, vgRecoveryPoint := range vgRecoveryPoints.([]interface{}) {
		vgRecoveryPoint := vgRecoveryPoint.(map[string]interface{})
		volumeGroupRecoveryPointRestoreOverride := config.VolumeGroupRecoveryPointRestoreOverride{
			VolumeGroupRecoveryPointExtId: utils.StringPtr(vgRecoveryPoint["volume_group_recovery_point_ext_id"].(string)),
		}
		if volumeGroupOverrideSpec, ok := vgRecoveryPoint["volume_group_override_spec"]; ok {
			volumeGroupRecoveryPointRestoreOverride.VolumeGroupOverrideSpec = expandVolumeGroupOverrideSpec(volumeGroupOverrideSpec.([]interface{}))
		}
		volumeGroupRecoveryPointRestoreOverrides = append(volumeGroupRecoveryPointRestoreOverrides, volumeGroupRecoveryPointRestoreOverride)
	}
	return volumeGroupRecoveryPointRestoreOverrides
}

func expandVolumeGroupOverrideSpec(volumeGroupSpec []interface{}) *config.VolumeGroupOverrideSpec {
	volumeGroupOverrideSpec := config.VolumeGroupOverrideSpec{}
	if len(volumeGroupSpec) > 0 {
		volumeGroupOverrideSpec.Name = utils.StringPtr(volumeGroupSpec[0].(map[string]interface{})["name"].(string))
	}
	return &volumeGroupOverrideSpec
}

func expandVMRecoveryPointRestoreOverrides(vmRecoveryPoints interface{}) []config.VmRecoveryPointRestoreOverride {
	vmRecoveryPointRestoreOverrides := make([]config.VmRecoveryPointRestoreOverride, 0)
	for _, vmRecoveryPoint := range vmRecoveryPoints.([]interface{}) {
		vmRecoveryPoint := vmRecoveryPoint.(map[string]interface{})
		vmRecoveryPointRestoreOverride := config.VmRecoveryPointRestoreOverride{
			VmRecoveryPointExtId: utils.StringPtr(vmRecoveryPoint["vm_recovery_point_ext_id"].(string)),
		}
		vmRecoveryPointRestoreOverrides = append(vmRecoveryPointRestoreOverrides, vmRecoveryPointRestoreOverride)
	}
	return vmRecoveryPointRestoreOverrides
}

func ResourceNutanixRecoveryPointRestoreV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRecoveryPointRestoreV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixRecoveryPointRestoreV2Read(ctx, d, meta)
}

func ResourceNutanixRecoveryPointRestoreV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
