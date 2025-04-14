package selfservice

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
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixCalmAppCustomAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixCalmAppCustomActionCreate,
		ReadContext:   ResourceNutanixCalmAppCustomActionRead,
		UpdateContext: ResourceNutanixCalmAppCustomActionUpdate,
		DeleteContext: ResourceNutanixCalmAppCustomActionDelete,
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
			"runlog_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixCalmAppCustomActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm

	var appUUID string

	appName := d.Get("app_name").(string)

	appFilter := &calm.ApplicationListInput{}

	appFilter.Filter = fmt.Sprintf("name==%s;_state!=deleted", appName)

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

	if appUUIDRead, ok := d.GetOk("app_uuid"); ok {
		appUUID = appUUIDRead.(string)
	}

	actionName := d.Get("action_name").(string)

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

	fetchInput := &calm.ActionInput{}
	fetchInput.APIVersion = appResp.APIVersion
	fetchInput.Metadata = objMetadata

	var actionUUID string
	// fetch patch for spec
	fetchSpec := &calm.TaskSpec{}
	fetchSpec.TargetUUID = appUUID
	fetchSpec.TargetKind = "Application"
	fetchSpec.Args = []*calm.VariableList{}
	_, actionUUID = expandCustomActionSpec(objSpec, actionName)

	fetchInput.Spec = *fetchSpec

	fetchResp, err := conn.Service.PerformActionUUID(ctx, appUUID, actionUUID, fetchInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := fetchResp.Status.RunlogUUID

	fmt.Println("Response:", runlogUUID)
	d.SetId(runlogUUID)
	// poll till action is completed
	const delayDuration = 5 * time.Second
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "POLICY_EXEC", "ABORTING", "APPROVAL"},
		Target:  []string{"SUCCESS", "FAILURE", "WARNING", "ERROR", "SYS_FAILURE", "SYS_ERROR", "SYS_ABORTED", "TIMEOUT", "APPROVAL_FAILED"},
		Refresh: ActionStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   delayDuration,
	}

	if _, errWaitTask := appStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Error waiting for app to perform Restore Action: %s", errWaitTask)
	}
	if err := d.Set("runlog_uuid", runlogUUID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixCalmAppCustomActionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func ResourceNutanixCalmAppCustomActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func ResourceNutanixCalmAppCustomActionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandCustomActionSpec(pr map[string]interface{}, actionName string) (map[string]interface{}, string) {
	calmActionName := "action_" + strings.ToLower(actionName)
	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		// fmt.Println("RESOURCESSSSS")
		if actionList, ok := resource["action_list"].([]interface{}); ok {
			for _, action := range actionList {
				if dep, ok := action.(map[string]interface{}); ok {
					fmt.Println("DEP UUID::::", dep["uuid"])
					if dep["name"] == actionName || dep["name"] == calmActionName {
						fmt.Println("DEP UUID::::", dep["uuid"])
						return action.(map[string]interface{}), dep["uuid"].(string)
					}
				}
			}
		}
	}
	return nil, ""
}

func ActionStateRefreshFunc(ctx context.Context, client *calm.Client, appUUID, runlogUUID string) resource.StateRefreshFunc {
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
