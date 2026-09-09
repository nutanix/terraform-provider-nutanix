package lifecyclev2

import (
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/response"
	import3 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// flattenIPAddress converts an SDK IPAddress into the Terraform list-of-map representation.
func flattenIPAddress(ip *import1.IPAddress) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	out := make(map[string]interface{})
	if ip.Ipv4 != nil {
		out["ipv4"] = flattenIPv4Address(ip.Ipv4)
	}
	if ip.Ipv6 != nil {
		out["ipv6"] = flattenIPv6Address(ip.Ipv6)
	}
	return []map[string]interface{}{out}
}

func flattenIPv4Address(ip *import1.IPv4Address) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	out := make(map[string]interface{})
	out["value"] = utils.StringValue(ip.Value)
	if ip.PrefixLength != nil {
		out["prefix_length"] = utils.IntValue(ip.PrefixLength)
	}
	return []map[string]interface{}{out}
}

func flattenIPv6Address(ip *import1.IPv6Address) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	out := make(map[string]interface{})
	out["value"] = utils.StringValue(ip.Value)
	if ip.PrefixLength != nil {
		out["prefix_length"] = utils.IntValue(ip.PrefixLength)
	}
	return []map[string]interface{}{out}
}

// expandIPAddress builds an SDK IPAddress from the Terraform list-of-map representation.
func expandIPAddress(pr []interface{}) *import1.IPAddress {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ip := &import1.IPAddress{}
	if v, ok := val["ipv4"]; ok && len(v.([]interface{})) > 0 {
		ip.Ipv4 = expandIPv4Address(v.([]interface{}))
	}
	if v, ok := val["ipv6"]; ok && len(v.([]interface{})) > 0 {
		ip.Ipv6 = expandIPv6Address(v.([]interface{}))
	}
	if ip.Ipv4 == nil && ip.Ipv6 == nil {
		return nil
	}
	return ip
}

func expandIPv4Address(pr []interface{}) *import1.IPv4Address {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ip := &import1.IPv4Address{}
	if v, ok := val["value"]; ok && v.(string) != "" {
		ip.Value = utils.StringPtr(v.(string))
	}
	if v, ok := val["prefix_length"]; ok && v.(int) != 0 {
		ip.PrefixLength = utils.IntPtr(v.(int))
	}
	if ip.Value == nil {
		return nil
	}
	return ip
}

func expandIPv6Address(pr []interface{}) *import1.IPv6Address {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ip := &import1.IPv6Address{}
	if v, ok := val["value"]; ok && v.(string) != "" {
		ip.Value = utils.StringPtr(v.(string))
	}
	if v, ok := val["prefix_length"]; ok && v.(int) != 0 {
		ip.PrefixLength = utils.IntPtr(v.(int))
	}
	if ip.Value == nil {
		return nil
	}
	return ip
}

// flattenNodeLinks flattens the HATEOAS links of the node into Terraform state.
func flattenNodeLinks(links []import2.ApiLink) []map[string]interface{} {
	return common.FlattenLinks(links)
}

// flattenKVStringPairs flattens a list of key-value string pairs.
func flattenKVStringPairs(pairs []import1.KVStringPair) []map[string]interface{} {
	if len(pairs) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(pairs))
	for i, p := range pairs {
		out[i] = map[string]interface{}{
			"name":  utils.StringValue(p.Name),
			"value": utils.StringValue(p.Value),
		}
	}
	return out
}

// nodeStateName returns the string name of a node state enum (empty if nil).
func nodeStateName(s *import3.NodeState) string {
	return common.FlattenPtrEnum(s)
}
