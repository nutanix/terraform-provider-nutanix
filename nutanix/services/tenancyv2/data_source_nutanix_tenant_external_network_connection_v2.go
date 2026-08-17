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

// DatasourceNutanixTenantExternalNetworkConnectionV2 reads the external network
// connection settings for a tenant domain allocation.
func DatasourceNutanixTenantExternalNetworkConnectionV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixTenantExternalNetworkConnectionV2Read,
		Schema: map[string]*schema.Schema{
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_allocation_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     externalNetworkConnConfigDataSourceElem(),
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

// DatasourceNutanixTenantExternalNetworkConnectionV2Read fetches the external network
// connection settings and populates state.
func DatasourceNutanixTenantExternalNetworkConnectionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	tenantExtID := d.Get("tenant_ext_id").(string)
	domainAllocationExtID := d.Get("domain_allocation_ext_id").(string)

	getRequest := import4.GetTenantExternalNetworkConnectionRequest{
		TenantExtId:           utils.StringPtr(tenantExtID),
		DomainAllocationExtId: utils.StringPtr(domainAllocationExtID),
	}
	resp, err := conn.TenantManagementAPIInstance.GetTenantExternalNetworkConnection(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching tenant external network connection: %v", err)
	}

	var extID, tenantID *string
	var config *import1.ExternalNetworkConnConfig
	var stateInfo *import1.ExternalNetworkConfigStateInfo
	var links interface{}

	switch data := resp.Data.GetValue().(type) {
	case import1.ExternalNetworkConnectionProjection:
		extID = data.ExtId
		tenantID = data.TenantId
		config = data.Config
		stateInfo = data.StateInfo
		links = data.Links
	case import1.ExternalNetworkConnection:
		extID = data.ExtId
		tenantID = data.TenantId
		config = data.Config
		stateInfo = data.StateInfo
		links = data.Links
	default:
		return diag.Errorf("unexpected response type while fetching tenant external network connection")
	}

	if err := d.Set("ext_id", extID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config", flattenExternalNetworkConnConfig(config)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state_info", flattenExternalNetworkConfigStateInfo(stateInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", tenantID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tenantExtID + "/" + domainAllocationExtID)
	return nil
}
