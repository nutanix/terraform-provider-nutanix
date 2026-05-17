package cluster_managementv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clustermgmtConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	commonResponse "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/response"
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

func schemaForIPAddress() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"prefix_length": {
								Type:     schema.TypeInt,
								Computed: true,
							},
						},
					},
				},
				"ipv6": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"prefix_length": {
								Type:     schema.TypeInt,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func schemaForIPAddressInput() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
							"prefix_length": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
				"ipv6": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
							"prefix_length": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func flattenLinks(links []commonResponse.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		result = append(result, map[string]interface{}{
			"href": utils.StringValue(link.Href),
			"rel":  utils.StringValue(link.Rel),
		})
	}
	return result
}

func flattenIPAddress(ip *commonConfig.IPAddress) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	m := map[string]interface{}{}
	if ip.Ipv4 != nil {
		m["ipv4"] = []map[string]interface{}{
			{
				"value":         utils.StringValue(ip.Ipv4.Value),
				"prefix_length": utils.IntValue(ip.Ipv4.PrefixLength),
			},
		}
	} else {
		m["ipv4"] = []map[string]interface{}{}
	}
	if ip.Ipv6 != nil {
		m["ipv6"] = []map[string]interface{}{
			{
				"value":         utils.StringValue(ip.Ipv6.Value),
				"prefix_length": utils.IntValue(ip.Ipv6.PrefixLength),
			},
		}
	} else {
		m["ipv6"] = []map[string]interface{}{}
	}
	return []map[string]interface{}{m}
}

func flattenRsyslogModules(modules []clustermgmtConfig.RsyslogModuleItem) []map[string]interface{} {
	if len(modules) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(modules))
	for _, module := range modules {
		m := map[string]interface{}{
			"should_log_monitor_files": utils.BoolValue(module.ShouldLogMonitorFiles),
		}
		if module.Name != nil {
			m["name"] = module.Name.GetName()
		} else {
			m["name"] = ""
		}
		if module.LogSeverityLevel != nil {
			m["log_severity_level"] = module.LogSeverityLevel.GetName()
		} else {
			m["log_severity_level"] = ""
		}
		result = append(result, m)
	}
	return result
}

func expandIPAddress(ipList []interface{}) *commonConfig.IPAddress {
	if len(ipList) == 0 {
		return nil
	}
	ipMap := ipList[0].(map[string]interface{})
	ip := &commonConfig.IPAddress{}

	if ipv4List, ok := ipMap["ipv4"].([]interface{}); ok && len(ipv4List) > 0 {
		ipv4Map := ipv4List[0].(map[string]interface{})
		ip.Ipv4 = &commonConfig.IPv4Address{
			Value: utils.StringPtr(ipv4Map["value"].(string)),
		}
		if pl, ok := ipv4Map["prefix_length"].(int); ok && pl > 0 {
			ip.Ipv4.PrefixLength = utils.IntPtr(pl)
		}
	}

	if ipv6List, ok := ipMap["ipv6"].([]interface{}); ok && len(ipv6List) > 0 {
		ipv6Map := ipv6List[0].(map[string]interface{})
		ip.Ipv6 = &commonConfig.IPv6Address{
			Value: utils.StringPtr(ipv6Map["value"].(string)),
		}
		if pl, ok := ipv6Map["prefix_length"].(int); ok && pl > 0 {
			ip.Ipv6.PrefixLength = utils.IntPtr(pl)
		}
	}

	return ip
}

func expandRsyslogModules(modulesList []interface{}) []clustermgmtConfig.RsyslogModuleItem {
	if len(modulesList) == 0 {
		return nil
	}
	modules := make([]clustermgmtConfig.RsyslogModuleItem, 0, len(modulesList))
	for _, m := range modulesList {
		moduleMap := m.(map[string]interface{})
		item := clustermgmtConfig.RsyslogModuleItem{}

		if name, ok := moduleMap["name"].(string); ok && name != "" {
			item.Name = rsyslogModuleNameFromString(name)
		}
		if severity, ok := moduleMap["log_severity_level"].(string); ok && severity != "" {
			item.LogSeverityLevel = rsyslogModuleLogSeverityLevelFromString(severity)
		}
		if shouldLog, ok := moduleMap["should_log_monitor_files"].(bool); ok {
			item.ShouldLogMonitorFiles = utils.BoolPtr(shouldLog)
		}
		modules = append(modules, item)
	}
	return modules
}

func rsyslogNetworkProtocolFromString(s string) *clustermgmtConfig.RsyslogNetworkProtocol {
	protoMap := map[string]clustermgmtConfig.RsyslogNetworkProtocol{
		"UDP":  clustermgmtConfig.RSYSLOGNETWORKPROTOCOL_UDP,
		"TCP":  clustermgmtConfig.RSYSLOGNETWORKPROTOCOL_TCP,
		"RELP": clustermgmtConfig.RSYSLOGNETWORKPROTOCOL_RELP,
	}
	if v, ok := protoMap[s]; ok {
		return &v
	}
	return nil
}

func rsyslogModuleNameFromString(s string) *clustermgmtConfig.RsyslogModuleName {
	nameMap := map[string]clustermgmtConfig.RsyslogModuleName{
		"CASSANDRA":         clustermgmtConfig.RSYSLOGMODULENAME_CASSANDRA,
		"CEREBRO":           clustermgmtConfig.RSYSLOGMODULENAME_CEREBRO,
		"CURATOR":           clustermgmtConfig.RSYSLOGMODULENAME_CURATOR,
		"GENESIS":           clustermgmtConfig.RSYSLOGMODULENAME_GENESIS,
		"PRISM":             clustermgmtConfig.RSYSLOGMODULENAME_PRISM,
		"STARGATE":          clustermgmtConfig.RSYSLOGMODULENAME_STARGATE,
		"SYSLOG_MODULE":     clustermgmtConfig.RSYSLOGMODULENAME_SYSLOG_MODULE,
		"ZOOKEEPER":         clustermgmtConfig.RSYSLOGMODULENAME_ZOOKEEPER,
		"UHARA":             clustermgmtConfig.RSYSLOGMODULENAME_UHARA,
		"LAZAN":             clustermgmtConfig.RSYSLOGMODULENAME_LAZAN,
		"API_AUDIT":         clustermgmtConfig.RSYSLOGMODULENAME_API_AUDIT,
		"AUDIT":             clustermgmtConfig.RSYSLOGMODULENAME_AUDIT,
		"CALM":              clustermgmtConfig.RSYSLOGMODULENAME_CALM,
		"EPSILON":           clustermgmtConfig.RSYSLOGMODULENAME_EPSILON,
		"ACROPOLIS":         clustermgmtConfig.RSYSLOGMODULENAME_ACROPOLIS,
		"MINERVA_CVM":       clustermgmtConfig.RSYSLOGMODULENAME_MINERVA_CVM,
		"FLOW":              clustermgmtConfig.RSYSLOGMODULENAME_FLOW,
		"FLOW_SERVICE_LOGS": clustermgmtConfig.RSYSLOGMODULENAME_FLOW_SERVICE_LOGS,
		"LCM":               clustermgmtConfig.RSYSLOGMODULENAME_LCM,
		"APLOS":             clustermgmtConfig.RSYSLOGMODULENAME_APLOS,
		"NCM_AIOPS":         clustermgmtConfig.RSYSLOGMODULENAME_NCM_AIOPS,
	}
	if v, ok := nameMap[s]; ok {
		return &v
	}
	return nil
}

func rsyslogModuleLogSeverityLevelFromString(s string) *clustermgmtConfig.RsyslogModuleLogSeverityLevel {
	levelMap := map[string]clustermgmtConfig.RsyslogModuleLogSeverityLevel{
		"EMERGENCY": clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_EMERGENCY,
		"ALERT":     clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_ALERT,
		"CRITICAL":  clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_CRITICAL,
		"ERROR":     clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_ERROR,
		"WARNING":   clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_WARNING,
		"NOTICE":    clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_NOTICE,
		"INFO":      clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_INFO,
		"DEBUG":     clustermgmtConfig.RSYSLOGMODULELOGSEVERITYLEVEL_DEBUG,
	}
	if v, ok := levelMap[s]; ok {
		return &v
	}
	return nil
}
