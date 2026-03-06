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
				Elem:     schemaResourceGroupPlacementTargets(),
			},
			"links": schemaForLinks(),
		},
	}
}

func schemaResourceGroupPlacementTargets() *schema.Resource {
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

	if resp.Data == nil {
		return diag.Errorf("no resource group data in response")
	}

	var rg config.ResourceGroup
	switch v := resp.Data.GetValue().(type) {
	case config.ResourceGroup:
		rg = v
	case *config.ResourceGroup:
		if v != nil {
			rg = *v
		} else {
			return diag.Errorf("error parsing GetResourceGroupById response data")
		}
	default:
		return diag.Errorf("error parsing GetResourceGroupById response data")
	}

	_ = d.Set("name", utils.StringValue(rg.Name))
	_ = d.Set("project_ext_id", utils.StringValue(rg.ProjectExtId))
	_ = d.Set("tenant_id", utils.StringValue(rg.TenantId))
	_ = d.Set("created_by", utils.StringValue(rg.CreatedBy))
	_ = d.Set("last_updated_by", utils.StringValue(rg.LastUpdatedBy))
	if rg.CreateTime != nil {
		_ = d.Set("create_time", rg.CreateTime.Format("2006-01-02T15:04:05Z07:00"))
	}
	if rg.LastUpdateTime != nil {
		_ = d.Set("last_update_time", rg.LastUpdateTime.Format("2006-01-02T15:04:05Z07:00"))
	}
	_ = d.Set("placement_targets", flattenResourceGroupPlacementTargets(rg.PlacementTargets))
	_ = d.Set("links", flattenLinks(rg.Links))

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
