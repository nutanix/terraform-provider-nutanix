package selfservice

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"strconv"
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

func ResourceNutanixCalmAppRecoveryPoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCalmAppRecoveryPointCreate,
		ReadContext:   resourceNutanixCalmAppRecoveryPointRead,
		UpdateContext: resourceNutanixCalmAppRecoveryPointUpdate,
		DeleteContext: resourceNutanixCalmAppRecoveryPointDelete,
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"recovery_point_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNutanixCalmAppRecoveryPointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	snapshotActionName := d.Get("action_name").(string)
	snapshotName := d.Get("recovery_point_name").(string)

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

	uuid, _ := uuid.GenerateUUID()

	appMetadata["uuid"] = uuid
	delete(appMetadata, "owner_reference")

	// create spec
	snapshotSpec := &selfservice.TaskSpec{}
	snapshotSpec.TargetUUID = appUUID
	snapshotSpec.TargetKind = "Application"
	snapshotSpec.Args = []*selfservice.VariableList{}

	snapshotConfig := &selfservice.VariableList{}

	snapshotConfig.Name = "snapshot_name"
	snapshotConfig.Value = snapshotName
	snapshotActionUUID, snapshotActionTaskUUID := fetchSnapshotActionUUID(appStatus, snapshotActionName)
	if snapshotActionUUID == "" {
		return diag.Errorf("UUID for snapshot action with name %s not found.", snapshotActionName)
	}
	if snapshotActionTaskUUID == "" {
		return diag.Errorf("UUID for snapshot action task with name %s not found.", snapshotActionName)
	}
	snapshotConfig.TaskUUID = snapshotActionTaskUUID

	snapshotSpec.Args = append(snapshotSpec.Args, snapshotConfig)

	snapshotInput := &selfservice.ActionInput{}
	snapshotInput.APIVersion = appResp.APIVersion
	snapshotInput.Metadata = appMetadata
	snapshotInput.Spec = *snapshotSpec

	snapshotResp, err := conn.Service.PerformActionUUID(ctx, appUUID, snapshotActionUUID, snapshotInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := snapshotResp.Status.RunlogUUID

	log.Println("[DEBUG] Runlog UUID:", runlogUUID)
	d.SetId(runlogUUID)
	// poll till action is completed
	const delayDuration = 5 * time.Second
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "POLICY_EXEC", "ABORTING", "APPROVAL"},
		Target:  []string{"SUCCESS", "FAILURE", "WARNING", "ERROR", "SYS_FAILURE", "SYS_ERROR", "SYS_ABORTED", "TIMEOUT", "APPROVAL_FAILED"},
		Refresh: SnapshotStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   delayDuration,
	}
	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Error waiting for app to perform Restore Action: %s", errWaitTask)
	}

	return nil
}

func resourceNutanixCalmAppRecoveryPointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmAppRecoveryPointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixCalmAppRecoveryPointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		log.Println("Error unmarshalling AppName:", err)
		return diag.FromErr(err)
	}

	entities := AppNameStatus[0].(map[string]interface{})

	if entity, ok := entities["metadata"].(map[string]interface{}); ok {
		appUUID = entity["uuid"].(string)
	}
	log.Println("[Debug] App uuid: ", appUUID)

	snapshotName := d.Get("recovery_point_name").(string)
	log.Printf("DELETE CALLED FOR %s %s", appUUID, snapshotName)
	length := 250
	offset := 0
	appResp, err := conn.Service.GetApp(ctx, appUUID)
	if err != nil {
		return diag.FromErr(err)
	}
	var appStatus map[string]interface{}
	if err = json.Unmarshal(appResp.Status, &appStatus); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec to get status:", err)
		return diag.FromErr(err)
	}

	var appMetadata map[string]interface{}
	if err = json.Unmarshal(appResp.Metadata, &appMetadata); err != nil {
		log.Println("[DEBUG] Error unmarshalling Spec to get metadata:", err)
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

	var snapshotGroupID string

	foundRecoveryPoint := false

	for _, entity := range listResp.Entities {
		if status, ok := entity["status"].(map[string]interface{}); ok {
			if recoveryPointInfoList, ok := status["recovery_point_info_list"].([]interface{}); ok {
				for _, recoveryPoint := range recoveryPointInfoList {
					if snapshotName == recoveryPoint.(map[string]interface{})["name"].(string) {
						snapshotGroupID = status["uuid"].(string)
						foundRecoveryPoint = true
						break
					}
				}
			}
			if foundRecoveryPoint {
				break
			}
		}
	}

	snapshotSpec := &selfservice.TaskSpec{}
	snapshotSpec.TargetUUID = appUUID
	snapshotSpec.TargetKind = "Application"
	snapshotSpec.Args = []*selfservice.VariableList{}

	snapshotConfig := &selfservice.VariableList{}
	snapshotConfig.Name = "snapshot_group_id"
	snapshotConfig.Value = snapshotGroupID
	snapshotSpec.Args = append(snapshotSpec.Args, snapshotConfig)

	snapshotInput := &selfservice.ActionInput{}
	snapshotInput.APIVersion = appResp.APIVersion
	snapshotInput.Metadata = appMetadata
	snapshotInput.Spec = *snapshotSpec

	snapshotResp, err := conn.Service.RecoveryPointsDelete(ctx, appUUID, snapshotInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := snapshotResp.Status.RunlogUUID

	log.Println("[DEBUG] Trigger delete of snapshot with Runlog UUID:", runlogUUID)

	return nil
}

func fetchSnapshotActionUUID(appStatus map[string]interface{}, snapshotActionName string) (string, string) {
	var snapshotActionTaskUUID string
	var snapshotActionUUID string
	if resources, ok := appStatus["resources"].(map[string]interface{}); ok {
		if actionList, ok := resources["action_list"].([]interface{}); ok {
			for _, action := range actionList {
				if act, ok := action.(map[string]interface{}); ok {
					if act["name"].(string) == snapshotActionName {
						snapshotActionUUID = act["uuid"].(string)
						if runbook, ok := act["runbook"].(map[string]interface{}); ok {
							if taskDefinitionList, ok := runbook["task_definition_list"].([]interface{}); ok {
								for _, taskDef := range taskDefinitionList {
									if task, ok := taskDef.(map[string]interface{}); ok {
										if task["type"].(string) == "CALL_CONFIG" {
											snapshotActionTaskUUID = task["uuid"].(string)
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
	return snapshotActionUUID, snapshotActionTaskUUID
}

func SnapshotStateRefreshFunc(ctx context.Context, client *selfservice.Client, appUUID, runlogUUID string) resource.StateRefreshFunc {
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
