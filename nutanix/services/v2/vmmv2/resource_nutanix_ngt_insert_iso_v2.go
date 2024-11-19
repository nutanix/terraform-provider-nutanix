package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v16/models/prism/v4/config"
	vmmPrism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"
	vmmConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixNGTInsertIsoV4 TF schema for NGT install/uninstall
func ResourceNutanixNGTInsertIsoV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixNGTInsertIsoV4Create,
		ReadContext:   ResourceNutanixNGTInsertIsoV4Read,
		UpdateContext: ResourceNutanixNGTInsertIsoV4Update,
		DeleteContext: ResourceNutanixNGTInsertIsoV4Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"capablities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"SELF_SERVICE_RESTORE", "VSS_SNAPSHOT"}, false),
				},
			},
			"is_config_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

// Install NGT on Vm
func ResourceNutanixNGTInsertIsoV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	readResp, err := conn.VMAPIInstance.GetGuestToolsById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching Vm : %v", err)
	}
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	body := &vmmConfig.GuestToolsInsertConfig{}

	// prepare the body
	if capabilities, ok := d.GetOk("capablities"); ok {
		capabilitiesList := capabilities.([]interface{})
		capabilitiesSet := make(map[vmmConfig.NgtCapability]bool)
		// capabilitiesListStr := make([]string, len(capabilitiesList))
		for _, v := range capabilitiesList {
			var cap vmmConfig.NgtCapability
			if v.(string) == "SELF_SERVICE_RESTORE" {
				cap = 2 // Assuming 2 represents SELF_SERVICE_RESTORE
			} else if v.(string) == "VSS_SNAPSHOT" {
				cap = 3 // Assuming 3 represents VSS_SNAPSHOT
			}
			// Step 3: Add capability to the set
			capabilitiesSet[cap] = true
		}
		// Convert the set back to a slice for the API call
		capabilitiesListStr := make([]vmmConfig.NgtCapability, 0, len(capabilitiesSet))
		for cap := range capabilitiesSet {
			capabilitiesListStr = append(capabilitiesListStr, cap)
		}

		body.Capabilities = capabilitiesListStr
	}

	if isConfigOnly, ok := d.GetOk("is_config_only"); ok {
		body.IsConfigOnly = utils.BoolPtr(isConfigOnly.(bool))
	}

	resp, err := conn.VMAPIInstance.InsertVmGuestTools(utils.StringPtr(extID.(string)), body, args)

	if err != nil {
		return diag.Errorf("error while Inserting  gest tools ISO : %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template (%s) to Insert gest tools ISO: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while Inserting  gest tools ISO  : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(taskPoll.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId

	d.SetId(*uuid)

	return ResourceNutanixNGTInsertIsoV4Read(ctx, d, meta)
}

// Read NGT Configuration
func ResourceNutanixNGTInsertIsoV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extId := d.Id()
	resp, err := conn.VMAPIInstance.GetGuestToolsById(utils.StringPtr(extId))
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

// ResourceNutanixNGTInsertIsoV4Update  Not supported
func ResourceNutanixNGTInsertIsoV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixNGTInsertIsoV4Create(ctx, d, meta)
}

// ResourceNutanixNGTInsertIsoV4Delete  Not supported
func ResourceNutanixNGTInsertIsoV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
