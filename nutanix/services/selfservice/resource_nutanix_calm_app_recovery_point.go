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
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
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
	conn := meta.(*conns.Client).Calm

	var appUUID string

	app_name := d.Get("app_name").(string)

	appFilter := &calm.ApplicationListInput{}

	appFilter.Filter = fmt.Sprintf("name==%s;_state!=deleted", app_name)

	log.Printf("[Debug] Qeurying apps/list API with filter %s", appFilter)

	appNameResp, err := conn.Service.ListApplication(ctx, appFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[Debug] Getting app uuid from app response: %s", appNameResp)

	var AppNameStatus []interface{}
	if err := json.Unmarshal([]byte(appNameResp.Entities), &AppNameStatus); err != nil {
		fmt.Println("Error unmarshalling AppName:", err)
	}

	entities := AppNameStatus[0].(map[string]interface{})

	if entity, ok := entities["metadata"].(map[string]interface{}); ok {
		appUUID = entity["uuid"].(string)
	}

	if appUUID, ok := d.GetOk("app_uuid"); ok {
		appUUID = appUUID.(string)
	}

	snapshotActionName := d.Get("action_name").(string)
	snapshotName := d.Get("recovery_point_name").(string)

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
	snapshotSpec := &calm.TaskSpec{}
	snapshotSpec.TargetUUID = appUUID
	snapshotSpec.TargetKind = "Application"
	snapshotSpec.Args = []*calm.VariableList{}

	snapshotConfig := &calm.VariableList{}

	snapshotConfig.Name = "snapshot_name"
	snapshotConfig.Value = snapshotName
	snapshotActionUUID, snapshotActionTaskUuid := fetchSnapshotActionUUID(appStatus, snapshotActionName)
	if snapshotActionUUID == "" {
		return diag.Errorf("UUID for snapshot action with name %s not found.", snapshotActionName)
	}
	if snapshotActionTaskUuid == "" {
		return diag.Errorf("UUID for snapshot action task with name %s not found.", snapshotActionName)
	}
	snapshotConfig.TaskUUID = snapshotActionTaskUuid

	snapshotSpec.Args = append(snapshotSpec.Args, snapshotConfig)

	snapshotInput := &calm.ActionInput{}
	snapshotInput.APIVersion = appResp.APIVersion
	snapshotInput.Metadata = appMetadata
	snapshotInput.Spec = *snapshotSpec

	snapshotResp, err := conn.Service.PerformActionUuid(ctx, appUUID, snapshotActionUUID, snapshotInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := snapshotResp.Status.RunlogUUID

	fmt.Println("Runlog UUID:", runlogUUID)
	d.SetId(runlogUUID)
	// poll till action is completed
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "POLICY_EXEC", "ABORTING", "APPROVAL"},
		Target:  []string{"SUCCESS", "FAILURE", "WARNING", "ERROR", "SYS_FAILURE", "SYS_ERROR", "SYS_ABORTED", "TIMEOUT", "APPROVAL_FAILED"},
		Refresh: SnapshotStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   5 * time.Second,
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

	conn := meta.(*conns.Client).Calm

	var appUUID string

	app_name := d.Get("app_name").(string)

	appFilter := &calm.ApplicationListInput{}

	appFilter.Filter = fmt.Sprintf("name==%s;_state!=deleted", app_name)

	log.Printf("[Debug] Qeurying apps/list API with filter %s", appFilter)

	appNameResp, err := conn.Service.ListApplication(ctx, appFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[Debug] Getting app uuid from app response: %s", appNameResp)

	var AppNameStatus []interface{}
	if err := json.Unmarshal([]byte(appNameResp.Entities), &AppNameStatus); err != nil {
		log.Println("Error unmarshalling AppName:", err)
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
	if err := json.Unmarshal(appResp.Status, &appStatus); err != nil {
		fmt.Println("Error unmarshalling Spec to get status:", err)
	}

	var appMetadata map[string]interface{}
	if err := json.Unmarshal(appResp.Metadata, &appMetadata); err != nil {
		fmt.Println("Error unmarshalling Spec to get metadata:", err)
	}

	substrateReference := fetchSubstrateReference(appStatus)

	currTime := strconv.FormatInt(time.Now().Unix(), 10)

	listInput := &calm.RecoveryPointsListInput{}

	listInput.Filter = fmt.Sprintf("substrate_reference==%s;expiration_time=ge=%s", substrateReference, currTime)
	listInput.Length = length
	listInput.Offset = offset

	listResp, err := conn.Service.RecoveryPointsList(ctx, appUUID, listInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var snapshotGroupId string

	foundRecoveryPoint := false

	for _, entity := range listResp.Entities {
		if status, ok := entity["status"].(map[string]interface{}); ok {
			if recovery_point_info_list, ok := status["recovery_point_info_list"].([]interface{}); ok {
				for _, recovery_point := range recovery_point_info_list {
					if snapshotName == recovery_point.(map[string]interface{})["name"].(string) {
						snapshotGroupId = status["uuid"].(string)
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

	snapshotSpec := &calm.TaskSpec{}
	snapshotSpec.TargetUUID = appUUID
	snapshotSpec.TargetKind = "Application"
	snapshotSpec.Args = []*calm.VariableList{}

	snapshotConfig := &calm.VariableList{}
	snapshotConfig.Name = "snapshot_group_id"
	snapshotConfig.Value = snapshotGroupId
	snapshotSpec.Args = append(snapshotSpec.Args, snapshotConfig)

	snapshotInput := &calm.ActionInput{}
	snapshotInput.APIVersion = appResp.APIVersion
	snapshotInput.Metadata = appMetadata
	snapshotInput.Spec = *snapshotSpec

	snapshotResp, err := conn.Service.RecoveryPointsDelete(ctx, appUUID, snapshotInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := snapshotResp.Status.RunlogUUID

	fmt.Println("[DEBUG] Trigger delete of snapshot with Runlog UUID:", runlogUUID)

	return nil
}

func fetchSnapshotActionUUID(appStatus map[string]interface{}, snapshotActionName string) (string, string) {
	var snapshotActionTaskUuid string
	var snapshotActionUuid string
	if resources, ok := appStatus["resources"].(map[string]interface{}); ok {
		if actionList, ok := resources["action_list"].([]interface{}); ok {
			for _, action := range actionList {
				if act, ok := action.(map[string]interface{}); ok {
					if act["name"].(string) == snapshotActionName {
						snapshotActionUuid = act["uuid"].(string)
						if runbook, ok := act["runbook"].(map[string]interface{}); ok {
							if taskDefinitionList, ok := runbook["task_definition_list"].([]interface{}); ok {
								for _, taskDef := range taskDefinitionList {
									if task, ok := taskDef.(map[string]interface{}); ok {
										if task["type"].(string) == "CALL_CONFIG" {
											snapshotActionTaskUuid = task["uuid"].(string)
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
	return snapshotActionUuid, snapshotActionTaskUuid
}

func SnapshotStateRefreshFunc(ctx context.Context, client *calm.Client, appUUID, runlogUUID string) resource.StateRefreshFunc {
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
