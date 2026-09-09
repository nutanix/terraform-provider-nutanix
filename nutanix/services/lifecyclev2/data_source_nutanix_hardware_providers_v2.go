package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

// DatasourceNutanixHardwareProvidersV2 returns a paginated list of all available
// hardware providers.
func DatasourceNutanixHardwareProvidersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixHardwareProvidersV2Read,
		Schema: map[string]*schema.Schema{
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
			"hardware_providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixHardwareProviderV2(),
			},
		},
	}
}

func DatasourceNutanixHardwareProvidersV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	page, limit, filter, orderBy, selects := expandListQueryParams(d)

	resp, err := conn.HardwareProvidersAPIInstance.ListHardwareProviders(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching hardware providers : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("hardware_providers", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of hardware providers.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.HardwareProvider)

	providers := make([]map[string]interface{}, len(getResp))
	for i := range getResp {
		providers[i] = flattenHardwareProvider(&getResp[i])
	}

	if err := d.Set("hardware_providers", providers); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
