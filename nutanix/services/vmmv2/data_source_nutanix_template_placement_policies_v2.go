package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vmmConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/templateplacementpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixTemplatePlacementPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixTemplatePlacementPoliciesV2Read,
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
			"template_placement_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixTemplatePlacementPolicyV2(),
			},
		},
	}
}

func datasourceNutanixTemplatePlacementPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	req := import3.ListTemplatePlacementPoliciesRequest{}

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

	resp, err := conn.TemplatePlacementPoliciesAPIInstance.ListTemplatePlacementPolicies(ctx, &req)
	if err != nil {
		return diag.Errorf("error listing template placement policies: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("template_placement_policies", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found",
			Detail:   "The API returned an empty list of template placement policies.",
		}}
	}

	policies := resp.Data.GetValue().([]vmmConfig.TemplatePlacementPolicy)

	if err := d.Set("template_placement_policies", flattenTemplatePlacementPolicies(policies)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenTemplatePlacementPolicies(policies []vmmConfig.TemplatePlacementPolicy) []interface{} {
	if len(policies) == 0 {
		return nil
	}
	result := make([]interface{}, len(policies))
	for i, p := range policies {
		policy := make(map[string]interface{})

		if p.ExtId != nil {
			policy["ext_id"] = utils.StringValue(p.ExtId)
		}
		if p.Name != nil {
			policy["name"] = utils.StringValue(p.Name)
		}
		if p.Description != nil {
			policy["description"] = utils.StringValue(p.Description)
		}
		if p.PlacementType != nil {
			policy["placement_type"] = common.FlattenPtrEnum(p.PlacementType)
		}
		if p.ClusterFilter != nil {
			policy["cluster_filter"] = flattenCategoriesFilter(p.ClusterFilter)
		}
		if p.ContentFilter != nil {
			policy["content_filter"] = flattenCategoriesFilter(p.ContentFilter)
		}
		if p.CreateTime != nil {
			policy["create_time"] = p.CreateTime.String()
		}
		if p.CreatedBy != nil {
			policy["created_by"] = utils.StringValue(p.CreatedBy)
		}
		if p.UpdateTime != nil {
			policy["update_time"] = p.UpdateTime.String()
		}
		if p.UpdatedBy != nil {
			policy["updated_by"] = utils.StringValue(p.UpdatedBy)
		}
		if p.TenantId != nil {
			policy["tenant_id"] = utils.StringValue(p.TenantId)
		}
		result[i] = policy
	}
	return result
}
