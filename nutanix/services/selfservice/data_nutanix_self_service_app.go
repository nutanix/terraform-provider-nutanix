package selfservice

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
)

func DatsourceNutanixCalmApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: datsourceNutanixCalmAppRead,
		Schema: map[string]*schema.Schema{
			"app_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vcpus": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"cores": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"memory": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"vm_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"nics": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mac_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"categories": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"app_summary": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"blueprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_profile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"actions": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datsourceNutanixCalmAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).CalmAPI

	appID := d.Get("app_uuid").(string)
	resp, err := conn.Service.GetApp(ctx, appID)
	if err != nil {
		return diag.FromErr(err)
	}

	AppResp := &selfservice.AppResponse{}
	if err = json.Unmarshal([]byte(resp.Status), &AppResp.Status); err != nil {
		log.Println("[DEBUG] Error unmarshalling App:", err)
		return diag.FromErr(err)
	}
	if specErr := json.Unmarshal([]byte(resp.Spec), &AppResp.Spec); specErr != nil {
		log.Println("[DEBUG] Error unmarshalling App:", specErr)
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	// unMarshall to get state of an APP
	var objStatus map[string]interface{}
	if err := json.Unmarshal(resp.Status, &objStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}

	var objMetadata map[string]interface{}
	if err := json.Unmarshal(resp.Metadata, &objMetadata); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}
	var appState string

	if state, ok := objStatus["state"].(string); ok {
		appState = state
	}

	if err := d.Set("state", appState); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("app_name", objStatus["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("app_description", objStatus["description"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm", flattenVM(objStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_summary", flattenAppSummary(objStatus, objMetadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("actions", flattenActions(objStatus)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(appID)
	return nil
}

func flattenAppSummary(pr map[string]interface{}, meta map[string]interface{}) []interface{} {
	appSummaryMap := make(map[string]interface{})
	appSummaryList := make([]interface{}, 0)

	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		// Access the list "appProfile"
		if appProfile, ok := resource["app_profile_config_reference"].(map[string]interface{}); ok {
			appSummaryMap["application_profile"] = appProfile["name"]
		}
		if bpReference, ok := resource["app_blueprint_reference"].(map[string]interface{}); ok {
			appSummaryMap["blueprint"] = bpReference["name"]
		}
	}
	if project, ok := meta["project_reference"].(map[string]interface{}); ok {
		appSummaryMap["project"] = project["name"]
	}
	if owner, ok := meta["owner_reference"].(map[string]interface{}); ok {
		appSummaryMap["owner"] = owner["name"]
	}
	if createdOn, ok := meta["creation_time"].(string); ok {
		appSummaryMap["created_on"] = createdOn
	}
	if lastUpdatedOn, ok := meta["last_update_time"].(string); ok {
		appSummaryMap["last_updated_on"] = lastUpdatedOn
	}
	if appUUUID, ok := meta["uuid"].(string); ok {
		appSummaryMap["application_uuid"] = appUUUID
	}

	appSummaryList = append(appSummaryList, appSummaryMap)
	return appSummaryList
}

func flattenActions(pr map[string]interface{}) []interface{} {
	actionsOutput := make([]interface{}, 0)
	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		if actionsList, ok := resource["action_list"].([]interface{}); ok {
			for _, action := range actionsList {
				actionMap := make(map[string]interface{})
				if action, ok := action.(map[string]interface{}); ok {
					actionMap["name"] = func(parts []string) string {
						if parts[0] == "action" {
							return strings.Join(parts[1:], " ")
						}
						return strings.Join(parts, " ")
					}(strings.Split(action["name"].(string), "_"))
					actionMap["uuid"] = action["uuid"]
					actionMap["description"] = action["description"]
				}
				actionsOutput = append(actionsOutput, actionMap)
			}
			return actionsOutput
		}
	}
	return nil
}
