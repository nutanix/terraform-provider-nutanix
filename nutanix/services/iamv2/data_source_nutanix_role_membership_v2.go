package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/rolemembership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRoleMembershipV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRoleMembershipV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "External identifier of the role membership.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"authorization_policy_ext_id": {
				Description: "External identifier of the authorization policy.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"role_ext_id": {
				Description: "External identifier of the role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"identity_ext_id": {
				Description: "External identifier of the identity (user or group) associated with the role membership.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"identity_type": {
				Description: "Type of identity. Valid values are USER, GROUP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"identity_value": {
				Description: "Value of the identity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"idp_ext_id": {
				Description: "External identifier of the identity provider.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"scope_template_name": {
				Description: "Name of the scope template.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"scope_template_name_values": schemaForScopeTemplateNameValues(),
			"project_ext_id": {
				Description: "External identifier of the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key_value_pairs": schemaForKeyValuePairs(),
			"created_by": {
				Description: "User or service name that created the role membership.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_time": {
				Description: "The creation time of the role membership.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated_time": {
				Description: "The time when the role membership was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixRoleMembershipV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id").(string)
	getRequest := import1.GetRoleMembershipByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.RoleMembershipAPIInstance.GetRoleMembershipById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching role membership: %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.RoleMembership)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("authorization_policy_ext_id", getResp.AuthorizationPolicyExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role_ext_id", getResp.RoleExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identity_ext_id", getResp.IdentityExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identity_type", flattenRmIdentityType(getResp.IdentityType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identity_value", getResp.IdentityValue); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("idp_ext_id", getResp.IdpExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scope_template_name", getResp.ScopeTemplateName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scope_template_name_values", flattenScopeTemplateNameValues(getResp.ScopeTemplateNameValues)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", getResp.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key_value_pairs", flattenKeyValuePairs(getResp.KeyValuePairs)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdatedTime != nil {
		if err := d.Set("last_updated_time", getResp.LastUpdatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
