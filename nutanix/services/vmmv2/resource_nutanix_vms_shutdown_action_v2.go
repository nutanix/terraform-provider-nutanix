package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsShutdownActionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsShutdownActionV2Create,
		ReadContext:   ResourceNutanixVmsShutdownActionV2Read,
		UpdateContext: ResourceNutanixVmsShutdownActionV2Update,
		DeleteContext: ResourceNutanixVmsShutdownActionV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"shutdown", "guest_shutdown", "reboot", "guest_reboot"}, false),
			},
			"guest_power_state_transition_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"should_enable_script_exec": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"should_fail_on_script_failure": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixVmsShutdownActionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("ext_id")

	var action string
	body := config.GuestPowerOptions{}
	if actionType, ok := d.GetOk("action"); ok {
		action = actionType.(string)

		if action == "shutdown" || action == "reboot" {
			if _, ok := d.GetOk("guest_power_state_transition_config"); ok {
				return diag.Errorf("guest_power_state_transition_config  attribute is not optional for ['shutdown','reboot'] actions.")
			}
		}
	}

	if gst, ok := d.GetOk("guest_power_state_transition_config"); ok && len(gst.([]interface{})) > 0 {
		prI := gst.([]interface{})
		gstData := prI[0].(map[string]interface{})
		gstVal := config.GuestPowerStateTransitionConfig{}
		if enableScript, ok := gstData["should_enable_script_exec"]; ok {
			gstVal.ShouldEnableScriptExec = utils.BoolPtr(enableScript.(bool))
		}
		if scriptFailure, ok := gstData["should_fail_on_script_failure"]; ok {
			gstVal.ShouldFailOnScriptFailure = utils.BoolPtr(scriptFailure.(bool))
		}
		body.GuestPowerStateTransitionConfig = &gstVal
	}

	readResp, errR := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if errR != nil {
		return diag.Errorf("error while reading vm : %v", errR)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	var TaskRef import1.TaskReference
	//nolint:gocritic // Keeping if-else for clarity in this specific case
	if action == "shutdown" {
		resp, err := conn.VMAPIInstance.ShutdownVm(utils.StringPtr(vmExtID.(string)), args)
		if err != nil {
			return diag.Errorf("error while Shutdown VM : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	} else if action == "guest_shutdown" {
		resp, err := conn.VMAPIInstance.ShutdownGuestVm(utils.StringPtr(vmExtID.(string)), &body, args)
		if err != nil {
			return diag.Errorf("error while Shutdown Guest VM : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	} else if action == "reboot" {
		resp, err := conn.VMAPIInstance.RebootVm(utils.StringPtr(vmExtID.(string)), args)
		if err != nil {
			return diag.Errorf("error while performing Reboot VM  : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	} else if action == "guest_reboot" {
		resp, err := conn.VMAPIInstance.RebootGuestVm(utils.StringPtr(vmExtID.(string)), &body, args)
		if err != nil {
			return diag.Errorf("error while performing Reboot Guest VM : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	}

	// TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vm (%s) to perform (%s): %s", utils.StringValue(taskUUID), action, errWaitTask)
	}
	d.SetId(resource.UniqueId())
	return nil
}

func ResourceNutanixVmsShutdownActionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVmsShutdownActionV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixVmsShutdownActionV2Create(ctx, d, meta)
}

func ResourceNutanixVmsShutdownActionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
