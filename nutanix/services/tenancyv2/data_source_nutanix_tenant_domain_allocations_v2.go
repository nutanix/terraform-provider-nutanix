package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixTenantDomainAllocationsV2 lists all domain allocations for a tenant.
func DatasourceNutanixTenantDomainAllocationsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixTenantDomainAllocationsV2Read,
		Schema: map[string]*schema.Schema{
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_allocations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixTenantDomainAllocationV2(),
			},
		},
	}
}

// DatasourceNutanixTenantDomainAllocationsV2Read fetches the list of domain allocations and populates state.
func DatasourceNutanixTenantDomainAllocationsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)

	listRequest := import4.ListTenantDomainAllocationsRequest{
		TenantExtId: utils.StringPtr(tenantExtID),
	}
	if v, ok := d.GetOk("expand"); ok {
		listRequest.Expand_ = utils.StringPtr(v.(string))
	}
	resp, err := conn.TenantManagementAPIInstance.ListTenantDomainAllocations(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant domain allocations: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("domain_allocations", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of tenant domain allocations.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.DomainAllocation)

	if err := d.Set("domain_allocations", flattenDomainAllocationsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// flattenDomainAllocationsEntities maps a list of DomainAllocation responses to datasource entries.
func flattenDomainAllocationsEntities(allocations []import1.DomainAllocation) []map[string]interface{} {
	if len(allocations) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(allocations))
	for k, v := range allocations {
		allocation := make(map[string]interface{})
		allocation["ext_id"] = utils.StringValue(v.ExtId)
		allocation["domain_ext_id"] = utils.StringValue(v.DomainExtId)
		allocation["external_network_config"] = flattenExternalNetworkConnection(v.ExternalNetworkConfig)
		allocation["projects"] = flattenProjects(v.Projects)
		allocation["state_info"] = flattenDomainAllocationStateInfo(v.StateInfo)
		allocation["tenant_id"] = utils.StringValue(v.TenantId)
		allocation["links"] = common.FlattenLinks(v.Links)
		result[k] = allocation
	}
	return result
}
