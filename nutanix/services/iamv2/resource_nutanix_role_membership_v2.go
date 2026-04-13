package iamv2

import (
	"context"
	"encoding/json"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/rolemembership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRoleMembershipV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRoleMembershipV2Create,
		ReadContext:   ResourceNutanixRoleMembershipV2Read,
		DeleteContext: ResourceNutanixRoleMembershipV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "External identifier of the role membership.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"role_ext_id": {
				Description: "External identifier of the role.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"identity_ext_id": {
				Description: "External identifier of the identity (user or group) associated with the role membership.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"identity_type": {
				Description: "Type of identity. Valid values are USER, GROUP.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"idp_ext_id": {
				Description: "External identifier of the identity provider.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"scope_template_name": {
				Description: "Name of the scope template.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"scope_template_name_values": {
				Description: "Name value pairs to substitute in the scope template variables.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"project_ext_id": {
				Description: "External identifier of the project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
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
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixRoleMembershipV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	body := iamConfig.RoleMembership{}

	if v, ok := d.GetOk("role_ext_id"); ok {
		body.RoleExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("identity_ext_id"); ok {
		body.IdentityExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("identity_type"); ok {
		body.IdentityType = expandRmIdentityType(v.(string))
	}
	if v, ok := d.GetOk("idp_ext_id"); ok {
		body.IdpExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("scope_template_name"); ok {
		body.ScopeTemplateName = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("scope_template_name_values"); ok {
		body.ScopeTemplateNameValues = expandScopeTemplateNameValues(v.([]interface{}))
	}
	if v, ok := d.GetOk("project_ext_id"); ok {
		body.ProjectExtId = utils.StringPtr(v.(string))
	}

	createRequest := import1.CreateRoleMembershipRequest{
		Body: &body,
	}
  Json, _ := json.MarshalIndent(createRequest, "", "  ")
  log.Printf("[DEBUG] Create Role Membership Request Body: %s", string(Json))
	resp, err := conn.RoleMembershipAPIInstance.CreateRoleMembership(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error while creating role membership: %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.RoleMembership)
	d.SetId(utils.StringValue(getResp.ExtId))

	return ResourceNutanixRoleMembershipV2Read(ctx, d, meta)
}

func ResourceNutanixRoleMembershipV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	getRequest := import1.GetRoleMembershipByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	resp, err := conn.RoleMembershipAPIInstance.GetRoleMembershipById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching role membership: %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.RoleMembership)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
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

	return nil
}

func ResourceNutanixRoleMembershipV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI
  getRequest := import1.GetRoleMembershipByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	resp, err := conn.RoleMembershipAPIInstance.GetRoleMembershipById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching role membership: %v", err)
	}
	getResp := resp.Data.GetValue().(iamConfig.RoleMembership)
	etagValue := conn.RoleMembershipAPIInstance.ApiClient.GetEtag(getResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	deleteRequest := import1.DeleteRoleMembershipByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}

	_, err = conn.RoleMembershipAPIInstance.DeleteRoleMembershipById(ctx, &deleteRequest, headers)
	if err != nil {
		return diag.Errorf("error while deleting role membership: %v", err)
	}
	
	return nil
}
