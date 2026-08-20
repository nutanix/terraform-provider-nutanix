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

// DatasourceNutanixCategoryMappingV2 reads the category to directory
// configuration information by its external identifier.
func DatasourceNutanixCategoryMappingV2() *schema.Resource {
	dsSchema := datasourceCategoryMappingSchema()
	dsSchema["ext_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{
		ReadContext: DatasourceNutanixCategoryMappingV2Read,
		Schema:      dsSchema,
	}
}

func DatasourceNutanixCategoryMappingV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id").(string)
	getRequest := import3.GetDsCategoryMappingByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.DirectoryServerConfigsAPIInstance.GetDsCategoryMappingById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching Category Mapping: %s", err)
	}

	if resp.Data == nil {
		return diag.Errorf("no data in GetDsCategoryMappingById response")
	}

	getResp, ok := resp.Data.GetValue().(import2.CategoryMapping)
	if !ok {
		return diag.Errorf("invalid CategoryMapping type in response")
	}

	if err := d.Set("ext_id", utils.StringValue(getResp.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(getResp.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_name", utils.StringValue(getResp.CategoryName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_value", utils.StringValue(getResp.CategoryValue)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ad_info", flattenAdInfo(getResp.AdInfo)); err != nil {
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
