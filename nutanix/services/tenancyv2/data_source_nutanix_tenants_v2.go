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

// DatasourceNutanixTenantsV2 lists all tenants.
func DatasourceNutanixTenantsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixTenantsV2Read,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tenants": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixTenantV2(),
			},
		},
	}
}

// DatasourceNutanixTenantsV2Read fetches the list of tenants and populates state.
func DatasourceNutanixTenantsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	listTenantsRequest := import4.ListTenantsRequest{}
	if v, ok := d.GetOk("filter"); ok {
		listTenantsRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("page"); ok {
		listTenantsRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listTenantsRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listTenantsRequest.Orderby_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.TenantManagementAPIInstance.ListTenants(ctx, &listTenantsRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenants: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("tenants", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of tenants.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.Tenant)

	if err := d.Set("tenants", flattenTenantsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// flattenTenantsEntities maps a list of Tenant responses to datasource entries.
func flattenTenantsEntities(tenants []import1.Tenant) []map[string]interface{} {
	if len(tenants) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(tenants))
	for k, v := range tenants {
		tenant := make(map[string]interface{})
		tenant["ext_id"] = utils.StringValue(v.ExtId)
		tenant["name"] = utils.StringValue(v.Name)
		tenant["description"] = utils.StringValue(v.Description)
		tenant["url"] = utils.StringValue(v.Url)
		tenant["added_by"] = utils.StringValue(v.AddedBy)
		tenant["domains"] = v.Domains
		tenant["admin_info"] = flattenLocalUserInfo(v.AdminInfo)
		tenant["aggregated_quota"] = flattenQuota(v.AggregatedQuota)
		if v.CreatedTime != nil {
			tenant["created_time"] = v.CreatedTime.String()
		}
		tenant["state_info"] = flattenTenantStateInfo(v.StateInfo)
		tenant["tenant_id"] = utils.StringValue(v.TenantId)
		tenant["links"] = common.FlattenLinks(v.Links)
		result[k] = tenant
	}
	return result
}
