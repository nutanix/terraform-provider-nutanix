package selfservice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
)

func DatsourceNutanixSelfServiceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: datsourceNutanixSelfServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,	
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}


func datsourceNutanixSelfServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(*conns.Client).CalmAPI
  
	acFilter := &selfservice.AccountsListInput{}
	if filter, ok := d.GetOk("filter"); ok {
		acFilter.Filter = filter.(string)
	}

	accountResp, err := conn.Service.ListAccounts(ctx, acFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(accountResp.Entities) == 0 {
		if err := d.Set("accounts", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resource.UniqueId())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No Data found",
			Detail:   "The API returned an empty list of Accounts.",
		}}
	}

	if err := d.Set("api_version", accountResp.APIVersion); err != nil {
		return diag.FromErr(err)
	}
  
	if err := d.Set("accounts", flattenAccountEntities(accountResp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenAccountEntities(entities []map[string]interface{}) []map[string]interface{} {
	EntityList := make([]map[string]interface{}, 0)
	for _, entity := range entities {
		EntityMap := make(map[string]interface{})
		if status, ok := entity["status"].(map[string]interface{}); ok {
			// Dig into resources -> data -> server, port
			if resources, ok := status["resources"].(map[string]interface{}); ok {
				if typeName, ok := resources["type"].(string); ok {
					if typeName == "custom_provider"{
						continue // Skip custom provider types
					}
					EntityMap["type"] = typeName
				}
				if state, ok := resources["state"].(string); ok {
					EntityMap["state"] = state
				}
				if data, ok := resources["data"].(map[string]interface{}); ok {
					if server, ok := data["server"].(string); ok {
						EntityMap["server"] = server
					}
					if port, ok := data["port"].(float64); ok {
						EntityMap["port"] = port
					}
				}
			}
			if description, ok := status["description"].(string); ok {
				EntityMap["description"] = description
			}
		}
		if metadata, ok := entity["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				EntityMap["name"] = name
			}
			if uuid, ok := metadata["uuid"].(string); ok {
				EntityMap["uuid"] = uuid
			}
		}
		EntityList = append(EntityList, EntityMap)
	}
	return EntityList
}


