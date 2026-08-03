package microsegv2

import (
	commonCfg "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// expandDirectoryServerConfigMatchingCriterias converts the HCL matching
// criteria blocks into SDK MatchingCriteria values.
func expandDirectoryServerConfigMatchingCriterias(l []interface{}) []import2.MatchingCriteria {
	if len(l) == 0 {
		return nil
	}
	out := make([]import2.MatchingCriteria, 0, len(l))
	for _, e := range l {
		m, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		mc := import2.NewMatchingCriteria()
		if v, ok := m["criteria"].(string); ok && v != "" {
			mc.Criteria = utils.StringPtr(v)
		}
		if v, ok := m["match_entity"].(string); ok && v != "" {
			mc.MatchEntity = common.ExpandEnum[import2.MatchEntity](v)
		}
		if v, ok := m["match_field"].(string); ok && v != "" {
			mc.MatchField = common.ExpandEnum[import2.MatchField](v)
		}
		if v, ok := m["match_type"].(string); ok && v != "" {
			mc.MatchType = common.ExpandEnum[import2.MatchType](v)
		}
		out = append(out, *mc)
	}
	return out
}

// expandIPAddressOrFQDNList converts a list of HCL domain controller blocks
// into SDK IPAddressOrFQDN values.
func expandIPAddressOrFQDNList(l []interface{}) []commonCfg.IPAddressOrFQDN {
	if len(l) == 0 {
		return nil
	}
	out := make([]commonCfg.IPAddressOrFQDN, 0, len(l))
	for _, e := range l {
		m, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		dc := commonCfg.NewIPAddressOrFQDN()
		if v, ok := m["ipv4"].([]interface{}); ok && len(v) > 0 {
			dc.Ipv4 = expandIPv4AddressValue(v)
		}
		if v, ok := m["ipv6"].([]interface{}); ok && len(v) > 0 {
			dc.Ipv6 = expandIPv6AddressValue(v)
		}
		if v, ok := m["fqdn"].([]interface{}); ok && len(v) > 0 {
			dc.Fqdn = expandFQDNValue(v)
		}
		out = append(out, *dc)
	}
	return out
}

func expandIPv4AddressValue(l []interface{}) *commonCfg.IPv4Address {
	if len(l) == 0 {
		return nil
	}
	m, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}
	ip := commonCfg.NewIPv4Address()
	if v, ok := m["value"].(string); ok && v != "" {
		ip.Value = utils.StringPtr(v)
	}
	if v, ok := m["prefix_length"].(int); ok && v != 0 {
		ip.PrefixLength = utils.IntPtr(v)
	}
	return ip
}

func expandIPv6AddressValue(l []interface{}) *commonCfg.IPv6Address {
	if len(l) == 0 {
		return nil
	}
	m, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}
	ip := commonCfg.NewIPv6Address()
	if v, ok := m["value"].(string); ok && v != "" {
		ip.Value = utils.StringPtr(v)
	}
	if v, ok := m["prefix_length"].(int); ok && v != 0 {
		ip.PrefixLength = utils.IntPtr(v)
	}
	return ip
}

func expandFQDNValue(l []interface{}) *commonCfg.FQDN {
	if len(l) == 0 {
		return nil
	}
	m, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}
	f := commonCfg.NewFQDN()
	if v, ok := m["value"].(string); ok && v != "" {
		f.Value = utils.StringPtr(v)
	}
	return f
}

// expandAdInfo converts the HCL ad_info block into an SDK AdInfo value.
func expandAdInfo(l []interface{}) *import2.AdInfo {
	if len(l) == 0 {
		return nil
	}
	m, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}
	adInfo := import2.NewAdInfo()
	if v, ok := m["directory_service_reference"].(string); ok && v != "" {
		adInfo.DirectoryServiceReference = utils.StringPtr(v)
	}
	if v, ok := m["object_identifier"].(string); ok && v != "" {
		adInfo.ObjectIdentifier = utils.StringPtr(v)
	}
	if v, ok := m["status"].(string); ok && v != "" {
		adInfo.Status = common.ExpandEnum[import2.AdStatus](v)
	}
	return adInfo
}

// flattenMatchingCriterias converts SDK MatchingCriteria values into the HCL
// representation.
func flattenMatchingCriterias(criterias []import2.MatchingCriteria) []map[string]interface{} {
	if len(criterias) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(criterias))
	for _, c := range criterias {
		m := map[string]interface{}{
			"criteria": utils.StringValue(c.Criteria),
		}
		if c.MatchEntity != nil {
			m["match_entity"] = c.MatchEntity.GetName()
		}
		if c.MatchField != nil {
			m["match_field"] = c.MatchField.GetName()
		}
		if c.MatchType != nil {
			m["match_type"] = c.MatchType.GetName()
		}
		out = append(out, m)
	}
	return out
}

// flattenIPAddressOrFQDNList converts SDK IPAddressOrFQDN values into the HCL
// representation.
func flattenIPAddressOrFQDNList(dcs []commonCfg.IPAddressOrFQDN) []map[string]interface{} {
	if len(dcs) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(dcs))
	for _, dc := range dcs {
		m := map[string]interface{}{}
		if dc.Ipv4 != nil {
			m["ipv4"] = []map[string]interface{}{
				{
					"value":         utils.StringValue(dc.Ipv4.Value),
					"prefix_length": utils.IntValue(dc.Ipv4.PrefixLength),
				},
			}
		}
		if dc.Ipv6 != nil {
			m["ipv6"] = []map[string]interface{}{
				{
					"value":         utils.StringValue(dc.Ipv6.Value),
					"prefix_length": utils.IntValue(dc.Ipv6.PrefixLength),
				},
			}
		}
		if dc.Fqdn != nil {
			m["fqdn"] = []map[string]interface{}{
				{
					"value": utils.StringValue(dc.Fqdn.Value),
				},
			}
		}
		out = append(out, m)
	}
	return out
}

// flattenAdInfo converts an SDK AdInfo value into the HCL representation.
func flattenAdInfo(adInfo *import2.AdInfo) []map[string]interface{} {
	if adInfo == nil {
		return nil
	}
	m := map[string]interface{}{
		"directory_service_reference": utils.StringValue(adInfo.DirectoryServiceReference),
		"object_identifier":           utils.StringValue(adInfo.ObjectIdentifier),
		"object_path":                 utils.StringValue(adInfo.ObjectPath),
	}
	if adInfo.Status != nil {
		m["status"] = adInfo.Status.GetName()
	}
	return []map[string]interface{}{m}
}
