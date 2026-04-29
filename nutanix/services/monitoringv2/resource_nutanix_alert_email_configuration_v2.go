package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitoringCommon "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/common"
	monitoringService "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixAlertEmailConfigurationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixAlertEmailConfigurationV2Create,
		ReadContext:   ResourceNutanixAlertEmailConfigurationV2Read,
		UpdateContext: ResourceNutanixAlertEmailConfigurationV2Update,
		DeleteContext: ResourceNutanixAlertEmailConfigurationV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"links": schemaForLinks(),
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"alert_email_digest_send_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Time in HH:mm format when the alert email digest is sent daily.",
			},
			"alert_email_digest_send_timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timezone for the time at which the alert email digest is sent daily.",
			},
			"default_nutanix_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The default Nutanix email ID to which alert emails are sent.",
			},
			"email_config_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Rules for email configuration.",
				Elem:        schemaForEmailConfigRule(),
			},
			"email_contact_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of email contacts.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"email_template": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_suffix": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Suffix for the email body.",
						},
						"subject_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Prefix for the email subject.",
						},
					},
				},
			},
			"has_default_nutanix_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether alert emails are enabled or not on default Nutanix email ID.",
			},
			"is_email_digest_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether alert email digest is enabled or not.",
			},
			"is_empty_alert_email_digest_skipped": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Send alert email digest only if there are one or more alerts.",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether alert emails are enabled or not.",
			},
			"tunnel_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForRemoteTunnelDetails(),
			},
		},
	}
}

func schemaForEmailConfigRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_uuids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Cluster UUIDs to which this rule applies.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"has_global_email_contact_list": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether to include a global email contact list.",
			},
			"impact_types": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the configuration rule is enabled or not.",
			},
			"match_phrases": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of phrases to match the alert.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"recipients": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of recipients who will receive emails.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"severities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func schemaForRemoteTunnelDetails() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"connection_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForCommunicationStatus(),
			},
			"http_proxy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForHTTPProxy(),
			},
			"service_center": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForServiceCenter(),
			},
			"transport_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForCommunicationStatus(),
			},
		},
	}
}

func schemaForCommunicationStatus() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"last_changed_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last changed time.",
			},
			"last_checked_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last checked time.",
			},
			"last_successful_transmission_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last successful transmission time.",
			},
			"message": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message.",
						},
						"attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaForHTTPProxy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Proxy name.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port on which proxy is binding.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User name for proxy authentication.",
			},
			"proxy_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func schemaForServiceCenter() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the service center.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of service center.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port number.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username.",
			},
		},
	}
}

func ResourceNutanixAlertEmailConfigurationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	getResp, err := conn.AlertEmailConfiguration.GetAlertEmailConfiguration()
	if err != nil {
		return diag.Errorf("error while fetching alert email configuration for ETag: %s", err)
	}

	etagValue := conn.AlertEmailConfiguration.ApiClient.GetEtag(getResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := expandAlertEmailConfiguration(d)

	_, err = conn.AlertEmailConfiguration.UpdateAlertEmailConfiguration(body, args)
	if err != nil {
		return diag.Errorf("error while creating alert email configuration: %s", err)
	}

	return ResourceNutanixAlertEmailConfigurationV2Read(ctx, d, meta)
}

func ResourceNutanixAlertEmailConfigurationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	resp, err := conn.AlertEmailConfiguration.GetAlertEmailConfiguration()
	if err != nil {
		return diag.Errorf("error while fetching alert email configuration: %s", err)
	}

	config := resp.Data.GetValue().(monitoringService.AlertEmailConfiguration)

	if err := d.Set("ext_id", utils.StringValue(config.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(config.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(config.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("alert_email_digest_send_time", utils.StringValue(config.AlertEmailDigestSendTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("alert_email_digest_send_timezone", utils.StringValue(config.AlertEmailDigestSendTimezone)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_nutanix_email", utils.StringValue(config.DefaultNutanixEmail)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_config_rules", flattenEmailConfigRules(config.EmailConfigRules)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_contact_list", config.EmailContactList); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_template", flattenEmailTemplate(config.EmailTemplate)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_default_nutanix_email", utils.BoolValue(config.HasDefaultNutanixEmail)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_email_digest_enabled", utils.BoolValue(config.IsEmailDigestEnabled)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_empty_alert_email_digest_skipped", utils.BoolValue(config.IsEmptyAlertEmailDigestSkipped)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", utils.BoolValue(config.IsEnabled)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tunnel_details", flattenRemoteTunnelDetails(config.TunnelDetails)); err != nil {
		return diag.FromErr(err)
	}

	if config.ExtId != nil {
		d.SetId(utils.StringValue(config.ExtId))
	} else {
		d.SetId("alert-email-configuration")
	}

	return nil
}

func ResourceNutanixAlertEmailConfigurationV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	getResp, err := conn.AlertEmailConfiguration.GetAlertEmailConfiguration()
	if err != nil {
		return diag.Errorf("error while fetching alert email configuration for ETag: %s", err)
	}

	etagValue := conn.AlertEmailConfiguration.ApiClient.GetEtag(getResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	body := expandAlertEmailConfiguration(d)

	_, err = conn.AlertEmailConfiguration.UpdateAlertEmailConfiguration(body, args)
	if err != nil {
		return diag.Errorf("error while updating alert email configuration: %s", err)
	}

	return ResourceNutanixAlertEmailConfigurationV2Read(ctx, d, meta)
}

func ResourceNutanixAlertEmailConfigurationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func expandAlertEmailConfiguration(d *schema.ResourceData) *monitoringService.AlertEmailConfiguration {
	body := &monitoringService.AlertEmailConfiguration{}

	if v, ok := d.GetOk("alert_email_digest_send_time"); ok {
		body.AlertEmailDigestSendTime = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("alert_email_digest_send_timezone"); ok {
		body.AlertEmailDigestSendTimezone = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("default_nutanix_email"); ok {
		body.DefaultNutanixEmail = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("email_config_rules"); ok {
		body.EmailConfigRules = expandEmailConfigRules(v.([]interface{}))
	}
	if v, ok := d.GetOk("email_contact_list"); ok {
		body.EmailContactList = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("email_template"); ok {
		templateList := v.([]interface{})
		if len(templateList) > 0 && templateList[0] != nil {
			body.EmailTemplate = expandEmailTemplate(templateList[0].(map[string]interface{}))
		}
	}
	if v, ok := d.GetOk("has_default_nutanix_email"); ok {
		body.HasDefaultNutanixEmail = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("is_email_digest_enabled"); ok {
		body.IsEmailDigestEnabled = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("is_empty_alert_email_digest_skipped"); ok {
		body.IsEmptyAlertEmailDigestSkipped = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("is_enabled"); ok {
		body.IsEnabled = utils.BoolPtr(v.(bool))
	}

	return body
}

func expandEmailConfigRules(rules []interface{}) []monitoringService.EmailConfigurationRule {
	if len(rules) == 0 {
		return nil
	}
	result := make([]monitoringService.EmailConfigurationRule, len(rules))
	for i, r := range rules {
		ruleMap := r.(map[string]interface{})
		rule := monitoringService.EmailConfigurationRule{}

		if v, ok := ruleMap["cluster_uuids"]; ok {
			rule.ClusterUuids = expandStringList(v.([]interface{}))
		}
		if v, ok := ruleMap["has_global_email_contact_list"]; ok {
			rule.HasGlobalEmailContactList = utils.BoolPtr(v.(bool))
		}
		if v, ok := ruleMap["impact_types"]; ok {
			rule.ImpactTypes = expandImpactTypes(v.([]interface{}))
		}
		if v, ok := ruleMap["is_enabled"]; ok {
			rule.IsEnabled = utils.BoolPtr(v.(bool))
		}
		if v, ok := ruleMap["match_phrases"]; ok {
			rule.MatchPhrases = expandStringList(v.([]interface{}))
		}
		if v, ok := ruleMap["recipients"]; ok {
			rule.Recipients = expandStringList(v.([]interface{}))
		}
		if v, ok := ruleMap["severities"]; ok {
			rule.Severities = expandSeverities(v.([]interface{}))
		}

		result[i] = rule
	}
	return result
}

func expandImpactTypes(items []interface{}) []monitoringCommon.ImpactType {
	if len(items) == 0 {
		return nil
	}
	impactTypeMap := map[string]monitoringCommon.ImpactType{
		"AVAILABILITY":     monitoringCommon.IMPACTTYPE_AVAILABILITY,
		"CAPACITY":         monitoringCommon.IMPACTTYPE_CAPACITY,
		"CONFIGURATION":    monitoringCommon.IMPACTTYPE_CONFIGURATION,
		"PERFORMANCE":      monitoringCommon.IMPACTTYPE_PERFORMANCE,
		"SYSTEM_INDICATOR": monitoringCommon.IMPACTTYPE_SYSTEM_INDICATOR,
		"CPU_CAPACITY":     monitoringCommon.IMPACTTYPE_CPU_CAPACITY,
		"MEMORY_CAPACITY":  monitoringCommon.IMPACTTYPE_MEMORY_CAPACITY,
		"STORAGE_CAPACITY": monitoringCommon.IMPACTTYPE_STORAGE_CAPACITY,
	}
	result := make([]monitoringCommon.ImpactType, 0, len(items))
	for _, item := range items {
		val := item.(string)
		if it, ok := impactTypeMap[val]; ok {
			result = append(result, it)
		}
	}
	return result
}

func expandSeverities(items []interface{}) []monitoringCommon.Severity {
	if len(items) == 0 {
		return nil
	}
	severityMap := map[string]monitoringCommon.Severity{
		"INFO":     monitoringCommon.SEVERITY_INFO,
		"WARNING":  monitoringCommon.SEVERITY_WARNING,
		"CRITICAL": monitoringCommon.SEVERITY_CRITICAL,
	}
	result := make([]monitoringCommon.Severity, 0, len(items))
	for _, item := range items {
		val := item.(string)
		if s, ok := severityMap[val]; ok {
			result = append(result, s)
		}
	}
	return result
}

func expandEmailTemplate(m map[string]interface{}) *monitoringService.EmailTemplate {
	if m == nil {
		return nil
	}
	tmpl := &monitoringService.EmailTemplate{}
	if v, ok := m["body_suffix"]; ok && v.(string) != "" {
		tmpl.BodySuffix = utils.StringPtr(v.(string))
	}
	if v, ok := m["subject_prefix"]; ok && v.(string) != "" {
		tmpl.SubjectPrefix = utils.StringPtr(v.(string))
	}
	return tmpl
}

func expandStringList(items []interface{}) []string {
	if len(items) == 0 {
		return nil
	}
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.(string)
	}
	return result
}

func flattenEmailConfigRules(rules []monitoringService.EmailConfigurationRule) []map[string]interface{} {
	if len(rules) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		m := map[string]interface{}{
			"cluster_uuids":                rule.ClusterUuids,
			"has_global_email_contact_list": utils.BoolValue(rule.HasGlobalEmailContactList),
			"is_enabled":                   utils.BoolValue(rule.IsEnabled),
			"match_phrases":                rule.MatchPhrases,
			"recipients":                   rule.Recipients,
		}

		if len(rule.ImpactTypes) > 0 {
			types := make([]string, len(rule.ImpactTypes))
			for j, it := range rule.ImpactTypes {
				types[j] = flattenEnumValue(&it)
			}
			m["impact_types"] = types
		}

		if len(rule.Severities) > 0 {
			sevs := make([]string, len(rule.Severities))
			for j, s := range rule.Severities {
				sevs[j] = flattenEnumValue(&s)
			}
			m["severities"] = sevs
		}

		result[i] = m
	}
	return result
}

func flattenEmailTemplate(tmpl *monitoringService.EmailTemplate) []map[string]interface{} {
	if tmpl == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"body_suffix":    utils.StringValue(tmpl.BodySuffix),
			"subject_prefix": utils.StringValue(tmpl.SubjectPrefix),
		},
	}
}

func flattenRemoteTunnelDetails(details *monitoringService.RemoteTunnelDetails) []map[string]interface{} {
	if details == nil {
		return nil
	}
	m := map[string]interface{}{
		"connection_status": flattenCommunicationStatus(details.ConnectionStatus),
		"http_proxy":        flattenHTTPProxy(details.HttpProxy),
		"service_center":    flattenServiceCenter(details.ServiceCenter),
		"transport_status":  flattenCommunicationStatus(details.TransportStatus),
	}
	return []map[string]interface{}{m}
}

func flattenCommunicationStatus(status *monitoringService.CommunicationStatus) []map[string]interface{} {
	if status == nil {
		return nil
	}
	m := map[string]interface{}{
		"last_changed_time":                 utils.TimeStringValue(status.LastChangedTime),
		"last_checked_time":                 utils.TimeStringValue(status.LastCheckedTime),
		"last_successful_transmission_time": utils.TimeStringValue(status.LastSuccessfulTransmissionTime),
	}
	if status.Status != nil {
		m["status"] = flattenEnumValue(status.Status)
	}
	if status.Message != nil {
		m["message"] = flattenParameterizedMessage(status.Message)
	}
	return []map[string]interface{}{m}
}

func flattenParameterizedMessage(msg *monitoringService.ParameterizedMessage) []map[string]interface{} {
	if msg == nil {
		return nil
	}
	m := map[string]interface{}{
		"message": utils.StringValue(msg.Message),
	}
	if len(msg.Attributes) > 0 {
		attrs := make([]map[string]interface{}, len(msg.Attributes))
		for i, attr := range msg.Attributes {
			attrs[i] = map[string]interface{}{
				"name":  utils.StringValue(attr.Name),
				"value": utils.StringValue(attr.Value),
			}
		}
		m["attributes"] = attrs
	}
	return []map[string]interface{}{m}
}

func flattenHTTPProxy(proxy *monitoringService.HttpProxy) []map[string]interface{} {
	if proxy == nil {
		return nil
	}
	m := map[string]interface{}{
		"name":     utils.StringValue(proxy.Name),
		"port":     utils.IntValue(proxy.Port),
		"username": utils.StringValue(proxy.Username),
	}
	if len(proxy.ProxyTypes) > 0 {
		types := make([]string, len(proxy.ProxyTypes))
		for i, pt := range proxy.ProxyTypes {
			types[i] = flattenEnumValue(&pt)
		}
		m["proxy_types"] = types
	}
	return []map[string]interface{}{m}
}

func flattenServiceCenter(sc *monitoringService.ServiceCenter) []map[string]interface{} {
	if sc == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ip_address": utils.StringValue(sc.IpAddress),
			"name":       utils.StringValue(sc.Name),
			"port":       utils.IntValue(sc.Port),
			"username":   utils.StringValue(sc.Username),
		},
	}
}
