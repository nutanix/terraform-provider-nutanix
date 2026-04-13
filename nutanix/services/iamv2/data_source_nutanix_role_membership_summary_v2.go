package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/rolemembership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRoleMembershipSummaryV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRoleMembershipSummaryV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Description: "A URL query parameter that specifies the page number of the result set.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"limit": {
				Description: "A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.",
				Type:        schema.TypeInt,
				Optional:    true,
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
			"summaries": {
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
						"links": schemaForLinks(),
						"users_count": {
							Description: "Count of distinct users.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"groups_count": {
							Description: "Count of distinct groups.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"roles_count": {
							Description: "Count of distinct roles.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"total_identities_count": {
							Description: "Total count of identities.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixRoleMembershipSummaryV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	listRequest := import1.ListRoleMembershipSummaryRequest{}
	if v, ok := d.GetOk("page"); ok {
		listRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.RoleMembershipAPIInstance.ListRoleMembershipSummary(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while fetching role membership summaries: %v", err)
	}

	summariesRaw := resp.Data.GetValue()
	summariesList, ok := summariesRaw.([]iamConfig.RoleMembershipSummary)
	if !ok || len(summariesList) == 0 {
		if err := d.Set("summaries", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of role membership summaries.",
		}}
	}

	if err := d.Set("summaries", flattenRoleMembershipSummaries(summariesList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRoleMembershipSummaries(summaries []iamConfig.RoleMembershipSummary) []interface{} {
	if len(summaries) == 0 {
		return nil
	}
	result := make([]interface{}, len(summaries))
	for i, s := range summaries {
		summary := make(map[string]interface{})
		if s.ExtId != nil {
			summary["ext_id"] = utils.StringValue(s.ExtId)
		}
		if s.TenantId != nil {
			summary["tenant_id"] = utils.StringValue(s.TenantId)
		}
		if s.Links != nil {
			summary["links"] = flattenLinks(s.Links)
		}
		if s.UsersCount != nil {
			summary["users_count"] = *s.UsersCount
		}
		if s.GroupsCount != nil {
			summary["groups_count"] = *s.GroupsCount
		}
		if s.RolesCount != nil {
			summary["roles_count"] = *s.RolesCount
		}
		if s.TotalIdentitiesCount != nil {
			summary["total_identities_count"] = *s.TotalIdentitiesCount
		}
		result[i] = summary
	}
	return result
}
