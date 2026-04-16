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
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
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
	conn := meta.(*conns.Client).CalmAPI

	var appUUID string

	appName := d.Get("app_name").(string)

	appFilter := &selfservice.ApplicationListInput{}

	appFilter.Filter = fmt.Sprintf("name==%s;_state!=deleted", appName)

	log.Printf("[Debug] Qeurying apps/list API with filter %s", appFilter)

	appNameResp, err := conn.Service.ListApplication(ctx, appFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[Debug] Getting app uuid from app response: %s", appNameResp)

	var AppNameStatus []interface{}
	if err = json.Unmarshal([]byte(appNameResp.Entities), &AppNameStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling AppName:", err)
		return diag.FromErr(err)
	}

	entities := AppNameStatus[0].(map[string]interface{})

	if entity, ok := entities["metadata"].(map[string]interface{}); ok {
		appUUID = entity["uuid"].(string)
	}

	if appUUIDRead, ok := d.GetOk("app_uuid"); ok {
		appUUID = appUUIDRead.(string)
	}

	length := d.Get("length").(int)
	offset := d.Get("offset").(int)
	appResp, err := conn.Service.GetApp(ctx, appUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	var appStatus map[string]interface{}
	if err = json.Unmarshal(appResp.Status, &appStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec to get status:", err)
		return diag.FromErr(err)
	}

	substrateReference := fetchSubstrateReference(appStatus)

	currTime := strconv.FormatInt(time.Now().Unix(), 10)

	listInput := &selfservice.RecoveryPointsListInput{}

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
			if actionName, ok := status["action_name"].(string); ok {
				EntityMap["action_name"] = actionName
			}
			if recoveryPointInfoList, ok := status["recovery_point_info_list"].([]interface{}); ok {
				EntityMap["recovery_point_info_list"] = recoveryPointInfoList
			}
		}

		if apiVersion, ok := entity["api_version"].(string); ok {
			EntityMap["api_version"] = apiVersion
		}

		if spec, ok := entity["spec"].(map[string]interface{}); ok {
			EntityMap["spec"] = spec
		}

		if meta, ok := entity["metadata"].(map[string]interface{}); ok {
			if creationTime, ok := meta["creation_time"].(string); ok {
				EntityMap["creation_time"] = creationTime
			}
			if lastUpdateTime, ok := meta["last_update_time"].(string); ok {
				EntityMap["last_update_time"] = lastUpdateTime
			}
			if specVersion, ok := meta["spec_version"].(int); ok {
				EntityMap["spec_version"] = specVersion
			}
			if kind, ok := meta["kind"].(string); ok {
				EntityMap["kind"] = kind
			}
		}

		EntityList = append(EntityList, EntityMap)
	}

	return EntityList
}
