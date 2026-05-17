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

func DatasourceNutanixRsyslogServerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixRsyslogServerV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func datasourceNutanixRsyslogServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterMgmtAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Get("ext_id").(string)

	request := &clustermgmtRequest.GetRsyslogServerByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}

	resp, err := conn.ClustersServiceAPI.GetRsyslogServerById(ctx, request)
	if err != nil {
		return diag.Errorf("error while reading Rsyslog Server: %v", err)
	}
	if resp == nil || resp.Data == nil {
		return diag.Errorf("no Rsyslog Server found with ext_id: %s", extID)
	}

	server := resp.Data.GetValue().(clustermgmtConfig.RsyslogServer)
	aJSON, _ := json.MarshalIndent(server, "", "  ")
	log.Printf("[DEBUG] Get Rsyslog Server Response: %s", string(aJSON))

	d.SetId(utils.StringValue(server.ExtId))

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
