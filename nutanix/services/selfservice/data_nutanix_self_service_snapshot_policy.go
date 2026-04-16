package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
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
	conn := meta.(*conns.Client).CalmAPI
	length := d.Get("length").(int)
	offset := d.Get("offset").(int)

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
	if err := json.Unmarshal([]byte(bpNameResp.Entities), &BpNameStatus); err != nil {
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
	if err := json.Unmarshal([]byte(bpOut.Spec), &bpResp); err != nil {
		log.Println("[DEBUG] Error unmarshalling BPOut:", err)
		return diag.FromErr(err)
	}

	var objStatus map[string]interface{}
	if err := json.Unmarshal(bpOut.Spec, &objStatus); err != nil {
		log.Println("[DEBUG]Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}
	var appUUID, envRefUUID string
	var snapshotConfigUUIDList *[]string
	var snapshotConfigNameList *[]string
	if resource, ok := objStatus["resources"].(map[string]interface{}); ok {
		// Access the list "app_profile"
		if appProfileList, ok := resource["app_profile_list"].([]interface{}); ok {
			for i, item := range appProfileList {
				if appProfile, ok := item.(map[string]interface{}); ok {
					// Print values in each "app_profile" map
					appUUID = appProfile["uuid"].(string)

					//Get env uuid
					if envList, ok := appProfile["environment_reference_list"].([]interface{}); ok {
						for _, envUUID := range envList {
							envRefUUID = envUUID.(string)
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
					snapshotConfigUUIDList = &snapshotUUIDList
					snapshotConfigNameList = &snapshotNameList
					log.Printf("[DEBUG] app_profile %d: name = %s, uuid = %s\n", i+1, appProfile["name"], appProfile["uuid"])
				} else {
					log.Printf("[DEBUG] app_profile %d is not a map\n", i+1)
				}
			}
		}
	}

	log.Printf("[DEBUG] envRefUUID is %s\n", envRefUUID)

	PolicyList := make([]map[string]interface{}, 0)

	for idx, configUUID := range *snapshotConfigUUIDList {
		policyListInput := &selfservice.PolicyListInput{}
		policyListInput.Length = length
		policyListInput.Offset = offset
		policyListInput.Filter = fmt.Sprintf("environment_references==%s", envRefUUID)

		policyResp, err := conn.Service.GetAppProtectionPolicyList(ctx, bpUUID, appUUID, configUUID, policyListInput)
		if err != nil {
			log.Println("[DEBUG] Error GetAppProtectionPolicyList:", err)
			return diag.FromErr(err)
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
										PolicyMap["snapshot_config_name"] = (*snapshotConfigNameList)[idx]
										PolicyMap["snapshot_config_uuid"] = configUUID
										PolicyMap["policy_expiry_days"] = policy.(map[string]interface{})["multiple"].(float64)
										PolicyList = append(PolicyList, PolicyMap)
										log.Printf("[DEBUG] Added policy with param %s %s", status["uuid"].(string), configUUID)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	log.Println("[DEBUG] Final PolicyList:", PolicyList)

	d.SetId(bpUUID)

	if err := d.Set("policy_list", PolicyList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
