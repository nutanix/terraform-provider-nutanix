package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/directoryserverconfigs"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixDirectoryServerConfigV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixDirectoryServerConfigV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"links": schemaForLinks(),
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_service_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_controllers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForIPAddressOrFQDNComputed(),
			},
			"is_default_category_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"matching_criterias": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForMatchingCriteriaComputed(),
			},
			"should_keep_default_category_on_login": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixDirectoryServerConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id").(string)
	req := import3.GetDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.GetDirectoryServerConfigById(ctx, &req)
	if err != nil {
		return diag.Errorf("error while fetching Directory Server Config: %s", err)
	}

	getResp := resp.Data.GetValue().(import2.DirectoryServerConfig)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksDSC(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("directory_service_reference", getResp.DirectoryServiceReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_controllers", flattenDomainControllers(getResp.DomainControllers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_default_category_enabled", getResp.IsDefaultCategoryEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("matching_criterias", flattenMatchingCriterias(getResp.MatchingCriterias)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("should_keep_default_category_on_login", getResp.ShouldKeepDefaultCategoryOnLogin); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", getResp.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}
