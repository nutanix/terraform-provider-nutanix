package monitoringv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	monitoringCommon "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The URL at which the entity described by the link can be accessed.",
				},
				"rel": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of \"self\" identifies the URL for the object.",
				},
			},
		},
	}
}

func schemaForEntityReference() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "UUID of the entity.",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the entity.",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The type of entity. For example, VM, node, or cluster.",
		},
	}
}

func schemaForOneOfValue() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"string_value": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Denotes a value of type string.",
		},
		"bool_value": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Denotes a value of type boolean.",
		},
		"int_value": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Denotes a value of type integer.",
		},
		"double_value": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "Denotes a value of type double.",
		},
	}
}

func schemaForMetricDetails() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Details of the metric for a metric-based event.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"comparison_operator": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The comparison operator used for the condition evaluation.",
				},
				"condition_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Indicating if this symptom is caused by a static threshold or anomaly (dynamic threshold) evaluation.",
				},
				"data_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Data type of the metric value as stored in the database.",
				},
				"metric_category": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Broad category under which this metric falls. For example, disk, CPU, or memory.",
				},
				"metric_display_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Readable name of the metric in English.",
				},
				"metric_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The metric key.",
				},
				"metric_value": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "The raw value of the metric when the condition threshold was exceeded.",
					Elem: &schema.Resource{
						Schema: schemaForOneOfValue(),
					},
				},
				"threshold_value": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "The threshold value that was used for the condition evaluation.",
					Elem: &schema.Resource{
						Schema: schemaForOneOfValue(),
					},
				},
				"trigger_time": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The time in ISO 8601 format when the event was triggered.",
				},
				"trigger_wait_time_seconds": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "How long the metric breached the given condition before raising an event.",
				},
				"unit": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unit of the metric. For example, percentage, ms or usecs.",
				},
			},
		},
	}
}

func schemaForParameters() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Additional parameters associated with the event.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"param_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Name or key of additional parameter for an instance.",
				},
				"param_value": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Value of additional parameter for an instance.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"string_value": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Denotes a value of type string.",
							},
							"bool_value": {
								Type:        schema.TypeBool,
								Computed:    true,
								Description: "Denotes a value of type boolean.",
							},
							"int_value": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "Denotes a value of type integer.",
							},
						},
					},
				},
			},
		},
	}
}

func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return []map[string]interface{}{}
	}

	linkList := make([]map[string]interface{}, 0)
	for _, link := range links {
		linkMap := make(map[string]interface{})
		if link.Rel != nil {
			linkMap["rel"] = utils.StringValue(link.Rel)
		}
		if link.Href != nil {
			linkMap["href"] = utils.StringValue(link.Href)
		}
		linkList = append(linkList, linkMap)
	}
	return linkList
}

func flattenEntityReferences(entities []monitoringCommon.EntityReference) []map[string]interface{} {
	if len(entities) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		entityMap := make(map[string]interface{})
		entityMap["ext_id"] = utils.StringValue(entity.ExtId)
		entityMap["name"] = utils.StringValue(entity.Name)
		entityMap["type"] = utils.StringValue(entity.Type)
		result[i] = entityMap
	}
	return result
}

func flattenEventEntityReference(entity *serviceability.EventEntityReference) []map[string]interface{} {
	if entity == nil {
		return []map[string]interface{}{}
	}

	entityMap := make(map[string]interface{})
	entityMap["ext_id"] = utils.StringValue(entity.ExtId)
	entityMap["name"] = utils.StringValue(entity.Name)
	entityMap["type"] = utils.StringValue(entity.Type)
	return []map[string]interface{}{entityMap}
}

func flattenMetricDetails(metricDetails []monitoringCommon.MetricDetail) []map[string]interface{} {
	if len(metricDetails) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, len(metricDetails))
	for i, md := range metricDetails {
		mdMap := make(map[string]interface{})

		if md.ComparisonOperator != nil {
			mdMap["comparison_operator"] = md.ComparisonOperator.GetName()
		}
		if md.ConditionType != nil {
			mdMap["condition_type"] = md.ConditionType.GetName()
		}
		if md.DataType != nil {
			mdMap["data_type"] = md.DataType.GetName()
		}
		mdMap["metric_category"] = utils.StringValue(md.MetricCategory)
		mdMap["metric_display_name"] = utils.StringValue(md.MetricDisplayName)
		mdMap["metric_name"] = utils.StringValue(md.MetricName)
		mdMap["metric_value"] = flattenOneOfMetricValue(md.MetricValue)
		mdMap["threshold_value"] = flattenOneOfThresholdValue(md.ThresholdValue)

		if md.TriggerTime != nil {
			mdMap["trigger_time"] = md.TriggerTime.String()
		}
		if md.TriggerWaitTimeSeconds != nil {
			mdMap["trigger_wait_time_seconds"] = utils.Int64Value(md.TriggerWaitTimeSeconds)
		}
		mdMap["unit"] = utils.StringValue(md.Unit)

		result[i] = mdMap
	}
	return result
}

func flattenOneOfMetricValue(oneOfValue *monitoringCommon.OneOfMetricDetailMetricValue) []map[string]interface{} {
	if oneOfValue == nil || oneOfValue.ObjectType_ == nil {
		return []map[string]interface{}{}
	}

	valueMap := make(map[string]interface{})
	value := oneOfValue.GetValue()
	if value == nil {
		return []map[string]interface{}{}
	}

	switch *oneOfValue.ObjectType_ {
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
	}

	return []map[string]interface{}{valueMap}
}

func flattenOneOfThresholdValue(oneOfValue *monitoringCommon.OneOfMetricDetailThresholdValue) []map[string]interface{} {
	if oneOfValue == nil || oneOfValue.ObjectType_ == nil {
		return []map[string]interface{}{}
	}

	valueMap := make(map[string]interface{})
	value := oneOfValue.GetValue()
	if value == nil {
		return []map[string]interface{}{}
	}

	switch *oneOfValue.ObjectType_ {
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
	}

	return []map[string]interface{}{valueMap}
}

func flattenParameters(params []monitoringCommon.Parameter) []map[string]interface{} {
	if len(params) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, len(params))
	for i, param := range params {
		paramMap := make(map[string]interface{})
		paramMap["param_name"] = utils.StringValue(param.ParamName)
		paramMap["param_value"] = flattenOneOfParamValue(param.ParamValue)
		result[i] = paramMap
	}
	return result
}

func flattenOneOfParamValue(oneOfValue *monitoringCommon.OneOfParameterParamValue) []map[string]interface{} {
	if oneOfValue == nil || oneOfValue.ObjectType_ == nil {
		return []map[string]interface{}{}
	}

	valueMap := make(map[string]interface{})
	value := oneOfValue.GetValue()
	if value == nil {
		return []map[string]interface{}{}
	}

	switch *oneOfValue.ObjectType_ {
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
	}

	return []map[string]interface{}{valueMap}
}
