package selfservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func DataSourceNutanixSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixCalmSnapshotPolicyRead,
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
			"length": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_expiry_days": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"snapshot_config_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_config_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixCalmSnapshotPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm
	length := d.Get("length").(int)
	offset := d.Get("offset").(int)

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
	var app_uuid, env_uuid string
	var snapshot_config_uuid_list *[]string
	var snapshot_config_name_list *[]string
	if resource, ok := objStatus["resources"].(map[string]interface{}); ok {
		// Access the list "app_profile"
		if appProfileList, ok := resource["app_profile_list"].([]interface{}); ok {
			for i, item := range appProfileList {
				if appProfile, ok := item.(map[string]interface{}); ok {
					// Print values in each "app_profile" map
					app_uuid = appProfile["uuid"].(string)

					//Get env uuid
					if envList, ok := appProfile["environment_reference_list"].([]interface{}); ok {
						for _, envUUID := range envList {
							env_uuid = envUUID.(string)
						}
					}

					//Get snapshot configlists
					var snapshotUUIDList []string
					var snapshotNameList []string
					if snapshotList, ok := appProfile["snapshot_config_list"].([]interface{}); ok {
						for _, snapshotItem := range snapshotList {
							if snapshot, ok := snapshotItem.(map[string]interface{}); ok {
								snapshotUUIDList = append(snapshotUUIDList, snapshot["uuid"].(string))
								snapshotNameList = append(snapshotNameList, snapshot["name"].(string))
							}
						}
					}
					snapshot_config_uuid_list = &snapshotUUIDList
					snapshot_config_name_list = &snapshotNameList
					fmt.Printf("app_profile %d: name = %s, uuid = %s\n", i+1, appProfile["name"], appProfile["uuid"])
				} else {
					fmt.Printf("app_profile %d is not a map\n", i+1)
				}
			}
		}
	}

	fmt.Printf("env_uuid is %s\n", env_uuid)

	// var proj_uuid string

	// var objMetadata map[string]interface{}
	// if err := json.Unmarshal(bpOut.Metadata, &objMetadata); err != nil {
	// 	fmt.Println("Error unmarshalling metadata:", err)
	// }

	// if proj_reference, ok := objMetadata["project_refernce"].(map[string]string); ok {
	// 	proj_uuid = proj_reference["uuid"]
	// 	fmt.Printf("proj_uuid is %s p\n", proj_uuid)
	// } else {
	// 	fmt.Println("Error getting project uuid")
	// }

	// projOut, er := conn.Service.GetProject(ctx, proj_uuid)

	// var projSpec map[string]interface{}
	// if err := json.Unmarshal(projOut.Spec, &projSpec); err != nil {
	// 	fmt.Println("Error unmarshalling Spec:", err)
	// }

	// var env_uuid string

	// if detail, ok := projSpec["project_detail"].(map[string]interface{}); ok {
	// 	if resources, ok := detail["resources"].(map[string]interface{}); ok {
	// 		if env_ref, ok := resources["default_environment_reference"].(map[string]string); ok {
	// 			env_uuid = env_ref["uuid"]
	// 			fmt.Printf("env_uuid is %s p\n", env_uuid)
	// 		}
	// 	}
	// }

	PolicyList := make([]map[string]interface{}, 0)

	for idx, config_uuid := range *snapshot_config_uuid_list {
		policyListInput := &calm.PolicyListInput{}
		policyListInput.Length = length
		policyListInput.Offset = offset
		policyListInput.Filter = fmt.Sprintf("environment_references==%s", env_uuid)

		policyResp, err := conn.Service.GetAppProtectionPolicyList(ctx, bp_uuid, app_uuid, config_uuid, policyListInput)
		if err != nil {
			fmt.Println("Error GetAppProtectionPolicyList:", err)
		}

		for _, entity := range policyResp.Entities {
			if status, ok := entity["status"].(map[string]interface{}); ok {
				if resource, ok := status["resources"].(map[string]interface{}); ok {
					if appProtectionRuleList, ok := resource["app_protection_rule_list"].([]interface{}); ok {
						for _, appProtectionRuleItem := range appProtectionRuleList {
							if appProtectionRule, ok := appProtectionRuleItem.(map[string]interface{}); ok {
								if local, ok := appProtectionRule["local_snapshot_retention_policy"].(map[string]interface{}); ok {
									if policy, ok := local["snapshot_expiry_policy"]; ok {
										PolicyMap := make(map[string]interface{})
										PolicyMap["policy_name"] = status["name"].(string)
										PolicyMap["policy_uuid"] = status["uuid"].(string)
										PolicyMap["snapshot_config_name"] = (*snapshot_config_name_list)[idx]
										PolicyMap["snapshot_config_uuid"] = config_uuid
										PolicyMap["policy_expiry_days"] = policy.(map[string]interface{})["multiple"].(float64)
										PolicyList = append(PolicyList, PolicyMap)
										fmt.Println("Added policy with param %s %s", status["uuid"].(string), config_uuid)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Println("Final PolicyList:", PolicyList)

	d.SetId(bp_uuid)

	if err := d.Set("policy_list", PolicyList); err != nil {
		fmt.Println("GGGGGG")
		return diag.FromErr(err)
	}

	return nil
}
