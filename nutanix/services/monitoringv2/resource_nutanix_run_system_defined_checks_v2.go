package monitoringv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/config"
	monitoringRequest "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/request/systemdefinedchecks"
	monitoringServiceability "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	monitoringTaskRef "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRunSystemDefinedChecksV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixRunSystemDefinedChecksV2Create,
		ReadContext:   resourceNutanixRunSystemDefinedChecksV2Read,
		UpdateContext: resourceNutanixRunSystemDefinedChecksV2Update,
		DeleteContext: resourceNutanixRunSystemDefinedChecksV2Delete,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for the cluster for which run System-Defined Checks is requested.",
			},
			"additional_recipients": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of additional email addresses for sending the run summary. Either this should be set or should_send_report_to_configured_recipients should be true. If both are set then email would be sent to all the recipients.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"node_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of node IP addresses where the Check will run. This field will be ignored if the check scope is a cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The prefix length of the network to which this host IPv4 address belongs.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IPv4 address of the host.",
						},
					},
				},
			},
			"sda_ext_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Check IDs to be executed. This field cannot be set simultaneously with should_run_all_checks; only one of them should be specified.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"should_anonymize": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to mask sensitive data in the check run summary.",
			},
			"should_run_all_checks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether all System-Defined Checks applicable to the specified cluster should be executed. This field is mutually exclusive with the sda_ext_ids parameter, meaning that only one of these should be set at a time. Please use this field with caution, as it is resource-intensive.",
			},
			"should_send_report_to_configured_recipients": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if the run summary should be sent to the configured email address associated with the cluster. Either this should be true or additional_recipients should be provided. If both are set then email would be sent to all the recipients.",
			},
			"task_ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier for the task created by running system-defined checks.",
			},
		},
	}
}

func resourceNutanixRunSystemDefinedChecksV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	body := monitoringServiceability.NewRunSystemDefinedChecksSpec()

	if v, ok := d.GetOk("additional_recipients"); ok {
		body.AdditionalRecipients = common.ExpandListOfString(v.([]interface{}))
	}

	if v, ok := d.GetOk("node_ips"); ok {
		body.NodeIps = expandNodeIps(v.([]interface{}))
	}

	if v, ok := d.GetOk("sda_ext_ids"); ok {
		body.SdaExtIds = common.ExpandListOfString(v.([]interface{}))
	}

	if v, ok := d.GetOk("should_anonymize"); ok {
		body.ShouldAnonymize = utils.BoolPtr(v.(bool))
	}

	if v, ok := d.GetOk("should_run_all_checks"); ok {
		body.ShouldRunAllChecks = utils.BoolPtr(v.(bool))
	}

	if v, ok := d.GetOk("should_send_report_to_configured_recipients"); ok {
		body.ShouldSendReportToConfiguredRecipients = utils.BoolPtr(v.(bool))
	}

	request := &monitoringRequest.RunSystemDefinedChecksRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		Body:         body,
	}

	resp, err := conn.SystemDefinedChecksAPI.RunSystemDefinedChecks(ctx, request)
	if err != nil {
		return diag.Errorf("error while running System-Defined Checks: %v", err)
	}

	if resp.Data == nil {
		return diag.Errorf("error: empty response data from RunSystemDefinedChecks API")
	}

	TaskRef := resp.Data.GetValue().(monitoringTaskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	aJSON, _ := json.MarshalIndent(TaskRef, "", "  ")
	log.Printf("[DEBUG] RunSystemDefinedChecks TaskReference: %s", string(aJSON))

	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for System-Defined Checks task (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching System-Defined Checks task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Run System-Defined Checks Task Details: %s", string(aJSON))

	d.SetId(utils.StringValue(taskDetails.ExtId))
	if err := d.Set("task_ext_id", utils.StringValue(taskDetails.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNutanixRunSystemDefinedChecksV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixRunSystemDefinedChecksV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixRunSystemDefinedChecksV2Create(ctx, d, meta)
}

func resourceNutanixRunSystemDefinedChecksV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandNodeIps(input []interface{}) []commonConfig.IPv4Address {
	if len(input) == 0 {
		return nil
	}
	nodeIps := make([]commonConfig.IPv4Address, 0, len(input))
	for _, v := range input {
		if v == nil {
			continue
		}
		item := v.(map[string]interface{})
		ipAddr := commonConfig.IPv4Address{}
		if val, ok := item["prefix_length"]; ok && val.(int) != 0 {
			ipAddr.PrefixLength = utils.IntPtr(val.(int))
		}
		if val, ok := item["value"]; ok && val.(string) != "" {
			ipAddr.Value = utils.StringPtr(val.(string))
		}
		nodeIps = append(nodeIps, ipAddr)
	}
	return nodeIps
}
