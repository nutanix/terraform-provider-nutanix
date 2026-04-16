package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authz"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/scopetemplates"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixScopeTemplatesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixScopeTemplatesV2Read,
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
			"scope_templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func DatasourceNutanixScopeTemplatesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	listRequest := import1.ListScopeTemplatesRequest{}
	if v, ok := d.GetOk("filter"); ok {
		listRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.ScopeTemplatesAPIInstance.ListScopeTemplates(ctx, &listRequest)
	if err != nil {
		return diag.Errorf("error while fetching scope templates: %v", err)
	}

	scopeTemplatesRaw := resp.Data.GetValue()
	scopeTemplatesList, ok := scopeTemplatesRaw.([]iamConfig.ScopeTemplate)
	if !ok || len(scopeTemplatesList) == 0 {
		if err := d.Set("scope_templates", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of scope templates.",
		}}
	}

	if err := d.Set("scope_templates", flattenScopeTemplatesEntities(scopeTemplatesList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenScopeTemplatesEntities(scopeTemplates []iamConfig.ScopeTemplate) []interface{} {
	if len(scopeTemplates) == 0 {
		return nil
	}

	result := make([]interface{}, len(scopeTemplates))

	for k, v := range scopeTemplates {
		st := make(map[string]interface{})

		if v.ExtId != nil {
			st["ext_id"] = v.ExtId
		}
		if v.TenantId != nil {
			st["tenant_id"] = v.TenantId
		}
		if v.Links != nil {
			st["links"] = flattenLinks(v.Links)
		}
		if v.DisplayName != nil {
			st["display_name"] = v.DisplayName
		}
		if v.Description != nil {
			st["description"] = v.Description
		}
		if v.Entities != nil {
			st["entities"] = flattenScopeTemplateEntityFilters(v.Entities)
		}
		if v.CreatedBy != nil {
			st["created_by"] = v.CreatedBy
		}
		if v.CreatedTime != nil {
			st["created_time"] = v.CreatedTime.String()
		}

		result[k] = st
	}
	return result
}
