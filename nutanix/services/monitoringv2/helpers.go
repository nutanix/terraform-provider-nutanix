package monitoringv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
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
					Description: "A name that identifies the relationship of the link to the object that is returned by the URL.",
				},
			},
		},
	}
}

func schemaForTriggerConditionsComputed() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Trigger conditions for the policy.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"condition": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"metric_name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The metric key.",
							},
							"operator": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Comparison operator.",
							},
							"threshold_value": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
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
							},
						},
					},
				},
				"condition_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"severity_level": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func schemaForTriggerConditionsInput() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "Trigger conditions for the policy. If there are multiple trigger conditions, all of them will be considered during the operation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"condition": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"metric_name": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The metric key.",
							},
							"operator": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Comparison operator.",
							},
							"threshold_value": {
								Type:     schema.TypeList,
								Required: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"int_value": {
											Type:        schema.TypeInt,
											Optional:    true,
											Description: "Denotes a value of type integer.",
										},
										"double_value": {
											Type:        schema.TypeFloat,
											Optional:    true,
											Description: "Denotes a value of type double.",
										},
									},
								},
							},
						},
					},
				},
				"condition_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Condition type.",
				},
				"severity_level": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Severity level.",
				},
			},
		},
	}
}

func schemaForFiltersComputed() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Filter criteria for narrowing down the entities on which User-Defined Alert policies can be set up.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"entity_filter": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Entity UUID on which the User-Defined Alert policy should be set up.",
							},
						},
					},
				},
				"group_filter": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Entity UUID of the group entity type on which the User-Defined Alert policy should be set up.",
							},
							"type": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func schemaForFiltersInput() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: "Filter criteria for narrowing down the entities on which User-Defined Alert policies can be set up.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"entity_filter": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Entity UUID on which the User-Defined Alert policy should be set up.",
							},
						},
					},
				},
				"group_filter": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ext_id": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Entity UUID of the group entity type on which the User-Defined Alert policy should be set up.",
							},
							"type": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

func schemaForRelatedPoliciesComputed() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of alert policies that are related to the entities of the current policy.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"entity_uuid": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "UUID of the entity the User-Defined Alert policy is associated with.",
				},
				"policy_ids": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Policy IDs associated with the specified entity.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

// Flatten functions

func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if links == nil {
		return nil
	}
	result := make([]map[string]interface{}, len(links))
	for i, link := range links {
		result[i] = map[string]interface{}{
			"href": link.Href,
			"rel":  link.Rel,
		}
	}
	return result
}

func flattenTriggerConditions(conditions []serviceability.TriggerCondition) []map[string]interface{} {
	if conditions == nil {
		return nil
	}
	result := make([]map[string]interface{}, len(conditions))
	for i, tc := range conditions {
		condMap := map[string]interface{}{}

		if tc.Condition != nil {
			conditionList := []map[string]interface{}{flattenCondition(tc.Condition)}
			condMap["condition"] = conditionList
		} else {
			condMap["condition"] = []map[string]interface{}{}
		}

		if tc.ConditionType != nil {
			condMap["condition_type"] = tc.ConditionType.GetName()
		} else {
			condMap["condition_type"] = ""
		}

		if tc.SeverityLevel != nil {
			condMap["severity_level"] = tc.SeverityLevel.GetName()
		} else {
			condMap["severity_level"] = ""
		}

		result[i] = condMap
	}
	return result
}

func flattenCondition(cond *serviceability.Condition) map[string]interface{} {
	condMap := map[string]interface{}{
		"metric_name": "",
		"operator":    "",
	}

	if cond.MetricName != nil {
		condMap["metric_name"] = utils.StringValue(cond.MetricName)
	}
	if cond.Operator != nil {
		condMap["operator"] = cond.Operator.GetName()
	}

	condMap["threshold_value"] = flattenThresholdValue(cond.ThresholdValue)
	return condMap
}

func flattenThresholdValue(tv *serviceability.OneOfConditionThresholdValue) []map[string]interface{} {
	if tv == nil || tv.ObjectType_ == nil {
		return []map[string]interface{}{}
	}

	valueMap := map[string]interface{}{
		"int_value":    0,
		"double_value": 0.0,
	}

	val := tv.GetValue()
	if val != nil {
		switch *tv.ObjectType_ {
		case "monitoring.v4.common.IntValue":
			if intVal, ok := val.(common.IntValue); ok && intVal.IntValue != nil {
				valueMap["int_value"] = int(utils.Int64Value(intVal.IntValue))
			}
		case "monitoring.v4.common.DoubleValue":
			if doubleVal, ok := val.(common.DoubleValue); ok && doubleVal.DoubleValue != nil {
				valueMap["double_value"] = utils.Float64Value(doubleVal.DoubleValue)
			}
		}
	}

	return []map[string]interface{}{valueMap}
}

func flattenFilters(filters *serviceability.OneOfUserDefinedPolicyFilters) []map[string]interface{} {
	if filters == nil || filters.ObjectType_ == nil {
		return []map[string]interface{}{}
	}

	filterMap := map[string]interface{}{
		"entity_filter": []map[string]interface{}{},
		"group_filter":  []map[string]interface{}{},
	}

	val := filters.GetValue()
	if val != nil {
		switch v := val.(type) {
		case []serviceability.EntityFilter:
			entityFilters := make([]map[string]interface{}, len(v))
			for i, ef := range v {
				entityFilters[i] = map[string]interface{}{
					"ext_id": utils.StringValue(ef.ExtId),
				}
			}
			filterMap["entity_filter"] = entityFilters
		case []serviceability.GroupFilter:
			groupFilters := make([]map[string]interface{}, len(v))
			for i, gf := range v {
				gfMap := map[string]interface{}{
					"ext_id": utils.StringValue(gf.ExtId),
					"type":   "",
				}
				if gf.Type != nil {
					gfMap["type"] = gf.Type.GetName()
				}
				groupFilters[i] = gfMap
			}
			filterMap["group_filter"] = groupFilters
		}
	}

	return []map[string]interface{}{filterMap}
}

func flattenImpactTypes(impactTypes []common.ImpactType) []string {
	if impactTypes == nil {
		return nil
	}
	result := make([]string, len(impactTypes))
	for i, it := range impactTypes {
		result[i] = it.GetName()
	}
	return result
}

func flattenRelatedPolicies(policies []serviceability.RelatedPolicy) []map[string]interface{} {
	if policies == nil {
		return nil
	}
	result := make([]map[string]interface{}, len(policies))
	for i, rp := range policies {
		result[i] = map[string]interface{}{
			"entity_uuid": utils.StringValue(rp.EntityUuid),
			"policy_ids":  rp.PolicyIds,
		}
	}
	return result
}

func flattenUdaPolicies(policies []serviceability.UserDefinedPolicy) []map[string]interface{} {
	if policies == nil {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		pMap := map[string]interface{}{
			"ext_id":                           utils.StringValue(policy.ExtId),
			"tenant_id":                        policy.TenantId,
			"links":                            flattenLinks(policy.Links),
			"title":                            policy.Title,
			"description":                      policy.Description,
			"entity_type":                      policy.EntityType,
			"trigger_conditions":               flattenTriggerConditions(policy.TriggerConditions),
			"filters":                          flattenFilters(policy.Filters),
			"impact_types":                     flattenImpactTypes(policy.ImpactTypes),
			"is_auto_resolved":                 policy.IsAutoResolved,
			"is_enabled":                       policy.IsEnabled,
			"trigger_wait_period":              policy.TriggerWaitPeriod,
			"created_by":                       policy.CreatedBy,
			"policies_to_override":             policy.PoliciesToOverride,
			"related_policies":                 flattenRelatedPolicies(policy.RelatedPolicies),
			"is_expected_to_error_on_conflict": policy.IsExpectedToErrorOnConflict,
		}
		if policy.LastUpdatedTime != nil {
			pMap["last_updated_time"] = policy.LastUpdatedTime.String()
		} else {
			pMap["last_updated_time"] = ""
		}
		result[i] = pMap
	}
	return result
}

// Expand functions

func expandTriggerConditions(tcList []interface{}) []serviceability.TriggerCondition {
	if len(tcList) == 0 {
		return nil
	}
	result := make([]serviceability.TriggerCondition, len(tcList))
	for i, v := range tcList {
		tcMap := v.(map[string]interface{})
		tc := *serviceability.NewTriggerCondition()

		if condList, ok := tcMap["condition"].([]interface{}); ok && len(condList) > 0 {
			tc.Condition = expandCondition(condList[0].(map[string]interface{}))
		}

		if ct, ok := tcMap["condition_type"].(string); ok && ct != "" {
			ctVal := conditionTypeFromString(ct)
			tc.ConditionType = ctVal.Ref()
		}

		if sl, ok := tcMap["severity_level"].(string); ok && sl != "" {
			slVal := policySeverityLevelFromString(sl)
			tc.SeverityLevel = slVal.Ref()
		}

		result[i] = tc
	}
	return result
}

func expandCondition(condMap map[string]interface{}) *serviceability.Condition {
	cond := serviceability.NewCondition()

	if mn, ok := condMap["metric_name"].(string); ok && mn != "" {
		cond.MetricName = utils.StringPtr(mn)
	}

	if op, ok := condMap["operator"].(string); ok && op != "" {
		opVal := comparisonOperatorFromString(op)
		cond.Operator = opVal.Ref()
	}

	if tvList, ok := condMap["threshold_value"].([]interface{}); ok && len(tvList) > 0 {
		cond.ThresholdValue = expandThresholdValue(tvList[0].(map[string]interface{}))
	}

	return cond
}

func expandThresholdValue(tvMap map[string]interface{}) *serviceability.OneOfConditionThresholdValue {
	tv := serviceability.NewOneOfConditionThresholdValue()

	if intVal, ok := tvMap["int_value"].(int); ok && intVal != 0 {
		intV := common.NewIntValue()
		intV.IntValue = utils.Int64Ptr(int64(intVal))
		tv.SetValue(*intV)
	} else if doubleVal, ok := tvMap["double_value"].(float64); ok && doubleVal != 0 {
		doubleV := common.NewDoubleValue()
		doubleV.DoubleValue = utils.Float64Ptr(doubleVal)
		tv.SetValue(*doubleV)
	}

	return tv
}

func expandFilters(filterList []interface{}) *serviceability.OneOfUserDefinedPolicyFilters {
	if len(filterList) == 0 {
		return nil
	}

	filterMap := filterList[0].(map[string]interface{})
	filters := serviceability.NewOneOfUserDefinedPolicyFilters()

	if entityFilters, ok := filterMap["entity_filter"].([]interface{}); ok && len(entityFilters) > 0 {
		efs := make([]serviceability.EntityFilter, len(entityFilters))
		for i, ef := range entityFilters {
			efMap := ef.(map[string]interface{})
			efs[i] = serviceability.EntityFilter{
				ExtId: utils.StringPtr(efMap["ext_id"].(string)),
			}
		}
		filters.SetValue(efs)
	} else if groupFilters, ok := filterMap["group_filter"].([]interface{}); ok && len(groupFilters) > 0 {
		gfs := make([]serviceability.GroupFilter, len(groupFilters))
		for i, gf := range groupFilters {
			gfMap := gf.(map[string]interface{})
			gfs[i] = serviceability.GroupFilter{
				ExtId: utils.StringPtr(gfMap["ext_id"].(string)),
			}
			if t, tok := gfMap["type"].(string); tok && t != "" {
				getVal := groupEntityTypeFromString(t)
				gfs[i].Type = getVal.Ref()
			}
		}
		filters.SetValue(gfs)
	}

	return filters
}

func expandImpactTypes(impactList []interface{}) []common.ImpactType {
	if len(impactList) == 0 {
		return nil
	}
	result := make([]common.ImpactType, len(impactList))
	for i, v := range impactList {
		result[i] = impactTypeFromString(v.(string))
	}
	return result
}

// Enum conversion helpers

func conditionTypeFromString(s string) common.ConditionType {
	switch s {
	case "STATIC_THRESHOLD":
		return common.CONDITIONTYPE_STATIC_THRESHOLD
	default:
		return common.CONDITIONTYPE_UNKNOWN
	}
}

func policySeverityLevelFromString(s string) serviceability.PolicySeverityLevel {
	switch s {
	case "WARNING":
		return serviceability.POLICYSEVERITYLEVEL_WARNING
	case "CRITICAL":
		return serviceability.POLICYSEVERITYLEVEL_CRITICAL
	default:
		return serviceability.POLICYSEVERITYLEVEL_UNKNOWN
	}
}

func comparisonOperatorFromString(s string) common.ComparisonOperator {
	switch s {
	case "EQUAL_TO":
		return common.COMPARISONOPERATOR_EQUAL_TO
	case "GREATER_THAN":
		return common.COMPARISONOPERATOR_GREATER_THAN
	case "GREATER_THAN_OR_EQUAL_TO":
		return common.COMPARISONOPERATOR_GREATER_THAN_OR_EQUAL_TO
	case "LESS_THAN":
		return common.COMPARISONOPERATOR_LESS_THAN
	case "LESS_THAN_OR_EQUAL_TO":
		return common.COMPARISONOPERATOR_LESS_THAN_OR_EQUAL_TO
	default:
		return common.COMPARISONOPERATOR_UNKNOWN
	}
}

func groupEntityTypeFromString(s string) serviceability.GroupEntityType {
	switch s {
	case "CATEGORY":
		return serviceability.GROUPENTITYTYPE_CATEGORY
	case "CLUSTER":
		return serviceability.GROUPENTITYTYPE_CLUSTER
	default:
		return serviceability.GROUPENTITYTYPE_UNKNOWN
	}
}

func impactTypeFromString(s string) common.ImpactType {
	switch s {
	case "AVAILABILITY":
		return common.IMPACTTYPE_AVAILABILITY
	case "CAPACITY":
		return common.IMPACTTYPE_CAPACITY
	case "CONFIGURATION":
		return common.IMPACTTYPE_CONFIGURATION
	case "PERFORMANCE":
		return common.IMPACTTYPE_PERFORMANCE
	case "SYSTEM_INDICATOR":
		return common.IMPACTTYPE_SYSTEM_INDICATOR
	case "CPU_CAPACITY":
		return common.IMPACTTYPE_CPU_CAPACITY
	default:
		return common.IMPACTTYPE_UNKNOWN
	}
}
