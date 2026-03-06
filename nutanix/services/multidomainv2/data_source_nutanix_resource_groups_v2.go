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
			"resource_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_by": {
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
						"placement_targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaResourceGroupPlacementTargets(),
						},
						"links": schemaForLinks(),
					},
				},
			},
		},
	}
}

func DatasourceNutanixResourceGroupsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	listReq := import1.ListResourceGroupsRequest{}
	resp, err := conn.ResourceGroups.ListResourceGroups(ctx, &listReq)
	if err != nil {
		return diag.Errorf("error while listing ResourceGroups: %s", err)
	}

	if resp.Data == nil {
		_ = d.Set("resource_groups", []map[string]interface{}{})
		d.SetId(utils.GenUUID())
		return nil
	}

	projections, ok := resp.Data.GetValue().([]config.ResourceGroupProjection)
	if !ok {
		_ = d.Set("resource_groups", []map[string]interface{}{})
		d.SetId(utils.GenUUID())
		return nil
	}

	if err := d.Set("resource_groups", flattenResourceGroupProjections(projections)); err != nil {
		return diag.Errorf("error setting resource_groups: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenResourceGroupProjections(projections []config.ResourceGroupProjection) []map[string]interface{} {
	if len(projections) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(projections))
	for _, p := range projections {
		m := map[string]interface{}{
			"ext_id":             utils.StringValue(p.ExtId),
			"name":               utils.StringValue(p.Name),
			"project_ext_id":     utils.StringValue(p.ProjectExtId),
			"tenant_id":          utils.StringValue(p.TenantId),
			"created_by":         utils.StringValue(p.CreatedBy),
			"last_updated_by":    utils.StringValue(p.LastUpdatedBy),
			"placement_targets":  flattenResourceGroupPlacementTargets(p.PlacementTargets),
			"links":              flattenLinks(p.Links),
		}
		if p.CreateTime != nil {
			m["create_time"] = p.CreateTime.Format("2006-01-02T15:04:05Z07:00")
		}
		if p.LastUpdateTime != nil {
			m["last_update_time"] = p.LastUpdateTime.Format("2006-01-02T15:04:05Z07:00")
		}
		out = append(out, m)
	}
	return out
}
