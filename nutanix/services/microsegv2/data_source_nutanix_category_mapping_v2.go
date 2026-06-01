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

func DatasourceNutanixCategoryMappingV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixCategoryMappingV2Read,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ad_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForAdInfoComputed(),
			},
			"project_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixCategoryMappingV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id").(string)
	req := import3.GetDsCategoryMappingByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.DirectoryServerConfigsAPIInstance.GetDsCategoryMappingById(ctx, &req)
	if err != nil {
		return diag.Errorf("error while fetching Category Mapping: %s", err)
	}

	getResp := resp.Data.GetValue().(import2.CategoryMapping)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinksDSC(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_name", getResp.CategoryName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_value", getResp.CategoryValue); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ad_info", flattenAdInfo(getResp.AdInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", getResp.ProjectExtId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}
