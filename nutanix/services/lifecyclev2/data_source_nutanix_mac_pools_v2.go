package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixMacPoolsV2 returns a list of MAC address pools available in a
// hardware provider connection.
func DatasourceNutanixMacPoolsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixMacPoolsV2Read,
		Schema: map[string]*schema.Schema{
			"hardware_provider_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External ID of the hardware provider",
			},
			"connection_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External ID of the hardware provider connection",
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
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mac_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixMacPoolV2(),
			},
		},
	}
}

func DatasourceNutanixMacPoolsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)
	connExtID := d.Get("connection_ext_id").(string)
	page, limit, filter, orderBy, selects := expandListQueryParams(d)

	resp, err := conn.HardwareProvidersAPIInstance.ListMacPoolsByConnectionId(utils.StringPtr(hpExtID), utils.StringPtr(connExtID), page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching MAC pools : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("mac_pools", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of MAC pools.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.MacPool)

	pools := make([]map[string]interface{}, len(getResp))
	for i := range getResp {
		pools[i] = flattenMacPool(&getResp[i])
	}

	if err := d.Set("mac_pools", pools); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
