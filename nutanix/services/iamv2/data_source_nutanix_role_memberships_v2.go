package iamv2

import (
	"context"
	"encoding/json"
	"log"

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
				Elem:     DatasourceNutanixRoleMembershipV2(),
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
  Json, _ := json.MarshalIndent(listRequest, "", "  ")
  log.Printf("[DEBUG] List Role Memberships Request Body: %s", string(Json))
	resp, err := conn.RoleMembershipAPIInstance.ListRoleMemberships(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while fetching role memberships: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("role_memberships", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of role memberships.",
		}}
	}

  membershipsList := resp.Data.GetValue().([]iamConfig.RoleMembershipProjection)
	if err := d.Set("role_memberships", flattenRoleMembershipEntities(membershipsList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRoleMembershipEntities(memberships []iamConfig.RoleMembershipProjection) []interface{} {
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
