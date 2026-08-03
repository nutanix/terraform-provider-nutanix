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
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/tasks"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/prism/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/vm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
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

	var TaskRef import1.TaskReference
	//nolint:gocritic // Keeping if-else for clarity in this specific case
	if action == "shutdown" {
		shutdownVmRequest := import3.ShutdownVmRequest{
			ExtId: utils.StringPtr(vmExtID.(string)),
		}
		resp, err := conn.VMAPIInstance.ShutdownVm(ctx, &shutdownVmRequest)
		if err != nil {
			return diag.Errorf("error while Shutdown VM : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	} else if action == "guest_shutdown" {
		shutdownGuestVmRequest := import3.ShutdownGuestVmRequest{
			ExtId: utils.StringPtr(vmExtID.(string)),
			Body:  &body,
		}
		resp, err := conn.VMAPIInstance.ShutdownGuestVm(ctx, &shutdownGuestVmRequest)
		if err != nil {
			return diag.Errorf("error while Shutdown Guest VM : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	} else if action == "reboot" {
		rebootVmRequest := import3.RebootVmRequest{
			ExtId: utils.StringPtr(vmExtID.(string)),
		}
		resp, err := conn.VMAPIInstance.RebootVm(ctx, &rebootVmRequest)
		if err != nil {
			return diag.Errorf("error while performing Reboot VM  : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	} else if action == "guest_reboot" {
		rebootGuestVmRequest := import3.RebootGuestVmRequest{
			ExtId: utils.StringPtr(vmExtID.(string)),
			Body:  &body,
		}
		resp, err := conn.VMAPIInstance.RebootGuestVm(ctx, &rebootGuestVmRequest)
		if err != nil {
			return diag.Errorf("error while performing Reboot Guest VM : %v", err)
		}
		TaskRef = resp.Data.GetValue().(import1.TaskReference)
	}

	// TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM action to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM action (%s) (%s) to complete: %s", action, utils.StringValue(taskUUID), errWaitTask)
	}
	log.Printf("[DEBUG] VM power action (%s) task (%s) reported SUCCEEDED", action, utils.StringValue(taskUUID))

	// The task reaching SUCCEEDED only means the power action was accepted/issued.
	// Both the hard "shutdown" and the guest actions are handled asynchronously by
	// AHV/the guest OS, so the VM's observed power state can lag behind the task
	// completion. Wait for the VM to actually reach its terminal power state so an
	// immediate re-read of the VM is deterministic (otherwise callers, e.g.
	// acceptance tests asserting power drift, race the shutdown/boot).
	expectedPowerState := ""
	switch action {
	case "shutdown", "guest_shutdown":
		expectedPowerState = "OFF"
	case "reboot", "guest_reboot":
		expectedPowerState = "ON"
	}

	if expectedPowerState != "" {
		// Log the power state observed immediately after the task succeeded to make
		// the async lag visible in debug output.
		if immediateResp, immErr := conn.VMAPIInstance.GetVmById(ctx, &import3.GetVmByIdRequest{ExtId: utils.StringPtr(vmExtID.(string))}); immErr == nil {
			if vmNow, ok := immediateResp.Data.GetValue().(config.Vm); ok && vmNow.PowerState != nil {
				log.Printf("[DEBUG] VM (%s) power state immediately after action (%s) task success: %s (expecting %s)",
					vmExtID.(string), action, vmNow.PowerState.GetName(), expectedPowerState)
			}
		} else {
			log.Printf("[DEBUG] VM (%s) immediate power-state read after action (%s) failed: %v", vmExtID.(string), action, immErr)
		}

		waitStart := time.Now()
		powerStateConf := &resource.StateChangeConf{
			Pending: []string{"WAITING"},
			Target:  []string{"DONE"},
			Refresh: func() (interface{}, string, error) {
				getVMReq := import3.GetVmByIdRequest{ExtId: utils.StringPtr(vmExtID.(string))}
				readResp, err := conn.VMAPIInstance.GetVmById(ctx, &getVMReq)
				if err != nil {
					log.Printf("[DEBUG] VM (%s) power-state poll after action (%s) errored: %v", vmExtID.(string), action, err)
					return nil, "", err
				}
				vmData, ok := readResp.Data.GetValue().(config.Vm)
				if !ok {
					return readResp, "WAITING", nil
				}
				current := "<nil>"
				if vmData.PowerState != nil {
					current = vmData.PowerState.GetName()
				}
				log.Printf("[DEBUG] VM (%s) power-state poll after action (%s): current=%s expected=%s elapsed=%s",
					vmExtID.(string), action, current, expectedPowerState, time.Since(waitStart).Round(time.Second))
				if current == expectedPowerState {
					return readResp, "DONE", nil
				}
				return readResp, "WAITING", nil
			},
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 5 * time.Second,
		}
		if _, errWaitPower := powerStateConf.WaitForStateContext(ctx); errWaitPower != nil {
			return diag.Errorf("error waiting for VM (%s) to reach power state %q after action (%s): %s", vmExtID.(string), expectedPowerState, action, errWaitPower)
		}
		log.Printf("[DEBUG] VM (%s) reached expected power state %s after action (%s) in %s",
			vmExtID.(string), expectedPowerState, action, time.Since(waitStart).Round(time.Second))
	}

	// Get UUID from TASK API
	getTaskByIdRequest := import4.GetTaskByIdRequest{
		ExtId: utils.StringPtr(*taskUUID),
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching VM action task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Shutdown Action Task Details: %s", string(aJSON))

	// This is an action resource that does not maintain state.
	// The resource ID is set to the task ExtId for traceability.
	d.SetId(utils.StringValue(taskDetails.ExtId))
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
