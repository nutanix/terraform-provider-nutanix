package lcmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/common"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceLcmUpgradeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceLcmUpgradeV2Create,
		ReadContext:   ResourceLcmUpgradeV2Read,
		UpdateContext: ResourceLcmUpgradeV2Update,
		DeleteContext: ResourceLcmUpgradeV2Delete,
		Schema: map[string]*schema.Schema{
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"management_server": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hypervisor_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"HYPERV", "ESX", "AHV"}, false),
						},
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"entity_update_specs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"to_version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"skipped_precheck_flags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"auto_handle_flags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"max_wait_time_in_secs": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(60, 86400), //nolint:gomnd
			},
		},
	}
}

func ResourceLcmUpgradeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	var clusterID *string
	if id := d.Get("x_cluster_id").(string); id != "" {
		clusterID = &id
	}

	body := common.NewUpgradeSpec()

	if managementServer, ok := d.GetOk("management_server"); ok && len(managementServer.([]interface{})) > 0 {
		body.ManagementServer = expandManagementServer(managementServer)
	}
	if entityUpdateSpecs, ok := d.GetOk("entity_update_specs"); ok && len(entityUpdateSpecs.([]interface{})) > 0 {
		body.EntityUpdateSpecs = expandEntityUpdateSpecs(entityUpdateSpecs.([]interface{}))
	}
	if skippedPrecheckFlags, ok := d.GetOk("skipped_precheck_flags"); ok && len(skippedPrecheckFlags.([]interface{})) > 0 {
		body.SkippedPrecheckFlags = expandSystemAutoMgmtFlag(skippedPrecheckFlags.([]interface{}))
	}
	if autoHandleFlags, ok := d.GetOk("auto_handle_flags"); ok && len(autoHandleFlags.([]interface{})) > 0 {
		body.AutoHandleFlags = expandSystemAutoMgmtFlag(autoHandleFlags.([]interface{}))
	}
	if maxWaitTimeInSecs, ok := d.GetOk("max_wait_time_in_secs"); ok && maxWaitTimeInSecs.(int) > 0 {
		body.MaxWaitTimeInSecs = utils.IntPtr(maxWaitTimeInSecs.(int))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] LCM Upgrade Request Spec: %s", string(aJSON))
	// pass nil for the new dyRun flag
	resp, err := conn.LcmUpgradeAPIInstance.PerformUpgrade(body, clusterID, nil)
	if err != nil {
		return diag.Errorf("error while Perform Upgrade the LCM config: %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the LCM upgrade to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for LCM upgrade (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching LCM upgrade task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] LCM Upgrade Task Details: %s", string(aJSON))

	// This is an action resource that does not maintain state.
	// The resource ID is set to the task ExtId for traceability.
	d.SetId(utils.StringValue(taskDetails.ExtId))
	return nil
}

func ResourceLcmUpgradeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceLcmUpgradeV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceLcmUpgradeV2Create(ctx, d, meta)
}

func ResourceLcmUpgradeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandManagementServer(managementServer interface{}) *common.ManagementServer {
	managementServerList := managementServer.([]interface{})

	if len(managementServerList) == 0 {
		return nil
	}
	managementServerMap := managementServerList[0].(map[string]interface{})
	managementServerObj := common.NewManagementServer()

	managementServerObj.HypervisorType = expandHypervisorType(managementServerMap["hypervisor_type"].(string))
	managementServerObj.Ip = utils.StringPtr(managementServerMap["ip"].(string))
	managementServerObj.Username = utils.StringPtr(managementServerMap["username"].(string))
	managementServerObj.Password = utils.StringPtr(managementServerMap["password"].(string))

	return managementServerObj
}

func expandHypervisorType(hypervisorType string) *common.HypervisorType {
	switch hypervisorType {
	case "HYPERV":
		p := common.HYPERVISORTYPE_HYPERV
		return &p
	case "ESX":
		p := common.HYPERVISORTYPE_ESX
		return &p
	case "AHV":
		p := common.HYPERVISORTYPE_AHV
		return &p
	}
	return nil
}

func expandEntityUpdateSpecs(entityUpdateSpec []interface{}) []common.EntityUpdateSpec {
	if len(entityUpdateSpec) == 0 {
		return nil
	}

	entityUpdateSpecsList := make([]common.EntityUpdateSpec, 0)

	for _, entityUpdateSpecItem := range entityUpdateSpec {
		entityUpdateSpecMap := entityUpdateSpecItem.(map[string]interface{})
		entityUpdateSpecObj := common.NewEntityUpdateSpec()
		entityUpdateSpecObj.EntityUuid = utils.StringPtr(entityUpdateSpecMap["entity_uuid"].(string))
		entityUpdateSpecObj.ToVersion = utils.StringPtr(entityUpdateSpecMap["to_version"].(string))
		entityUpdateSpecsList = append(entityUpdateSpecsList, *entityUpdateSpecObj)
	}
	return entityUpdateSpecsList
}

func expandSystemAutoMgmtFlag(flags []interface{}) []common.SystemAutoMgmtFlag {
	if len(flags) == 0 {
		return nil
	}

	systemAutoMgmtFlags := make([]common.SystemAutoMgmtFlag, 0)
	for _, flag := range flags {
		if flag == "POWER_OFF_UVMS" {
			p := common.SYSTEMAUTOMGMTFLAG_POWER_OFF_UVMS
			systemAutoMgmtFlags = append(systemAutoMgmtFlags, p)
		}
	}
	return systemAutoMgmtFlags
}
