package monitoringv2

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonResp "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	monCommon "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
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

func schemaForMetricDetails() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Details of the metric for a metric-based event.",
		Elem: &schema.Resource{
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
				"metric_value": schemaForOneOfValue(),
				"threshold_value": schemaForOneOfValue(),
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

func schemaForOneOfValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
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
				"double_value": {
					Type:        schema.TypeFloat,
					Computed:    true,
					Description: "Denotes a value of type double.",
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

// Flatten functions

func flattenLinks(links []commonResp.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return []map[string]interface{}{}
	}
	linkList := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		linkMap := make(map[string]interface{})
		if link.Href != nil {
			linkMap["href"] = utils.StringValue(link.Href)
		}
		if link.Rel != nil {
			linkMap["rel"] = utils.StringValue(link.Rel)
		}
		linkList = append(linkList, linkMap)
	}
	return linkList
}

func flattenTime(t *time.Time) string {
	if t != nil {
		return t.Format(time.RFC3339)
	}
	return ""
}

func flattenOperationType(opType *monCommon.OperationType) string {
	if opType != nil {
		return opType.GetName()
	}
	return ""
}

func flattenEventEntityReference(ref *serviceability.EventEntityReference) []map[string]interface{} {
	if ref == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"ext_id": utils.StringValue(ref.ExtId),
			"name":   utils.StringValue(ref.Name),
			"type":   utils.StringValue(ref.Type),
		},
	}
}

func flattenEntityReferences(refs []monCommon.EntityReference) []map[string]interface{} {
	if len(refs) == 0 {
		return []map[string]interface{}{}
	}
	refList := make([]map[string]interface{}, 0, len(refs))
	for _, ref := range refs {
		refMap := map[string]interface{}{
			"ext_id": utils.StringValue(ref.ExtId),
			"name":   utils.StringValue(ref.Name),
			"type":   utils.StringValue(ref.Type),
		}
		refList = append(refList, refMap)
	}
	return refList
}

func flattenMetricDetails(details []monCommon.MetricDetail) []map[string]interface{} {
	if len(details) == 0 {
		return []map[string]interface{}{}
	}
	detailList := make([]map[string]interface{}, 0, len(details))
	for _, detail := range details {
		detailMap := make(map[string]interface{})

		if detail.ComparisonOperator != nil {
			detailMap["comparison_operator"] = detail.ComparisonOperator.GetName()
		}
		if detail.ConditionType != nil {
			detailMap["condition_type"] = detail.ConditionType.GetName()
		}
		if detail.DataType != nil {
			detailMap["data_type"] = detail.DataType.GetName()
		}
		detailMap["metric_category"] = utils.StringValue(detail.MetricCategory)
		detailMap["metric_display_name"] = utils.StringValue(detail.MetricDisplayName)
		detailMap["metric_name"] = utils.StringValue(detail.MetricName)
		detailMap["metric_value"] = flattenOneOfMetricValue(detail.MetricValue)
		detailMap["threshold_value"] = flattenOneOfThresholdValue(detail.ThresholdValue)
		detailMap["trigger_time"] = flattenTime(detail.TriggerTime)
		detailMap["trigger_wait_time_seconds"] = utils.Int64Value(detail.TriggerWaitTimeSeconds)
		detailMap["unit"] = utils.StringValue(detail.Unit)

		detailList = append(detailList, detailMap)
	}
	return detailList
}

func flattenOneOfMetricValue(oneOfValue *monCommon.OneOfMetricDetailMetricValue) []map[string]interface{} {
	if oneOfValue != nil && oneOfValue.ObjectType_ != nil {
		valueMap := make(map[string]interface{})
		value := oneOfValue.GetValue()
		if value != nil {
			switch *oneOfValue.ObjectType_ {
			case "monitoring.v4.common.StringValue":
				if strVal, ok := value.(monCommon.StringValue); ok && strVal.StringValue != nil {
					valueMap["string_value"] = utils.StringValue(strVal.StringValue)
				}
			case "monitoring.v4.common.BoolValue":
				if boolVal, ok := value.(monCommon.BoolValue); ok && boolVal.BoolValue != nil {
					valueMap["bool_value"] = utils.BoolValue(boolVal.BoolValue)
				}
			case "monitoring.v4.common.IntValue":
				if intVal, ok := value.(monCommon.IntValue); ok && intVal.IntValue != nil {
					valueMap["int_value"] = utils.Int64Value(intVal.IntValue)
				}
			case "monitoring.v4.common.DoubleValue":
				if doubleVal, ok := value.(monCommon.DoubleValue); ok && doubleVal.DoubleValue != nil {
					valueMap["double_value"] = utils.Float64Value(doubleVal.DoubleValue)
				}
			}
		}
		return []map[string]interface{}{valueMap}
	}
	return []map[string]interface{}{}
}

func flattenOneOfThresholdValue(oneOfValue *monCommon.OneOfMetricDetailThresholdValue) []map[string]interface{} {
	if oneOfValue != nil && oneOfValue.ObjectType_ != nil {
		valueMap := make(map[string]interface{})
		value := oneOfValue.GetValue()
		if value != nil {
			switch *oneOfValue.ObjectType_ {
			case "monitoring.v4.common.StringValue":
				if strVal, ok := value.(monCommon.StringValue); ok && strVal.StringValue != nil {
					valueMap["string_value"] = utils.StringValue(strVal.StringValue)
				}
			case "monitoring.v4.common.BoolValue":
				if boolVal, ok := value.(monCommon.BoolValue); ok && boolVal.BoolValue != nil {
					valueMap["bool_value"] = utils.BoolValue(boolVal.BoolValue)
				}
			case "monitoring.v4.common.IntValue":
				if intVal, ok := value.(monCommon.IntValue); ok && intVal.IntValue != nil {
					valueMap["int_value"] = utils.Int64Value(intVal.IntValue)
				}
			case "monitoring.v4.common.DoubleValue":
				if doubleVal, ok := value.(monCommon.DoubleValue); ok && doubleVal.DoubleValue != nil {
					valueMap["double_value"] = utils.Float64Value(doubleVal.DoubleValue)
				}
			}
		}
		return []map[string]interface{}{valueMap}
	}
	return []map[string]interface{}{}
}

func flattenParameters(params []monCommon.Parameter) []map[string]interface{} {
	if len(params) == 0 {
		return []map[string]interface{}{}
	}
	paramList := make([]map[string]interface{}, 0, len(params))
	for _, param := range params {
		paramMap := make(map[string]interface{})
		paramMap["param_name"] = utils.StringValue(param.ParamName)
		paramMap["param_value"] = flattenOneOfParamValue(param.ParamValue)
		paramList = append(paramList, paramMap)
	}
	return paramList
}

func flattenOneOfParamValue(oneOfValue *monCommon.OneOfParameterParamValue) []map[string]interface{} {
	if oneOfValue != nil && oneOfValue.ObjectType_ != nil {
		valueMap := make(map[string]interface{})
		value := oneOfValue.GetValue()
		if value != nil {
			switch *oneOfValue.ObjectType_ {
			case "monitoring.v4.common.StringValue":
				if strVal, ok := value.(monCommon.StringValue); ok && strVal.StringValue != nil {
					valueMap["string_value"] = utils.StringValue(strVal.StringValue)
				}
			case "monitoring.v4.common.BoolValue":
				if boolVal, ok := value.(monCommon.BoolValue); ok && boolVal.BoolValue != nil {
					valueMap["bool_value"] = utils.BoolValue(boolVal.BoolValue)
				}
			case "monitoring.v4.common.IntValue":
				if intVal, ok := value.(monCommon.IntValue); ok && intVal.IntValue != nil {
					valueMap["int_value"] = utils.Int64Value(intVal.IntValue)
				}
			}
		}
		return []map[string]interface{}{valueMap}
	}
	return []map[string]interface{}{}
}
