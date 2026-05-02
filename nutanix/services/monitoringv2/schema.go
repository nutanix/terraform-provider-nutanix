package monitoringv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaForParamValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"int_value": {
					Type:     schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_int_value": {
								Type:     schema.TypeInt,
								Computed: true,
								Optional: true,
							},
							"default_int_value": {
								Type:     schema.TypeInt,
								Computed: true,
							},
						},
					},
				},
				"float_value": {
					Type:     schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_float_value": {
								Type:     schema.TypeFloat,
								Computed: true,
								Optional: true,
							},
							"default_float_value": {
								Type:     schema.TypeFloat,
								Computed: true,
							},
						},
					},
				},
				"bool_value": {
					Type:     schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_bool_value": {
								Type:     schema.TypeBool,
								Computed: true,
								Optional: true,
							},
							"default_bool_value": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
				"str_value": {
					Type:     schema.TypeList,
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"current_str_value": {
								Type:     schema.TypeString,
								Computed: true,
								Optional: true,
							},
							"default_str_value": {
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

func schemaForThresholdParameters(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: !computed,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: !computed,
				},
				"unit": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"param_value": schemaForParamValue(),
			},
		},
	}
}

func schemaForSeverityConfig(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: !computed,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"state": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: !computed,
				},
				"threshold_parameters": schemaForThresholdParameters(computed),
			},
		},
	}
}

func schemaForAlertConfig(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Optional: !computed,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_resolve": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: !computed,
				},
				"critical_severity": schemaForSeverityConfig(computed),
				"info_severity":     schemaForSeverityConfig(computed),
				"warning_severity":  schemaForSeverityConfig(computed),
			},
		},
	}
}

func schemaForClusterConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: clusterConfigSchemaMap(),
		},
	}
}

func clusterConfigSchemaMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alert_config":              schemaForAlertConfig(true),
		"configurable_parameters":   schemaForThresholdParameters(true),
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"last_modified_by_user": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_modified_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": {
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
		},
		"schedule_interval_seconds": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func linksSchema() *schema.Schema {
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
