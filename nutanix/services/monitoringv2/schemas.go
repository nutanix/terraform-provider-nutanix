package monitoringv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	monitoringCommon "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
	monitoringService "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(links))
	for i, link := range links {
		result[i] = map[string]interface{}{
			"href": utils.StringValue(link.Href),
			"rel":  utils.StringValue(link.Rel),
		}
	}
	return result
}

func schemaForEntityReference() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaForAlertEntityReference() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaForMetricDetail() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"comparison_operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"condition_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"string_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bool_value": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"int_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"double_value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
			"threshold_value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"string_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bool_value": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"int_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"double_value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
			"trigger_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trigger_wait_time_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaForParameter() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"param_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"param_value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"string_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bool_value": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"int_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func schemaForRootCauseAnalysis() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cause": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resolution": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaForSeverityTrail() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity_change_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func flattenEntityReferences(refs []monitoringCommon.EntityReference) []map[string]interface{} {
	if len(refs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(refs))
	for i, ref := range refs {
		result[i] = map[string]interface{}{
			"ext_id": utils.StringValue(ref.ExtId),
			"name":   utils.StringValue(ref.Name),
			"type":   utils.StringValue(ref.Type),
		}
	}
	return result
}

func flattenAlertEntityReference(ref *monitoringCommon.AlertEntityReference) []map[string]interface{} {
	if ref == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ext_id": utils.StringValue(ref.ExtId),
			"name":   utils.StringValue(ref.Name),
			"type":   utils.StringValue(ref.Type),
		},
	}
}

func flattenOneOfValue(oneOfValue interface{ GetValue() interface{} }, objectType *string) []map[string]interface{} {
	if oneOfValue == nil || objectType == nil {
		return nil
	}
	valueMap := make(map[string]interface{})
	value := oneOfValue.GetValue()
	if value == nil {
		return nil
	}
	switch *objectType {
	case "monitoring.v4.common.StringValue":
		if strVal, ok := value.(monitoringCommon.StringValue); ok && strVal.StringValue != nil {
			valueMap["string_value"] = utils.StringValue(strVal.StringValue)
		}
	case "monitoring.v4.common.BoolValue":
		if boolVal, ok := value.(monitoringCommon.BoolValue); ok && boolVal.BoolValue != nil {
			valueMap["bool_value"] = utils.BoolValue(boolVal.BoolValue)
		}
	case "monitoring.v4.common.IntValue":
		if intVal, ok := value.(monitoringCommon.IntValue); ok && intVal.IntValue != nil {
			valueMap["int_value"] = utils.Int64Value(intVal.IntValue)
		}
	case "monitoring.v4.common.DoubleValue":
		if doubleVal, ok := value.(monitoringCommon.DoubleValue); ok && doubleVal.DoubleValue != nil {
			valueMap["double_value"] = utils.Float64Value(doubleVal.DoubleValue)
		}
	default:
		return nil
	}
	return []map[string]interface{}{valueMap}
}

func flattenMetricDetails(details []monitoringCommon.MetricDetail) []map[string]interface{} {
	if len(details) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(details))
	for i, d := range details {
		m := map[string]interface{}{
			"metric_category":           utils.StringValue(d.MetricCategory),
			"metric_display_name":       utils.StringValue(d.MetricDisplayName),
			"metric_name":               utils.StringValue(d.MetricName),
			"trigger_time":              utils.TimeStringValue(d.TriggerTime),
			"trigger_wait_time_seconds": utils.Int64Value(d.TriggerWaitTimeSeconds),
			"unit":                      utils.StringValue(d.Unit),
		}

		if d.ComparisonOperator != nil {
			m["comparison_operator"] = flattenEnumValue(d.ComparisonOperator)
		}
		if d.ConditionType != nil {
			m["condition_type"] = flattenEnumValue(d.ConditionType)
		}
		if d.DataType != nil {
			m["data_type"] = flattenEnumValue(d.DataType)
		}

		if d.MetricValue != nil && d.MetricValue.ObjectType_ != nil {
			m["metric_value"] = flattenOneOfValue(d.MetricValue, d.MetricValue.ObjectType_)
		}
		if d.ThresholdValue != nil && d.ThresholdValue.ObjectType_ != nil {
			m["threshold_value"] = flattenOneOfValue(d.ThresholdValue, d.ThresholdValue.ObjectType_)
		}

		result[i] = m
	}
	return result
}

func flattenParameters(params []monitoringCommon.Parameter) []map[string]interface{} {
	if len(params) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(params))
	for i, p := range params {
		m := map[string]interface{}{
			"param_name": utils.StringValue(p.ParamName),
		}
		if p.ParamValue != nil && p.ParamValue.ObjectType_ != nil {
			m["param_value"] = flattenOneOfValue(p.ParamValue, p.ParamValue.ObjectType_)
		}
		result[i] = m
	}
	return result
}

func flattenRootCauseAnalysis(rcas []monitoringService.RootCauseAnalysis) []map[string]interface{} {
	if len(rcas) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(rcas))
	for i, rca := range rcas {
		result[i] = map[string]interface{}{
			"cause":      utils.StringValue(rca.Cause),
			"detail":     utils.StringValue(rca.Detail),
			"resolution": utils.StringValue(rca.Resolution),
		}
	}
	return result
}

func flattenSeverityTrails(trails []monitoringService.SeverityTrail) []map[string]interface{} {
	if len(trails) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(trails))
	for i, t := range trails {
		m := map[string]interface{}{
			"severity_change_time": utils.TimeStringValue(t.SeverityChangeTime),
		}
		if t.Severity != nil {
			m["severity"] = flattenEnumValue(t.Severity)
		}
		result[i] = m
	}
	return result
}

func flattenEnumValue(v interface{}) string {
	type stringer interface {
		GetName() string
	}
	if s, ok := v.(stringer); ok {
		return s.GetName()
	}
	return ""
}

func flattenImpactTypes(impactTypes []monitoringCommon.ImpactType) []string {
	if len(impactTypes) == 0 {
		return nil
	}
	result := make([]string, len(impactTypes))
	for i, it := range impactTypes {
		result[i] = flattenEnumValue(&it)
	}
	return result
}

func flattenAlert(alert monitoringService.Alert) map[string]interface{} {
	m := map[string]interface{}{
		"ext_id":                    utils.StringValue(alert.ExtId),
		"tenant_id":                 utils.StringValue(alert.TenantId),
		"acknowledged_by_username":  utils.StringValue(alert.AcknowledgedByUsername),
		"acknowledged_time":         utils.TimeStringValue(alert.AcknowledgedTime),
		"alert_type":                utils.StringValue(alert.AlertType),
		"classifications":           alert.Classifications,
		"cluster_name":              utils.StringValue(alert.ClusterName),
		"cluster_uuid":              utils.StringValue(alert.ClusterUUID),
		"creation_time":             utils.TimeStringValue(alert.CreationTime),
		"is_acknowledged":           utils.BoolValue(alert.IsAcknowledged),
		"is_auto_resolved":          utils.BoolValue(alert.IsAutoResolved),
		"is_resolved":               utils.BoolValue(alert.IsResolved),
		"is_runnable":               utils.BoolValue(alert.IsRunnable),
		"is_user_defined":           utils.BoolValue(alert.IsUserDefined),
		"kb_articles":               alert.KbArticles,
		"last_updated_time":         utils.TimeStringValue(alert.LastUpdatedTime),
		"message":                   utils.StringValue(alert.Message),
		"originating_cluster_uuid":  utils.StringValue(alert.OriginatingClusterUUID),
		"resolved_by_username":      utils.StringValue(alert.ResolvedByUsername),
		"resolved_time":             utils.TimeStringValue(alert.ResolvedTime),
		"service_name":              utils.StringValue(alert.ServiceName),
		"title":                     utils.StringValue(alert.Title),
		"affected_entities":         flattenEntityReferences(alert.AffectedEntities),
		"source_entity":             flattenAlertEntityReference(alert.SourceEntity),
		"impact_types":              flattenImpactTypes(alert.ImpactTypes),
		"metric_details":            flattenMetricDetails(alert.MetricDetails),
		"parameters":                flattenParameters(alert.Parameters),
		"root_cause_analysis":       flattenRootCauseAnalysis(alert.RootCauseAnalysis),
		"severity_trails":           flattenSeverityTrails(alert.SeverityTrails),
		"links":                     flattenLinks(alert.Links),
	}

	if alert.Severity != nil {
		m["severity"] = flattenEnumValue(alert.Severity)
	}

	return m
}

func schemaForAlertComputed() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"links":     schemaForLinks(),
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"acknowledged_by_username": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the user who acknowledged this alert.",
		},
		"acknowledged_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time in ISO 8601 format when the alert was acknowledged.",
		},
		"affected_entities": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of all the entities that are affected by the alert.",
			Elem:        schemaForEntityReference(),
		},
		"alert_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "A preconfigured or dynamically generated unique value for each alert type.",
		},
		"classifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Various categories into which this alert type can be classified.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cluster_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the cluster associated with the entity.",
		},
		"cluster_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster UUID associated with the source entity of the alert.",
		},
		"creation_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time in ISO 8601 format when the alert was created.",
		},
		"impact_types": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The impact this alert or event will have on the system.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"is_acknowledged": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates whether the alert is acknowledged or not.",
		},
		"is_auto_resolved": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates whether the alert is auto-resolved or not.",
		},
		"is_resolved": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates whether the alert is resolved or not.",
		},
		"is_runnable": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates whether the policy associated with the alert is runnable or not.",
		},
		"is_user_defined": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Flag to indicate if the alert was generated from a User-Defined Alert policy.",
		},
		"kb_articles": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of knowledge base article links.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"last_updated_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time in ISO 8601 format when the alert was last updated.",
		},
		"message": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Additional message associated with the alert.",
		},
		"metric_details": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Details of the metric for a metric-based event.",
			Elem:        schemaForMetricDetail(),
		},
		"originating_cluster_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster UUID associated with the cluster where the alert was first raised.",
		},
		"parameters": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Additional parameters associated with the alert.",
			Elem:        schemaForParameter(),
		},
		"resolved_by_username": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the user who resolved this alert.",
		},
		"resolved_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The time in ISO 8601 format when the alert was resolved.",
		},
		"root_cause_analysis": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Possible causes, resolutions and additional details to troubleshoot this alert.",
			Elem:        schemaForRootCauseAnalysis(),
		},
		"service_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The service that raised the alert.",
		},
		"severity": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"severity_trails": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Contains information on the severity change history for alerts.",
			Elem:        schemaForSeverityTrail(),
		},
		"source_entity": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     schemaForAlertEntityReference(),
		},
		"title": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Title of the alert.",
		},
	}
}
