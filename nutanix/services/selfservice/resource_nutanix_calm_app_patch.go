package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func ResourceNutanixCalmAppPatch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmAppPatchCreate,
		ReadContext:   resourceNutanixCalmAppPatchRead,
		UpdateContext: resourceNutanixCalmAppPatchUpdate,
		DeleteContext: resourceNutanixCalmAppPatchDelete,
		Schema: map[string]*schema.Schema{
			"app_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"patch_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"runlog_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory_size_mib": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"num_sockets": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"num_vcpus_per_socket": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"nics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"categories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixCalmAppPatchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm

	appUUID := d.Get("app_uuid").(string)
	patchName := d.Get("patch_name").(string)

	// fetch app for spec

	appResp, err := conn.Service.GetApp(ctx, appUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	var objSpec map[string]interface{}
	if err := json.Unmarshal(appResp.Spec, &objSpec); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}

	var objMetadata map[string]interface{}
	if err := json.Unmarshal(appResp.Metadata, &objMetadata); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}

	var objStatus map[string]interface{}
	if err := json.Unmarshal(appResp.Status, &objStatus); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}

	//fetch input

	fetchInput := &calm.PatchInput{}
	fetchInput.APIVersion = appResp.APIVersion
	fetchInput.Metadata = objMetadata

	var patchUUID string
	// fetch patch for spec
	fetchSpec := &calm.PatchSpec{}
	fetchSpec.TargetUUID = appUUID
	fetchSpec.TargetKind = "Application"
	fetchSpec.Args.Variables = []*calm.VariableList{}
	fetchSpec.Args.Patch, patchUUID = expandPatchSpec(objSpec, patchName)
	// fetchSpec.Args.Variables = []

	if vmConfigRuntimeEditable, ok := d.GetOk("vm_config"); ok {
		vmConfigRuntimeEditable := vmConfigRuntimeEditable.([]interface{})
		for _, vmConfig := range vmConfigRuntimeEditable {
			vmConfigMap := vmConfig.(map[string]interface{})
			// log.Println("VM CONFIG MAP::::", vmConfigMap)
			if numSockets, ok := vmConfigMap["num_sockets"].(int); ok {
				log.Println("NUM SOCKETS::::", numSockets)
				fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["num_sockets_ruleset"].(map[string]interface{})["value"] = numSockets
			}
			if memorySizeMib, ok := vmConfigMap["memory_size_mib"].(int); ok {
				log.Println("MEMORY SIZE::::", memorySizeMib)
				fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["memory_size_mib_ruleset"].(map[string]interface{})["value"] = memorySizeMib
			}
			if numVcpusPerSocket, ok := vmConfigMap["num_vcpus_per_socket"].(int); ok {
				log.Println("NUM VCPUS PER SOCKET::::", numVcpusPerSocket)
				fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["num_vcpus_per_socket_ruleset"].(map[string]interface{})["value"] = numVcpusPerSocket
			}
		}
	}

	if categoriesRuntimeEditable, ok := d.GetOk("categories"); ok {
		categoriesRuntimeEditable := categoriesRuntimeEditable.([]interface{})
		for _, category := range categoriesRuntimeEditable {
			categoryMap := category.(map[string]interface{})
			log.Println("CATEGORY MAP::::", categoryMap)
			categoryList := fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["pre_defined_categories"].([]interface{})
			if operation, ok := categoryMap["operation"].(string); ok {
				if operation == "add" {
					categoryList = append(categoryList, map[string]interface{}{
						"value":     categoryMap["value"],
						"operation": "add",
					})
				} else {
					categoryList = append(categoryList, map[string]interface{}{
						"value":     categoryMap["value"],
						"operation": "delete",
					})
				}
			}
			fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["pre_defined_categories"] = categoryList
		}
	}

	// if runtimeEditables, ok := d.GetOk("run_action"); ok {
	// 	// get the path list from  objSpec

	// 	// runtimeVMConfigMap := runtime.([]interface{})[0].(map[string]interface{})["vm_config"].(map[string]interface{})
	// 	// num_sockets := runtimeVMConfigMap["num_sockets"].(int)

	// 	// print("NUM SOCKETS::::", num_sockets)

	// 	// func to return attrs_list from patch_list

	// 	attsDataMap := getAttrsListFromPatchList(objSpec, patchName)
	// 	log.Println("ATTRS LIST::::", attsDataMap)

	// 	runtimeEditablesList := runtimeEditables.([]interface{})

	// 	for _, runtimeEditable := range runtimeEditablesList {
	// 		runtimeEditableMap := runtimeEditable.(map[string]interface{})

	// 		// // fetch the current nic present in app
	// 		// getNicList := fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["pre_defined_nic_list"].([]interface{})
	// 		// for _, getNicMap := range getNicList {
	// 		// 	// now get the nic in config spec
	// 		// 	fmt.Println("Length of NIC LIST::::", len(getNicList))
	// 		// 	if nics, ok := runtimeEditableMap["nics"].([]interface{}); ok {
	// 		// 		for _, nic := range nics {
	// 		// 			nicMap := nic.(map[string]interface{})
	// 		// 			idx := nicMap["index"].(int)
	// 		// 			ops := nicMap["operation"].(string)
	// 		// 			fmt.Println("IDX::::", idx)
	// 		// 			fmt.Println("OPS::::", ops)

	// 		// 			getNicList = append(getNicList, map[string]interface{}{
	// 		// 				"identifier": idx,
	// 		// 				"operation":  ops,
	// 		// 			})
	// 		// 			// if getNicMap.(map[string]interface{})["identifier"].(string) == string(idx) {
	// 		// 			// 	getNicMap.(map[string]interface{})["operation"] = ops
	// 		// 			// 	fmt.Println("INSIDE NIC MAP")
	// 		// 			// }
	// 		// 		}
	// 		// 	}
	// 		// }

	// 		// if resource, ok := objStatus["resources"].(map[string]interface{}); ok {
	// 		// 	fmt.Println("INSIDE RESOURCE")
	// 		// 	// Access the list "app_profile"
	// 		// 	if deployList, ok := resource["deployment_list"].([]interface{}); ok {
	// 		// 		fmt.Println("INSIDE DEPLOYMENT")
	// 		// 		for _, deploy := range deployList {
	// 		// 			deployMap := deploy.(map[string]interface{})
	// 		// 			log.Println("DEPLOYYYYY MAPPPPPPPP")
	// 		// 			if subs, ok := deployMap["substrate_configuration"].(map[string]interface{}); ok {
	// 		// 				fmt.Println("INSIDE SUBSTRATE")
	// 		// 				if element, ok := subs["element_list"].([]interface{}); ok {
	// 		// 					for _, elem := range element {
	// 		// 						fmt.Println("INSIDE ELEMENT")
	// 		// 						if nics, ok := elem.(map[string]interface{})["create_spec"].(map[string]interface{}); ok {
	// 		// 							fmt.Println("create_spec")
	// 		// 							if resources, ok := nics["resources"].(map[string]interface{}); ok {
	// 		// 								if nicList, ok := resources["nic_list"].([]interface{}); ok {
	// 		// 									fmt.Println("INSIDE NICS LIST")
	// 		// 									for _, nic := range nicList {
	// 		// 										nicMap := nic.(map[string]interface{})
	// 		// 										identifier := nicMap["nic_type"].(string)
	// 		// 										fmt.Println("NIC TYPE::::", identifier)
	// 		// 										// if nics, ok := runtimeEditableMap["nics"].([]interface{}); ok {
	// 		// 										// 	for _, nic := range nics {
	// 		// 										// 		fmt.Println("NIC IDENTIFIER::::", nic.(map[string]interface{}))
	// 		// 										// 	}
	// 		// 										// }
	// 		// 									}
	// 		// 								}
	// 		// 							}
	// 		// 						}
	// 		// 					}
	// 		// 				}
	// 		// 			}
	// 		// 		}
	// 		// 	}
	// 		// }
	// 	}

	// }
	fetchInput.Spec = *fetchSpec

	// log.Println("HELLLLLOOOOOO2222")
	// aJSON, _ := json.Marshal(fetchSpec)
	// fmt.Printf("JSON Print - \n%s\n", string(aJSON))

	// return nil

	fetchResp, err := conn.Service.PatchApp(ctx, appUUID, patchUUID, fetchInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := fetchResp.Status.RunlogUUID

	fmt.Println("Response:", runlogUUID)

	// poll till action is completed
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"POLICY_EXEC"},
		Refresh: RunlogStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   5 * time.Second,
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for app to perform Patch(%s): %s", errWaitTask)
	}

	if err := d.Set("runlog_uuid", runlogUUID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(runlogUUID)
	return nil
}

func resourceNutanixCalmAppPatchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixCalmAppPatchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixCalmAppPatchCreate(ctx, d, meta)
}
func resourceNutanixCalmAppPatchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandPatchSpec(pr map[string]interface{}, patchName string) (map[string]interface{}, string) {
	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		// fmt.Println("RESOURCESSSSS")
		if patchList, ok := resource["patch_list"].([]interface{}); ok {
			for _, patch := range patchList {
				if dep, ok := patch.(map[string]interface{}); ok {
					fmt.Println("DEP UUID::::", dep["uuid"])
					if dep["name"] == patchName {
						fmt.Println("DEP UUID::::", dep["uuid"])
						return patch.(map[string]interface{}), dep["uuid"].(string)
					}
				}
			}
		}
	}
	return nil, ""
}

func getAttrsListFromPatchList(pr map[string]interface{}, patchName string) map[string]interface{} {
	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		if patchList, ok := resource["patch_list"].([]interface{}); ok {
			for _, patch := range patchList {
				if dep, ok := patch.(map[string]interface{}); ok {
					if dep["name"] == patchName {
						if attrs, ok := dep["attrs_list"].([]interface{}); ok {
							for _, attr := range attrs {
								if data, ok := attr.(map[string]interface{}); ok {
									return data
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
