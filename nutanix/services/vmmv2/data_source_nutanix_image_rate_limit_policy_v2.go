package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/images/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/imageratelimitpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixImageRateLimitPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixImageRateLimitPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The external identifier of image rate limit policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the image rate limit policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image rate limit policy specification.",
			},
			"rate_limit_kbps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network bandwidth in KBps that the rate limited image operation can utilize.",
			},
			"cluster_entity_filter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"matching_cluster_ext_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "External identifier of the Prism Elements where a rate limit is the effective rate limit policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External identifier of the owner of the rate limit policy.",
			},
			"owner_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the owner of the rate limit policy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image rate limit policy creation time.",
			},
			"last_update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time of an image rate limit policy.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
		},
	}
}

func datasourceNutanixImageRateLimitPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)
	req := import3.GetRateLimitPolicyByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.ImageRateLimitPoliciesAPIInstance.GetRateLimitPolicyById(ctx, &req)
	if err != nil {
		return diag.Errorf("error reading image rate limit policy: %v", err)
	}

	policy := resp.Data.GetValue().(import7.RateLimitPolicy)

	if err := d.Set("name", policy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", policy.Description); err != nil {
		return diag.FromErr(err)
	}
	if policy.RateLimitKbps != nil {
		if err := d.Set("rate_limit_kbps", int(*policy.RateLimitKbps)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("cluster_entity_filter", flattenRateLimitClusterEntityFilter(policy.ClusterEntityFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("matching_cluster_ext_ids", utils.StringSlice(policy.MatchingClusterExtIds)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", policy.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_name", policy.OwnerName); err != nil {
		return diag.FromErr(err)
	}
	if policy.CreateTime != nil {
		if err := d.Set("create_time", policy.CreateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if policy.LastUpdateTime != nil {
		if err := d.Set("last_update_time", policy.LastUpdateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("tenant_id", policy.TenantId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(policy.ExtId))
	return nil
}
