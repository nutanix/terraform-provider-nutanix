package monitoringv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixSdaClusterConfigV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSdaClusterConfigV2Create,
		ReadContext:   ResourceNutanixSdaClusterConfigV2Read,
		UpdateContext: ResourceNutanixSdaClusterConfigV2Update,
		DeleteContext: ResourceNutanixSdaClusterConfigV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Update: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Delete: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"system_defined_policy_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique ID of the System-Defined Alert Policy.",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster UUID.",
			},
			"alert_config": schemaForAlertConfig(false),
			"configurable_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"param_value": schemaForParamValue(),
					},
				},
			},
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

func ResourceNutanixSdaClusterConfigV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)
	extID := d.Get("ext_id").(string)

	d.SetId(fmt.Sprintf("%s/%s", sdaPolicyExtID, extID))

	if d.HasChange("alert_config") {
		return ResourceNutanixSdaClusterConfigV2Update(ctx, d, meta)
	}

	return ResourceNutanixSdaClusterConfigV2Read(ctx, d, meta)
}

func ResourceNutanixSdaClusterConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return diag.Errorf("invalid resource ID format, expected <system_defined_policy_ext_id>/<ext_id>")
	}
	sdaPolicyExtID := idParts[0]
	extID := idParts[1]

	if err := d.Set("system_defined_policy_ext_id", sdaPolicyExtID); err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.SystemDefinedPoliciesAPI.GetClusterConfigById(utils.StringPtr(sdaPolicyExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching cluster config: %v", err)
	}

	getResp, ok := resp.Data.GetValue().(serviceability.ClusterConfig)
	if !ok {
		return diag.Errorf("error: unexpected response type, expected ClusterConfig")
	}

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("alert_config", flattenAlertConfig(getResp.AlertConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("configurable_parameters", flattenAlertPolicyConfigurableParameters(getResp.ConfigurableParameters)); err != nil {
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

	return nil
}

func ResourceNutanixSdaClusterConfigV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return diag.Errorf("invalid resource ID format, expected <system_defined_policy_ext_id>/<ext_id>")
	}
	sdaPolicyExtID := idParts[0]
	extID := idParts[1]

	getResp, err := conn.SystemDefinedPoliciesAPI.GetClusterConfigById(utils.StringPtr(sdaPolicyExtID), utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching cluster config for update: %v", err)
	}

	etagValue := conn.SystemDefinedPoliciesAPI.ApiClient.GetEtag(getResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updateSpec, ok := getResp.Data.GetValue().(serviceability.ClusterConfig)
	if !ok {
		return diag.Errorf("error: unexpected response type, expected ClusterConfig")
	}

	if d.HasChange("alert_config") {
		if v, ok := d.GetOk("alert_config"); ok {
			alertConfig := expandAlertConfig(v.([]interface{}))
			updateSpec.AlertConfig = alertConfig
		}
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] UpdateClusterConfigById payload: %s", aJSON)

	_, err = conn.SystemDefinedPoliciesAPI.UpdateClusterConfigById(utils.StringPtr(sdaPolicyExtID), utils.StringPtr(extID), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating cluster config: %v", err)
	}

	return ResourceNutanixSdaClusterConfigV2Read(ctx, d, meta)
}

func ResourceNutanixSdaClusterConfigV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func expandAlertConfig(pr []interface{}) *serviceability.AlertConfig {
	if len(pr) == 0 {
		return nil
	}
	m, ok := pr[0].(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := serviceability.NewAlertConfig()

	if v, ok := m["auto_resolve"].(string); ok && v != "" {
		autoResolve := autoResolveStateFromString(v)
		cfg.AutoResolve = &autoResolve
	}
	if v, ok := m["critical_severity"].([]interface{}); ok && len(v) > 0 {
		cfg.CriticalSeverity = expandSeverityConfig(v)
	}
	if v, ok := m["info_severity"].([]interface{}); ok && len(v) > 0 {
		cfg.InfoSeverity = expandSeverityConfig(v)
	}
	if v, ok := m["warning_severity"].([]interface{}); ok && len(v) > 0 {
		cfg.WarningSeverity = expandSeverityConfig(v)
	}
	return cfg
}

func expandSeverityConfig(pr []interface{}) *serviceability.SeverityConfig {
	if len(pr) == 0 {
		return nil
	}
	m, ok := pr[0].(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := serviceability.NewSeverityConfig()

	if v, ok := m["state"].(string); ok && v != "" {
		state := propertyStateFromString(v)
		cfg.State = &state
	}
	if v, ok := m["threshold_parameters"].([]interface{}); ok && len(v) > 0 {
		cfg.ThresholdParameters = expandThresholdParameters(v)
	}
	return cfg
}

func expandThresholdParameters(pr []interface{}) []serviceability.AlertPolicyConfigurableParameter {
	if len(pr) == 0 {
		return nil
	}
	result := make([]serviceability.AlertPolicyConfigurableParameter, 0, len(pr))
	for _, item := range pr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		param := serviceability.AlertPolicyConfigurableParameter{}
		if v, ok := m["name"].(string); ok && v != "" {
			param.Name = utils.StringPtr(v)
		}
		if v, ok := m["display_name"].(string); ok && v != "" {
			param.DisplayName = utils.StringPtr(v)
		}
		if v, ok := m["unit"].(string); ok && v != "" {
			param.Unit = utils.StringPtr(v)
		}
		if v, ok := m["param_value"].([]interface{}); ok && len(v) > 0 {
			pv := expandParamValue(v)
			if pv != nil {
				param.ParamValue = pv
			}
		}
		result = append(result, param)
	}
	return result
}

func expandParamValue(pr []interface{}) *serviceability.OneOfAlertPolicyConfigurableParameterParamValue {
	if len(pr) == 0 {
		return nil
	}
	m, ok := pr[0].(map[string]interface{})
	if !ok {
		return nil
	}

	pv := serviceability.NewOneOfAlertPolicyConfigurableParameterParamValue()

	if v, ok := m["int_value"].([]interface{}); ok && len(v) > 0 {
		intM, ok := v[0].(map[string]interface{})
		if ok {
			intVal := serviceability.IntConfigurableParamValue{}
			if cv, ok := intM["current_int_value"].(int); ok {
				val := int64(cv)
				intVal.CurrentIntValue = &val
			}
			if err := pv.SetValue(intVal); err != nil {
				log.Printf("[WARN] failed to set int param value: %v", err)
			}
		}
	} else if v, ok := m["float_value"].([]interface{}); ok && len(v) > 0 {
		floatM, ok := v[0].(map[string]interface{})
		if ok {
			floatVal := serviceability.FloatConfigurableParamValue{}
			if cv, ok := floatM["current_float_value"].(float64); ok {
				fv := float32(cv)
				floatVal.CurrentFloatValue = &fv
			}
			if err := pv.SetValue(floatVal); err != nil {
				log.Printf("[WARN] failed to set float param value: %v", err)
			}
		}
	} else if v, ok := m["bool_value"].([]interface{}); ok && len(v) > 0 {
		boolM, ok := v[0].(map[string]interface{})
		if ok {
			boolVal := serviceability.BooleanConfigurableParamValue{}
			if cv, ok := boolM["current_bool_value"].(bool); ok {
				boolVal.CurrentBoolValue = &cv
			}
			if err := pv.SetValue(boolVal); err != nil {
				log.Printf("[WARN] failed to set bool param value: %v", err)
			}
		}
	} else if v, ok := m["str_value"].([]interface{}); ok && len(v) > 0 {
		strM, ok := v[0].(map[string]interface{})
		if ok {
			strVal := serviceability.StringConfigurableParamValue{}
			if cv, ok := strM["current_str_value"].(string); ok {
				strVal.CurrentStrValue = utils.StringPtr(cv)
			}
			if err := pv.SetValue(strVal); err != nil {
				log.Printf("[WARN] failed to set str param value: %v", err)
			}
		}
	}

	return pv
}
