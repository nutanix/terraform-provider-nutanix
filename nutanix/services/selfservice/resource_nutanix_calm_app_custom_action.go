package selfservice

import (
	"context"
	"encoding/json"
	"fmt"
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
			"app_uuid": {
				Type:     schema.TypeString,
				Required: true,
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

	appUUID := d.Get("app_uuid").(string)
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

	fetchResp, err := conn.Service.PerformActionUuid(ctx, appUUID, actionUUID, fetchInput)
	if err != nil {
		return diag.FromErr(err)
	}

	runlogUUID := fetchResp.Status.RunlogUUID

	fmt.Println("Response:", runlogUUID)
	d.SetId(runlogUUID)
	// poll till action is completed
	appStateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "POLICY_EXEC", "ABORTING", "APPROVAL"},
		Target:  []string{"SUCCESS", "FAILURE", "WARNING", "ERROR", "SYS_FAILURE", "SYS_ERROR", "SYS_ABORTED", "TIMEOUT", "APPROVAL_FAILED"},
		Refresh: ActionStateRefreshFunc(ctx, conn, appUUID, runlogUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
		Delay:   5 * time.Second,
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
	calm_action_name := "action_" + strings.ToLower(actionName)
	if resource, ok := pr["resources"].(map[string]interface{}); ok {
		// fmt.Println("RESOURCESSSSS")
		if actionList, ok := resource["action_list"].([]interface{}); ok {
			for _, action := range actionList {
				if dep, ok := action.(map[string]interface{}); ok {
					fmt.Println("DEP UUID::::", dep["uuid"])
					if dep["name"] == actionName || dep["name"] == calm_action_name {
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
