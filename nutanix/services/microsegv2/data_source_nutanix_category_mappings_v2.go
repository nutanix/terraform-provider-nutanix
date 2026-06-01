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
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func DatasourceNutanixCategoryMappingsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	req := import3.ListCategoryMappingsRequest{}

	if v, ok := d.GetOk("page"); ok {
		req.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		req.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		req.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		req.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		req.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.DirectoryServerConfigsAPIInstance.ListCategoryMappings(ctx, &req)
	if err != nil {
		return diag.Errorf("error while listing Category Mappings: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("category_mappings", []map[string]interface{}{}); err != nil {
			return diag.Errorf("error setting Category Mappings: %s", err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of category mappings.",
		}}
	}

	getResp := resp.Data.GetValue().([]import2.CategoryMapping)

	if err := d.Set("category_mappings", flattenCategoryMappings(getResp)); err != nil {
		return diag.Errorf("error setting Category Mappings: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}
