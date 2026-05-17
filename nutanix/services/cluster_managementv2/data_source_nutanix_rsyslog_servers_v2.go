package cluster_managementv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
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
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": schemaForIPAddressComputed(),
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
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": schemaForLinks(),
					},
				},
			},
		},
	}
}

func datasourceNutanixRsyslogServersV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	resp, err := conn.ClusterEntityAPI.ListRsyslogServersByClusterId(utils.StringPtr(clusterExtID))
	if err != nil {
		return diag.Errorf("error while listing Rsyslog Servers: %v", err)
	}
	if resp.Data == nil {
		if err := d.Set("rsyslog_servers", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No Rsyslog Servers found.",
			Detail:   "The API returned an empty list of rsyslog servers.",
		}}
	}

	servers := resp.Data.GetValue().([]clusterConfig.RsyslogServer)
	if err := d.Set("rsyslog_servers", flattenRsyslogServers(servers)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenRsyslogServers(servers []clusterConfig.RsyslogServer) []map[string]interface{} {
	if len(servers) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(servers))
	for _, server := range servers {
		result = append(result, flattenRsyslogServer(server))
	}
	return result
}

func flattenRsyslogServer(server clusterConfig.RsyslogServer) map[string]interface{} {
	return map[string]interface{}{
		"ext_id":           utils.StringValue(server.ExtId),
		"server_name":      utils.StringValue(server.ServerName),
		"ip_address":       flattenIPAddress(server.IpAddress),
		"port":             utils.IntValue(server.Port),
		"network_protocol": common.FlattenPtrEnum(server.NetworkProtocol),
		"modules":          flattenRsyslogModules(server.Modules),
		"tenant_id":        utils.StringValue(server.TenantId),
		"links":            flattenApiLinks(server.Links),
	}
}
