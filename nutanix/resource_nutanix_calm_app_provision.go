package nutanix

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
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixCalmAppProvision() *schema.Resource {
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
			// "system_action": {
			// 	Type:         schema.TypeString,
			// 	Optional:     true,
			// 	ValidateFunc: validation.StringInSlice([]string{"action_stop", "action_restart", "action_start"}, false),
			// },
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
		},
	}
}

func resourceNutanixCalmAppProvisionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Calm

	var bp_uuid string
	// fetch bp_uuid from bp_name
	bp_name := d.Get("bp_name").(string)

	bpFilter := &calm.BlueprintListInput{}

	bpFilter.Filter = fmt.Sprintf("name==%s;state!=DELETED", bp_name)

	bpNameResp, err := conn.Service.ListBlueprint(ctx, bpFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	var BpNameStatus []interface{}
	if err := json.Unmarshal([]byte(bpNameResp.Entities), &BpNameStatus); err != nil {
		fmt.Println("Error unmarshalling BPName:", err)
	}

	entities := BpNameStatus[0].(map[string]interface{})

	if entity, ok := entities["metadata"].(map[string]interface{}); ok {
		bp_uuid = entity["uuid"].(string)
	}

	if bpUUID, ok := d.GetOk("bp_uuid"); ok {
		bp_uuid = bpUUID.(string)
	}

	// call bp

	bpOut, er := conn.Service.GetBlueprint(ctx, bp_uuid)
	if er != nil {
		return diag.Errorf("Error getting Blueprint: %s", er)
	}

	bpResp := &calm.BlueprintResponse{}
	if err := json.Unmarshal([]byte(bpOut.Spec), &bpResp); err != nil {
		fmt.Println("Error unmarshalling BPOut:", err)
	}

	var objStatus map[string]interface{}
	if err := json.Unmarshal(bpOut.Spec, &objStatus); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}

	var app_uuid, app_name string
	if resource, ok := objStatus["resources"].(map[string]interface{}); ok {
		// Access the list "app_profile"
		if appProfileList, ok := resource["app_profile_list"].([]interface{}); ok {
			for i, item := range appProfileList {
				if appProfile, ok := item.(map[string]interface{}); ok {
					// Print values in each "app_profile" map
					app_name = appProfile["name"].(string)
					app_uuid = appProfile["uuid"].(string)
					fmt.Printf("app_profile %d: name = %s, uuid = %s\n", i+1, appProfile["name"], appProfile["uuid"])
				} else {
					fmt.Printf("app_profile %d is not a map\n", i+1)
				}
			}
		}
	}

	// check for runtime editables
	runtimeSpec := &calm.RuntimeEditables{}
	if runtime, ok := d.GetOk("runtime_editables"); ok {
		getRuntime, err := conn.Service.GetRuntimeEditables(ctx, bp_uuid)
		if err != nil {
			return diag.Errorf("Error getting Runtime Editables: %s", err)
		}

		runtimeSpec = getRuntime.Resources[0].RuntimeEditables

		fmt.Println("Runtime Editables: ", runtimeSpec)
		// log.Println("HELLLLLOOOOOO")
		// aJSON, _ := json.Marshal(runtimeSpec)
		// fmt.Printf("JSON Print - \n%s\n", string(aJSON))

		runtimeList := runtime.([]interface{})

		for k, item := range runtimeList {
			itemMap := item.(map[string]interface{})

			if variable_list, ok := itemMap["variable_list"].([]interface{}); ok {
				for _, variable := range variable_list {
					variableMap := variable.(map[string]interface{})
					// fmt.Println("Variable Name: ", variableMap["name"])
					// fmt.Println("Variable Value: ", variableMap["value"])
					val := variableMap["value"].(string)

					rawMsg := json.RawMessage(val)
					runtimeSpec.VariableList[k].Value = &rawMsg
				}
			}
			if substrate_list, ok := itemMap["substrate_list"].([]interface{}); ok {
				for _, substrate := range substrate_list {
					substrateMap := substrate.(map[string]interface{})
					// fmt.Println("Substrate Name: ", substrateMap["name"])
					// fmt.Println("Substrate Value: ", substrateMap["value"])
					val := substrateMap["value"].(string)

					rawMsg := json.RawMessage(val)
					runtimeSpec.SubstrateList[k].Value = &rawMsg
				}
			}
			if deployment_list, ok := itemMap["deployment_list"].([]interface{}); ok {
				for _, deployment := range deployment_list {
					deploymentMap := deployment.(map[string]interface{})
					// fmt.Println("Deployment Name: ", deploymentMap["name"])
					// fmt.Println("Deployment Value: ", deploymentMap["value"])
					val := deploymentMap["value"].(string)

					rawMsg := json.RawMessage(val)
					runtimeSpec.DeploymentList[k].Value = &rawMsg
				}
			}
		}
		// log.Println("HELLLLLOOOOOO22")
		// bJSON, _ := json.Marshal(runtimeSpec)
		// fmt.Printf("JSON Print - \n%s\n", string(bJSON))
	}

	// return nil

	bpSpec := &calm.BPspec{
		AppName: d.Get("app_name").(string),
		AppDesc: d.Get("app_description").(string),
		AppProfileReference: calm.AppProfileReference{
			Kind: "app_profile",
			Name: app_name,
			UUID: app_uuid,
		},
		RuntimeEditables: runtimeSpec,
	}

	input := &calm.BlueprintProvisionInput{
		Spec: *bpSpec,
	}

	log.Println("HELLLLLOOOOOO22333333")
	bJSON, _ := json.Marshal(input)
	fmt.Printf("JSON Print - \n%s\n", string(bJSON))

	output, err := conn.Service.ProvisionBlueprint(ctx, bp_uuid, input)
	if err != nil {
		return diag.Errorf("Error creating App: %s", err)
	}

	// var objStatusResp map[string]interface{}
	// if err := json.Unmarshal(output.Spec, &objStatusResp); err != nil {
	// 	fmt.Println("Error unmarshalling Spec:", err)
	// }
	// var objStatus map[string]interface{}
	// if err := json.Unmarshal(bpOut.Spec, &objStatus); err == nil {
	// 	fmt.Println("Status as object:", objStatus)
	// }

	// fmt.Println("app_reference", objStatus["resource"])
	// d.Set("status", objStatus["resource"])
	// d.Set("spec", objStatus)
	fmt.Println("Status as object:", output.Status)

	// Set the values in the resource data
	// if err := d.Set("status", output.Status.RequestID); err != nil {
	// 	return diag.FromErr(err)
	// }

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
		Refresh: calmtaskStateRefreshFunc(ctx, conn, bp_uuid, taskUUID),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	var obj interface{}
	var errWaitTask error
	if obj, errWaitTask = stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for app (%s) to create: %s", obj, errWaitTask)
	}

	applicationUUID := ""
	if taskStatus, ok := obj.(*calm.PollResponse); ok {
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
		return diag.Errorf("error waiting for app (%s) to be running: %s", obj, errWaitTask)
	}
	return resourceNutanixCalmAppProvisionRead(ctx, d, meta)
}

func resourceNutanixCalmAppProvisionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Calm

	resp, err := conn.Service.GetApp(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	AppResp := &calm.AppResponse{}
	if err := json.Unmarshal([]byte(resp.Status), &AppResp.Status); err != nil {
		fmt.Println("Error unmarshalling App:", err)
	}

	// Convert JSON object to JSON string
	jsonData, err := json.MarshalIndent(AppResp.Status, " ", "  ")
	if err != nil {
		return diag.FromErr(err)
	}
	// Set the JSON data as a string in the Terraform resource data
	if err := d.Set("status", string(jsonData)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", AppResp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	// unMarshall to get state of an APP
	var objStatus map[string]interface{}
	if err := json.Unmarshal(resp.Status, &objStatus); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}
	var app_state string

	if state, ok := objStatus["state"].(string); ok {
		app_state = state
		fmt.Printf("State of APPP: %s\n", state)
	}

	if err := d.Set("state", app_state); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm", flattenVM(objStatus)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixCalmAppProvisionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Calm

	actionSpec := &calm.ActionSpec{}

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
	appStateConf := &resource.StateChangeConf{
		Pending:    []string{"RUNNING"},
		Target:     []string{"SUCCESS"},
		Refresh:    RunlogStateRefreshFunc(ctx, conn, d.Id(), resp.RunlogUUID),
		MinTimeout: 2 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for app (%s): %s", errWaitTask)
	}
	return resourceNutanixCalmAppProvisionRead(ctx, d, meta)
}

func resourceNutanixCalmAppProvisionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Calm

	log.Printf("[Debug] Destroying the app with the ID %s", d.Id())

	if _, err := conn.Service.DeleteApp(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func calmtaskStateRefreshFunc(ctx context.Context, client *calm.Client, bpID, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.Service.TaskPoll(ctx, bpID, taskUUID)
		fmt.Println("V: ", *v)
		fmt.Println("V.state: ", *v.Status.State)
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
		return v, *v.Status.State, nil
	}
}

func calmappStateRefreshFunc(ctx context.Context, client *calm.Client, appUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.Service.GetApp(ctx, appUUID)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, ERROR, nil
			}
			return nil, "", err
		}

		var objStatus map[string]interface{}
		if err := json.Unmarshal(v.Status, &objStatus); err != nil {
			fmt.Println("Error unmarshalling Spec:", err)
		}

		var app_state string
		if state, ok := objStatus["state"].(string); ok {
			app_state = state
			fmt.Printf("State of APPP: %s\n", state)
		}

		if utils.StringValue(&app_state) == "failed" {
			return v, appUUID,
				fmt.Errorf("error_detail: %s, progress_message: %s", *utils.StringPtr(appUUID), app_state)
		}
		return v, app_state, nil
	}
}

func RunlogStateRefreshFunc(ctx context.Context, client *calm.Client, appUUID, runlogUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.Service.AppRunlogs(ctx, appUUID, runlogUUID)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return v, ERROR, nil
			}
			return nil, "", err
		}
		fmt.Println("V State: ", v.Status.RunlogState)
		fmt.Println("V: ", *v)

		runlogstate := utils.StringValue(v.Status.RunlogState)

		fmt.Printf("Runlog State: %s\n", runlogstate)

		// if runlogstate == "ERROR" || v.Status.ExitCode != utils.IntPtr(-1) {
		// 	return v, runlogstate,
		// 		fmt.Errorf("error_detail: %s, progress_message: %s", *v.OutputList[0].Output, v.Status.RunlogState)
		// }
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
				configMap := map[string]interface{}{}
				configMapList := make([]map[string]interface{}, 0)
				nicsList := make([]map[string]interface{}, 0)
				nicsMap := map[string]interface{}{}

				itemList := item.(map[string]interface{})

				if subs, ok := itemList["substrate_configuration"].(map[string]interface{}); ok {
					fmt.Println("AHV_VM: ", subs["type"])
					fmt.Println("uuid:", subs["uuid"])

					if elemList, ok := subs["element_list"].([]interface{}); ok {

						for _, elem := range elemList {
							elemMap := elem.(map[string]interface{})

							configMap["name"] = elemMap["instance_name"]
							configMap["ip_address"] = elemMap["address"]
							configMap["vm_uuid"] = elemMap["instance_id"]

							fmt.Println("Address: ", elemMap["address"])
							fmt.Println("Instance_id: ", elemMap["instance_id"])
							fmt.Println("Instance Name: ", elemMap["instance_name"])

							if createSpec, ok := elemMap["create_spec"].(map[string]interface{}); ok {
								if resources, ok := createSpec["resources"].(map[string]interface{}); ok {

									configMap["vcpus"] = resources["num_sockets"]
									configMap["cores"] = resources["num_vcpus_per_socket"]
									configMap["memory"] = resources["memory_size_mib"]

								}

								if resource, ok := createSpec["resources"].(map[string]interface{}); ok {
									if nics, ok := resource["nic_list"].([]interface{}); ok {
										for _, nic := range nics {
											nicMap := nic.(map[string]interface{})
											nicsMap["mac_address"] = nicMap["mac_address"]
											nicsMap["type"] = nicMap["nic_type"]
											nicsMap["subnet"] = nicMap["subnet_reference"].(map[string]interface{})["name"]
											nicsList = append(nicsList, nicsMap)
											configsVal["nics"] = nicsList
										}
									}
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
						}
					}

					if createSpec, ok := subs["create_spec"].(map[string]interface{}); ok {
						if resource, ok := createSpec["resources"].(map[string]interface{}); ok {
							if diskList, ok := resource["disk_list"].([]interface{}); ok {
								for _, disk := range diskList {
									fmt.Println("Disk: ", disk)
									fmt.Println("AAAAAAAAAAAAAAAAAAAA")
									diskMap := disk.(map[string]interface{})
									fmt.Println("DISK::::::", diskMap["data_source_reference"].(map[string]interface{})["name"])
									configMap["image"] = diskMap["data_source_reference"].(map[string]interface{})["name"]

									fmt.Println("MEMEORY:", diskMap["disk_size_mib"])
								}
							}
						}
					}

					if variableList, ok := subs["variable_list"].([]interface{}); ok {
						for _, elem := range variableList {
							elemMap := elem.(map[string]interface{})

							if elemMap["name"] == "mac_address" {
								nicsMap["mac_address"] = elemMap["value"]
							}
						}
					}

					configMapList = append(configMapList, configMap)
					configsVal["configuration"] = configMapList
					vms = append(vms, configsVal)
				}
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
