package cluster_managementv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	clusterConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	responseConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/response"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"
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
			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": schemaForIPAddress(true),
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(rsyslogModuleNameValues(), false),
						},
						"log_severity_level": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(rsyslogModuleLogSeverityLevelValues(), false),
						},
						"should_log_monitor_files": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
		},
	}
}

func resourceNutanixRsyslogServerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	body := clusterConfig.NewRsyslogServer()

	if v, ok := d.GetOk("server_name"); ok {
		body.ServerName = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("port"); ok {
		body.Port = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("network_protocol"); ok {
		body.NetworkProtocol = common.ExpandEnum[clusterConfig.RsyslogNetworkProtocol](v)
	}
	if v, ok := d.GetOk("ip_address"); ok {
		body.IpAddress = expandIPAddress(v.([]interface{}))
	}
	if v, ok := d.GetOk("modules"); ok {
		body.Modules = expandRsyslogModules(v.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create RsyslogServer Payload: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.CreateRsyslogServer(utils.StringPtr(clusterExtID), body)
	if err != nil {
		return diag.Errorf("error while creating Rsyslog Server: %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for rsyslog server (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	serverName := d.Get("server_name").(string)
	rsyslogExtID, findErr := findRsyslogServerExtIDByName(conn, clusterExtID, serverName)
	if findErr != nil {
		return diag.Errorf("rsyslog server was created but could not be found by name %q: %v", serverName, findErr)
	}

	d.SetId(fmt.Sprintf("%s:%s", clusterExtID, rsyslogExtID))
	return resourceNutanixRsyslogServerV2Read(ctx, d, meta)
}

func resourceNutanixRsyslogServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID, extID, err := parseRsyslogServerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.ClusterEntityAPI.GetRsyslogServerById(utils.StringPtr(clusterExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while reading Rsyslog Server: %v", err)
	}

	rsyslogServer := resp.Data.GetValue().(clusterConfig.RsyslogServer)
	aJSON, _ := json.MarshalIndent(rsyslogServer, "", "  ")
	log.Printf("[DEBUG] Read RsyslogServer Response: %s", string(aJSON))

	if err := d.Set("cluster_ext_id", clusterExtID); err != nil {
		return diag.FromErr(err)
	}
	return setRsyslogServerState(d, rsyslogServer)
}

func resourceNutanixRsyslogServerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID, extID, err := parseRsyslogServerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	readResp, err := conn.ClusterEntityAPI.GetRsyslogServerById(utils.StringPtr(clusterExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching Rsyslog Server for update: %v", err)
	}

	etagValue := conn.ClusterEntityAPI.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	updateSpec := readResp.Data.GetValue().(clusterConfig.RsyslogServer)

	if d.HasChange("port") {
		updateSpec.Port = utils.IntPtr(d.Get("port").(int))
	}
	if d.HasChange("network_protocol") {
		updateSpec.NetworkProtocol = common.ExpandEnum[clusterConfig.RsyslogNetworkProtocol](d.Get("network_protocol"))
	}
	if d.HasChange("ip_address") {
		updateSpec.IpAddress = expandIPAddress(d.Get("ip_address").([]interface{}))
	}
	if d.HasChange("modules") {
		updateSpec.Modules = expandRsyslogModules(d.Get("modules").([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Update RsyslogServer Payload: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.UpdateRsyslogServerById(
		utils.StringPtr(clusterExtID), utils.StringPtr(extID), &updateSpec, headers,
	)
	if err != nil {
		return diag.Errorf("error while updating Rsyslog Server: %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for rsyslog server (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return resourceNutanixRsyslogServerV2Read(ctx, d, meta)
}

func resourceNutanixRsyslogServerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID, extID, err := parseRsyslogServerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.ClusterEntityAPI.DeleteRsyslogServerById(utils.StringPtr(clusterExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while deleting Rsyslog Server: %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for rsyslog server (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func findRsyslogServerExtIDByName(conn *clusters.Client, clusterExtID, serverName string) (string, error) {
	resp, err := conn.ClusterEntityAPI.ListRsyslogServersByClusterId(utils.StringPtr(clusterExtID))
	if err != nil {
		return "", fmt.Errorf("error listing rsyslog servers: %v", err)
	}
	if resp.Data == nil {
		return "", fmt.Errorf("no rsyslog servers found in cluster %s", clusterExtID)
	}
	servers := resp.Data.GetValue().([]clusterConfig.RsyslogServer)
	for _, s := range servers {
		if s.ServerName != nil && utils.StringValue(s.ServerName) == serverName {
			return utils.StringValue(s.ExtId), nil
		}
	}
	return "", fmt.Errorf("rsyslog server %q not found in cluster %s", serverName, clusterExtID)
}

func parseRsyslogServerID(id string) (string, string, error) {
	parts := splitCompositeID(id)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid rsyslog server ID format: expected cluster_ext_id:ext_id, got %s", id)
	}
	return parts[0], parts[1], nil
}

func splitCompositeID(id string) []string {
	result := make([]string, 0)
	current := ""
	for _, c := range id {
		if c == ':' {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	result = append(result, current)
	return result
}

func setRsyslogServerState(d *schema.ResourceData, server clusterConfig.RsyslogServer) diag.Diagnostics {
	if server.ExtId != nil {
		if err := d.Set("ext_id", utils.StringValue(server.ExtId)); err != nil {
			return diag.FromErr(err)
		}
	}
	if server.ServerName != nil {
		if err := d.Set("server_name", utils.StringValue(server.ServerName)); err != nil {
			return diag.FromErr(err)
		}
	}
	if server.Port != nil {
		if err := d.Set("port", *server.Port); err != nil {
			return diag.FromErr(err)
		}
	}
	if server.NetworkProtocol != nil {
		if err := d.Set("network_protocol", server.NetworkProtocol.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}
	if server.IpAddress != nil {
		if err := d.Set("ip_address", flattenIPAddress(server.IpAddress)); err != nil {
			return diag.FromErr(err)
		}
	}
	if server.Modules != nil {
		if err := d.Set("modules", flattenRsyslogModules(server.Modules)); err != nil {
			return diag.FromErr(err)
		}
	}
	if server.TenantId != nil {
		if err := d.Set("tenant_id", utils.StringValue(server.TenantId)); err != nil {
			return diag.FromErr(err)
		}
	}
	links := flattenApiLinks(server.Links)
	if links == nil {
		links = []map[string]interface{}{}
	}
	if err := d.Set("links", links); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func expandIPAddress(pr []interface{}) *commonConfig.IPAddress {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ipAddress := commonConfig.NewIPAddress()

	if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
		ipAddress.Ipv4 = expandIPv4Address(ipv4.([]interface{}))
	}
	if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
		ipAddress.Ipv6 = expandIPv6Address(ipv6.([]interface{}))
	}
	return ipAddress
}

func expandIPv4Address(pr []interface{}) *commonConfig.IPv4Address {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ipv4 := commonConfig.NewIPv4Address()

	if v, ok := val["value"]; ok {
		ipv4.Value = utils.StringPtr(v.(string))
	}
	if p, ok := val["prefix_length"]; ok {
		ipv4.PrefixLength = utils.IntPtr(p.(int))
	}
	return ipv4
}

func expandIPv6Address(pr []interface{}) *commonConfig.IPv6Address {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ipv6 := commonConfig.NewIPv6Address()

	if v, ok := val["value"]; ok {
		ipv6.Value = utils.StringPtr(v.(string))
	}
	if p, ok := val["prefix_length"]; ok {
		ipv6.PrefixLength = utils.IntPtr(p.(int))
	}
	return ipv6
}

func expandRsyslogModules(pr []interface{}) []clusterConfig.RsyslogModuleItem {
	if len(pr) == 0 {
		return nil
	}
	modules := make([]clusterConfig.RsyslogModuleItem, 0, len(pr))
	for _, v := range pr {
		val := v.(map[string]interface{})
		module := clusterConfig.RsyslogModuleItem{}

		if name, ok := val["name"]; ok {
			module.Name = common.ExpandEnum[clusterConfig.RsyslogModuleName](name)
		}
		if severity, ok := val["log_severity_level"]; ok {
			module.LogSeverityLevel = common.ExpandEnum[clusterConfig.RsyslogModuleLogSeverityLevel](severity)
		}
		if shouldLog, ok := val["should_log_monitor_files"]; ok {
			module.ShouldLogMonitorFiles = utils.BoolPtr(shouldLog.(bool))
		}
		modules = append(modules, module)
	}
	return modules
}

func flattenIPAddress(addr *commonConfig.IPAddress) []map[string]interface{} {
	if addr == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ipv4": flattenIPv4Address(addr.Ipv4),
			"ipv6": flattenIPv6Address(addr.Ipv6),
		},
	}
}

func flattenIPv4Address(ipv4 *commonConfig.IPv4Address) []map[string]interface{} {
	if ipv4 == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"value":         utils.StringValue(ipv4.Value),
			"prefix_length": utils.IntValue(ipv4.PrefixLength),
		},
	}
}

func flattenIPv6Address(ipv6 *commonConfig.IPv6Address) []map[string]interface{} {
	if ipv6 == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"value":         utils.StringValue(ipv6.Value),
			"prefix_length": utils.IntValue(ipv6.PrefixLength),
		},
	}
}

func flattenRsyslogModules(modules []clusterConfig.RsyslogModuleItem) []map[string]interface{} {
	if len(modules) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(modules))
	for _, m := range modules {
		moduleMap := map[string]interface{}{
			"name":                    common.FlattenPtrEnum(m.Name),
			"log_severity_level":      common.FlattenPtrEnum(m.LogSeverityLevel),
			"should_log_monitor_files": utils.BoolValue(m.ShouldLogMonitorFiles),
		}
		result = append(result, moduleMap)
	}
	return result
}

func flattenApiLinks(links []responseConfig.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		linkMap := map[string]interface{}{
			"rel":  utils.StringValue(link.Rel),
			"href": utils.StringValue(link.Href),
		}
		result = append(result, linkMap)
	}
	return result
}

func rsyslogModuleNameValues() []string {
	return []string{
		"CASSANDRA", "CEREBRO", "CURATOR", "GENESIS", "PRISM",
		"STARGATE", "SYSLOG_MODULE", "ZOOKEEPER", "UHARA", "LAZAN",
		"API_AUDIT", "AUDIT", "CALM", "EPSILON", "ACROPOLIS",
		"MINERVA_CVM", "FLOW", "FLOW_SERVICE_LOGS", "LCM", "APLOS",
		"NCM_AIOPS",
	}
}

func rsyslogModuleLogSeverityLevelValues() []string {
	return []string{
		"EMERGENCY", "ALERT", "CRITICAL", "ERROR", "WARNING",
		"NOTICE", "INFO", "DEBUG",
	}
}
