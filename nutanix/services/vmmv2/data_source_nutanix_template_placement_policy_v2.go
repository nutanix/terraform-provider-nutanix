package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vmmConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/templateplacementpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixTemplatePlacementPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixTemplatePlacementPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"placement_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_filter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"content_filter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceNutanixTemplatePlacementPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)
	req := import3.GetTemplatePlacementPolicyByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.TemplatePlacementPoliciesAPIInstance.GetTemplatePlacementPolicyById(ctx, &req)
	if err != nil {
		return diag.Errorf("error reading template placement policy: %v", err)
	}

	policy := resp.Data.GetValue().(vmmConfig.TemplatePlacementPolicy)

	if err := d.Set("name", policy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", policy.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("placement_type", common.FlattenPtrEnum(policy.PlacementType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_filter", flattenCategoriesFilter(policy.ClusterFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("content_filter", flattenCategoriesFilter(policy.ContentFilter)); err != nil {
		return diag.FromErr(err)
	}
	if policy.CreateTime != nil {
		if err := d.Set("create_time", policy.CreateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_by", policy.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if policy.UpdateTime != nil {
		if err := d.Set("update_time", policy.UpdateTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("updated_by", policy.UpdatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", policy.TenantId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(policy.ExtId))
	return nil
}
