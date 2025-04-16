package selfservice

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
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
							Optional: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_uuid": {
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
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size_mib": {
							Type:     schema.TypeInt,
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
	conn := meta.(*conns.Client).CalmAPI

	appUUID := d.Get("app_uuid").(string)
	patchName := d.Get("patch_name").(string)

	// fetch app for spec

	appResp, err := conn.Service.GetApp(ctx, appUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	var objSpec map[string]interface{}
	if err = json.Unmarshal(appResp.Spec, &objSpec); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}

	var objMetadata map[string]interface{}
	if err = json.Unmarshal(appResp.Metadata, &objMetadata); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}

	var objStatus map[string]interface{}
	if err = json.Unmarshal(appResp.Status, &objStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}

	//fetch input

	fetchInput := &selfservice.PatchInput{}
	fetchInput.APIVersion = appResp.APIVersion
	fetchInput.Metadata = objMetadata

	var patchUUID string
	// fetch patch for spec
	fetchSpec := &selfservice.PatchSpec{}
	fetchSpec.TargetUUID = appUUID
	fetchSpec.TargetKind = "Application"
	fetchSpec.Args.Variables = []*selfservice.VariableList{}
	fetchSpec.Args.Patch, patchUUID = expandPatchSpec(objSpec, patchName)

	if vmConfigRuntimeEditable, ok := d.GetOk("vm_config"); ok {
		vmConfigRuntimeEditable := vmConfigRuntimeEditable.([]interface{})
		for _, vmConfig := range vmConfigRuntimeEditable {
			vmConfigMap := vmConfig.(map[string]interface{})
			log.Println("[DEBUG] VM CONFIG MAP::::", vmConfigMap)
			if numSockets, ok := vmConfigMap["num_sockets"].(int); ok {
				log.Println("[DEBUG] NUM SOCKETS::::", numSockets)
				fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["num_sockets_ruleset"].(map[string]interface{})["value"] = numSockets
			}
			if memorySizeMib, ok := vmConfigMap["memory_size_mib"].(int); ok {
				log.Println("[DEBUG] MEMORY SIZE::::", memorySizeMib)
				fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["memory_size_mib_ruleset"].(map[string]interface{})["value"] = memorySizeMib
			}
			if numVcpusPerSocket, ok := vmConfigMap["num_vcpus_per_socket"].(int); ok {
				log.Println("[DEBUG] NUM VCPUS PER SOCKET::::", numVcpusPerSocket)
				fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["num_vcpus_per_socket_ruleset"].(map[string]interface{})["value"] = numVcpusPerSocket
			}
		}
	}

	if categoriesRuntimeEditable, ok := d.GetOk("categories"); ok {
		categoriesRuntimeEditable := categoriesRuntimeEditable.([]interface{})
		for _, category := range categoriesRuntimeEditable {
			categoryMap := category.(map[string]interface{})
			log.Println("[DEBUG] CATEGORY MAP::::", categoryMap)
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

	if nicsRuntimeEditable, ok := d.GetOk("nics"); ok {
		nicsRuntimeEditable := nicsRuntimeEditable.([]interface{})
		startIndex := 0
		for _, nic := range nicsRuntimeEditable {
			nicMap := nic.(map[string]interface{})
			log.Println("[DEBUG] NIC MAP::::", nicMap)
			nicList := fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["pre_defined_nic_list"].([]interface{})
			if operation, ok := nicMap["operation"].(string); ok {
				if operation == "add" {
					for indx := startIndex; indx < len(nicList); indx++ {
						nicListOperation := nicList[indx].(map[string]interface{})["operation"].(string)
						nicListEditable := nicList[indx].(map[string]interface{})["editable"].(bool)
						// Add a nic only if it's editable else proceed with addition of original nic present in nicList
						if nicListEditable && nicListOperation == operation {
							nicList[indx].(map[string]interface{})["subnet_reference"] = map[string]interface{}{
								"kind": "subnet",
								"type": "",
								"name": "",
								"uuid": nicMap["subnet_uuid"],
							}
							startIndex = indx + 1
							break
						}
					}
				} else {
					nicList = append(nicList, map[string]interface{}{
						"identifier": nicMap["index"],
						"operation":  "delete",
						"subnet_reference": map[string]interface{}{
							"kind": "subnet",
							"name": "",
							"uuid": nicMap["subnet_uuid"],
						},
					})
				}
			}
			fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["pre_defined_nic_list"] = nicList
		}
	}
	if disksRuntimeEditable, ok := d.GetOk("disks"); ok {
		disksRuntimeEditable := disksRuntimeEditable.([]interface{})
		startIndex := 0
		for _, disk := range disksRuntimeEditable {
			diskMap := disk.(map[string]interface{})
			log.Println("[DEBUG] DISK MAP::::", diskMap)
			diskList := fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["pre_defined_disk_list"].([]interface{})
			log.Println("[DEBUG] DISK LIST::::", diskList)
			if operation, ok := diskMap["operation"].(string); ok {
				if operation == "add" {
					// config_details = fetchSpec.Args.Patch["resources"]
					for indx := startIndex; indx < len(diskList); indx++ {
						if diskList[indx].(map[string]interface{})["operation"].(string) == operation {
							diskSizeMib := diskList[indx].(map[string]interface{})["disk_size_mib"].(map[string]interface{})
							diskSizeMib["value"] = diskMap["disk_size_mib"]
							startIndex = indx + 1
							break
						}
					}
				} else {
					diskList = append(diskList, map[string]interface{}{
						"disk_size_mib": diskMap["disk_size_mib"],
						"operation":     "delete",
					})
				}
			}
			fetchSpec.Args.Patch["attrs_list"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["pre_defined_disk_list"] = diskList
		}
	}

	fetchInput.Spec = *fetchSpec

	fetchResp, err := conn.Service.PatchApp(ctx, appUUID, patchUUID, fetchInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := fetchResp.Status.RunlogUUID

	log.Println("[DEBUG] Response:", runlogUUID)

	// poll till action is completed
	const delayDuration = 5 * time.Second
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "POLICY_EXEC"},
		Target:  []string{"SUCCESS"},
		Refresh: RunlogStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   delayDuration,
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for app to perform Patch: %s", errWaitTask)
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
		if patchList, ok := resource["patch_list"].([]interface{}); ok {
			for _, patch := range patchList {
				if dep, ok := patch.(map[string]interface{}); ok {
					log.Println("[DEBUG] DEP UUID::::", dep["uuid"])
					if dep["name"] == patchName {
						log.Println("[DEBUG] DEP UUID::::", dep["uuid"])
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
