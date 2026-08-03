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

func DatasourceNutanixEffectiveImageRateLimitPoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixEffectiveImageRateLimitPoliciesV2Read,
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
			"effective_rate_limit_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
						},
						"cluster_ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster external identifier.",
						},
						"rate_limit_ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The external identifier of image rate limit policy.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A globally unique identifier that represents the tenant that owns this entity.",
						},
					},
				},
			},
		},
	}
}

func datasourceNutanixEffectiveImageRateLimitPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	req := import3.ListEffectiveRateLimitPoliciesRequest{}

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

	resp, err := conn.ImageRateLimitPoliciesAPIInstance.ListEffectiveRateLimitPolicies(ctx, &req)
	if err != nil {
		return diag.Errorf("error listing effective image rate limit policies: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("effective_rate_limit_policies", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return nil
	}

	policies := resp.Data.GetValue().([]import7.EffectiveRateLimitPolicy)

	if err := d.Set("effective_rate_limit_policies", flattenEffectiveRateLimitPolicies(policies)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenEffectiveRateLimitPolicies(policies []import7.EffectiveRateLimitPolicy) []interface{} {
	if len(policies) == 0 {
		return nil
	}
	result := make([]interface{}, len(policies))
	for i, policy := range policies {
		p := make(map[string]interface{})
		if policy.ExtId != nil {
			p["ext_id"] = utils.StringValue(policy.ExtId)
		}
		if policy.ClusterExtId != nil {
			p["cluster_ext_id"] = utils.StringValue(policy.ClusterExtId)
		}
		if policy.RateLimitExtId != nil {
			p["rate_limit_ext_id"] = utils.StringValue(policy.RateLimitExtId)
		}
		if policy.TenantId != nil {
			p["tenant_id"] = utils.StringValue(policy.TenantId)
		}
		result[i] = p
	}
	return result
}
