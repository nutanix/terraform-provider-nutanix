package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/common/v1/response"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/projects"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixProjectV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixProjectV2Read,
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
	}
}

func DatasourceNutanixProjectV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	extID := d.Get("ext_id").(string)
	getProjectByIdRequest := import1.GetProjectByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.Projects.GetProjectById(ctx, &getProjectByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching Project: %s", err)
	}

	if resp.Data == nil {
		return diag.Errorf("no project data in response")
	}

	project, ok := resp.Data.GetValue().(config.Project)
	if !ok {
		return diag.Errorf("error parsing GetProjectById response data")
	}

	_ = d.Set("name", utils.StringValue(project.Name))
	_ = d.Set("description", utils.StringValue(project.Description))
	_ = d.Set("tenant_id", utils.StringValue(project.TenantId))
	_ = d.Set("state", utils.StringValue(project.State))
	_ = d.Set("is_default", utils.BoolValue(project.IsDefault))
	_ = d.Set("is_system_defined", utils.BoolValue(project.IsSystemDefined))
	_ = d.Set("created_by", utils.StringValue(project.CreatedBy))
	_ = d.Set("updated_by", utils.StringValue(project.UpdatedBy))
	_ = d.Set("created_timestamp", utils.Int64Value(project.CreatedTimestamp))
	_ = d.Set("modified_timestamp", utils.Int64Value(project.ModifiedTimestamp))
	_ = d.Set("links", flattenLinks(project.Links))

	d.SetId(utils.StringValue(project.ExtId))
	return nil
}

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) > 0 {
		linkList := make([]map[string]interface{}, 0, len(links))
		for _, link := range links {
			linkMap := make(map[string]interface{})
			if link.Rel != nil {
				linkMap["rel"] = utils.StringValue(link.Rel)
			}
			if link.Href != nil {
				linkMap["href"] = utils.StringValue(link.Href)
			}
			linkList = append(linkList, linkMap)
		}
		return linkList
	}
	return []map[string]interface{}{}
}
