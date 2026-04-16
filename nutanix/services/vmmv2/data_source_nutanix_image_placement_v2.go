package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import7 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/images/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixImagePlacementV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixImagePlacementV4Read,
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
			"image_entity_filter": {
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
			"enforcement_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixImagePlacementV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")

	resp, err := conn.ImagesPlacementAPIInstance.GetPlacementPolicyById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching image placement : %v", err)
	}

	getResp := resp.Data.GetValue().(import7.PlacementPolicy)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("placement_type", flattenPlacementType(getResp.PlacementType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("image_entity_filter", flattenEntityFilter(getResp.ImageEntityFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_entity_filter", flattenEntityFilter(getResp.ClusterEntityFilter)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enforcement_state", flattenEnforcementState(getResp.EnforcementState)); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdateTime != nil {
		t := getResp.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenPlacementType(pr *import7.PlacementType) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import7.PlacementType(two) {
			return "SOFT"
		}
		if *pr == import7.PlacementType(three) {
			return "HARD"
		}
	}
	return "UNKNOWN"
}

func flattenEnforcementState(pr *import7.EnforcementState) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import7.EnforcementState(two) {
			return "ACTIVE"
		}
		if *pr == import7.EnforcementState(three) {
			return "SUSPENDED"
		}
	}
	return "UNKNOWN"
}

func flattenEntityFilter(pr *import7.Filter) []map[string]interface{} {
	if pr != nil {
		filters := make([]map[string]interface{}, 0)
		filter := make(map[string]interface{})

		filter["type"] = flattenType(pr.Type)
		filter["category_ext_ids"] = utils.StringSlice(pr.CategoryExtIds)

		filters = append(filters, filter)
		return filters
	}
	return nil
}

func flattenType(pr *import7.FilterMatchType) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == import7.FilterMatchType(two) {
			return "ALL"
		}
		if *pr == import7.FilterMatchType(three) {
			return "ANY"
		}
	}
	return "UNKNOWN"
}
