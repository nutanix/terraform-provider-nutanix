package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitoringService "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixAlertsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixAlertsV2Read,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alerts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links":     schemaForLinks(),
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acknowledged_by_username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acknowledged_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"affected_entities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForEntityReference(),
						},
						"alert_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"classifications": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"impact_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"is_acknowledged": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_auto_resolved": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_resolved": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_runnable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_user_defined": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"kb_articles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"last_updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForMetricDetail(),
						},
						"originating_cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForParameter(),
						},
						"resolved_by_username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resolved_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_cause_analysis": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForRootCauseAnalysis(),
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity_trails": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForSeverityTrail(),
						},
						"source_entity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForAlertEntityReference(),
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixAlertsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	var filter, orderBy, selectParam *string
	if v, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		selectParam = utils.StringPtr(v.(string))
	}

	resp, err := conn.Alerts.ListAlerts(nil, nil, filter, orderBy, selectParam)
	if err != nil {
		return diag.Errorf("error while fetching alerts: %s", err)
	}

	alerts := resp.Data.GetValue().([]monitoringService.Alert)

	alertList := make([]map[string]interface{}, len(alerts))
	for i, alert := range alerts {
		alertList[i] = flattenAlert(alert)
	}

	if err := d.Set("alerts", alertList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}
