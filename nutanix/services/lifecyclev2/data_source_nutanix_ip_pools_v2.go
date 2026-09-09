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

// DatasourceNutanixIPPoolsV2 returns a list of IP address pools available in a
// hardware provider connection.
func DatasourceNutanixIPPoolsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixIPPoolsV2Read,
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
			"ip_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixIPPoolV2(),
			},
		},
	}
}

func DatasourceNutanixIPPoolsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)
	connExtID := d.Get("connection_ext_id").(string)
	page, limit, filter, orderBy, selects := expandListQueryParams(d)

	resp, err := conn.HardwareProvidersAPIInstance.ListIpPoolsByConnectionId(utils.StringPtr(hpExtID), utils.StringPtr(connExtID), page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching IP pools : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("ip_pools", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of IP pools.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.IpPool)

	pools := make([]map[string]interface{}, len(getResp))
	for i := range getResp {
		pools[i] = flattenIPPool(&getResp[i])
	}

	if err := d.Set("ip_pools", pools); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
