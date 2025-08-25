package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	// ERROR ..
	ERROR              = "ERROR"
	DEFAULTWAITTIMEOUT = 60
)

func ResourceNutanixCalmAppProvision() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmAppProvisionCreate,
		ReadContext:   resourceNutanixCalmAppProvisionRead,
		UpdateContext: resourceNutanixCalmAppProvisionUpdate,
		DeleteContext: resourceNutanixCalmAppProvisionDelete,
		Schema: map[string]*schema.Schema{
			"bp_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bp_uuid"},
			},
			"bp_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bp_name"},
			},
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"stop", "start", "restart"}, false),
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
			"runtime_editables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"service_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"credential_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"substrate_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"package_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"snapshot_config_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"app_profile": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"task_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"restore_config_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"variable_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
						},
						"deployment_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpec(),
							},
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
			"soft_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNutanixCalmAppProvisionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).CalmAPI

	var bpUUID string
	// fetch bp_uuid from bp_name
	bpName := d.Get("bp_name").(string)

	bpFilter := &selfservice.BlueprintListInput{}

	bpFilter.Filter = fmt.Sprintf("name==%s;state!=DELETED", bpName)

	bpNameResp, err := conn.Service.ListBlueprint(ctx, bpFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	var BpNameStatus []interface{}
	if err = json.Unmarshal([]byte(bpNameResp.Entities), &BpNameStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling BPName:", err)
		return diag.FromErr(err)
	}

	entities := BpNameStatus[0].(map[string]interface{})

	if entity, ok := entities["metadata"].(map[string]interface{}); ok {
		bpUUID = entity["uuid"].(string)
	}

	if bpUUIDRead, ok := d.GetOk("bp_uuid"); ok {
		bpUUID = bpUUIDRead.(string)
	}

	// call bp

	bpOut, er := conn.Service.GetBlueprint(ctx, bpUUID)
	if er != nil {
		return diag.Errorf("Error getting Blueprint: %s", er)
	}

	bpResp := &selfservice.BlueprintResponse{}
	if err = json.Unmarshal([]byte(bpOut.Spec), &bpResp); err != nil {
		log.Println("[DEBUG] Error unmarshalling BPOut:", err)
		return diag.FromErr(err)
	}

	var objStatus map[string]interface{}
	if err = json.Unmarshal(bpOut.Spec, &objStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}

	var appUUID, appName string
	if resource, ok := objStatus["resources"].(map[string]interface{}); ok {
		// Access the list "app_profile"
		if appProfileList, ok := resource["app_profile_list"].([]interface{}); ok {
			for i, item := range appProfileList {
				if appProfile, ok := item.(map[string]interface{}); ok {
					// Print values in each "app_profile" map
					appName = appProfile["name"].(string)
					appUUID = appProfile["uuid"].(string)
					log.Printf("[DEBUG] app_profile %d: name = %s, uuid = %s\n", i+1, appProfile["name"], appProfile["uuid"])
				} else {
					log.Printf("[DEBUG] app_profile %d is not a map\n", i+1)
				}
			}
		}
	}

	// check for runtime editables
	runtimeSpec := &selfservice.RuntimeEditables{}
	if runtime, ok := d.GetOk("runtime_editables"); ok {
		getRuntime, errResp := conn.Service.GetRuntimeEditables(ctx, bpUUID)
		if errResp != nil {
			return diag.Errorf("Error getting Runtime Editables: %s", errResp)
		}

		runtimeSpec = getRuntime.Resources[0].RuntimeEditables

		runtimeList := runtime.([]interface{})

		for _, item := range runtimeList {
			itemMap := item.(map[string]interface{})

			if variableList, ok := itemMap["variable_list"].([]interface{}); ok {
				for _, variable := range variableList {
					variableMap := variable.(map[string]interface{})

					for vbs := range runtimeSpec.VariableList {
						if variableMap["name"] == utils.StringValue(runtimeSpec.VariableList[vbs].Name) {
							val := variableMap["value"].(string)

							rawMsg := json.RawMessage(val)
							runtimeSpec.VariableList[vbs].Value = &rawMsg
						}
					}
				}
			}
			if substrateList, ok := itemMap["substrate_list"].([]interface{}); ok {
				log.Println("[DEBUG] Substrate List::: ", len(substrateList))
				log.Println("[DEBUG] RUNTIME SUBSTRATE LIST::: ", runtimeSpec.SubstrateList)
				for _, substrate := range substrateList {
					substrateMap := substrate.(map[string]interface{})

					for sbs := range runtimeSpec.SubstrateList {
						if substrateMap["name"] == utils.StringValue(runtimeSpec.SubstrateList[sbs].Name) {
							val := substrateMap["value"].(string)
							rawMsg := json.RawMessage(val)
							runtimeSpec.SubstrateList[sbs].Value = &rawMsg
						}
					}
				}
			}
			if deploymentList, ok := itemMap["deployment_list"].([]interface{}); ok {
				for _, deployment := range deploymentList {
					deploymentMap := deployment.(map[string]interface{})

					for dps := range runtimeSpec.DeploymentList {
						if deploymentMap["name"] == utils.StringValue(runtimeSpec.DeploymentList[dps].Name) {
							val := deploymentMap["value"].(string)
							rawMsg := json.RawMessage(val)
							runtimeSpec.DeploymentList[dps].Value = &rawMsg
						}
					}
				}
			}
		}
	}

	bpSpec := &selfservice.BPspec{
		AppName: d.Get("app_name").(string),
		AppDesc: d.Get("app_description").(string),
		AppProfileReference: selfservice.AppProfileReference{
			Kind: "app_profile",
			Name: appName,
			UUID: appUUID,
		},
		RuntimeEditables: runtimeSpec,
	}

	input := &selfservice.BlueprintProvisionInput{
		Spec: *bpSpec,
	}

	output, err := conn.Service.ProvisionBlueprint(ctx, bpUUID, input)
	if err != nil {
		return diag.Errorf("Error creating App: %s", err)
	}

	// Convert JSON object to JSON string
	jsonData, err := json.Marshal(output.Spec)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the JSON data as a string in the Terraform resource data
	if err := d.Set("spec", string(jsonData)); err != nil {
		return diag.FromErr(err)
	}

	// call the poll API to get the status of the task
	taskUUID := output.Status.RequestID
	// Wait for the APP to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"running"},
		Target:  []string{"success"},
		Refresh: calmtaskStateRefreshFunc(ctx, conn, bpUUID, taskUUID),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	var obj interface{}
	var errWaitTask error
	if obj, errWaitTask = stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for app (%s) to create: %s", obj, errWaitTask)
	}

	applicationUUID := ""
	if taskStatus, ok := obj.(*selfservice.PollResponse); ok {
		applicationUUID = *taskStatus.Status.AppUUID
	} else {
		return diag.Errorf("error extracting UUID from task status")
	}
	d.SetId(applicationUUID)

	// poll till app state is running
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"provisioning"},
		Target:  []string{"running"},
		Refresh: calmappStateRefreshFunc(ctx, conn, applicationUUID),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask = appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for app (%s) to be running: %s", d.Id(), errWaitTask)
	}
	return resourceNutanixCalmAppProvisionRead(ctx, d, meta)
}

func resourceNutanixCalmAppProvisionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).CalmAPI

	resp, err := conn.Service.GetApp(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	AppResp := &selfservice.AppResponse{}
	if err = json.Unmarshal([]byte(resp.Status), &AppResp.Status); err != nil {
		log.Println("[DEBUG] Error unmarshalling App:", err)
		return diag.FromErr(err)
	}

	var objMetadata map[string]interface{}
	if err = json.Unmarshal(resp.Metadata, &objMetadata); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}

	if err = d.Set("api_version", AppResp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	// unMarshall to get state of an APP
	var objStatus map[string]interface{}
	if err = json.Unmarshal(resp.Status, &objStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}
	var appState string

	if state, ok := objStatus["state"].(string); ok {
		appState = state
		log.Printf("[DEBUG] State of APPP: %s\n", state)
	}

	if err := d.Set("state", appState); err != nil {
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

	return nil
}

func resourceNutanixCalmAppProvisionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).CalmAPI

	actionSpec := &selfservice.ActionSpec{}

	if d.HasChange("soft_delete") {
		log.Printf("[INFO] soft_delete changed to: %v", d.Get("soft_delete").(bool))
		return resourceNutanixCalmAppProvisionRead(ctx, d, meta)
	}

	if d.HasChange("action") {
		actionSpec.Name = d.Get("action").(string)
	}

	actionMap := map[string]string{
		"stop":    "action_stop",
		"start":   "action_start",
		"restart": "action_restart",
	}

	if action, ok := actionMap[actionSpec.Name]; ok {
		actionSpec.Name = action
	} else {
		return diag.Errorf("Invalid action %s", actionSpec.Name)
	}
	// Call action API

	resp, err := conn.Service.PerformAction(ctx, d.Id(), actionSpec)
	if err != nil {
		return diag.FromErr(err)
	}

	// poll till action is completed
	const delayDuration = 5 * time.Second
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCESS"},
		Refresh: RunlogStateRefreshFunc(ctx, conn, d.Id(), resp.RunlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   delayDuration,
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for app: %s", errWaitTask)
	}
	return resourceNutanixCalmAppProvisionRead(ctx, d, meta)
}

func resourceNutanixCalmAppProvisionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).CalmAPI

	log.Printf("[Debug] Destroying the app with the ID %s", d.Id())

	soft := d.Get("soft_delete").(bool)

	var err error
	if soft {
		log.Printf("[Debug] Performing soft delete on app: %s", d.Id())
		_, err = conn.Service.SoftDeleteApp(ctx, d.Id())
	} else {
		log.Printf("[Debug] Performing hard delete on app: %s", d.Id())
		_, err = conn.Service.DeleteApp(ctx, d.Id())
	}

	if err != nil {
		return diag.Errorf("Error deleting self-service application: %s", err)
	}

	// poll till app state is deleted
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"running", "deleting"},
		Target:  []string{"deleted"},
		Refresh: calmappStateRefreshFunc(ctx, conn, d.Id()),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Error waiting for app %s to be deleted: %s", d.Id(), errWaitTask)
	}

	return nil
}

func calmtaskStateRefreshFunc(ctx context.Context, client *selfservice.Client, bpID, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.Service.TaskPoll(ctx, bpID, taskUUID)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, ERROR, nil
			}
			return nil, "", err
		}

		if utils.StringValue(v.Status.State) == "failed" {
			return v, *v.Status.AppUUID,
				fmt.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.Status.AppUUID), utils.StringValue(v.Status.State))
		}
		return v, utils.StringValue(v.Status.State), nil
	}
}

func calmappStateRefreshFunc(ctx context.Context, client *selfservice.Client, appUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.Service.GetApp(ctx, appUUID)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, ERROR, nil
			}
			return nil, "", err
		}

		var objStatus map[string]interface{}
		if err = json.Unmarshal(v.Status, &objStatus); err != nil {
			log.Println("[DEBUG] Error unmarshalling Spec:", err)
		}

		var appState string
		if state, ok := objStatus["state"].(string); ok {
			appState = state
			log.Printf("[DEBUG] State of APP: %s\n", state)
		}

		if utils.StringValue(&appState) == "failed" {
			return v, appUUID,
				fmt.Errorf("error_detail: %s, progress_message: %s", *utils.StringPtr(appUUID), appState)
		}
		return v, appState, nil
	}
}

func RunlogStateRefreshFunc(ctx context.Context, client *selfservice.Client, appUUID, runlogUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.Service.AppRunlogs(ctx, appUUID, runlogUUID)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, ERROR, nil
			}
			return nil, "", err
		}
		log.Println("[DEBUG] V State: ", v.Status.RunlogState)
		log.Println("[DEBUG] V: ", *v)

		runlogstate := utils.StringValue(v.Status.RunlogState)

		log.Printf("[DEBUG] Runlog State: %s\n", runlogstate)

		return v, runlogstate, nil
	}
}

func flattenVM(pr map[string]interface{}) []interface{} {
	var vms []interface{}
	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		// Access the list "app_profile"
		if deploymentList, ok := resource["deployment_list"].([]interface{}); ok {
			for _, item := range deploymentList {
				configsVal := make(map[string]interface{})
				itemList := item.(map[string]interface{})

				if subs, ok := itemList["substrate_configuration"].(map[string]interface{}); ok {
					var configMap map[string]interface{}
					configMapList := make([]map[string]interface{}, 0)

					if elemList, ok := subs["element_list"].([]interface{}); ok {
						for _, elem := range elemList {
							elemMap := elem.(map[string]interface{})
							// elem for multiple services
							configMap = make(map[string]interface{})
							configMap["name"] = elemMap["instance_name"]
							configMap["ip_address"] = elemMap["address"]
							configMap["vm_uuid"] = elemMap["instance_id"]
							configMap["image"] = flattenImageConfigName(subs)

							if createSpec, ok := elemMap["create_spec"].(map[string]interface{}); ok {
								if resources, ok := createSpec["resources"].(map[string]interface{}); ok {
									configMap["vcpus"] = resources["num_sockets"]
									configMap["cores"] = resources["num_vcpus_per_socket"]
									configMap["memory"] = resources["memory_size_mib"]
								}

								configsVal["nics"] = flattenNicMacAddress(subs)

								if categories, ok := createSpec["categories"].(map[string]interface{}); ok {
									configsVal["categories"] = categories
								}

								if cluster, ok := createSpec["cluster_reference"].(map[string]interface{}); ok {
									clusterList := make([]map[string]interface{}, 0)
									clusterMap := map[string]interface{}{
										"cluster_name": cluster["name"],
										"cluster_uuid": cluster["uuid"],
									}
									clusterList = append(clusterList, clusterMap)
									configsVal["cluster_info"] = clusterList
								}
							}
							configMapList = append(configMapList, configMap)
							configsVal["configuration"] = configMapList
						}
					}
				}
				vms = append(vms, configsVal)
			}
		}
	}
	return vms
}

func RuntimeSpec() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"uuid": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"context": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func flattenImageConfigName(subs map[string]interface{}) string {
	if createSpec, ok := subs["create_spec"].(map[string]interface{}); ok {
		if resource, ok := createSpec["resources"].(map[string]interface{}); ok {
			if diskList, ok := resource["disk_list"].([]interface{}); ok {
				for _, disk := range diskList {
					diskMap := disk.(map[string]interface{})
					return diskMap["data_source_reference"].(map[string]interface{})["name"].(string)
				}
			}
		}
	}
	return ""
}

func flattenNicMacAddress(subs map[string]interface{}) []map[string]interface{} {
	nicMapList := make([]map[string]interface{}, 0)
	if variableList, ok := subs["variable_list"].([]interface{}); ok {
		for _, elem := range variableList {
			elemMap := elem.(map[string]interface{})

			if elemMap["name"] == "platform_data" {
				var result []interface{}

				// Unmarshal the JSON string into the []interface{} slice
				err := json.Unmarshal([]byte(elemMap["value"].(string)), &result)
				if err != nil {
					log.Println("[DEBUG] Error decoding JSON:", err)
				}

				for _, elemMap := range result {
					elemMap := elemMap.(map[string]interface{})

					if status, ok := elemMap["status"].(map[string]interface{}); ok {
						if resources, ok := status["resources"].(map[string]interface{}); ok {
							if elemMap, ok := resources["nic_list"].([]interface{}); ok {
								for _, elemMap := range elemMap {
									elemMap := elemMap.(map[string]interface{})
									nicMap := map[string]interface{}{}

									nicMap["mac_address"] = elemMap["mac_address"]
									nicMap["type"] = elemMap["nic_type"]
									nicMap["subnet"] = elemMap["subnet_reference"].(map[string]interface{})["name"]

									nicMapList = append(nicMapList, nicMap)
								}
							}
						}
					}
				}
			}
		}
	}
	return nicMapList
}
