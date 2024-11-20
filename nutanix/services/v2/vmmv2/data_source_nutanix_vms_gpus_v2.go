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

func DatasourceNutanixVmsGpusV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmsGpusV4Read,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gpus": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func DatasourceNutanixVmsGpusV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	// initialize query params
	var page, limit *int
	var filter *string

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	vmExtID := d.Get("vm_ext_id")

	resp, err := conn.VMAPIInstance.ListGpusByVmId(utils.StringPtr(vmExtID.(string)), page, limit, filter)
	if err != nil {
		return diag.Errorf("error while fetching gpus : %v", err)
	}

	serialResp := resp.Data

	if serialResp != nil {
		serialRespData := serialResp.GetValue().([]config.Gpu)

		if err := d.Set("gpus", flattenGpusEntities(serialRespData)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(resource.UniqueId())
	return nil
}

func flattenGpusEntities(pr []config.Gpu) []interface{} {
	if len(pr) > 0 {
		gpus := make([]interface{}, len(pr))

		for k, v := range pr {
			gpu := make(map[string]interface{})

			if v.ExtId != nil {
				gpu["ext_id"] = v.ExtId
			}
			if v.TenantId != nil {
				gpu["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				gpu["links"] = v.Links
			}
			if v.Mode != nil {
				gpu["mode"] = flattenGpuMode(v.Mode)
			}
			if v.DeviceId != nil {
				gpu["device_id"] = v.DeviceId
			}
			if v.Vendor != nil {
				gpu["vendor"] = flattenGpuVendor(v.Vendor)
			}
			if v.PciAddress != nil {
				gpu["pci_address"] = flattenSBDF(v.PciAddress)
			}
			if v.GuestDriverVersion != nil {
				gpu["guest_driver_version"] = v.GuestDriverVersion
			}
			if v.Name != nil {
				gpu["name"] = v.Name
			}
			if v.FrameBufferSizeBytes != nil {
				gpu["frame_buffer_size_bytes"] = v.FrameBufferSizeBytes
			}
			if v.NumVirtualDisplayHeads != nil {
				gpu["num_virtual_display_heads"] = v.NumVirtualDisplayHeads
			}
			if v.Fraction != nil {
				gpu["fraction"] = v.Fraction
			}
			gpus[k] = gpu
		}
		return gpus
	}
	return nil
}
