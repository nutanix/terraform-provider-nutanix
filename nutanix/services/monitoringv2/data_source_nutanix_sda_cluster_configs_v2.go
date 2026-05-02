package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSdaClusterConfigsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSdaClusterConfigsV2Read,
		Schema: map[string]*schema.Schema{
			"system_defined_policy_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of the System-Defined Alert Policy.",
			},
			"cluster_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cluster-specific configurations for the SDA policy.",
				Elem: &schema.Resource{
					Schema: clusterConfigSchemaMap(),
				},
			},
		},
	}
}

func DatasourceNutanixSdaClusterConfigsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)

	resp, err := conn.SystemDefinedPoliciesAPI.ListClusterConfigsBySdaId(utils.StringPtr(sdaPolicyExtID), nil, nil, nil, nil, nil)
	if err != nil {
		return diag.Errorf("error while listing cluster configs for SDA policy: %v", err)
	}

	if resp.Data == nil {
		if setErr := d.Set("cluster_configs", []map[string]interface{}{}); setErr != nil {
			return diag.FromErr(setErr)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	listResp, ok := resp.Data.GetValue().([]serviceability.ClusterConfig)
	if !ok {
		return diag.Errorf("error: unexpected response type from ListClusterConfigsBySdaId, expected []ClusterConfig")
	}

	if err := d.Set("cluster_configs", flattenClusterConfigs(listResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}
