package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iamResponse "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/common/v1/response"
	iamConfig "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEntityV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixEntityV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Description: "ExtId for the Entity.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Description: "Tenant ID for the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"links": {
				Description: "A HATEOAS style link for the response.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Description: "Name of the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"display_name": {
				Description: "Display name for the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"client_name": {
				Description: "Client that created the entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"search_url": {
				Description: "Search URL for the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_time": {
				Description: "Creation time of the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated_time": {
				Description: "Last updated time of the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_by": {
				Description: "User or Service that created the Entity.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attribute_list": {
				Description: "List of attributes for the Entity (used in authorization policy filters). Each item may include $reserved fields: acceptedValues, displayName, supportedOperators, uiDisplayName.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rel": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_operator": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"attribute_values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"is_logical_and_supported_for_attributes": {
				Description: "Whether logical AND is supported for attributes.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func DatasourceNutanixEntityV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.EntityAPIInstance.GetEntityById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching entity: %v", err)
	}

	if resp.Data == nil {
		return diag.Errorf("no data returned for entity %s", extID)
	}

	entityVal := resp.Data.GetValue()
	if entityVal == nil {
		return diag.Errorf("no entity data for ext_id %s", extID)
	}

	getResp, ok := entityVal.(iamConfig.Entity)
	if !ok {
		return diag.Errorf("unexpected entity response type for ext_id %s", extID)
	}

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenEntityLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_name", getResp.DisplayName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_name", getResp.ClientName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("search_url", getResp.SearchURL); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreatedTime != nil {
		if err := d.Set("created_time", utils.TimeStringValue(getResp.CreatedTime)); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdatedTime != nil {
		if err := d.Set("last_updated_time", utils.TimeStringValue(getResp.LastUpdatedTime)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("attribute_list", flattenAttributeList(getResp.AttributeList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_logical_and_supported_for_attributes", getResp.IsLogicalAndSupportedForAttributes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}

func flattenEntityLinks(apiLinks []iamResponse.ApiLink) []map[string]interface{} {
	if len(apiLinks) == 0 {
		return nil
	}
	linkList := make([]map[string]interface{}, len(apiLinks))
	for i, v := range apiLinks {
		link := map[string]interface{}{}
		if v.Href != nil {
			link["href"] = utils.StringValue(v.Href)
		}
		if v.Rel != nil {
			link["rel"] = utils.StringValue(v.Rel)
		}
		linkList[i] = link
	}
	return linkList
}

func flattenAttributeList(attrList []iamConfig.AttributeEntity) []map[string]interface{} {
	if len(attrList) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(attrList))
	for _, attr := range attrList {
		displayName := utils.StringValue(attr.DisplayName)
		supportedOps := common.FlattenEnumValueList(attr.SupportedOperators)

		if attr.Reserved_ != nil {
			if v, ok := attr.Reserved_["displayName"]; ok && v != nil && displayName == "" {
				if s, ok := v.(string); ok {
					displayName = s
				}
			}
			if v, ok := attr.Reserved_["supportedOperators"]; ok && v != nil && len(supportedOps) == 0 {
				if sl, ok := v.([]interface{}); ok {
					for _, item := range sl {
						if s, ok := item.(string); ok {
							supportedOps = append(supportedOps, s)
						}
					}
				}
			}
		}

		attrValues := attr.AttributeValues
		if attrValues == nil {
			attrValues = []string{}
		}
		m := map[string]interface{}{
			"tenant_id":          utils.StringValue(attr.TenantId),
			"ext_id":             utils.StringValue(attr.ExtId),
			"links":              flattenEntityLinks(attr.Links),
			"display_name":       displayName,
			"name":               utils.StringValue(attr.Name),
			"supported_operator": supportedOps,
			"attribute_values":   attrValues,
		}
		result = append(result, m)
	}
	return result
}
