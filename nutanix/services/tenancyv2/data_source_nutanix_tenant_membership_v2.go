package tenancyv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/membership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixTenantMembershipV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixTenantMembershipV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_provider_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceNutanixTenantMembershipV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	getRequest := import3.GetDomainMembershipByIdRequest{
		ExtId: utils.StringPtr(d.Get("ext_id").(string)),
	}
	resp, err := conn.TenantMembershipAPIInstance.GetDomainMembershipById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant membership: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Membership)
	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] Get Tenant Membership Response: %s", string(aJSON))

	d.SetId(utils.StringValue(getResp.ExtId))

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
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
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
