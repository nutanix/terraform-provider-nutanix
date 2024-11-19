package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVmsGpuV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmsGpuV4Read,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pci_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"func": {
							Type:     schema.TypeInt,
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

func DatasourceNutanixVmsGpuV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")

	resp, err := conn.VMAPIInstance.GetGpuById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching gpu : %v", err)
	}

	gpuResp := resp.Data

	if gpuResp != nil {
		gpuRespData := gpuResp.GetValue().(config.Gpu)
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
		d.SetId(*gpuRespData.ExtId)
		return nil
	}
	d.SetId(resource.UniqueId())
	return nil
}
