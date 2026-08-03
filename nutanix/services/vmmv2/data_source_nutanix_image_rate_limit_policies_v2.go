package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/images/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/imageratelimitpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixImageRateLimitPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixImageRateLimitPoliciesV2Read,
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
			"rate_limit_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixImageRateLimitPolicyV2(),
			},
		},
	}
}

func datasourceNutanixImageRateLimitPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	req := import3.ListRateLimitPoliciesRequest{}

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

	resp, err := conn.ImageRateLimitPoliciesAPIInstance.ListRateLimitPolicies(ctx, &req)
	if err != nil {
		return diag.Errorf("error listing image rate limit policies: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("rate_limit_policies", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of image rate limit policies.",
		}}
	}

	policies := resp.Data.GetValue().([]import7.RateLimitPolicy)

	if err := d.Set("rate_limit_policies", flattenRateLimitPolicies(policies)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRateLimitPolicies(policies []import7.RateLimitPolicy) []interface{} {
	if len(policies) == 0 {
		return nil
	}
	result := make([]interface{}, len(policies))
	for i, policy := range policies {
		p := make(map[string]interface{})
		if policy.ExtId != nil {
			p["ext_id"] = utils.StringValue(policy.ExtId)
		}
		if policy.Name != nil {
			p["name"] = utils.StringValue(policy.Name)
		}
		if policy.Description != nil {
			p["description"] = utils.StringValue(policy.Description)
		}
		if policy.RateLimitKbps != nil {
			p["rate_limit_kbps"] = int(*policy.RateLimitKbps)
		}
		if policy.ClusterEntityFilter != nil {
			p["cluster_entity_filter"] = flattenRateLimitClusterEntityFilter(policy.ClusterEntityFilter)
		}
		if policy.MatchingClusterExtIds != nil {
			p["matching_cluster_ext_ids"] = utils.StringSlice(policy.MatchingClusterExtIds)
		}
		if policy.OwnerExtId != nil {
			p["owner_ext_id"] = utils.StringValue(policy.OwnerExtId)
		}
		if policy.OwnerName != nil {
			p["owner_name"] = utils.StringValue(policy.OwnerName)
		}
		if policy.CreateTime != nil {
			p["create_time"] = policy.CreateTime.String()
		}
		if policy.LastUpdateTime != nil {
			p["last_update_time"] = policy.LastUpdateTime.String()
		}
		if policy.TenantId != nil {
			p["tenant_id"] = utils.StringValue(policy.TenantId)
		}
		result[i] = p
	}
	return result
}
