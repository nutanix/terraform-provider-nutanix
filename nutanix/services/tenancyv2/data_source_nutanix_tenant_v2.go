package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixTenantV2 reads a single tenant by its external identifier.
func DatasourceNutanixTenantV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixTenantV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"added_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"admin_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"aggregated_quota": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaDataSourceElem(),
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     stateInfoDataSourceElem(),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

// quotaDataSourceElem returns the fully-computed datasource schema for a quota block.
func quotaDataSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// stateInfoDataSourceElem returns the fully-computed datasource schema for a state_info block.
func stateInfoDataSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// DatasourceNutanixTenantV2Read fetches a tenant and populates state.
func DatasourceNutanixTenantV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	extID := d.Get("ext_id").(string)
	getTenantRequest := import4.GetTenantByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Tenant)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", getResp.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("added_by", getResp.AddedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domains", getResp.Domains); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("admin_info", flattenLocalUserInfo(getResp.AdminInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("aggregated_quota", flattenQuota(getResp.AggregatedQuota)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("state_info", flattenTenantStateInfo(getResp.StateInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
