package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsSerialPortsV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsSerialPortsV4Create,
		ReadContext:   ResourceNutanixVmsSerialPortsV4Read,
		UpdateContext: ResourceNutanixVmsSerialPortsV4Update,
		DeleteContext: ResourceNutanixVmsSerialPortsV4Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_connected": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"index": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixVmsSerialPortsV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	body := config.SerialPort{}
	var idxx int

	if isConn, ok := d.GetOk("is_connected"); ok {
		body.IsConnected = utils.BoolPtr(isConn.(bool))
	}
	if idx, ok := d.GetOk("index"); ok {
		body.Index = utils.IntPtr(idx.(int))
		idxx = utils.IntValue(body.Index)
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.CreateSerialPort(utils.StringPtr(vmExtID.(string)), &body, args)
	if err != nil {
		return diag.Errorf("error while creating serial ports : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for serial port (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// reading again VM to fetch serial port UUID

	SubResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	out := SubResp.Data.GetValue().(config.Vm)

	serialPortOut := out.SerialPorts
	serialUUID := ""
	for _, v := range serialPortOut {
		if utils.IntValue(v.Index) == idxx {
			serialUUID = *v.ExtId
		}
	}
	d.SetId(serialUUID)
	return ResourceNutanixVmsSerialPortsV4Read(ctx, d, meta)
}

func ResourceNutanixVmsSerialPortsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	resp, err := conn.VMAPIInstance.GetSerialPortById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching serial port : %v", err)
	}

	serialRespData := resp.Data.GetValue().(config.SerialPort)

	if err := d.Set("tenant_id", serialRespData.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenApiLink(serialRespData.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_connected", serialRespData.IsConnected); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("index", serialRespData.Index); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVmsSerialPortsV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	resp, err := conn.VMAPIInstance.GetSerialPortById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching serial port : %v", err)
	}

	getResp := resp.Data.GetValue().(config.SerialPort)

	updateSpec := getResp

	if d.HasChange("is_connected") {
		updateSpec.IsConnected = utils.BoolPtr(d.Get("is_connected").(bool))
	}
	if d.HasChange("index") {
		updateSpec.Index = utils.IntPtr(d.Get("index").(int))
	}

	updateResp, err := conn.VMAPIInstance.UpdateSerialPortById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating serial port : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for serial port (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixVmsSerialPortsV4Read(ctx, d, meta)
}

func ResourceNutanixVmsSerialPortsV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.DeleteSerialPortById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while fetching serial port : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for serial port (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
