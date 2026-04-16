package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/selfservice"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixCalmAppRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmAppRestoreCreate,
		ReadContext:   resourceNutanixCalmAppRestoreRead,
		UpdateContext: resourceNutanixCalmAppRestoreUpdate,
		DeleteContext: resourceNutanixCalmAppRestoreDelete,
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restore_action_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixCalmAppRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	restoreActionName := d.Get("restore_action_name").(string)
	recoveryPointUUID := d.Get("snapshot_uuid").(string)

	appResp, err := conn.Service.GetApp(ctx, appUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	var appSpec map[string]interface{}
	if err = json.Unmarshal(appResp.Spec, &appSpec); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec:", err)
		return diag.FromErr(err)
	}

	var appMetadata map[string]interface{}
	if err = json.Unmarshal(appResp.Metadata, &appMetadata); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec to get metadata:", err)
		return diag.FromErr(err)
	}

	var appStatus map[string]interface{}
	if err = json.Unmarshal(appResp.Status, &appStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec to get status:", err)
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("Could not generate uuid with error: %s", err.Error())
	}

	appMetadata["uuid"] = uuid
	delete(appMetadata, "owner_reference")

	// create spec
	restoreSpec := &selfservice.TaskSpec{}
	restoreSpec.TargetUUID = appUUID
	restoreSpec.TargetKind = "Application"
	restoreSpec.Args = []*selfservice.VariableList{}

	restoreConfig := &selfservice.VariableList{}

	restoreConfig.Name = "recovery_point_group_uuid"
	restoreConfig.Value = recoveryPointUUID
	restoreActionUUID, restoreActionTaskUUID := fetchRestoreActionUUID(appStatus, restoreActionName)
	if restoreActionUUID == "" {
		return diag.Errorf("UUID for restore action with name %s not found.", restoreActionName)
	}
	if restoreActionTaskUUID == "" {
		return diag.Errorf("UUID for restore action task with name %s not found.", restoreActionName)
	}
	restoreConfig.TaskUUID = restoreActionTaskUUID

	restoreSpec.Args = append(restoreSpec.Args, restoreConfig)

	restoreInput := &selfservice.ActionInput{}
	restoreInput.APIVersion = appResp.APIVersion
	restoreInput.Metadata = appMetadata
	restoreInput.Spec = *restoreSpec

	restoreResp, err := conn.Service.PerformActionUUID(ctx, appUUID, restoreActionUUID, restoreInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := restoreResp.Status.RunlogUUID

	log.Println("[DEBUG] Runlog UUID:", runlogUUID)

	d.SetId(runlogUUID)

	// poll till action is completed
	const delayDuration = 5 * time.Second
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "POLICY_EXEC", "ABORTING", "APPROVAL"},
		Target:  []string{"SUCCESS", "FAILURE", "WARNING", "ERROR", "SYS_FAILURE", "SYS_ERROR", "SYS_ABORTED", "TIMEOUT", "APPROVAL_FAILED"},
		Refresh: RestoreStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   delayDuration,
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Error waiting for app to perform Restore Action: %s", errWaitTask)
	}

	v, err := conn.Service.AppRunlogs(ctx, appUUID, runlogUUID)
	if err != nil {
		diag.Errorf("Error in getting runlog output: %s", err.Error())
	}

	runlogState := utils.StringValue(v.Status.RunlogState)

	if err := d.Set("snapshot_uuid", recoveryPointUUID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("state", runlogState); err != nil {
		return diag.FromErr(err)
	}

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
	var restoreActionTaskUUID string
	var restoreActionUUID string
	if resources, ok := appStatus["resources"].(map[string]interface{}); ok {
		if actionList, ok := resources["action_list"].([]interface{}); ok {
			for _, action := range actionList {
				if act, ok := action.(map[string]interface{}); ok {
					if act["name"].(string) == restoreActionName {
						restoreActionUUID = act["uuid"].(string)
						if runbook, ok := act["runbook"].(map[string]interface{}); ok {
							if taskDefinitionList, ok := runbook["task_definition_list"].([]interface{}); ok {
								for _, taskDef := range taskDefinitionList {
									if task, ok := taskDef.(map[string]interface{}); ok {
										if task["type"].(string) == "CALL_CONFIG" {
											restoreActionTaskUUID = task["uuid"].(string)
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
	return restoreActionUUID, restoreActionTaskUUID
}

func RestoreStateRefreshFunc(ctx context.Context, client *selfservice.Client, appUUID, runlogUUID string) resource.StateRefreshFunc {
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
