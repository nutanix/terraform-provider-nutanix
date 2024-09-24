package nutanix

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
)

func datsourceNutanixCalmApp() *schema.Resource {
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
						// "categories":   {
						// 	Type:    schema.TypeMap,
						// 	Computed: true,
						// },
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
						// "provider": {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// "environment": {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
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
			// "runtime_editables": {},
		},
	}
}

func datsourceNutanixCalmAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Calm

	appID := d.Get("app_uuid").(string)
	resp, err := conn.Service.GetApp(ctx, appID)
	if err != nil {
		return diag.FromErr(err)
	}

	AppResp := &calm.AppResponse{}
	if err := json.Unmarshal([]byte(resp.Status), &AppResp.Status); err != nil {
		fmt.Println("Error unmarshalling App:", err)
	}
	if specErr := json.Unmarshal([]byte(resp.Spec), &AppResp.Spec); specErr != nil {
		fmt.Println("Error unmarshalling App:", specErr)
	}

	// Convert JSON object to JSON string
	// jsonData, err := json.MarshalIndent(AppResp.Status, " ", "  ")
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// Set the JSON data as a string in the Terraform resource data
	// if err := d.Set("status", string(jsonData)); err != nil {
	// 	return diag.FromErr(err)
	// }

	// Convert JSON object to JSON string
	// jsonSpecData, err := json.MarshalIndent(AppResp.Spec, " ", "  ")
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// // Set the JSON data as a string in the Terraform resource data
	// if err := d.Set("spec", string(jsonSpecData)); err != nil {
	// 	return diag.FromErr(err)
	// }

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	// unMarshall to get state of an APP
	var objStatus map[string]interface{}
	if err := json.Unmarshal(resp.Status, &objStatus); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}

	var objMetadata map[string]interface{}
	if err := json.Unmarshal(resp.Metadata, &objMetadata); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}
	var app_state string

	if state, ok := objStatus["state"].(string); ok {
		app_state = state
		// fmt.Printf("State of APPP: %s\n", state)
	}

	if err := d.Set("state", app_state); err != nil {
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

	d.SetId(appID)
	return nil
}

func flattenAppSummary(pr map[string]interface{}, meta map[string]interface{}) []interface{} {
	appSummaryMap := make(map[string]interface{})
	appSummaryList := make([]interface{}, 0)

	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		// Access the list "app_profile"
		if app_profile, ok := resource["app_profile_config_reference"].(map[string]interface{}); ok {
			appSummaryMap["application_profile"] = app_profile["name"]
		}
		if bp_reference, ok := resource["app_blueprint_reference"].(map[string]interface{}); ok {
			appSummaryMap["blueprint"] = bp_reference["name"]
		}
	}
	if project, ok := meta["project_reference"].(map[string]interface{}); ok {
		appSummaryMap["project"] = project["name"]
	}
	if owner, ok := meta["owner_reference"].(map[string]interface{}); ok {
		appSummaryMap["owner"] = owner["name"]
	}
	if created_on, ok := meta["creation_time"].(string); ok {
		appSummaryMap["created_on"] = created_on
	}
	if last_updated_on, ok := meta["last_update_time"].(string); ok {
		appSummaryMap["last_updated_on"] = last_updated_on
	}
	if appUUUID, ok := meta["uuid"].(string); ok {
		appSummaryMap["application_uuid"] = appUUUID
	}

	appSummaryList = append(appSummaryList, appSummaryMap)
	return appSummaryList
}
