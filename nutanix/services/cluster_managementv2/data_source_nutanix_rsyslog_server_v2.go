package cluster_managementv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
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
	}
}

func datasourceNutanixRsyslogServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.ClusterEntityAPI.GetRsyslogServerById(utils.StringPtr(clusterExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while reading Rsyslog Server: %v", err)
	}
	if resp == nil || resp.Data == nil {
		return diag.Errorf("no Rsyslog Server found with ext_id: %s in cluster: %s", extID, clusterExtID)
	}

	rsyslogServer := resp.Data.GetValue().(clusterConfig.RsyslogServer)
	aJSON, _ := json.MarshalIndent(rsyslogServer, "", "  ")
	log.Printf("[DEBUG] Get RsyslogServer Response: %s", string(aJSON))

	d.SetId(utils.StringValue(rsyslogServer.ExtId))
	return setRsyslogServerState(d, rsyslogServer)
}
