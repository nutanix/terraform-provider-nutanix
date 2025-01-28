package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
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
	if err != nil {
		return diag.Errorf("Could not generate uuid with error: %s", err.Error())
	}

	appMetadata["uuid"] = uuid
	delete(appMetadata, "owner_reference")

	// create spec
	restoreSpec := &calm.TaskSpec{}
	restoreSpec.TargetUUID = appUUID
	restoreSpec.TargetKind = "Application"
	restoreSpec.Args = []*calm.VariableList{}

	restoreConfig := &calm.VariableList{}

	restoreConfig.Name = "recovery_point_group_uuid"
	restoreConfig.Value = recoveryPointUUID
	restoreActionUuid, restoreActionTaskUuid := fetchRestoreActionUUID(appStatus, restoreActionName)
	if restoreActionUuid == "" {
		return diag.Errorf("UUID for restore action with name %s not found.", restoreActionName)
	}
	if restoreActionTaskUuid == "" {
		return diag.Errorf("UUID for restore action task with name %s not found.", restoreActionName)
	}
	restoreConfig.TaskUUID = restoreActionTaskUuid

	restoreSpec.Args = append(restoreSpec.Args, restoreConfig)

	restoreInput := &calm.ActionInput{}
	restoreInput.APIVersion = appResp.APIVersion
	restoreInput.Metadata = appMetadata
	restoreInput.Spec = *restoreSpec

	restoreResp, err := conn.Service.PerformActionUuid(ctx, appUUID, restoreActionUuid, restoreInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := restoreResp.Status.RunlogUUID

	fmt.Println("Runlog UUID:", runlogUUID)

	d.SetId(runlogUUID)

	// poll till action is completed
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "POLICY_EXEC", "ABORTING", "APPROVAL"},
		Target:  []string{"SUCCESS", "FAILURE", "WARNING", "ERROR", "SYS_FAILURE", "SYS_ERROR", "SYS_ABORTED", "TIMEOUT", "APPROVAL_FAILED"},
		Refresh: RestoreStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   5 * time.Second,
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Error waiting for app to perform Restore Action: %s", errWaitTask)
	}

	v, err := conn.Service.AppRunlogs(ctx, appUUID, runlogUUID)
	if err != nil {
		diag.Errorf("Error in getting runlog output: %s", err.Error())
	}

	runlogState := utils.StringValue(v.Status.RunlogState)

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
	var restoreActionTaskUuid string
	var restoreActionUuid string
	if resources, ok := appStatus["resources"].(map[string]interface{}); ok {
		if actionList, ok := resources["action_list"].([]interface{}); ok {
			for _, action := range actionList {
				if act, ok := action.(map[string]interface{}); ok {
					if act["name"].(string) == restoreActionName {
						restoreActionUuid = act["uuid"].(string)
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

func RestoreStateRefreshFunc(ctx context.Context, client *calm.Client, appUUID, runlogUUID string) resource.StateRefreshFunc {
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

		return v, runlogstate, nil
	}
}
