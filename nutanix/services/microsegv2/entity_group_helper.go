package microsegv2

import (
	commonconfig "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/response"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func expandAllowedConfig(l []interface{}) *import2.AllowedConfig {
	if len(l) == 0 {
		return nil
	}
	m := l[0].(map[string]interface{})
	entities, ok := m["entities"].([]interface{})
	if !ok || len(entities) == 0 {
		return nil
	}
	cfg := import2.NewAllowedConfig()
	cfg.Entities = expandAllowedEntities(entities)
	return cfg
}

func expandAllowedEntities(l []interface{}) []import2.AllowedEntity {
	out := make([]import2.AllowedEntity, 0, len(l))
	for _, e := range l {
		m := e.(map[string]interface{})
		ent := import2.AllowedEntity{}
		if v, ok := m["selected_by"].(string); ok {
			ent.SelectBy = common.ExpandEnum[import2.AllowedSelectBy](v)
		}
		if v, ok := m["type"].(string); ok {
			ent.Type = common.ExpandEnum[import2.AllowedType](v)
		}
		if v, ok := m["addresses"].([]interface{}); ok && len(v) > 0 {
			ent.Addresses = expandAddresses(v)
		}
		if v, ok := m["ip_ranges"].([]interface{}); ok && len(v) > 0 {
			ent.IpRanges = expandIpRange(v)
		}
		if v, ok := m["kube_entities"].([]interface{}); ok {
			ent.KubeEntities = common.ExpandListOfString(v)
		}
		if v, ok := m["reference_ext_ids"].([]interface{}); ok {
			ent.ReferenceExtIds = common.ExpandListOfString(v)
		}
		out = append(out, ent)
	}
	return out
}

func expandExceptConfig(l []interface{}) *import2.ExceptConfig {
	if len(l) == 0 {
		return nil
	}
	m := l[0].(map[string]interface{})
	entities, ok := m["entities"].([]interface{})
	if !ok || len(entities) == 0 {
		return nil
	}
	cfg := import2.NewExceptConfig()
	cfg.Entities = expandExceptEntities(entities)
	return cfg
}

func expandExceptEntities(l []interface{}) []import2.ExceptEntity {
	out := make([]import2.ExceptEntity, 0, len(l))
	for _, e := range l {
		m := e.(map[string]interface{})
		ent := import2.ExceptEntity{}
		if v, ok := m["selected_by"].(string); ok {
			ent.SelectBy = common.ExpandEnum[import2.ExceptSelectBy](v)
		}
		if v, ok := m["type"].(string); ok {
			ent.Type = common.ExpandEnum[import2.ExceptType](v)
		}
		if v, ok := m["addresses"].([]interface{}); ok && len(v) > 0 {
			ent.Addresses = expandAddresses(v)
		}
		if v, ok := m["ip_ranges"].([]interface{}); ok && len(v) > 0 {
			ent.IpRanges = expandIpRange(v)
		}
		if v, ok := m["reference_ext_ids"].([]interface{}); ok {
			ent.ReferenceExtIds = common.ExpandListOfString(v)
		}
		out = append(out, ent)
	}
	return out
}

func expandAddresses(l []interface{}) *import2.Addresses {
	if len(l) == 0 {
		return nil
	}
	m := l[0].(map[string]interface{})
	ipList, ok := m["ipv4_addresses"].([]interface{})
	if !ok || len(ipList) == 0 {
		return nil
	}
	addr := import2.NewAddresses()
	addr.Ipv4Addresses = make([]commonconfig.IPv4Address, 0, len(ipList))
	for _, ip := range ipList {
		im := ip.(map[string]interface{})
		ipp := commonconfig.NewIPv4Address()
		if v, ok := im["value"].(string); ok {
			ipp.Value = utils.StringPtr(v)
		}
		if v, ok := im["prefix_length"].(int); ok {
			ipp.PrefixLength = utils.IntPtr(v)
		}
		addr.Ipv4Addresses = append(addr.Ipv4Addresses, *ipp)
	}
	return addr
}

func expandIpRange(l []interface{}) *import2.IpRange {
	if len(l) == 0 {
		return nil
	}
	m := l[0].(map[string]interface{})
	rangeList, ok := m["ipv4_ranges"].([]interface{})
	if !ok || len(rangeList) == 0 {
		return nil
	}
	ir := import2.NewIpRange()
	ir.Ipv4Ranges = make([]import2.IPv4Range, 0, len(rangeList))
	for _, r := range rangeList {
		rm := r.(map[string]interface{})
		rr := import2.NewIPv4Range()
		if v, ok := rm["start_ip"].(string); ok {
			rr.StartIp = utils.StringPtr(v)
		}
		if v, ok := rm["end_ip"].(string); ok {
			rr.EndIp = utils.StringPtr(v)
		}
		ir.Ipv4Ranges = append(ir.Ipv4Ranges, *rr)
	}
	return ir
}

// Flatten helpers

func flattenLinksEntityGroup(links []import1.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
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

func flattenAllowedConfig(cfg *import2.AllowedConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"entities": flattenAllowedEntities(cfg.Entities),
		},
	}
}

func flattenAllowedEntities(entities []import2.AllowedEntity) []map[string]interface{} {
	if len(entities) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(entities))
	for _, e := range entities {
		m := map[string]interface{}{
			"selected_by":       e.SelectBy.GetName(),
			"type":              e.Type.GetName(),
			"addresses":         flattenAddresses(e.Addresses),
			"ip_ranges":         flattenIpRanges(e.IpRanges),
			"kube_entities":     e.KubeEntities,
			"reference_ext_ids": e.ReferenceExtIds,
		}
		result = append(result, m)
	}
	return result
}

func flattenExceptConfig(cfg *import2.ExceptConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"entities": flattenExceptEntities(cfg.Entities),
		},
	}
}

func flattenExceptEntities(entities []import2.ExceptEntity) []map[string]interface{} {
	if len(entities) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(entities))
	for _, e := range entities {
		m := map[string]interface{}{
			"selected_by":       e.SelectBy.GetName(),
			"type":              e.Type.GetName(),
			"addresses":         flattenAddresses(e.Addresses),
			"ip_ranges":         flattenIpRanges(e.IpRanges),
			"reference_ext_ids": e.ReferenceExtIds,
		}
		result = append(result, m)
	}
	return result
}

func flattenAddresses(addr *import2.Addresses) []map[string]interface{} {
	if addr == nil || len(addr.Ipv4Addresses) == 0 {
		return nil
	}
	ipList := make([]map[string]interface{}, 0, len(addr.Ipv4Addresses))
	for _, ip := range addr.Ipv4Addresses {
		m := map[string]interface{}{
			"value": utils.StringValue(ip.Value),
		}
		if ip.PrefixLength != nil {
			m["prefix_length"] = *ip.PrefixLength
		}
		ipList = append(ipList, m)
	}
	return []map[string]interface{}{
		{"ipv4_addresses": ipList},
	}
}

func flattenIpRanges(ir *import2.IpRange) []map[string]interface{} {
	if ir == nil || len(ir.Ipv4Ranges) == 0 {
		return nil
	}
	rangeList := make([]map[string]interface{}, 0, len(ir.Ipv4Ranges))
	for _, r := range ir.Ipv4Ranges {
		rangeList = append(rangeList, map[string]interface{}{
			"start_ip": utils.StringValue(r.StartIp),
			"end_ip":   utils.StringValue(r.EndIp),
		})
	}
	return []map[string]interface{}{
		{"ipv4_ranges": rangeList},
	}
}
