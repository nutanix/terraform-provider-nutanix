package selfservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func ResourceNutanixCalmAppRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmAppRestoreCreate,
		ReadContext:   resourceNutanixCalmAppRestoreRead,
		UpdateContext: resourceNutanixCalmAppRestoreUpdate,
		DeleteContext: resourceNutanixCalmAppRestoreDelete,
		Schema: map[string]*schema.Schema{
			"app_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshot_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restore_action_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNutanixCalmAppRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm
	appUUID := d.Get("app_uuid").(string)
	restoreActionName := d.Get("restore_action_name").(string)
	recoveryPointUUID := d.Get("snapshot_uuid").(string)

	appResp, err := conn.Service.GetApp(ctx, appUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	var appSpec map[string]interface{}
	if err := json.Unmarshal(appResp.Spec, &appSpec); err != nil {
		fmt.Println("Error unmarshalling Spec:", err)
	}

	var appMetadata map[string]interface{}
	if err := json.Unmarshal(appResp.Metadata, &appMetadata); err != nil {
		fmt.Println("Error unmarshalling Spec to get metadata:", err)
	}

	var appStatus map[string]interface{}
	if err := json.Unmarshal(appResp.Status, &appStatus); err != nil {
		fmt.Println("Error unmarshalling Spec to get status:", err)
	}

	uuid, err := uuid.GenerateUUID()

	appMetadata["uuid"] = uuid
	delete(appMetadata, "owner_reference")

	// create spec
	fetchSpec := &calm.TaskSpec{}
	fetchSpec.TargetUUID = appUUID
	fetchSpec.TargetKind = "Application"
	fetchSpec.Args.Variables = []*calm.VariableList{}

	restoreConfig := &calm.VariableList{}

	restoreConfig.Name = "recovery_point_group_uuid"
	restoreConfig.Value = recoveryPointUUID
	restoreActionUUID, restoreActionTaskUuid := fetchRestoreActionUUID(appStatus, restoreActionName)
	if restoreActionUUID == "" {
		return diag.Errorf("UUID for restore action with name %s not found.", restoreActionName)
	}
	if restoreActionTaskUuid == "" {
		return diag.Errorf("UUID for restore action task with name %s not found.", restoreActionName)
	}
	restoreConfig.TaskUUID = restoreActionUUID

	fetchInput := &calm.ActionInput{}
	fetchInput.APIVersion = appResp.APIVersion
	fetchInput.Metadata = appMetadata
	fetchInput.Spec = *fetchSpec

	return nil
}

func resourceNutanixCalmAppRestoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmAppRestoreUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmAppRestoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func fetchRestoreActionUUID(appStatus map[string]interface{}, restoreActionName string) (string, string) {
	var restoreActionTaskUuid string
	var restoreActionUuid string
	if resources, ok := appStatus["resources"].(map[string]interface{}); ok {
		if actionList, ok := resources["action_list"].([]interface{}); ok {
			for _, action := range actionList {
				if act, ok := action.(map[string]interface{}); ok {
					if act["name"] == restoreActionName {
						if runbook, ok := act["runbook"].(map[string]interface{}); ok {
							if taskDefinitionList, ok := runbook["task_definition_list"].([]interface{}); ok {
								for _, taskDef := range taskDefinitionList {
									if task, ok := taskDef.(map[string]interface{}); ok {
										if task["type"].(string) == "CALL_CONFIG" {
											restoreActionTaskUuid = task["uuid"].(string)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return restoreActionUuid, restoreActionTaskUuid
}
