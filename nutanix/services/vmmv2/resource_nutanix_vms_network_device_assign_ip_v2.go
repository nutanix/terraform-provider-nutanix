package vmmv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsNetworkDeviceAssignIpV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsNetworkDeviceAssignIpV2Create,
		ReadContext:   ResourceNutanixVmsNetworkDeviceAssignIpV2Read,
		UpdateContext: ResourceNutanixVmsNetworkDeviceAssignIpV2Update,
		DeleteContext: ResourceNutanixVmsNetworkDeviceAssignIpV2Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefix_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 32),
							Default:      32,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixVmsNetworkDeviceAssignIpV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")
	body := config.AssignIpParams{}

	if ip_address, ok := d.GetOk("ip_address"); ok {
		body.IpAddress = expandIPv4Address(ip_address)
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.AssignIpById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)), &body, args)
	if err != nil {
		return diag.Errorf("error while assigning IP : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
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
		return diag.Errorf("error waiting for ip (%s) to assign: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func ResourceNutanixVmsNetworkDeviceAssignIpV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVmsNetworkDeviceAssignIpV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVmsNetworkDeviceAssignIpV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.ReleaseIpById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)), args)
	if err != nil {
		return diag.Errorf("error while releasing IP : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
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
		return diag.Errorf("error waiting for IP (%s) to release: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
