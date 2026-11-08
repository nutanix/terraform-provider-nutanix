package tenancyv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/common/v1/response"
	import5 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/iam/v4/authz"
	import1 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/request/membership"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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

	// Memberships (especially GROUP) can linger in PENDING; DELETE is rejected
	// until ACTIVE (TEN-33008 INVALID_STATE_TRANSITION).
	if err := waitForTenantMembershipState(ctx, d, meta,
		[]string{"PENDING", "UNKNOWN", ""},
		[]string{"ACTIVE"},
		d.Timeout(schema.TimeoutCreate),
	); err != nil {
		return diag.FromErr(err)
	}

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
		log.Println("[DEBUG] Tenant membership delete accepted.")
	}

	// In Future, when the task api is supported, we can wait for the delete to complete.

	d.SetId("")
	return nil
}

func waitForTenantMembershipState(
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	pending []string,
	target []string,
	timeout time.Duration,
) error {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI
	stateConf := &resource.StateChangeConf{
		Pending:      pending,
		Target:       target,
		Timeout:      timeout,
		PollInterval: 3 * time.Second,
		Refresh: func() (interface{}, string, error) {
			getRequest := import3.GetDomainMembershipByIdRequest{
				ExtId: utils.StringPtr(d.Id()),
			}
			resp, err := conn.TenantMembershipAPIInstance.GetDomainMembershipById(ctx, &getRequest)
			if err != nil {
				return nil, "", fmt.Errorf("error while checking tenant membership state: %w", err)
			}
			if resp == nil || resp.Data == nil {
				return nil, "", fmt.Errorf("error while checking tenant membership state: empty response")
			}

			membership, ok := resp.Data.GetValue().(import1.Membership)
			if !ok {
				return nil, "", fmt.Errorf("error while checking tenant membership state: unexpected type %T", resp.Data.GetValue())
			}
			if membership.StateInfo == nil || membership.StateInfo.State == nil {
				log.Printf("[DEBUG] Tenant membership lifecycle state: UNKNOWN")
				return &membership, "", nil
			}

			state := membership.StateInfo.State.GetName()
			log.Printf("[DEBUG] Tenant membership %s lifecycle state: %s", d.Id(), state)
			if state == "ERROR" {
				return &membership, state, fmt.Errorf("tenant membership reached ERROR state")
			}
			return &membership, state, nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for tenant membership %q to reach %s: %w", d.Id(), target[0], err)
	}
	return nil
}

// waitForTenantMembershipDeleted polls GET until the membership is gone (API
// error / empty) or reports a terminal deleted state.
func waitForTenantMembershipDeleted(
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	timeout time.Duration,
) error {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE", "PENDING", "DELETING", "UNKNOWN", ""},
		Target:       []string{"DELETED"},
		Timeout:      timeout,
		PollInterval: 3 * time.Second,
		Refresh: func() (interface{}, string, error) {
			getRequest := import3.GetDomainMembershipByIdRequest{
				ExtId: utils.StringPtr(d.Id()),
			}
			resp, err := conn.TenantMembershipAPIInstance.GetDomainMembershipById(ctx, &getRequest)
			if err != nil {
				// Not found / transport error from deleted resource: treat as gone.
				log.Printf("[DEBUG] Tenant membership %s GET after delete: %v (treating as DELETED)", d.Id(), err)
				return struct{}{}, "DELETED", nil
			}
			if resp == nil || resp.Data == nil {
				log.Printf("[DEBUG] Tenant membership %s empty GET after delete (treating as DELETED)", d.Id())
				return struct{}{}, "DELETED", nil
			}

			membership, ok := resp.Data.GetValue().(import1.Membership)
			if !ok {
				return nil, "", fmt.Errorf("error while checking tenant membership deletion: unexpected type %T", resp.Data.GetValue())
			}
			if membership.StateInfo == nil || membership.StateInfo.State == nil {
				log.Printf("[DEBUG] Tenant membership %s lifecycle state after delete: UNKNOWN", d.Id())
				return &membership, "", nil
			}

			state := membership.StateInfo.State.GetName()
			log.Printf("[DEBUG] Tenant membership %s lifecycle state after delete: %s", d.Id(), state)
			if state == "ERROR" {
				return &membership, state, fmt.Errorf("tenant membership %q reached ERROR while deleting", d.Id())
			}
			// Some APIs surface a terminal DELETED enum; map it to Target.
			if state == "DELETED" {
				return &membership, "DELETED", nil
			}
			return &membership, state, nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for tenant membership %q to be deleted: %w", d.Id(), err)
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
