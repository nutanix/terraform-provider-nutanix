package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/resourcegroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixResourceGroupsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixResourceGroupsV2Read,
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
			"resource_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: DatasourceNutanixResourceGroupV2(),
			},
		},
	}
}

func DatasourceNutanixResourceGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	listReq := import1.ListResourceGroupsRequest{}
	if v, ok := d.GetOk("page"); ok {
		listReq.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listReq.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listReq.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listReq.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listReq.Select_ = utils.StringPtr(v.(string))
	}
	resp, err := conn.ResourceGroups.ListResourceGroups(ctx, &listReq)
	if err != nil {
		return diag.Errorf("error while listing ResourceGroups: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("resource_groups", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of resource groups.",
		}}
	}
  
	resourceGroups := resp.Data.GetValue().([]config.ResourceGroup)
	if err := d.Set("resource_groups", flattenResourceGroups(resourceGroups)); err != nil {
		return diag.Errorf("error setting resource_groups: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenResourceGroups(resourceGroups []config.ResourceGroup) []map[string]interface{} {
	if len(resourceGroups) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(resourceGroups))
	for _, rg := range resourceGroups {
		m := map[string]interface{}{
			"ext_id":             utils.StringValue(rg.ExtId),
			"name":               utils.StringValue(rg.Name),
			"project_ext_id":     utils.StringValue(rg.ProjectExtId),
			"tenant_id":          utils.StringValue(rg.TenantId),
			"created_by":         utils.StringValue(rg.CreatedBy),
			"last_updated_by":    utils.StringValue(rg.LastUpdatedBy),
			"placement_targets":  flattenResourceGroupPlacementTargetsIncludingCapabilities(rg.PlacementTargets),
			"links":              flattenLinks(rg.Links),
			"capabilities":       flattenResourceGroupCapabilities(rg.Capabilities),
		}
		if rg.CreateTime != nil {
			m["create_time"] = rg.CreateTime.Format("2006-01-02T15:04:05Z07:00")
		}
		if rg.LastUpdateTime != nil {
			m["last_update_time"] = rg.LastUpdateTime.Format("2006-01-02T15:04:05Z07:00")
		}
		out = append(out, m)
	}
	return out
}
