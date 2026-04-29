package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitoringService "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAlertV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixAlertV2Read,
		Schema:     schemaForAlertComputed(),
	}
}

func DatasourceNutanixAlertV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.Alerts.GetAlertById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching alert: %s", err)
	}

	alert := resp.Data.GetValue().(monitoringService.Alert)

	if err := d.Set("tenant_id", utils.StringValue(alert.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(alert.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("acknowledged_by_username", utils.StringValue(alert.AcknowledgedByUsername)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("acknowledged_time", utils.TimeStringValue(alert.AcknowledgedTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("affected_entities", flattenEntityReferences(alert.AffectedEntities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("alert_type", utils.StringValue(alert.AlertType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("classifications", alert.Classifications); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_name", utils.StringValue(alert.ClusterName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_uuid", utils.StringValue(alert.ClusterUUID)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_time", utils.TimeStringValue(alert.CreationTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("impact_types", flattenImpactTypes(alert.ImpactTypes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_acknowledged", utils.BoolValue(alert.IsAcknowledged)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_auto_resolved", utils.BoolValue(alert.IsAutoResolved)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_resolved", utils.BoolValue(alert.IsResolved)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_runnable", utils.BoolValue(alert.IsRunnable)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_user_defined", utils.BoolValue(alert.IsUserDefined)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kb_articles", alert.KbArticles); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated_time", utils.TimeStringValue(alert.LastUpdatedTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("message", utils.StringValue(alert.Message)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metric_details", flattenMetricDetails(alert.MetricDetails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("originating_cluster_uuid", utils.StringValue(alert.OriginatingClusterUUID)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parameters", flattenParameters(alert.Parameters)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resolved_by_username", utils.StringValue(alert.ResolvedByUsername)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resolved_time", utils.TimeStringValue(alert.ResolvedTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("root_cause_analysis", flattenRootCauseAnalysis(alert.RootCauseAnalysis)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("service_name", utils.StringValue(alert.ServiceName)); err != nil {
		return diag.FromErr(err)
	}
	if alert.Severity != nil {
		if err := d.Set("severity", flattenEnumValue(alert.Severity)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("severity_trails", flattenSeverityTrails(alert.SeverityTrails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_entity", flattenAlertEntityReference(alert.SourceEntity)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("title", utils.StringValue(alert.Title)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(alert.ExtId))
	return nil
}
