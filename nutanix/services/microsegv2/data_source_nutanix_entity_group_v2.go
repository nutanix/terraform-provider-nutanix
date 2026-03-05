package microsegv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/response"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/entitygroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixEntityGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixEntityGroupV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"links": schemaForLinks(),
			"tenant_id": {
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
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaAllowedConfig(),
			},
		},
	}
}

func DatasourceNutanixEntityGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Get("ext_id")
	getEntityGroupByIdRequest := import3.GetEntityGroupByIdRequest{
		ExtId: utils.StringPtr(extID.(string)),
	}
	resp, err := conn.EntityGroupsAPIInstance.GetEntityGroupById(ctx, &getEntityGroupByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching Entity Group: %s", err)
	}

	getResp := resp.Data.GetValue().(import1.EntityGroup)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_ext_ids", getResp.PolicyExtIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allowed_config", flattenAllowedConfig(getResp.AllowedConfig)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}

// schema
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

func schemaAllowedConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaAllowedEntity(),
			},
		},
	}
}

func schemaAllowedEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kube_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reference_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"select_by": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem:     schemaAllowedSelectBy(),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaAllowedSelectBy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Add fields based on AllowedSelectBy struct if needed
			// This is a placeholder as the exact structure may vary
		},
	}
}

// flatten funcs
func flattenLinks(links []import2.ApiLink) []map[string]interface{} {
	if len(links) > 0 {
		linkList := make([]map[string]interface{}, 0)
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
	return nil
}

func flattenAllowedConfig(allowedConfig *import1.AllowedConfig) []map[string]interface{} {
	if allowedConfig != nil {
		allowedConfigMap := make(map[string]interface{})
		if allowedConfig.Entities != nil {
			allowedConfigMap["entities"] = flattenAllowedEntities(allowedConfig.Entities)
		}
		return []map[string]interface{}{allowedConfigMap}
	}
	return nil
}

func flattenAllowedEntities(entities []import1.AllowedEntity) []map[string]interface{} {
	if len(entities) > 0 {
		entityList := make([]map[string]interface{}, 0)
		for _, entity := range entities {
			entityMap := make(map[string]interface{})
			if entity.KubeEntities != nil {
				entityMap["kube_entities"] = entity.KubeEntities
			}
			if entity.ReferenceExtIds != nil {
				entityMap["reference_ext_ids"] = entity.ReferenceExtIds
			}
			if entity.SelectBy != nil {
				entityMap["select_by"] = flattenAllowedSelectBy(entity.SelectBy)
			}
			if entity.Type != nil {
				entityMap["type"] = utils.StringValue(entity.Type)
			}
			entityList = append(entityList, entityMap)
		}
		return entityList
	}
	return nil
}

func flattenAllowedSelectBy(selectBy *import1.AllowedSelectBy) []map[string]interface{} {
	if selectBy != nil {
		selectByMap := make(map[string]interface{})
		// Add fields based on AllowedSelectBy struct if needed
		return []map[string]interface{}{selectByMap}
	}
	return nil
}
