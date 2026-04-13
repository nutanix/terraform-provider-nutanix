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
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: DatasourceNutanixProjectV2(),
			},
		},
	}
}

func DatasourceNutanixProjectsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	listProjectsRequest := import1.ListProjectsRequest{}
	defaultFilter := "isSystemDefined ne true and isDefault ne true"
	filter := defaultFilter
	if v, ok := d.GetOk("filter"); ok {
		filter = v.(string) + " and " + defaultFilter
	}
	listProjectsRequest.Filter_ = utils.StringPtr(filter)
	if v, ok := d.GetOk("order_by"); ok {
		listProjectsRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listProjectsRequest.Select_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("limit"); ok {
		listProjectsRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("page"); ok {
		listProjectsRequest.Page_ = utils.IntPtr(v.(int))
	}
	resp, err := conn.Projects.ListProjects(ctx, &listProjectsRequest)
	if err != nil {
		return diag.Errorf("error while listing Projects: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("projects", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of projects.",
		}}
	}

	projects := resp.Data.GetValue().([]config.Project)
	if err := d.Set("projects", flattenProjects(projects)); err != nil {
		return diag.Errorf("error setting projects: %s", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenProjects(projects []config.Project) []map[string]interface{} {
	if len(projects) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(projects))
	for _, p := range projects {
		m := map[string]interface{}{
			"ext_id":              p.ExtId,
			"name":                p.Name,
			"description":         p.Description,
			"project_id":          p.Id,
			"tenant_id":           p.TenantId,
			"state":               p.State.GetName(),
			"is_default":          p.IsDefault,
			"is_system_defined":   p.IsSystemDefined,
			"created_by":          p.CreatedBy,
			"updated_by":          p.UpdatedBy,
			"created_timestamp":   p.CreatedTimestamp,
			"modified_timestamp":  p.ModifiedTimestamp,
			"links":               flattenLinks(p.Links),
		}
		result = append(result, m)
	}
	return result
}
