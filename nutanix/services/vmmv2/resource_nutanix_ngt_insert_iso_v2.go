package vmmv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmPrism "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	vmmConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixNGTInsertIsoV2 TF schema for NGT install/uninstall
func ResourceNutanixNGTInsertIsoV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixNGTInsertIsoV2Create,
		ReadContext:   ResourceNutanixNGTInsertIsoV2Read,
		UpdateContext: ResourceNutanixNGTInsertIsoV2Update,
		DeleteContext: ResourceNutanixNGTInsertIsoV2Delete,
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
			"cdrom_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "insert",
				ValidateFunc: validation.StringInSlice([]string{"insert", "eject"}, false),
			},
		},
	}
}

// Install NGT on Vm
func ResourceNutanixNGTInsertIsoV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	if action, ok := d.GetOk("action"); ok && action.(string) == "insert" {
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
		for _, entity := range rUUID.EntitiesAffected {
			if utils.StringValue(entity.Rel) == "vmm:ahv:config:vm:cdrom" {
				uuid := entity.ExtId
				d.Set("cdrom_ext_id", *uuid)
			}
		}

		d.SetId(resource.UniqueId())

		return ResourceNutanixNGTInsertIsoV2Read(ctx, d, meta)
	}
	return diag.Errorf("Action %s is not supported for NGT ISO Insert", d.Get("action").(string))
}

// Read NGT Configuration
func ResourceNutanixNGTInsertIsoV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)
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
	// Check if cdrom_ext_id is present and not null in the state file
	cdromExtID, cdromExists := d.GetOk("cdrom_ext_id")
	if cdromExists && cdromExtID != nil && cdromExtID.(string) != "" {
		if err := d.Set("cdrom_ext_id", cdromExtID.(string)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		// We need to find the CD-ROM ext id with iso_type GUEST_TOOLS, if possible
		vmResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(extID))
		if err != nil {
			return diag.Errorf("error while fetching vm details which helps us to set cdrom_ext_id : %v", err)
		}
		vmObj := vmResp.Data.GetValue().(vmmConfig.Vm)
		// Check that CdRoms is not nil and loop through the CdRoms to find the GUEST_TOOLS ISO
		if len(vmObj.CdRoms) > 0 {
			for _, cdrom := range vmObj.CdRoms {
				if cdrom.IsoType.GetName() == "GUEST_TOOLS" {
					if err := d.Set("cdrom_ext_id", utils.StringValue(cdrom.ExtId)); err != nil {
						return diag.FromErr(err)
					}
					break
				}
			}
		}
	}
	// Set the vm_ext_id to the state file
	if err := d.Set("vm_ext_id", extID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// ResourceNutanixNGTInsertIsoV2Update  Not supported
func ResourceNutanixNGTInsertIsoV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if action, ok := d.GetOk("action"); ok && action.(string) == "eject" {
		log.Printf("[DEBUG] ResourceNutanixNGTInsertIsoV2Update : Action %s", action.(string))
		diags := ejectCdromISO(ctx, d, meta)
		if diags.HasError() {
			// Ejection failed, set the action to INSERT to avoid Terraform from saving "EJECT" in state
			d.Set("action", "insert")
			return diags
		}
		return ResourceNutanixNGTInsertIsoV2Read(ctx, d, meta)
	}
	return ResourceNutanixNGTInsertIsoV2Create(ctx, d, meta)
}

// ResourceNutanixNGTInsertIsoV2Delete eject the ngt iso from the cd-rom of the vm
func ResourceNutanixNGTInsertIsoV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixNGTInsertIsoV2Delete : Ejecting NGT ISO from the CD-ROM %s of the VM %s", d.Get("cdrom_ext_id").(string), d.Get("vm_ext_id").(string))
	if action, ok := d.GetOk("action"); ok && action.(string) == "eject" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "NGT ISO is not inserted on the CD-ROM of the VM or ejected earlier using an action, Ignoring the request to eject the NGT ISO",
		}}
	}
	return ejectCdromISO(ctx, d, meta)
}

func ejectCdromISO(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Ejecting ISO from the CD-ROM %s of the VM %s", d.Get("cdrom_ext_id").(string), d.Get("vm_ext_id").(string))
	conn := meta.(*conns.Client).VmmAPI
	vmExtID := d.Get("vm_ext_id").(string)
	extID := d.Get("cdrom_ext_id").(string)

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	// eject the ngt iso from the cd-rom of the vm
	resp, err := conn.VMAPIInstance.EjectCdRomById(utils.StringPtr(vmExtID), utils.StringPtr(extID), args)
	if err != nil {
		return diag.Errorf("error while ejecting cd-rom : %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the cd-rom to be ejected
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("ISO EJECTION FAILED: REASON: %s : Task UUID: %s", errWaitTask, utils.StringValue(taskUUID))
	}
	return nil
}
