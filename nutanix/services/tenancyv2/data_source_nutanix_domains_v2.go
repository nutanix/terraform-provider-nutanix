package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixDomainsV2 lists all allowed domains in Nutanix Central tenancy.
func DatasourceNutanixDomainsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixDomainsV2Read,
		Schema: map[string]*schema.Schema{
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: domainDataSourceSchema(),
				},
			},
		},
	}
}

func domainDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain_ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": common.LinksSchema(),
	}
}

// DatasourceNutanixDomainsV2Read fetches allowed domains and populates state.
func DatasourceNutanixDomainsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	resp, err := conn.TenantManagementAPIInstance.ListDomains(ctx, &import4.ListDomainsRequest{})
	if err != nil {
		return diag.Errorf("error while fetching domains: %v", err)
	}
	if resp == nil || resp.Data == nil {
		if err := d.Set("domains", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of domains.",
		}}
	}

	domains, ok := resp.Data.GetValue().([]import1.Domain)
	if !ok {
		return diag.Errorf("error while fetching domains: unexpected response type %T", resp.Data.GetValue())
	}
	if err := d.Set("domains", flattenDomains(domains)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenDomains(domains []import1.Domain) []map[string]interface{} {
	if len(domains) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(domains))
	for i, domain := range domains {
		result[i] = map[string]interface{}{
			"domain_ext_id": utils.StringValue(domain.DomainExtId),
			"ext_id":        utils.StringValue(domain.ExtId),
			"tenant_id":     utils.StringValue(domain.TenantId),
			"links":         common.FlattenLinks(domain.Links),
		}
	}
	return result
}
