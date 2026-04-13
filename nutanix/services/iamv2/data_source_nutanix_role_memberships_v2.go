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

func DatasourceNutanixRoleMembershipsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRoleMembershipsV2Read,
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
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_memberships": {
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
						"links":                      schemaForLinks(),
						"authorization_policy_ext_id": {Type: schema.TypeString, Computed: true},
						"role_ext_id":                 {Type: schema.TypeString, Computed: true},
						"identity_ext_id":             {Type: schema.TypeString, Computed: true},
						"identity_type":               {Type: schema.TypeString, Computed: true},
						"identity_value":              {Type: schema.TypeString, Computed: true},
						"idp_ext_id":                  {Type: schema.TypeString, Computed: true},
						"scope_template_name":         {Type: schema.TypeString, Computed: true},
						"scope_template_name_values":  schemaForScopeTemplateNameValues(),
						"project_ext_id":              {Type: schema.TypeString, Computed: true},
						"key_value_pairs":             schemaForKeyValuePairs(),
						"created_by":                  {Type: schema.TypeString, Computed: true},
						"created_time":                {Type: schema.TypeString, Computed: true},
						"last_updated_time":           {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func DatasourceNutanixRoleMembershipsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	listRequest := import1.ListRoleMembershipsRequest{}
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
	if v, ok := d.GetOk("expand"); ok {
		listRequest.Expand_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.RoleMembershipAPIInstance.ListRoleMemberships(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while fetching role memberships: %v", err)
	}

	membershipsRaw := resp.Data.GetValue()
	membershipsList, ok := membershipsRaw.([]iamConfig.RoleMembership)
	if !ok || len(membershipsList) == 0 {
		if err := d.Set("role_memberships", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of role memberships.",
		}}
	}

	if err := d.Set("role_memberships", flattenRoleMembershipEntities(membershipsList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRoleMembershipEntities(memberships []iamConfig.RoleMembership) []interface{} {
	if len(memberships) == 0 {
		return nil
	}
	result := make([]interface{}, len(memberships))
	for i, m := range memberships {
		membership := make(map[string]interface{})
		if m.ExtId != nil {
			membership["ext_id"] = utils.StringValue(m.ExtId)
		}
		if m.TenantId != nil {
			membership["tenant_id"] = utils.StringValue(m.TenantId)
		}
		if m.Links != nil {
			membership["links"] = flattenLinks(m.Links)
		}
		if m.AuthorizationPolicyExtId != nil {
			membership["authorization_policy_ext_id"] = utils.StringValue(m.AuthorizationPolicyExtId)
		}
		if m.RoleExtId != nil {
			membership["role_ext_id"] = utils.StringValue(m.RoleExtId)
		}
		if m.IdentityExtId != nil {
			membership["identity_ext_id"] = utils.StringValue(m.IdentityExtId)
		}
		membership["identity_type"] = flattenRmIdentityType(m.IdentityType)
		if m.IdentityValue != nil {
			membership["identity_value"] = utils.StringValue(m.IdentityValue)
		}
		if m.IdpExtId != nil {
			membership["idp_ext_id"] = utils.StringValue(m.IdpExtId)
		}
		if m.ScopeTemplateName != nil {
			membership["scope_template_name"] = utils.StringValue(m.ScopeTemplateName)
		}
		if m.ScopeTemplateNameValues != nil {
			membership["scope_template_name_values"] = flattenScopeTemplateNameValues(m.ScopeTemplateNameValues)
		}
		if m.ProjectExtId != nil {
			membership["project_ext_id"] = utils.StringValue(m.ProjectExtId)
		}
		if m.KeyValuePairs != nil {
			membership["key_value_pairs"] = flattenKeyValuePairs(m.KeyValuePairs)
		}
		if m.CreatedBy != nil {
			membership["created_by"] = utils.StringValue(m.CreatedBy)
		}
		if m.CreatedTime != nil {
			membership["created_time"] = m.CreatedTime.String()
		}
		if m.LastUpdatedTime != nil {
			membership["last_updated_time"] = m.LastUpdatedTime.String()
		}
		result[i] = membership
	}
	return result
}
