package vmmv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmPrism "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	vmmConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixNGTUpgradeV2 TF schema for NGT install/uninstall
func ResourceNutanixNGTUpgradeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixNGTUpgradeV2Create,
		ReadContext:   ResourceNutanixNGTUpgradeV2Read,
		UpdateContext: ResourceNutanixNGTUpgradeV2Update,
		DeleteContext: ResourceNutanixNGTUpgradeV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reboot_preference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SKIP", "IMMEDIATE", "LATER"}, false),
						},
						"schedule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_iso_inserted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"capablities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"available_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guest_os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_reachable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_vss_snapshot_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_vm_mobility_drivers_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// ResourceNutanixNGTUpgradeV2Create to Upgrade NGT on Vm
func ResourceNutanixNGTUpgradeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	readResp, err := conn.VMAPIInstance.GetGuestToolsById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching Vm : %v", err)
	}
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	body := &vmmConfig.GuestToolsUpgradeConfig{}

	if rebootPreference, ok := d.GetOk("reboot_preference"); ok {
		if len(rebootPreference.([]interface{})) > 0 {
			rp := rebootPreference.([]interface{})[0].(map[string]interface{})
			const two, three, four = 2, 3, 4
			scheduleTypesMap := map[string]int{
				"SKIP":      two,
				"IMMEDIATE": three,
				"LATER":     four,
			}
			body.RebootPreference = &vmmConfig.RebootPreference{
				ScheduleType: (*vmmConfig.ScheduleType)(utils.IntPtr(scheduleTypesMap[(rp["schedule_type"].(string))])),
			}
			if scheduleType, ok := rp["schedule_type"].(string); ok && scheduleType == "LATER" {
				if schedule, ok := rp["schedule"]; ok {
					s := schedule.([]interface{})[0].(map[string]interface{})
					t, errTime := time.Parse(time.RFC3339, s["start_time"].(string))
					if errTime != nil {
						return diag.Errorf("error while Upgrading gest tools : %v", errTime)
					}
					body.RebootPreference.Schedule = &vmmConfig.RebootPreferenceSchedule{
						StartTime: utils.Time(t),
					}
				}
			}
		}
	}

	resp, err := conn.VMAPIInstance.UpgradeVmGuestTools(utils.StringPtr(extID.(string)), body, args)
	if err != nil {
		return diag.Errorf("error while Upgrading gest tools  : %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the NGT upgrade to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for NGT upgrade (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching NGT upgrade task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", " ")
	log.Printf("[DEBUG] NGT Upgrade Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVM, "VM")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixNGTUpgradeV2Read(ctx, d, meta)
}

// Read NGT Configuration
func ResourceNutanixNGTUpgradeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Id()
	resp, err := conn.VMAPIInstance.GetGuestToolsById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching Gest Tool : %v", err)
	}
	getResp := resp.Data.GetValue().(vmmConfig.GuestTools)

	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_installed", getResp.IsInstalled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", getResp.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_iso_inserted", getResp.IsIsoInserted); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("capablities", flattenCapabilities(getResp.Capabilities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("available_version", getResp.AvailableVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_os_version", getResp.GuestOsVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_reachable", getResp.IsReachable); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_vss_snapshot_capable", getResp.IsVssSnapshotCapable); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_vm_mobility_drivers_installed", getResp.IsVmMobilityDriversInstalled); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// ResourceNutanixNGTUpgradeV2Update  Not supported
func ResourceNutanixNGTUpgradeV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixNGTUpgradeV2Delete  Not supported
func ResourceNutanixNGTUpgradeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
