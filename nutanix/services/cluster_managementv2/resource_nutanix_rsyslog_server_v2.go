package cluster_managementv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	clustermgmtConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clustermgmtRequest "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/request/clusters"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRsyslogServerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixRsyslogServerV2Create,
		ReadContext:   resourceNutanixRsyslogServerV2Read,
		UpdateContext: resourceNutanixRsyslogServerV2Update,
		DeleteContext: resourceNutanixRsyslogServerV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"ip_address": schemaForIPAddressInput(),
			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"UDP", "TCP", "RELP"}, false),
			},
			"modules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CASSANDRA", "CEREBRO", "CURATOR", "GENESIS", "PRISM",
								"STARGATE", "SYSLOG_MODULE", "ZOOKEEPER", "UHARA", "LAZAN",
								"API_AUDIT", "AUDIT", "CALM", "EPSILON", "ACROPOLIS",
								"MINERVA_CVM", "FLOW", "FLOW_SERVICE_LOGS", "LCM", "APLOS", "NCM_AIOPS",
							}, false),
						},
						"log_severity_level": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"EMERGENCY", "ALERT", "CRITICAL", "ERROR",
								"WARNING", "NOTICE", "INFO", "DEBUG",
							}, false),
						},
						"should_log_monitor_files": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixRsyslogServerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterMgmtAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	body := &clustermgmtConfig.RsyslogServer{}

	if v, ok := d.GetOk("server_name"); ok {
		body.ServerName = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("port"); ok {
		body.Port = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("network_protocol"); ok {
		body.NetworkProtocol = rsyslogNetworkProtocolFromString(v.(string))
	}
	if v, ok := d.GetOk("ip_address"); ok {
		body.IpAddress = expandIPAddress(v.([]interface{}))
	}
	if v, ok := d.GetOk("modules"); ok {
		body.Modules = expandRsyslogModules(v.([]interface{}))
	}

	request := &clustermgmtRequest.CreateRsyslogServerRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		Body:         body,
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Rsyslog Server Request Body: %s", string(aJSON))

	resp, err := conn.ClustersServiceAPI.CreateRsyslogServer(ctx, request)
	if err != nil {
		return diag.Errorf("error while creating Rsyslog Server: %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := taskRef.ExtId
	log.Printf("[DEBUG] Create Rsyslog Server Task UUID: %s", utils.StringValue(taskUUID))

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: 5 * time.Minute,
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for rsyslog server to create: %s", errWaitTask)
	}

	// After task completes, list and find the created server by name
	listRequest := &clustermgmtRequest.ListRsyslogServersByClusterIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
	}
	listResp, err := conn.ClustersServiceAPI.ListRsyslogServersByClusterId(ctx, listRequest)
	if err != nil {
		return diag.Errorf("error while listing Rsyslog Servers after creation: %v", err)
	}

	if listResp != nil && listResp.Data != nil {
		servers := listResp.Data.GetValue().([]clustermgmtConfig.RsyslogServer)
		for _, server := range servers {
			if utils.StringValue(server.ServerName) == d.Get("server_name").(string) {
				d.SetId(utils.StringValue(server.ExtId))
				return resourceNutanixRsyslogServerV2Read(ctx, d, meta)
			}
		}
	}

	return diag.Errorf("rsyslog server was created but could not be found in list response")
}

func resourceNutanixRsyslogServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterMgmtAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	request := &clustermgmtRequest.GetRsyslogServerByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}

	resp, err := conn.ClustersServiceAPI.GetRsyslogServerById(ctx, request)
	if err != nil {
		return diag.Errorf("error while reading Rsyslog Server: %v", err)
	}
	if resp == nil || resp.Data == nil {
		d.SetId("")
		return nil
	}

	server := resp.Data.GetValue().(clustermgmtConfig.RsyslogServer)
	aJSON, _ := json.MarshalIndent(server, "", "  ")
	log.Printf("[DEBUG] Read Rsyslog Server Response: %s", string(aJSON))

	if err := d.Set("ext_id", utils.StringValue(server.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(server.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(server.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_address", flattenIPAddress(server.IpAddress)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server_name", utils.StringValue(server.ServerName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", utils.IntValue(server.Port)); err != nil {
		return diag.FromErr(err)
	}
	if server.NetworkProtocol != nil {
		if err := d.Set("network_protocol", server.NetworkProtocol.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("modules", flattenRsyslogModules(server.Modules)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixRsyslogServerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterMgmtAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	// Read first to get the ETag
	getRequest := &clustermgmtRequest.GetRsyslogServerByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	getResp, err := conn.ClustersServiceAPI.GetRsyslogServerById(ctx, getRequest)
	if err != nil {
		return diag.Errorf("error reading Rsyslog Server for update: %v", err)
	}
	serverData := getResp.Data.GetValue().(clustermgmtConfig.RsyslogServer)
	etagValue := conn.ClustersServiceAPI.ApiClient.GetEtag(&serverData)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := &clustermgmtConfig.RsyslogServer{}

	if v, ok := d.GetOk("server_name"); ok {
		body.ServerName = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("port"); ok {
		body.Port = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("network_protocol"); ok {
		body.NetworkProtocol = rsyslogNetworkProtocolFromString(v.(string))
	}
	if v, ok := d.GetOk("ip_address"); ok {
		body.IpAddress = expandIPAddress(v.([]interface{}))
	}
	if v, ok := d.GetOk("modules"); ok {
		body.Modules = expandRsyslogModules(v.([]interface{}))
	}

	request := &clustermgmtRequest.UpdateRsyslogServerByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
		Body:         body,
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Update Rsyslog Server Request Body: %s", string(aJSON))

	updateResp, err := conn.ClustersServiceAPI.UpdateRsyslogServerById(ctx, request, args)
	if err != nil {
		return diag.Errorf("error while updating Rsyslog Server: %v", err)
	}

	taskRef := updateResp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := taskRef.ExtId
	log.Printf("[DEBUG] Update Rsyslog Server Task UUID: %s", utils.StringValue(taskUUID))

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: 5 * time.Minute,
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for rsyslog server to update: %s", errWaitTask)
	}

	return resourceNutanixRsyslogServerV2Read(ctx, d, meta)
}

func resourceNutanixRsyslogServerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterMgmtAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	request := &clustermgmtRequest.DeleteRsyslogServerByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}

	deleteResp, err := conn.ClustersServiceAPI.DeleteRsyslogServerById(ctx, request)
	if err != nil {
		return diag.Errorf("error while deleting Rsyslog Server: %v", err)
	}

	taskRef := deleteResp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := taskRef.ExtId
	log.Printf("[DEBUG] Delete Rsyslog Server Task UUID: %s", utils.StringValue(taskUUID))

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: 5 * time.Minute,
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for rsyslog server to delete: %s", errWaitTask)
	}

	d.SetId("")
	return nil
}
