package monitoringv2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	monitoringModel "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of \"self\" identifies the URL for the object.",
				},
				"href": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The URL at which the entity described by the link can be accessed.",
				},
			},
		},
	}
}

func schemaForParamValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"int_value": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_int_value": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_int_value": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
				"float_value": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_float_value": {
								Type:        schema.TypeFloat,
								Computed:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_float_value": {
								Type:        schema.TypeFloat,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
				"bool_value": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_bool_value": {
								Type:        schema.TypeBool,
								Computed:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_bool_value": {
								Type:        schema.TypeBool,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
				"string_value": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_str_value": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_str_value": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaForParamValueResource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"int_value": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_int_value": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_int_value": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
				"float_value": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_float_value": {
								Type:        schema.TypeFloat,
								Optional:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_float_value": {
								Type:        schema.TypeFloat,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
				"bool_value": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_bool_value": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_bool_value": {
								Type:        schema.TypeBool,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
				"string_value": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_str_value": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Captures the current value of the parameter.",
							},
							"default_str_value": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Captures the default value of the parameter.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaForThresholdParameters() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Captures alert-related thresholds that correspond to a particular severity.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Equivalent name for the parameter used to display it on Prism UI.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unique identifier name for the parameter.",
				},
				"param_value": schemaForParamValue(),
				"unit": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unit for the parameter. For example, sec, %, MB, GB, and so on.",
				},
			},
		},
	}
}

func schemaForThresholdParametersResource() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Captures alert-related thresholds that correspond to a particular severity.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Equivalent name for the parameter used to display it on Prism UI.",
				},
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Unique identifier name for the parameter.",
				},
				"param_value": schemaForParamValueResource(),
				"unit": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unit for the parameter. For example, sec, %, MB, GB, and so on.",
				},
			},
		},
	}
}

func schemaForSeverityConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"state": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"threshold_parameters": schemaForThresholdParameters(),
			},
		},
	}
}

func schemaForSeverityConfigResource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"state": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"threshold_parameters": schemaForThresholdParametersResource(),
			},
		},
	}
}

func schemaForAlertConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_resolve": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"critical_severity":  schemaForSeverityConfig(),
				"info_severity":      schemaForSeverityConfig(),
				"warning_severity":   schemaForSeverityConfig(),
			},
		},
	}
}

func schemaForAlertConfigResource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_resolve": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"critical_severity":  schemaForSeverityConfigResource(),
				"info_severity":      schemaForSeverityConfigResource(),
				"warning_severity":   schemaForSeverityConfigResource(),
			},
		},
	}
}

func schemaForConfigurableParameters() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Parameters of the SDA that are configurable by a user.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Equivalent name for the parameter used to display it on Prism UI.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unique identifier name for the parameter.",
				},
				"param_value": schemaForParamValue(),
				"unit": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unit for the parameter. For example, sec, %, MB, GB, and so on.",
				},
			},
		},
	}
}

func schemaForConfigurableParametersResource() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Parameters of the SDA that are configurable by a user.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Equivalent name for the parameter used to display it on Prism UI.",
				},
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Unique identifier name for the parameter.",
				},
				"param_value": schemaForParamValueResource(),
				"unit": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unit for the parameter. For example, sec, %, MB, GB, and so on.",
				},
			},
		},
	}
}

// Flatten functions

func flattenLinks(links []import2.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(links))
	for i, link := range links {
		result[i] = map[string]interface{}{
			"rel":  utils.StringValue(link.Rel),
			"href": utils.StringValue(link.Href),
		}
	}
	return result
}

func flattenClusterConfigToState(d *schema.ResourceData, config monitoringModel.ClusterConfig) error {
	if config.ExtId != nil {
		if err := d.Set("ext_id", utils.StringValue(config.ExtId)); err != nil {
			return err
		}
	}
	if config.TenantId != nil {
		if err := d.Set("tenant_id", utils.StringValue(config.TenantId)); err != nil {
			return err
		}
	}
	if config.IsEnabled != nil {
		if err := d.Set("is_enabled", utils.BoolValue(config.IsEnabled)); err != nil {
			return err
		}
	}
	if config.LastModifiedByUser != nil {
		if err := d.Set("last_modified_by_user", utils.StringValue(config.LastModifiedByUser)); err != nil {
			return err
		}
	}
	if config.LastModifiedTime != nil {
		if err := d.Set("last_modified_time", config.LastModifiedTime.String()); err != nil {
			return err
		}
	}
	if config.ScheduleIntervalSeconds != nil {
		if err := d.Set("schedule_interval_seconds", *config.ScheduleIntervalSeconds); err != nil {
			return err
		}
	}
	if config.Links != nil {
		if err := d.Set("links", flattenLinks(config.Links)); err != nil {
			return err
		}
	}
	if config.AlertConfig != nil {
		if err := d.Set("alert_config", flattenAlertConfig(config.AlertConfig)); err != nil {
			return err
		}
	}
	if config.ConfigurableParameters != nil {
		if err := d.Set("configurable_parameters", flattenConfigurableParams(config.ConfigurableParameters)); err != nil {
			return err
		}
	}
	return nil
}

func flattenClusterConfigs(configs []monitoringModel.ClusterConfig) []map[string]interface{} {
	if len(configs) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, len(configs))
	for i, config := range configs {
		result[i] = flattenClusterConfig(config)
	}
	return result
}

func flattenClusterConfig(config monitoringModel.ClusterConfig) map[string]interface{} {
	m := map[string]interface{}{
		"ext_id":                     utils.StringValue(config.ExtId),
		"tenant_id":                  utils.StringValue(config.TenantId),
		"is_enabled":                 utils.BoolValue(config.IsEnabled),
		"last_modified_by_user":      utils.StringValue(config.LastModifiedByUser),
		"schedule_interval_seconds":  0,
		"links":                      flattenLinks(config.Links),
		"alert_config":               flattenAlertConfig(config.AlertConfig),
		"configurable_parameters":    flattenConfigurableParams(config.ConfigurableParameters),
		"system_defined_policy_ext_id": "",
	}
	if config.LastModifiedTime != nil {
		m["last_modified_time"] = config.LastModifiedTime.String()
	} else {
		m["last_modified_time"] = ""
	}
	if config.ScheduleIntervalSeconds != nil {
		m["schedule_interval_seconds"] = *config.ScheduleIntervalSeconds
	}
	return m
}

func flattenAlertConfig(ac *monitoringModel.AlertConfig) []interface{} {
	if ac == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"critical_severity": flattenSeverityConfig(ac.CriticalSeverity),
		"info_severity":     flattenSeverityConfig(ac.InfoSeverity),
		"warning_severity":  flattenSeverityConfig(ac.WarningSeverity),
	}
	if ac.AutoResolve != nil {
		m["auto_resolve"] = ac.AutoResolve.GetName()
	} else {
		m["auto_resolve"] = ""
	}
	return []interface{}{m}
}

func flattenSeverityConfig(sc *monitoringModel.SeverityConfig) []interface{} {
	if sc == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"threshold_parameters": flattenConfigurableParams(sc.ThresholdParameters),
	}
	if sc.State != nil {
		m["state"] = sc.State.GetName()
	} else {
		m["state"] = ""
	}
	return []interface{}{m}
}

func flattenConfigurableParams(params []monitoringModel.AlertPolicyConfigurableParameter) []interface{} {
	if len(params) == 0 {
		return []interface{}{}
	}
	result := make([]interface{}, len(params))
	for i, param := range params {
		m := map[string]interface{}{
			"display_name": utils.StringValue(param.DisplayName),
			"name":         utils.StringValue(param.Name),
			"unit":         utils.StringValue(param.Unit),
			"param_value":  flattenParamValue(param.ParamValue),
		}
		result[i] = m
	}
	return result
}

func flattenParamValue(pv *monitoringModel.OneOfAlertPolicyConfigurableParameterParamValue) []interface{} {
	if pv == nil {
		return []interface{}{}
	}
	value := pv.GetValue()
	if value == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"int_value":    []interface{}{},
		"float_value":  []interface{}{},
		"bool_value":   []interface{}{},
		"string_value": []interface{}{},
	}

	switch v := value.(type) {
	case monitoringModel.IntConfigurableParamValue:
		intMap := map[string]interface{}{}
		if v.CurrentIntValue != nil {
			intMap["current_int_value"] = int(*v.CurrentIntValue)
		}
		if v.DefaultIntValue != nil {
			intMap["default_int_value"] = int(*v.DefaultIntValue)
		}
		m["int_value"] = []interface{}{intMap}
	case monitoringModel.FloatConfigurableParamValue:
		floatMap := map[string]interface{}{}
		if v.CurrentFloatValue != nil {
			floatMap["current_float_value"] = float64(*v.CurrentFloatValue)
		}
		if v.DefaultFloatValue != nil {
			floatMap["default_float_value"] = float64(*v.DefaultFloatValue)
		}
		m["float_value"] = []interface{}{floatMap}
	case monitoringModel.BooleanConfigurableParamValue:
		boolMap := map[string]interface{}{}
		if v.CurrentBoolValue != nil {
			boolMap["current_bool_value"] = *v.CurrentBoolValue
		}
		if v.DefaultBoolValue != nil {
			boolMap["default_bool_value"] = *v.DefaultBoolValue
		}
		m["bool_value"] = []interface{}{boolMap}
	case monitoringModel.StringConfigurableParamValue:
		strMap := map[string]interface{}{}
		if v.CurrentStrValue != nil {
			strMap["current_str_value"] = *v.CurrentStrValue
		}
		if v.DefaultStrValue != nil {
			strMap["default_str_value"] = *v.DefaultStrValue
		}
		m["string_value"] = []interface{}{strMap}
	}

	return []interface{}{m}
}

// Expand functions

func expandAlertConfig(m map[string]interface{}) *monitoringModel.AlertConfig {
	ac := &monitoringModel.AlertConfig{}
	if v, ok := m["auto_resolve"]; ok && v.(string) != "" {
		var ars monitoringModel.AutoResolveState
		err := ars.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, v.(string))))
		if err == nil {
			ac.AutoResolve = ars.Ref()
		}
	}
	if v, ok := m["critical_severity"]; ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			ac.CriticalSeverity = expandSeverityConfig(list[0].(map[string]interface{}))
		}
	}
	if v, ok := m["info_severity"]; ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			ac.InfoSeverity = expandSeverityConfig(list[0].(map[string]interface{}))
		}
	}
	if v, ok := m["warning_severity"]; ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			ac.WarningSeverity = expandSeverityConfig(list[0].(map[string]interface{}))
		}
	}
	return ac
}

func expandSeverityConfig(m map[string]interface{}) *monitoringModel.SeverityConfig {
	sc := &monitoringModel.SeverityConfig{}
	if v, ok := m["state"]; ok && v.(string) != "" {
		var ps monitoringModel.PropertyState
		err := ps.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, v.(string))))
		if err == nil {
			sc.State = ps.Ref()
		}
	}
	if v, ok := m["threshold_parameters"]; ok {
		sc.ThresholdParameters = expandConfigurableParameters(v.([]interface{}))
	}
	return sc
}

func expandConfigurableParameters(params []interface{}) []monitoringModel.AlertPolicyConfigurableParameter {
	if len(params) == 0 {
		return nil
	}
	result := make([]monitoringModel.AlertPolicyConfigurableParameter, len(params))
	for i, p := range params {
		pm := p.(map[string]interface{})
		param := monitoringModel.AlertPolicyConfigurableParameter{}
		if v, ok := pm["name"]; ok && v.(string) != "" {
			param.Name = utils.StringPtr(v.(string))
		}
		if v, ok := pm["param_value"]; ok {
			paramValueList := v.([]interface{})
			if len(paramValueList) > 0 && paramValueList[0] != nil {
				param.ParamValue = expandParamValue(paramValueList[0].(map[string]interface{}))
			}
		}
		result[i] = param
	}
	return result
}

func expandParamValue(m map[string]interface{}) *monitoringModel.OneOfAlertPolicyConfigurableParameterParamValue {
	pv := monitoringModel.NewOneOfAlertPolicyConfigurableParameterParamValue()

	if v, ok := m["int_value"]; ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			intMap := list[0].(map[string]interface{})
			intVal := monitoringModel.IntConfigurableParamValue{}
			if cv, exists := intMap["current_int_value"]; exists {
				val := int64(cv.(int))
				intVal.CurrentIntValue = &val
			}
			pv.SetValue(intVal)
			return pv
		}
	}
	if v, ok := m["float_value"]; ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			floatMap := list[0].(map[string]interface{})
			floatVal := monitoringModel.FloatConfigurableParamValue{}
			if cv, exists := floatMap["current_float_value"]; exists {
				val := float32(cv.(float64))
				floatVal.CurrentFloatValue = &val
			}
			pv.SetValue(floatVal)
			return pv
		}
	}
	if v, ok := m["bool_value"]; ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			boolMap := list[0].(map[string]interface{})
			boolVal := monitoringModel.BooleanConfigurableParamValue{}
			if cv, exists := boolMap["current_bool_value"]; exists {
				val := cv.(bool)
				boolVal.CurrentBoolValue = &val
			}
			pv.SetValue(boolVal)
			return pv
		}
	}
	if v, ok := m["string_value"]; ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			strMap := list[0].(map[string]interface{})
			strVal := monitoringModel.StringConfigurableParamValue{}
			if cv, exists := strMap["current_str_value"]; exists {
				strVal.CurrentStrValue = utils.StringPtr(cv.(string))
			}
			pv.SetValue(strVal)
			return pv
		}
	}
	return nil
}
