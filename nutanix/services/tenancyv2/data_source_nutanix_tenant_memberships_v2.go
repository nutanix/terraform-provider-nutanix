package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/request/membership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixTenantMembershipsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixTenantMembershipsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"memberships": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceNutanixTenantMembershipV2(),
			},
			"total_available_results": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixTenantMembershipsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	listRequest := import3.ListDomainMembershipsRequest{}

	if v, ok := d.GetOk("page"); ok {
		listRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listRequest.Filter_ = utils.StringPtr(v.(string))
	}
	resp, err := conn.TenantMembershipAPIInstance.ListDomainMemberships(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant memberships: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("memberships", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of Tenant memberships.",
		}}
	}

	memberships := resp.Data.GetValue().([]import1.Membership)
	if err := d.Set("memberships", flattenTenantMemberships(memberships)); err != nil {
		return diag.FromErr(err)
	}
	if resp.Metadata != nil && resp.Metadata.TotalAvailableResults != nil {
		if err := d.Set("total_available_results", *resp.Metadata.TotalAvailableResults); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(resource.UniqueId())
	return nil
}

func flattenTenantMemberships(memberships []import1.Membership) []map[string]interface{} {
	if len(memberships) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(memberships))
	for _, membership := range memberships {
		result = append(result, flattenTenantMembership(membership))
	}
	return result
}

func flattenTenantMembership(membership import1.Membership) map[string]interface{} {
	m := map[string]interface{}{
		"ext_id":           utils.StringValue(membership.ExtId),
		"domain_ext_id":    utils.StringValue(membership.DomainExtId),
		"role_ext_id":      utils.StringValue(membership.RoleExtId),
		"identity_ext_id":  utils.StringValue(membership.IdentityExtId),
		"identity_details": flattenExternalUserDetails(membership.IdentityDetails),
		"created_by":       utils.StringValue(membership.CreatedBy),
		"principal":        utils.StringValue(membership.Principal),
		"tenant_id":        utils.StringValue(membership.TenantId),
		"links":            flattenAPILinks(membership.Links),
		"state_info":       flattenMembershipStateInfo(membership.StateInfo),
		"created_time":     utils.TimeStringValue(membership.CreatedTime),
	}
	return m
}
