package networkingv2

import (
	config "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func expandFlexRuleSpec(l []interface{}) *import1.FlexRuleSpec {
	if len(l) == 0 {
		return nil
	}
	val := l[0].(map[string]interface{})
	flex := import1.NewFlexRuleSpec()

	if action, ok := val["action"].(string); ok && len(action) > 0 {
		flex.Action = common.ExpandEnum[import1.FlexRuleAction](action)
	}
	if direction, ok := val["direction"].(string); ok && len(direction) > 0 {
		flex.Direction = common.ExpandEnum[import1.FlexRuleDirection](direction)
	}
	if v, ok := val["applied_to_entity_group_references"].([]interface{}); ok && len(v) > 0 {
		flex.AppliedToEntityGroupReferences = common.ExpandListOfString(v)
	}
	if v, ok := val["dest_entity_group_references"].([]interface{}); ok && len(v) > 0 {
		flex.DestEntityGroupReferences = common.ExpandListOfString(v)
	}
	if v, ok := val["dest_subnet"].([]interface{}); ok && len(v) > 0 {
		flex.DestSubnet = expandIPv4AddressMicroseg(v)
	}
	if v, ok := val["src_entity_group_references"].([]interface{}); ok && len(v) > 0 {
		flex.SrcEntityGroupReferences = common.ExpandListOfString(v)
	}
	if v, ok := val["src_subnet"].([]interface{}); ok && len(v) > 0 {
		flex.SrcSubnet = expandIPv4AddressMicroseg(v)
	}
	if v, ok := val["is_all_protocol_allowed"]; ok {
		flex.IsAllProtocolAllowed = utils.BoolPtr(v.(bool))
	}
	if v, ok := val["service_group_references"].([]interface{}); ok && len(v) > 0 {
		flex.ServiceGroupReferences = common.ExpandListOfString(v)
	}
	if v, ok := val["tcp_services"].([]interface{}); ok && len(v) > 0 {
		flex.TcpServices = expandTCPPortRangeSpec(v)
	}
	if v, ok := val["udp_services"].([]interface{}); ok && len(v) > 0 {
		flex.UdpServices = expandUDPPortRangeSpec(v)
	}
	if v, ok := val["icmp_services"].([]interface{}); ok && len(v) > 0 {
		flex.IcmpServices = expandIcmpTypeCodeSpec(v)
	}
	if v, ok := val["icmpv6_services"].([]interface{}); ok && len(v) > 0 {
		flex.IcmpV6Services = expandIcmpV6TypeCodeSpec(v)
	}
	if v, ok := val["should_allow_any_dst"]; ok {
		flex.ShouldAllowAnyDst = utils.BoolPtr(v.(bool))
	}
	if v, ok := val["should_allow_any_src"]; ok {
		flex.ShouldAllowAnySrc = utils.BoolPtr(v.(bool))
	}
	if v, ok := val["network_function_reference"].(string); ok && len(v) > 0 {
		flex.NetworkFunctionReference = utils.StringPtr(v)
	}
	if v, ok := val["priority"].(int); ok && v > 0 {
		flex.Priority = utils.IntPtr(v)
	}
	if v, ok := val["ip_version"].(string); ok && len(v) > 0 {
		flex.IpVersion = common.ExpandEnum[config.IPAddressScope](v)
	}
	return flex
}

func expandIcmpV6TypeCodeSpec(pr []interface{}) []import1.IcmpV6TypeCodeSpec {
	if len(pr) == 0 {
		return nil
	}
	specs := make([]import1.IcmpV6TypeCodeSpec, len(pr))
	for k, v := range pr {
		val := v.(map[string]interface{})
		spec := import1.IcmpV6TypeCodeSpec{}
		if isAll, ok := val["is_all_allowed"]; ok {
			spec.IsAllAllowed = utils.BoolPtr(isAll.(bool))
		}
		if t, ok := val["type"]; ok {
			spec.Type = utils.IntPtr(t.(int))
		}
		if c, ok := val["code"]; ok {
			spec.Code = utils.IntPtr(c.(int))
		}
		specs[k] = spec
	}
	return specs
}

func flattenFlexRuleSpec(flex import1.FlexRuleSpec) []map[string]interface{} {
	m := map[string]interface{}{}

	if flex.Action != nil {
		m["action"] = common.FlattenPtrEnum(flex.Action)
	}
	if flex.Direction != nil {
		m["direction"] = common.FlattenPtrEnum(flex.Direction)
	}
	if flex.AppliedToEntityGroupReferences != nil {
		m["applied_to_entity_group_references"] = flex.AppliedToEntityGroupReferences
	}
	if flex.DestEntityGroupReferences != nil {
		m["dest_entity_group_references"] = flex.DestEntityGroupReferences
	}
	if flex.DestSubnet != nil {
		m["dest_subnet"] = flattenIPv4AddressMicroSegList(flex.DestSubnet)
	}
	if flex.SrcEntityGroupReferences != nil {
		m["src_entity_group_references"] = flex.SrcEntityGroupReferences
	}
	if flex.SrcSubnet != nil {
		m["src_subnet"] = flattenIPv4AddressMicroSegList(flex.SrcSubnet)
	}
	if flex.IsAllProtocolAllowed != nil {
		m["is_all_protocol_allowed"] = utils.BoolValue(flex.IsAllProtocolAllowed)
	}
	if flex.ServiceGroupReferences != nil {
		m["service_group_references"] = flex.ServiceGroupReferences
	}
	if flex.TcpServices != nil {
		m["tcp_services"] = flattenTCPPortRangeSpec(flex.TcpServices)
	}
	if flex.UdpServices != nil {
		m["udp_services"] = flattenUDPPortRangeSpec(flex.UdpServices)
	}
	if flex.IcmpServices != nil {
		m["icmp_services"] = flattenIcmpTypeCodeSpec(flex.IcmpServices)
	}
	if flex.IcmpV6Services != nil {
		m["icmpv6_services"] = flattenIcmpV6TypeCodeSpec(flex.IcmpV6Services)
	}
	if flex.ShouldAllowAnyDst != nil {
		m["should_allow_any_dst"] = utils.BoolValue(flex.ShouldAllowAnyDst)
	}
	if flex.ShouldAllowAnySrc != nil {
		m["should_allow_any_src"] = utils.BoolValue(flex.ShouldAllowAnySrc)
	}
	if flex.NetworkFunctionReference != nil {
		m["network_function_reference"] = utils.StringValue(flex.NetworkFunctionReference)
	}
	if flex.Priority != nil {
		m["priority"] = utils.IntValue(flex.Priority)
	}
	if flex.IpVersion != nil {
		m["ip_version"] = common.FlattenPtrEnum(flex.IpVersion)
	}
	if flex.IsSystemRule != nil {
		m["is_system_rule"] = utils.BoolValue(flex.IsSystemRule)
	}
	return []map[string]interface{}{m}
}

func flattenIcmpV6TypeCodeSpec(specs []import1.IcmpV6TypeCodeSpec) []map[string]interface{} {
	if len(specs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(specs))
	for k, v := range specs {
		m := map[string]interface{}{}
		if v.IsAllAllowed != nil {
			m["is_all_allowed"] = utils.BoolValue(v.IsAllAllowed)
		}
		if v.Type != nil {
			m["type"] = utils.IntValue(v.Type)
		}
		if v.Code != nil {
			m["code"] = utils.IntValue(v.Code)
		}
		result[k] = m
	}
	return result
}
