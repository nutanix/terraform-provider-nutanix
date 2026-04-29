package monitoringv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitoringModel "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixSdaClusterConfigV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixSdaClusterConfigV2Create,
		ReadContext:   resourceNutanixSdaClusterConfigV2Read,
		UpdateContext: resourceNutanixSdaClusterConfigV2Update,
		DeleteContext: resourceNutanixSdaClusterConfigV2Delete,
		Schema: map[string]*schema.Schema{
			"system_defined_policy_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique ID of the System-Defined Alert Policy.",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster UUID.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
			},
			"links": schemaForLinks(),
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the SDA policy is enabled or not on the cluster.",
			},
			"last_modified_by_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the user who made the latest update to this policy. Its value will be Nutanix if the last update is due to an upgrade event.",
			},
			"last_modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time in ISO 8601 format when the SDA policy was last modified. It gets automatically updated by the Nutanix service from the user context during an update event.",
			},
			"schedule_interval_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Interval in seconds for periodically executing the SDA policy. This will not be set for policies with the type NOT_SCHEDULED & EVENT_DRIVEN.",
			},
			"alert_config":              schemaForAlertConfigResource(),
			"configurable_parameters":   schemaForConfigurableParametersResource(),
		},
	}
}

func resourceNutanixSdaClusterConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)
	extID := d.Get("ext_id").(string)

	readResp, err := conn.SystemDefinedPoliciesAPI.GetClusterConfigById(
		utils.StringPtr(sdaPolicyExtID),
		utils.StringPtr(extID),
	)
	if err != nil {
		return diag.Errorf("error while fetching SDA cluster config for create: %v", err)
	}

	etagValue := conn.SystemDefinedPoliciesAPI.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := readResp.Data.GetValue().(monitoringModel.ClusterConfig)

	if v, ok := d.GetOk("is_enabled"); ok {
		body.IsEnabled = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("alert_config"); ok {
		alertConfigList := v.([]interface{})
		if len(alertConfigList) > 0 && alertConfigList[0] != nil {
			body.AlertConfig = expandAlertConfig(alertConfigList[0].(map[string]interface{}))
		}
	}
	if v, ok := d.GetOk("configurable_parameters"); ok {
		body.ConfigurableParameters = expandConfigurableParameters(v.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Update SDA Cluster Config (Create) Payload: %s", string(aJSON))

	resp, err := conn.SystemDefinedPoliciesAPI.UpdateClusterConfigById(
		utils.StringPtr(sdaPolicyExtID),
		utils.StringPtr(extID),
		&body,
		args,
	)
	if err != nil {
		return diag.Errorf("error while updating SDA cluster config: %v", err)
	}

	TaskRefVal := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRefVal.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for SDA cluster config (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching SDA cluster config update task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] SDA Cluster Config Update Task Details: %s", string(aJSON))

	d.SetId(extID)
	return resourceNutanixSdaClusterConfigV2Read(ctx, d, meta)
}

func resourceNutanixSdaClusterConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)
	extID := d.Id()

	resp, err := conn.SystemDefinedPoliciesAPI.GetClusterConfigById(
		utils.StringPtr(sdaPolicyExtID),
		utils.StringPtr(extID),
	)
	if err != nil {
		return diag.Errorf("error while reading SDA cluster config: %v", err)
	}

	body := resp.Data.GetValue().(monitoringModel.ClusterConfig)
	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Read SDA Cluster Config Response: %s", string(aJSON))

	if err := flattenClusterConfigToState(d, body); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNutanixSdaClusterConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)
	extID := d.Id()

	readResp, err := conn.SystemDefinedPoliciesAPI.GetClusterConfigById(
		utils.StringPtr(sdaPolicyExtID),
		utils.StringPtr(extID),
	)
	if err != nil {
		return diag.Errorf("error while fetching SDA cluster config for update: %v", err)
	}

	etagValue := conn.SystemDefinedPoliciesAPI.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := readResp.Data.GetValue().(monitoringModel.ClusterConfig)

	if d.HasChange("is_enabled") {
		body.IsEnabled = utils.BoolPtr(d.Get("is_enabled").(bool))
	}
	if d.HasChange("alert_config") {
		if v, ok := d.GetOk("alert_config"); ok {
			alertConfigList := v.([]interface{})
			if len(alertConfigList) > 0 && alertConfigList[0] != nil {
				body.AlertConfig = expandAlertConfig(alertConfigList[0].(map[string]interface{}))
			}
		}
	}
	if d.HasChange("configurable_parameters") {
		if v, ok := d.GetOk("configurable_parameters"); ok {
			body.ConfigurableParameters = expandConfigurableParameters(v.([]interface{}))
		}
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Update SDA Cluster Config Payload: %s", string(aJSON))

	resp, err := conn.SystemDefinedPoliciesAPI.UpdateClusterConfigById(
		utils.StringPtr(sdaPolicyExtID),
		utils.StringPtr(extID),
		&body,
		args,
	)
	if err != nil {
		return diag.Errorf("error while updating SDA cluster config: %v", err)
	}

	TaskRefVal := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRefVal.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for SDA cluster config (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return resourceNutanixSdaClusterConfigV2Read(ctx, d, meta)
}

func resourceNutanixSdaClusterConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
