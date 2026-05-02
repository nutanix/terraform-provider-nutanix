package monitoringv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonCfg "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	monitoringPrism "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRunSystemDefinedChecksV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRunSystemDefinedChecksV2Create,
		ReadContext:   ResourceNutanixRunSystemDefinedChecksV2Read,
		DeleteContext: ResourceNutanixRunSystemDefinedChecksV2Delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
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
				ForceNew:    true,
				Description: "A list of additional email addresses for sending the run summary.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"node_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of node IP addresses where the Check will run.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The prefix length of the network to which this host IPv4 address belongs.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The IPv4 address of the host.",
						},
					},
				},
			},
			"sda_ext_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of Check IDs to be executed.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"should_anonymize": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether to mask sensitive data in the check run summary.",
			},
			"should_run_all_checks": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether all System-Defined Checks applicable to the specified cluster should be executed.",
			},
			"should_send_report_to_configured_recipients": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Determines if the run summary should be sent to the configured email address associated with the cluster.",
			},
			"task_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixRunSystemDefinedChecksV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	body := serviceability.NewRunSystemDefinedChecksSpec()

	if v, ok := d.GetOk("additional_recipients"); ok {
		recipients := make([]string, 0)
		for _, r := range v.([]interface{}) {
			recipients = append(recipients, r.(string))
		}
		body.AdditionalRecipients = recipients
	}

	if v, ok := d.GetOk("node_ips"); ok {
		nodeIps := expandIPv4Addresses(v.([]interface{}))
		body.NodeIps = nodeIps
	}

	if v, ok := d.GetOk("sda_ext_ids"); ok {
		sdaExtIds := make([]string, 0)
		for _, s := range v.([]interface{}) {
			sdaExtIds = append(sdaExtIds, s.(string))
		}
		body.SdaExtIds = sdaExtIds
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

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] RunSystemDefinedChecks payload: %s", aJSON)

	resp, err := conn.SystemDefinedChecksAPI.RunSystemDefinedChecks(utils.StringPtr(clusterExtID), body)
	if err != nil {
		return diag.Errorf("error while running System-Defined Checks: %v", err)
	}

	taskRefValue, ok := resp.Data.GetValue().(monitoringPrism.TaskReference)
	if !ok {
		return diag.Errorf("error: unexpected response type, expected TaskReference")
	}
	taskUUID := taskRefValue.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for System-Defined Checks task (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching System-Defined Checks task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] RunSystemDefinedChecks Task Details: %s", string(aJSON))

	d.SetId(utils.StringValue(taskUUID))
	if err := d.Set("task_ext_id", utils.StringValue(taskUUID)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixRunSystemDefinedChecksV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRunSystemDefinedChecksV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func expandIPv4Addresses(pr []interface{}) []commonCfg.IPv4Address {
	if len(pr) == 0 {
		return nil
	}
	result := make([]commonCfg.IPv4Address, 0, len(pr))
	for _, item := range pr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		ip := commonCfg.IPv4Address{}
		if v, ok := m["value"].(string); ok && v != "" {
			ip.Value = utils.StringPtr(v)
		}
		if v, ok := m["prefix_length"].(int); ok && v != 0 {
			ip.PrefixLength = utils.IntPtr(v)
		}
		result = append(result, ip)
	}
	return result
}
