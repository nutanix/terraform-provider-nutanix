package cluster_managementv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clustermgmtConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clustermgmtRequest "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/request/clusters"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRsyslogServersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixRsyslogServersV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rsyslog_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": schemaForLinks(),
						"ip_address": schemaForIPAddress(),
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_severity_level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"should_log_monitor_files": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func datasourceNutanixRsyslogServersV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterMgmtAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	request := &clustermgmtRequest.ListRsyslogServersByClusterIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
	}

	resp, err := conn.ClustersServiceAPI.ListRsyslogServersByClusterId(ctx, request)
	if err != nil {
		return diag.Errorf("error while listing Rsyslog Servers: %v", err)
	}
	if resp == nil || resp.Data == nil {
		if err := d.Set("rsyslog_servers", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	servers := resp.Data.GetValue().([]clustermgmtConfig.RsyslogServer)
	aJSON, _ := json.MarshalIndent(servers, "", "  ")
	log.Printf("[DEBUG] List Rsyslog Servers Response: %s", string(aJSON))

	if err := d.Set("rsyslog_servers", flattenRsyslogServers(servers)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenRsyslogServers(servers []clustermgmtConfig.RsyslogServer) []map[string]interface{} {
	if len(servers) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(servers))
	for _, server := range servers {
		result = append(result, flattenRsyslogServer(server))
	}
	return result
}

func flattenRsyslogServer(server clustermgmtConfig.RsyslogServer) map[string]interface{} {
	m := map[string]interface{}{
		"ext_id":      utils.StringValue(server.ExtId),
		"tenant_id":   utils.StringValue(server.TenantId),
		"links":       flattenLinks(server.Links),
		"ip_address":  flattenIPAddress(server.IpAddress),
		"server_name": utils.StringValue(server.ServerName),
		"port":        utils.IntValue(server.Port),
		"modules":     flattenRsyslogModules(server.Modules),
	}
	if server.NetworkProtocol != nil {
		m["network_protocol"] = server.NetworkProtocol.GetName()
	} else {
		m["network_protocol"] = ""
	}
	return m
}
