package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSdaClusterConfigV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSdaClusterConfigV2Read,
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
			"alert_config":            schemaForAlertConfig(true),
			"configurable_parameters": schemaForThresholdParameters(true),
			"is_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the SDA policy is enabled or not on the cluster.",
			},
			"last_modified_by_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the user who made the latest update to this policy.",
			},
			"last_modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time in ISO 8601 format when the SDA policy was last modified.",
			},
			"links": linksSchema(),
			"schedule_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Interval in seconds for periodically executing the SDA policy.",
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixSdaClusterConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.SystemDefinedPoliciesAPI.GetClusterConfigById(utils.StringPtr(sdaPolicyExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching cluster config for SDA policy: %v", err)
	}

	getResp, ok := resp.Data.GetValue().(serviceability.ClusterConfig)
	if !ok {
		return diag.Errorf("error: unexpected response type from GetClusterConfigById, expected ClusterConfig")
	}

	if err := d.Set("alert_config", flattenAlertConfig(getResp.AlertConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("configurable_parameters", flattenAlertPolicyConfigurableParameters(getResp.ConfigurableParameters)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", getResp.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_modified_by_user", getResp.LastModifiedByUser); err != nil {
		return diag.FromErr(err)
	}
	if getResp.LastModifiedTime != nil {
		if err := d.Set("last_modified_time", getResp.LastModifiedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schedule_interval_seconds", getResp.ScheduleIntervalSeconds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
