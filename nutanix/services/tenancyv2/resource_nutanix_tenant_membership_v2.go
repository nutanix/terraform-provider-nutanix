package tenancyv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/common/v1/response"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/iam/v4/authz"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/membership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixTenantMembershipV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTenantMembershipV2Create,
		ReadContext:   ResourceNutanixTenantMembershipV2Read,
		DeleteContext: ResourceNutanixTenantMembershipV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identity_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"identity_details": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_provider_ext_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"principal_value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"USER", "GROUP"}, false),
						},
					},
				},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"principal": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixTenantMembershipV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	inputSpec := import1.Membership{}

	if v, ok := d.GetOk("domain_ext_id"); ok {
		inputSpec.DomainExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("role_ext_id"); ok {
		inputSpec.RoleExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("identity_ext_id"); ok {
		inputSpec.IdentityExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("identity_details"); ok {
		inputSpec.IdentityDetails = expandExternalUserDetails(v.([]interface{}))
	}

	createRequest := import3.CreateDomainMembershipRequest{
		Body: &inputSpec,
	}
	aJSON, _ := json.MarshalIndent(createRequest, "", " ")
	log.Printf("[DEBUG] Tenant Membership create payload: %s", string(aJSON))

	resp, err := conn.TenantMembershipAPIInstance.CreateDomainMembership(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error while creating tenant membership: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Membership)
	d.SetId(utils.StringValue(getResp.ExtId))
	return ResourceNutanixTenantMembershipV2Read(ctx, d, meta)
}

func ResourceNutanixTenantMembershipV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	getRequest := import3.GetDomainMembershipByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantMembershipAPIInstance.GetDomainMembershipById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant membership: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Membership)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_ext_id", getResp.DomainExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role_ext_id", getResp.RoleExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identity_ext_id", getResp.IdentityExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identity_details", flattenExternalUserDetails(getResp.IdentityDetails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_time", utils.TimeStringValue(getResp.CreatedTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("principal", getResp.Principal); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenAPILinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state_info", flattenMembershipStateInfo(getResp.StateInfo)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixTenantMembershipV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	deleteRequest := import3.DeleteDomainMembershipByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.TenantMembershipAPIInstance.DeleteDomainMembershipById(ctx, &deleteRequest)
	if err != nil {
		return diag.Errorf("error while deleting Tenant membership: %v", err)
	}

	if resp == nil {
		log.Println("[DEBUG] Tenant membership deleted successfully.")
	}

	return nil
}

func expandExternalUserDetails(input []interface{}) *import1.ExternalUserDetails {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	details := &import1.ExternalUserDetails{}
	if v, ok := m["identity_provider_ext_id"]; ok && v.(string) != "" {
		details.IdentityProviderExtId = utils.StringPtr(v.(string))
	}
	if v, ok := m["principal_value"]; ok && v.(string) != "" {
		details.PrincipalValue = utils.StringPtr(v.(string))
	}
	if v, ok := m["type"]; ok && v.(string) != "" {
		details.Type = common.ExpandEnum[import5.RmIdentityType](v.(string))
	}
	return details
}

func flattenExternalUserDetails(details *import1.ExternalUserDetails) []map[string]interface{} {
	if details == nil {
		return nil
	}
	m := map[string]interface{}{
		"identity_provider_ext_id": utils.StringValue(details.IdentityProviderExtId),
		"principal_value":          utils.StringValue(details.PrincipalValue),
		"type":                     common.FlattenPtrEnum(details.Type),
	}
	return []map[string]interface{}{m}
}

func flattenMembershipStateInfo(stateInfo *import1.MembershipStateInfo) []map[string]interface{} {
	if stateInfo == nil {
		return nil
	}
	m := map[string]interface{}{
		"state": common.FlattenPtrEnum(stateInfo.State),
	}
	return []map[string]interface{}{m}
}

func flattenAPILinks(links []import2.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		result = append(result, map[string]interface{}{
			"href": utils.StringValue(link.Href),
			"rel":  utils.StringValue(link.Rel),
		})
	}
	return result
}
