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

// DatasourceNutanixConnectionsV2 returns a paginated list of connections for a
// specific hardware provider.
func DatasourceNutanixConnectionsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixConnectionsV2Read,
		Schema: map[string]*schema.Schema{
			"hardware_provider_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External ID of the hardware provider",
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
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixConnectionV2(),
			},
		},
	}
}

func DatasourceNutanixConnectionsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)
	page, limit, filter, orderBy, selects := expandListQueryParams(d)

	resp, err := conn.HardwareProvidersAPIInstance.ListConnectionsByHardwareProviderId(utils.StringPtr(hpExtID), page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching connections : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("connections", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of connections.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.Connection)

	connections := make([]map[string]interface{}, len(getResp))
	for i := range getResp {
		connections[i] = flattenConnection(&getResp[i])
	}

	if err := d.Set("connections", connections); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// expandListQueryParams reads the shared page/limit/filter/order_by/select attributes
// from schema into SDK-ready pointers.
func expandListQueryParams(d *schema.ResourceData) (page, limit *int, filter, orderBy, selects *string) {
	if v, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(v.(string))
	}
	return page, limit, filter, orderBy, selects
}
