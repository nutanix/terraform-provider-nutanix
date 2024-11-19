package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsGpusV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsGpusV4Create,
		ReadContext:   ResourceNutanixVmsGpusV4Read,
		UpdateContext: ResourceNutanixVmsGpusV4Update,
		DeleteContext: ResourceNutanixVmsGpusV4Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PASSTHROUGH_GRAPHICS", "PASSTHROUGH_COMPUTE", "VIRTUAL"}, false),
			},
			"device_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vendor": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NVIDIA", "AMD", "INTEL"}, false),
			},
			"pci_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"bus": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"func": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
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
			"guest_driver_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frame_buffer_size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_virtual_display_heads": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fraction": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixVmsGpusV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	gpuInput := config.Gpu{}

	if extID, ok := d.GetOk("ext_id"); ok {
		gpuInput.ExtId = utils.StringPtr(extID.(string))
	}
	if mode, ok := d.GetOk("mode"); ok {
		subMap := map[string]interface{}{
			"PASSTHROUGH_GRAPHICS": 2,
			"PASSTHROUGH_COMPUTE":  3,
			"VIRTUAL":              4,
		}
		pVal := subMap[mode.(string)]
		p := config.GpuMode(pVal.(int))
		gpuInput.Mode = &p
	}
	if deviceID, ok := d.GetOk("device_id"); ok {
		gpuInput.DeviceId = utils.IntPtr(deviceID.(int))
	}
	if vendor, ok := d.GetOk("vendor"); ok {
		subMap := map[string]interface{}{
			"NVIDIA": 2,
			"INTEL":  3,
			"AMD":    4,
		}
		pVal := subMap[vendor.(string)]
		p := config.GpuVendor(pVal.(int))
		gpuInput.Vendor = &p
	}
	if pci, ok := d.GetOk("pci_address"); ok {
		gpuInput.PciAddress = expandSBDF(pci)
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.CreateGpu(utils.StringPtr(vmExtID.(string)), &gpuInput, args)
	if err != nil {
		return diag.Errorf("error while creating gpu : %v", err)
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
		return diag.Errorf("error waiting for gpu (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// implement logic to get GPU UUID
	return ResourceNutanixVmsGpusV4Read(ctx, d, meta)
}

func ResourceNutanixVmsGpusV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	resp, err := conn.VMAPIInstance.GetGpuById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching gpu : %v", err)
	}
	gpuRespData := resp.Data.GetValue().(config.Gpu)

	if err := d.Set("tenant_id", gpuRespData.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenApiLink(gpuRespData.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mode", flattenGpuMode(gpuRespData.Mode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("device_id", gpuRespData.DeviceId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pci_address", flattenSBDF(gpuRespData.PciAddress)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_driver_version", gpuRespData.GuestDriverVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("frame_buffer_size_bytes", gpuRespData.FrameBufferSizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", gpuRespData.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_virtual_display_heads", gpuRespData.NumVirtualDisplayHeads); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("fraction", gpuRespData.Fraction); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVmsGpusV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVmsGpusV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	etagValue := conn.VMAPIInstance.ApiClient.GetEtag(readResp)

	args := make(map[string]interface{})
	args["If-Match"] = etagValue

	resp, err := conn.VMAPIInstance.DeleteGpuById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting gpu : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for gpu (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandSBDF(pr interface{}) *config.SBDF {
	if pr != nil {
		pci := &config.SBDF{}

		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if seg, ok := val["segment"]; ok {
			pci.Segment = utils.IntPtr(seg.(int))
		}
		if bus, ok := val["bus"]; ok {
			pci.Bus = utils.IntPtr(bus.(int))
		}
		if device, ok := val["device"]; ok {
			pci.Device = utils.IntPtr(device.(int))
		}
		if funct, ok := val["func"]; ok {
			pci.Func = utils.IntPtr(funct.(int))
		}
		return pci
	}
	return nil
}
