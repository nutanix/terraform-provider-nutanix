package iamv2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/scopetemplates"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixScopeTemplateV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixScopeTemplateV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "External identifier of the scope template.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"display_name": {
				Description: "The display name for the scope template.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the scope template.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"entities": {
				Description: "List of entities being scoped for the template.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_by": {
				Description: "Service name that created the scope template.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_time": {
				Description: "The creation time of the scope template.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixScopeTemplateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id").(string)
	getRequest := import1.GetScopeTemplateByIdRequest{
		ExtId: utils.StringPtr(extID),
	}

	resp, err := conn.ScopeTemplatesAPIInstance.GetScopeTemplateById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching scope template: %v", err)
	}

	getResp := resp.Data.GetValue().(iamConfig.ScopeTemplate)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_name", getResp.DisplayName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entities", flattenScopeTemplateEntityFilters(getResp.Entities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", getResp.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenScopeTemplateEntityFilters(entities []iamConfig.EntityFilter) []map[string]interface{} {
	if len(entities) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		entityMap := map[string]interface{}{}
		if entity.EntityFilter != nil {
			b, err := json.Marshal(*entity.EntityFilter)
			if err != nil {
				entityMap["entity_filter"] = fmt.Sprintf("%v", *entity.EntityFilter)
			} else {
				entityMap["entity_filter"] = string(b)
			}
		}
		result[i] = entityMap
	}
	return result
}
