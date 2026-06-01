package microsegv2

import (
	commonconfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/response"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func expandDomainControllers(l []interface{}) []commonconfig.IPAddressOrFQDN {
	if len(l) == 0 {
		return nil
	}
	out := make([]commonconfig.IPAddressOrFQDN, 0, len(l))
	for _, item := range l {
		m := item.(map[string]interface{})
		dc := commonconfig.IPAddressOrFQDN{}
		if v, ok := m["fqdn"].([]interface{}); ok && len(v) > 0 {
			fqdnMap := v[0].(map[string]interface{})
			fqdn := commonconfig.NewFQDN()
			if val, ok := fqdnMap["value"].(string); ok && val != "" {
				fqdn.Value = utils.StringPtr(val)
			}
			dc.Fqdn = fqdn
		}
		if v, ok := m["ipv4"].([]interface{}); ok && len(v) > 0 {
			ipv4Map := v[0].(map[string]interface{})
			ipv4 := commonconfig.NewIPv4Address()
			if val, ok := ipv4Map["value"].(string); ok && val != "" {
				ipv4.Value = utils.StringPtr(val)
			}
			if val, ok := ipv4Map["prefix_length"].(int); ok && val > 0 {
				ipv4.PrefixLength = utils.IntPtr(val)
			}
			dc.Ipv4 = ipv4
		}
		if v, ok := m["ipv6"].([]interface{}); ok && len(v) > 0 {
			ipv6Map := v[0].(map[string]interface{})
			ipv6 := commonconfig.NewIPv6Address()
			if val, ok := ipv6Map["value"].(string); ok && val != "" {
				ipv6.Value = utils.StringPtr(val)
			}
			if val, ok := ipv6Map["prefix_length"].(int); ok && val > 0 {
				ipv6.PrefixLength = utils.IntPtr(val)
			}
			dc.Ipv6 = ipv6
		}
		out = append(out, dc)
	}
	return out
}

func expandMatchingCriterias(l []interface{}) []import2.MatchingCriteria {
	if len(l) == 0 {
		return nil
	}
	out := make([]import2.MatchingCriteria, 0, len(l))
	for _, item := range l {
		m := item.(map[string]interface{})
		mc := import2.MatchingCriteria{}
		if v, ok := m["criteria"].(string); ok && v != "" {
			mc.Criteria = utils.StringPtr(v)
		}
		if v, ok := m["match_entity"].(string); ok && v != "" {
			mc.MatchEntity = commonUtils.ExpandEnum[import2.MatchEntity](v)
		}
		if v, ok := m["match_field"].(string); ok && v != "" {
			mc.MatchField = commonUtils.ExpandEnum[import2.MatchField](v)
		}
		if v, ok := m["match_type"].(string); ok && v != "" {
			mc.MatchType = commonUtils.ExpandEnum[import2.MatchType](v)
		}
		out = append(out, mc)
	}
	return out
}

func expandAdInfo(l []interface{}) *import2.AdInfo {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	m := l[0].(map[string]interface{})
	ai := import2.NewAdInfo()
	if v, ok := m["directory_service_reference"].(string); ok && v != "" {
		ai.DirectoryServiceReference = utils.StringPtr(v)
	}
	if v, ok := m["object_identifier"].(string); ok && v != "" {
		ai.ObjectIdentifier = utils.StringPtr(v)
	}
	if v, ok := m["object_path"].(string); ok && v != "" {
		ai.ObjectPath = utils.StringPtr(v)
	}
	if v, ok := m["status"].(string); ok && v != "" {
		ai.Status = commonUtils.ExpandEnum[import2.AdStatus](v)
	}
	return ai
}

func flattenDomainControllers(dcs []commonconfig.IPAddressOrFQDN) []map[string]interface{} {
	if len(dcs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(dcs))
	for _, dc := range dcs {
		m := map[string]interface{}{}
		if dc.Fqdn != nil {
			m["fqdn"] = []map[string]interface{}{
				{"value": utils.StringValue(dc.Fqdn.Value)},
			}
		} else {
			m["fqdn"] = []map[string]interface{}{}
		}
		if dc.Ipv4 != nil {
			ipv4Map := map[string]interface{}{
				"value": utils.StringValue(dc.Ipv4.Value),
			}
			if dc.Ipv4.PrefixLength != nil {
				ipv4Map["prefix_length"] = *dc.Ipv4.PrefixLength
			} else {
				ipv4Map["prefix_length"] = 0
			}
			m["ipv4"] = []map[string]interface{}{ipv4Map}
		} else {
			m["ipv4"] = []map[string]interface{}{}
		}
		if dc.Ipv6 != nil {
			ipv6Map := map[string]interface{}{
				"value": utils.StringValue(dc.Ipv6.Value),
			}
			if dc.Ipv6.PrefixLength != nil {
				ipv6Map["prefix_length"] = *dc.Ipv6.PrefixLength
			} else {
				ipv6Map["prefix_length"] = 0
			}
			m["ipv6"] = []map[string]interface{}{ipv6Map}
		} else {
			m["ipv6"] = []map[string]interface{}{}
		}
		result = append(result, m)
	}
	return result
}

func flattenMatchingCriterias(mcs []import2.MatchingCriteria) []map[string]interface{} {
	if len(mcs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(mcs))
	for _, mc := range mcs {
		m := map[string]interface{}{
			"criteria": utils.StringValue(mc.Criteria),
		}
		if mc.MatchEntity != nil {
			m["match_entity"] = mc.MatchEntity.GetName()
		} else {
			m["match_entity"] = ""
		}
		if mc.MatchField != nil {
			m["match_field"] = mc.MatchField.GetName()
		} else {
			m["match_field"] = ""
		}
		if mc.MatchType != nil {
			m["match_type"] = mc.MatchType.GetName()
		} else {
			m["match_type"] = ""
		}
		result = append(result, m)
	}
	return result
}

func flattenAdInfo(ai *import2.AdInfo) []map[string]interface{} {
	if ai == nil {
		return nil
	}
	m := map[string]interface{}{
		"directory_service_reference": utils.StringValue(ai.DirectoryServiceReference),
		"object_identifier":           utils.StringValue(ai.ObjectIdentifier),
		"object_path":                 utils.StringValue(ai.ObjectPath),
	}
	if ai.Status != nil {
		m["status"] = ai.Status.GetName()
	} else {
		m["status"] = ""
	}
	return []map[string]interface{}{m}
}

func flattenLinksDSC(links []import1.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		m := map[string]interface{}{
			"href": utils.StringValue(link.Href),
			"rel":  utils.StringValue(link.Rel),
		}
		result = append(result, m)
	}
	return result
}

func flattenDirectoryServerConfigs(configs []import2.DirectoryServerConfig) []map[string]interface{} {
	if len(configs) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		m := map[string]interface{}{
			"ext_id":                                utils.StringValue(c.ExtId),
			"directory_service_reference":           utils.StringValue(c.DirectoryServiceReference),
			"domain_controllers":                    flattenDomainControllers(c.DomainControllers),
			"is_default_category_enabled":           utils.BoolValue(c.IsDefaultCategoryEnabled),
			"matching_criterias":                    flattenMatchingCriterias(c.MatchingCriterias),
			"should_keep_default_category_on_login": utils.BoolValue(c.ShouldKeepDefaultCategoryOnLogin),
			"tenant_id":                             utils.StringValue(c.TenantId),
			"links":                                 flattenLinksDSC(c.Links),
			"project_ext_id":                        utils.StringValue(c.ProjectExtId),
		}
		result = append(result, m)
	}
	return result
}

func flattenCategoryMappings(mappings []import2.CategoryMapping) []map[string]interface{} {
	if len(mappings) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(mappings))
	for _, cm := range mappings {
		m := map[string]interface{}{
			"ext_id":         utils.StringValue(cm.ExtId),
			"name":           utils.StringValue(cm.Name),
			"category_name":  utils.StringValue(cm.CategoryName),
			"category_value": utils.StringValue(cm.CategoryValue),
			"ad_info":        flattenAdInfo(cm.AdInfo),
			"tenant_id":      utils.StringValue(cm.TenantId),
			"links":          flattenLinksDSC(cm.Links),
			"project_ext_id": utils.StringValue(cm.ProjectExtId),
		}
		result = append(result, m)
	}
	return result
}
