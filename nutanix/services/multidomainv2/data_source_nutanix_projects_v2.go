package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/projects"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixProjectsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixProjectsV2Read,
		Schema: map[string]*schema.Schema{
			"projects": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_system_defined": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"modified_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"links": schemaForLinks(),
					},
				},
			},
		},
	}
}

func DatasourceNutanixProjectsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	listProjectsRequest := import1.ListProjectsRequest{}
	resp, err := conn.Projects.ListProjects(ctx, &listProjectsRequest)
	if err != nil {
		return diag.Errorf("error while listing Projects: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("projects", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	projects, ok := resp.Data.GetValue().([]config.ProjectProjection)
	if !ok {
		if err := d.Set("projects", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	if err := d.Set("projects", flattenProjectProjections(projects)); err != nil {
		return diag.Errorf("error setting projects: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenProjectProjections(projects []config.ProjectProjection) []map[string]interface{} {
	if len(projects) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(projects))
	for _, p := range projects {
		m := map[string]interface{}{
			"ext_id":              utils.StringValue(p.ExtId),
			"name":                utils.StringValue(p.Name),
			"description":         utils.StringValue(p.Description),
			"tenant_id":           utils.StringValue(p.TenantId),
			"state":               utils.StringValue(p.State),
			"is_default":          utils.BoolValue(p.IsDefault),
			"is_system_defined":   utils.BoolValue(p.IsSystemDefined),
			"created_by":          utils.StringValue(p.CreatedBy),
			"updated_by":          utils.StringValue(p.UpdatedBy),
			"created_timestamp":   utils.Int64Value(p.CreatedTimestamp),
			"modified_timestamp":  utils.Int64Value(p.ModifiedTimestamp),
			"links":               flattenLinks(p.Links),
		}
		result = append(result, m)
	}
	return result
}
