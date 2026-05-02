package monitoringv2

import (
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	DEFAULTWAITTIMEOUT = 15
)

func autoResolveStateFromString(s string) serviceability.AutoResolveState {
	switch s {
	case "ENABLED":
		return serviceability.AUTORESOLVESTATE_ENABLED
	case "DISABLED":
		return serviceability.AUTORESOLVESTATE_DISABLED
	case "NOT_SUPPORTED":
		return serviceability.AUTORESOLVESTATE_NOT_SUPPORTED
	default:
		return serviceability.AUTORESOLVESTATE_UNKNOWN
	}
}

func propertyStateFromString(s string) serviceability.PropertyState {
	switch s {
	case "ENABLED":
		return serviceability.PROPERTYSTATE_ENABLED
	case "DISABLED":
		return serviceability.PROPERTYSTATE_DISABLED
	case "NOT_SUPPORTED":
		return serviceability.PROPERTYSTATE_NOT_SUPPORTED
	default:
		return serviceability.PROPERTYSTATE_UNKNOWN
	}
}

func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(links))
	for i, v := range links {
		link := map[string]interface{}{}
		if v.Href != nil {
			link["href"] = utils.StringValue(v.Href)
		}
		if v.Rel != nil {
			link["rel"] = utils.StringValue(v.Rel)
		}
		result[i] = link
	}
	return result
}

func flattenSeverityConfig(cfg *serviceability.SeverityConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}
	m := map[string]interface{}{}
	if cfg.State != nil {
		m["state"] = cfg.State.GetName()
	}
	if cfg.ThresholdParameters != nil {
		m["threshold_parameters"] = flattenAlertPolicyConfigurableParameters(cfg.ThresholdParameters)
	}
	return []map[string]interface{}{m}
}

func flattenAlertPolicyConfigurableParameters(params []serviceability.AlertPolicyConfigurableParameter) []map[string]interface{} {
	if len(params) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(params))
	for i, p := range params {
		m := map[string]interface{}{}
		if p.DisplayName != nil {
			m["display_name"] = utils.StringValue(p.DisplayName)
		}
		if p.Name != nil {
			m["name"] = utils.StringValue(p.Name)
		}
		if p.Unit != nil {
			m["unit"] = utils.StringValue(p.Unit)
		}
		if p.ParamValue != nil {
			m["param_value"] = flattenParamValue(p.ParamValue)
		}
		result[i] = m
	}
	return result
}

func flattenParamValue(pv *serviceability.OneOfAlertPolicyConfigurableParameterParamValue) []map[string]interface{} {
	if pv == nil || pv.ObjectType_ == nil {
		return nil
	}
	m := map[string]interface{}{}
	val := pv.GetValue()
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case serviceability.IntConfigurableParamValue:
		intMap := map[string]interface{}{}
		if v.CurrentIntValue != nil {
			intMap["current_int_value"] = *v.CurrentIntValue
		}
		if v.DefaultIntValue != nil {
			intMap["default_int_value"] = *v.DefaultIntValue
		}
		m["int_value"] = []map[string]interface{}{intMap}
	case serviceability.FloatConfigurableParamValue:
		floatMap := map[string]interface{}{}
		if v.CurrentFloatValue != nil {
			floatMap["current_float_value"] = float64(*v.CurrentFloatValue)
		}
		if v.DefaultFloatValue != nil {
			floatMap["default_float_value"] = float64(*v.DefaultFloatValue)
		}
		m["float_value"] = []map[string]interface{}{floatMap}
	case serviceability.BooleanConfigurableParamValue:
		boolMap := map[string]interface{}{}
		if v.CurrentBoolValue != nil {
			boolMap["current_bool_value"] = *v.CurrentBoolValue
		}
		if v.DefaultBoolValue != nil {
			boolMap["default_bool_value"] = *v.DefaultBoolValue
		}
		m["bool_value"] = []map[string]interface{}{boolMap}
	case serviceability.StringConfigurableParamValue:
		strMap := map[string]interface{}{}
		if v.CurrentStrValue != nil {
			strMap["current_str_value"] = utils.StringValue(v.CurrentStrValue)
		}
		if v.DefaultStrValue != nil {
			strMap["default_str_value"] = utils.StringValue(v.DefaultStrValue)
		}
		m["str_value"] = []map[string]interface{}{strMap}
	}
	return []map[string]interface{}{m}
}

func flattenAlertConfig(cfg *serviceability.AlertConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}
	m := map[string]interface{}{}
	if cfg.AutoResolve != nil {
		m["auto_resolve"] = cfg.AutoResolve.GetName()
	}
	if cfg.CriticalSeverity != nil {
		m["critical_severity"] = flattenSeverityConfig(cfg.CriticalSeverity)
	}
	if cfg.InfoSeverity != nil {
		m["info_severity"] = flattenSeverityConfig(cfg.InfoSeverity)
	}
	if cfg.WarningSeverity != nil {
		m["warning_severity"] = flattenSeverityConfig(cfg.WarningSeverity)
	}
	return []map[string]interface{}{m}
}

func flattenClusterConfigs(configs []serviceability.ClusterConfig) []map[string]interface{} {
	if len(configs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(configs))
	for i, c := range configs {
		result[i] = flattenClusterConfig(&c)
	}
	return result
}

func flattenClusterConfig(c *serviceability.ClusterConfig) map[string]interface{} {
	if c == nil {
		return nil
	}
	m := map[string]interface{}{}
	if c.AlertConfig != nil {
		m["alert_config"] = flattenAlertConfig(c.AlertConfig)
	}
	if c.ConfigurableParameters != nil {
		m["configurable_parameters"] = flattenAlertPolicyConfigurableParameters(c.ConfigurableParameters)
	}
	if c.ExtId != nil {
		m["ext_id"] = utils.StringValue(c.ExtId)
	}
	if c.IsEnabled != nil {
		m["is_enabled"] = *c.IsEnabled
	}
	if c.LastModifiedByUser != nil {
		m["last_modified_by_user"] = utils.StringValue(c.LastModifiedByUser)
	}
	if c.LastModifiedTime != nil {
		m["last_modified_time"] = c.LastModifiedTime.String()
	}
	if c.Links != nil {
		m["links"] = flattenLinks(c.Links)
	}
	if c.ScheduleIntervalSeconds != nil {
		m["schedule_interval_seconds"] = utils.IntValue(c.ScheduleIntervalSeconds)
	}
	if c.TenantId != nil {
		m["tenant_id"] = utils.StringValue(c.TenantId)
	}
	return m
}
