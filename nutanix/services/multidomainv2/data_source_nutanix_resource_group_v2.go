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

func DatasourceNutanixResourceGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixResourceGroupV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Elem:     schemaDatasourceResourceGroupPlacementTargets(),
			},
			"links": schemaForLinks(),
		},
	}
}

func schemaDatasourceResourceGroupPlacementTargets() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_containers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixResourceGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	extID := d.Get("ext_id").(string)
	getReq := import1.GetResourceGroupByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.ResourceGroups.GetResourceGroupById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching ResourceGroup: %s", err)
	}

	rg := resp.Data.GetValue().(config.ResourceGroup)
	if err := d.Set("name", utils.StringValue(rg.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", utils.StringValue(rg.ProjectExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", utils.StringValue(rg.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(rg.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", utils.StringValue(rg.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated_by", utils.StringValue(rg.LastUpdatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if rg.CreateTime != nil {
		if err := d.Set("create_time", rg.CreateTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
			return diag.FromErr(err)
		}
	}
	if rg.LastUpdateTime != nil {
		if err := d.Set("last_update_time", rg.LastUpdateTime.Format("2006-01-02T15:04:05Z07:00")); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("placement_targets", flattenResourceGroupPlacementTargets(rg.PlacementTargets)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(rg.Links)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(rg.ExtId))
	return nil
}

func flattenResourceGroupPlacementTargets(targets []config.TargetDetails) []map[string]interface{} {
	if len(targets) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(targets))
	for _, t := range targets {
		m := map[string]interface{}{
			"cluster_ext_id":     utils.StringValue(t.ClusterExtId),
			"storage_containers": flattenResourceGroupStorageContainers(t.StorageContainers),
		}
		out = append(out, m)
	}
	return out
}

func flattenResourceGroupStorageContainers(containers []config.StorageContainerDetails) []map[string]interface{} {
	if len(containers) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(containers))
	for _, c := range containers {
		out = append(out, map[string]interface{}{
			"ext_id": utils.StringValue(c.ExtId),
		})
	}
	return out
}
