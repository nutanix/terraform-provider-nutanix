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

func DatasourceNutanixVmsSerialPortsV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmsSerialPortsV4Read,
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
			"serial_ports": {
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
						"is_connected": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVmsSerialPortsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	// initialize query params
	var page, limit *int

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
	vmExtID := d.Get("vm_ext_id")

	resp, err := conn.VMAPIInstance.ListSerialPortsByVmId(utils.StringPtr(vmExtID.(string)), page, limit)
	if err != nil {
		return diag.Errorf("error while fetching serial ports : %v", err)
	}

	serialResp := resp.Data

	if serialResp != nil {
		serialRespData := serialResp.GetValue().([]config.SerialPort)

		if err := d.Set("serial_ports", flattenSerialPortsEntities(serialRespData)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(resource.UniqueId())
	return nil
}

func flattenSerialPortsEntities(pr []config.SerialPort) []interface{} {
	if len(pr) > 0 {
		ports := make([]interface{}, len(pr))

		for k, v := range pr {
			port := make(map[string]interface{})

			if v.ExtId != nil {
				port["ext_id"] = v.ExtId
			}
			if v.TenantId != nil {
				port["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				port["links"] = flattenApiLink(v.Links)
			}
			if v.IsConnected != nil {
				port["is_connected"] = v.IsConnected
			}
			if v.Index != nil {
				port["index"] = v.Index
			}
			ports[k] = port
		}
		return ports
	}
	return nil
}
