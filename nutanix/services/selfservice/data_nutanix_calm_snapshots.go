package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func DataSourceNutanixCalmSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixCalmSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"action_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recovery_point_info_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"recovery_point_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expiration_time": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"location_agnostic_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_references": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"config_spec_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_matches": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixCalmSnapshotsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm

	var appUUID string

	app_name := d.Get("app_name").(string)

	appFilter := &calm.ApplicationListInput{}

	appFilter.Filter = fmt.Sprintf("name==%s;_state!=deleted", app_name)

	log.Printf("[Debug] Qeurying apps/list API with filter %s", appFilter)

	appNameResp, err := conn.Service.ListApplication(ctx, appFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[Debug] Getting app uuid from app response: %s", appNameResp)

	var AppNameStatus []interface{}
	if err := json.Unmarshal([]byte(appNameResp.Entities), &AppNameStatus); err != nil {
		fmt.Println("Error unmarshalling AppName:", err)
	}

	entities := AppNameStatus[0].(map[string]interface{})

	if entity, ok := entities["metadata"].(map[string]interface{}); ok {
		appUUID = entity["uuid"].(string)
	}

	if appUUID, ok := d.GetOk("app_uuid"); ok {
		appUUID = appUUID.(string)
	}

	length := d.Get("length").(int)
	offset := d.Get("offset").(int)
	appResp, err := conn.Service.GetApp(ctx, appUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	var appStatus map[string]interface{}
	if err := json.Unmarshal(appResp.Status, &appStatus); err != nil {
		fmt.Println("Error unmarshalling Spec to get status:", err)
	}

	substrateReference := fetchSubstrateReference(appStatus)

	currTime := strconv.FormatInt(time.Now().Unix(), 10)

	listInput := &calm.RecoveryPointsListInput{}

	listInput.Filter = fmt.Sprintf("substrate_reference==%s;expiration_time=ge=%s", substrateReference, currTime)
	listInput.Length = length
	listInput.Offset = offset

	listResp, err := conn.Service.RecoveryPointsList(ctx, appUUID, listInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", listResp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("total_matches", int(listResp.Metadata["total_matches"].(float64))); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("kind", listResp.Metadata["kind"].(string)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entities", flattenSnapshotEntities(listResp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func fetchSubstrateReference(appStatus map[string]interface{}) string {
	var substrateReference string
	if resources, ok := appStatus["resources"].(map[string]interface{}); ok {
		if depList, ok := resources["deployment_list"].([]interface{}); ok {
			if dep, ok := depList[0].(map[string]interface{}); ok {
				if substrateConfig, ok := dep["substrate_configuration"].(map[string]interface{}); ok {
					substrateReference = substrateConfig["uuid"].(string)
				}
			}
		}
	}
	return substrateReference
}

func flattenSnapshotEntities(entities []map[string]interface{}) []map[string]interface{} {
	EntityList := make([]map[string]interface{}, 0)
	for _, entity := range entities {
		EntityMap := make(map[string]interface{})
		if status, ok := entity["status"].(map[string]interface{}); ok {
			if vmType, ok := status["type"].(string); ok {
				EntityMap["type"] = vmType
			}
			if name, ok := status["name"].(string); ok {
				EntityMap["name"] = name
			}
			if uuid, ok := status["uuid"].(string); ok {
				EntityMap["uuid"] = uuid
			}
			if description, ok := status["description"].(string); ok {
				EntityMap["description"] = description
			}
			if action_name, ok := status["action_name"].(string); ok {
				EntityMap["action_name"] = action_name
			}
			if recovery_point_info_list, ok := status["recovery_point_info_list"].([]interface{}); ok {
				EntityMap["recovery_point_info_list"] = recovery_point_info_list
			}

		}

		if api_version, ok := entity["api_version"].(string); ok {
			EntityMap["api_version"] = api_version
		}

		if spec, ok := entity["spec"].(map[string]interface{}); ok {
			EntityMap["spec"] = spec
		}

		if meta, ok := entity["metadata"].(map[string]interface{}); ok {
			if creation_time, ok := meta["creation_time"].(string); ok {
				EntityMap["creation_time"] = creation_time
			}
			if last_update_time, ok := meta["last_update_time"].(string); ok {
				EntityMap["last_update_time"] = last_update_time
			}
			if spec_version, ok := meta["spec_version"].(int); ok {
				EntityMap["spec_version"] = spec_version
			}
			if kind, ok := meta["kind"].(string); ok {
				EntityMap["kind"] = kind
			}
		}

		EntityList = append(EntityList, EntityMap)
	}

	return EntityList
}
