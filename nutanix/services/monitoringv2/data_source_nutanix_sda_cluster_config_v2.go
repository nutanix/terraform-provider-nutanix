package monitoringv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitoringModel "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSdaClusterConfigV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixSdaClusterConfigV2Read,
		Schema: map[string]*schema.Schema{
			"system_defined_policy_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of the System-Defined Alert Policy.",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
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
				Computed:    true,
				Description: "Interval in seconds for periodically executing the SDA policy. This will not be set for policies with the type NOT_SCHEDULED & EVENT_DRIVEN.",
			},
			"alert_config": schemaForAlertConfig(),
			"configurable_parameters": schemaForConfigurableParameters(),
		},
	}
}

func datasourceNutanixSdaClusterConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.SystemDefinedPoliciesAPI.GetClusterConfigById(
		utils.StringPtr(sdaPolicyExtID),
		utils.StringPtr(extID),
	)
	if err != nil {
		return diag.Errorf("error while reading SDA cluster config: %v", err)
	}
	if resp == nil || resp.Data == nil {
		return diag.Errorf("no SDA cluster config found for system_defined_policy_ext_id: %s, ext_id: %s", sdaPolicyExtID, extID)
	}

	body := resp.Data.GetValue().(monitoringModel.ClusterConfig)
	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Get SDA Cluster Config Response: %s", string(aJSON))

	if err := flattenClusterConfigToState(d, body); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(extID)
	return nil
}
