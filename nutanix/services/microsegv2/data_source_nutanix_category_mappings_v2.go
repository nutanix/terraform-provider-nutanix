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

// DatasourceNutanixCategoryMappingsV2 lists the Directory Server Category
// Mappings.
func DatasourceNutanixCategoryMappingsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixCategoryMappingsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: datasourceCategoryMappingSchema(),
				},
			},
		},
	}
}

func DatasourceNutanixCategoryMappingsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	listRequest := import3.ListCategoryMappingsRequest{}
	if v, ok := d.GetOk("page"); ok {
		listRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.DirectoryServerConfigsAPIInstance.ListCategoryMappings(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while listing Category Mappings: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("category_mappings", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of Category Mappings.",
		}}
	}

	listVal, ok := resp.Data.GetValue().([]import2.CategoryMapping)
	if !ok {
		if err := d.Set("category_mappings", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of Category Mappings.",
		}}
	}

	if err := d.Set("category_mappings", flattenCategoryMappings(listVal)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenCategoryMappings(mappings []import2.CategoryMapping) []map[string]interface{} {
	if len(mappings) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(mappings))
	for _, c := range mappings {
		m := map[string]interface{}{
			"ext_id":         utils.StringValue(c.ExtId),
			"name":           utils.StringValue(c.Name),
			"category_name":  utils.StringValue(c.CategoryName),
			"category_value": utils.StringValue(c.CategoryValue),
			"ad_info":        flattenAdInfo(c.AdInfo),
			"project_ext_id": utils.StringValue(c.ProjectExtId),
			"tenant_id":      utils.StringValue(c.TenantId),
			"links":          common.FlattenLinks(c.Links),
		}
		result = append(result, m)
	}
	return result
}
