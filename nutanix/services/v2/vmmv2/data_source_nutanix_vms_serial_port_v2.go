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

func DatasourceNutanixVmsSerialPortV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmsSerialPortV4Read,
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
			"is_connected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixVmsSerialPortV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")

	resp, err := conn.VMAPIInstance.GetSerialPortById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching serial port : %v", err)
	}

	serialResp := resp.Data

	if serialResp != nil {
		serialRespData := serialResp.GetValue().(config.SerialPort)
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
		d.SetId(*serialRespData.ExtId)
		return nil
	}
	d.SetId(resource.UniqueId())
	return nil
}
