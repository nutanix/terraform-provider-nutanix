package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/request/directoryserverconfigs"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixDirectoryServerConfigV2 reads a single Directory Server
// configuration by its external identifier.
func DatasourceNutanixDirectoryServerConfigV2() *schema.Resource {
	dsSchema := datasourceDirectoryServerConfigSchema()
	dsSchema["ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: DatasourceNutanixDirectoryServerConfigV2Read,
		Schema:      dsSchema,
	}
}

func DatasourceNutanixDirectoryServerConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id").(string)
	getRequest := import3.GetDirectoryServerConfigByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.DirectoryServerConfigsAPIInstance.GetDirectoryServerConfigById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching Directory Server Config: %s", err)
	}

	if resp.Data == nil {
		return diag.Errorf("no data in GetDirectoryServerConfigById response")
	}

	getResp, ok := resp.Data.GetValue().(import2.DirectoryServerConfig)
	if !ok {
		return diag.Errorf("invalid DirectoryServerConfig type in response")
	}

	if err := d.Set("ext_id", utils.StringValue(getResp.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("directory_service_reference", utils.StringValue(getResp.DirectoryServiceReference)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_controllers", flattenIPAddressOrFQDNList(getResp.DomainControllers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("matching_criterias", flattenMatchingCriterias(getResp.MatchingCriterias)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_default_category_enabled", utils.BoolValue(getResp.IsDefaultCategoryEnabled)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("should_keep_default_category_on_login", utils.BoolValue(getResp.ShouldKeepDefaultCategoryOnLogin)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", utils.StringValue(getResp.ProjectExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
