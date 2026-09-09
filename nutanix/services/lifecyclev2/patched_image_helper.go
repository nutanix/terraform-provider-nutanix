package lifecyclev2

import (
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// expandLocalHostImageDetails builds the OneOf image details for a patched image.
// Only the LocalHostImageDetails variant (an uploaded host image ext_id) is
// supported by the SDK.
func expandImageDetails(in []interface{}) *import1.OneOfPatchedImageImageDetails {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	oneOf := import1.NewOneOfPatchedImageImageDetails()
	details := *import1.NewLocalHostImageDetails()
	if v, ok := val["local_host_image_ext_id"]; ok && v.(string) != "" {
		details.ExtId = utils.StringPtr(v.(string))
	}
	if err := oneOf.SetValue(details); err != nil {
		return nil
	}
	return oneOf
}

func flattenImageDetails(oneOf *import1.OneOfPatchedImageImageDetails) []map[string]interface{} {
	if oneOf == nil {
		return nil
	}
	value := oneOf.GetValue()
	if details, ok := value.(import1.LocalHostImageDetails); ok {
		return []map[string]interface{}{
			{
				"local_host_image_ext_id": utils.StringValue(details.ExtId),
			},
		}
	}
	return nil
}

func expandNodeConfigurations(in []interface{}) []import1.NodeConfiguration {
	if len(in) == 0 {
		return nil
	}
	out := make([]import1.NodeConfiguration, 0, len(in))
	for _, raw := range in {
		val := raw.(map[string]interface{})
		nc := *import1.NewNodeConfiguration()
		if v, ok := val["node_ext_id"]; ok && v.(string) != "" {
			nc.NodeExtId = utils.StringPtr(v.(string))
		}
		if v, ok := val["host_configuration"]; ok {
			nc.HostConfiguration = expandHostConfiguration(v.([]interface{}))
		}
		out = append(out, nc)
	}
	return out
}

func flattenNodeConfigurations(in []import1.NodeConfiguration) []map[string]interface{} {
	if len(in) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(in))
	for i := range in {
		nc := in[i]
		out = append(out, map[string]interface{}{
			"node_ext_id":        utils.StringValue(nc.NodeExtId),
			"host_configuration": flattenHostConfiguration(nc.HostConfiguration),
		})
	}
	return out
}

func expandHostConfiguration(in []interface{}) *import1.HostConfiguration {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	hc := *import1.NewHostConfiguration()
	if v, ok := val["hostname"]; ok && v.(string) != "" {
		hc.Hostname = utils.StringPtr(v.(string))
	}
	if v, ok := val["nameservers"]; ok {
		hc.Nameservers = expandIPAddressList(v.([]interface{}))
	}
	if v, ok := val["ntpservers"]; ok {
		hc.Ntpservers = expandIPAddressOrFQDNList(v.([]interface{}))
	}
	if v, ok := val["network_details"]; ok {
		hc.NetworkDetails = expandHostNetworkDetails(v.([]interface{}))
	}
	return &hc
}

func flattenHostConfiguration(hc *import1.HostConfiguration) []map[string]interface{} {
	if hc == nil {
		return nil
	}
	m := map[string]interface{}{
		"hostname":        utils.StringValue(hc.Hostname),
		"nameservers":     flattenIPAddressList(hc.Nameservers),
		"ntpservers":      flattenIPAddressOrFQDNList(hc.Ntpservers),
		"network_details": flattenHostNetworkDetails(hc.NetworkDetails),
	}
	return []map[string]interface{}{m}
}

func expandHostNetworkDetails(in []interface{}) *import1.HostNetworkDetails {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	hnd := *import1.NewHostNetworkDetails()
	if v, ok := val["management"]; ok {
		hnd.Management = expandManagementNetwork(v.([]interface{}))
	}
	return &hnd
}

func flattenHostNetworkDetails(hnd *import1.HostNetworkDetails) []map[string]interface{} {
	if hnd == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"management": flattenManagementNetwork(hnd.Management),
		},
	}
}

func expandManagementNetwork(in []interface{}) *import1.ManagementNetwork {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	mn := *import1.NewManagementNetwork()
	if v, ok := val["ip"]; ok {
		mn.Ip = expandSingleIPAddress(v.([]interface{}))
	}
	if v, ok := val["gateway"]; ok {
		mn.Gateway = expandSingleIPAddress(v.([]interface{}))
	}
	if v, ok := val["vlan_id"]; ok && v.(int) != 0 {
		mn.VlanId = utils.IntPtr(v.(int))
	}
	if v, ok := val["mtu_bytes"]; ok && v.(int) != 0 {
		mn.MtuBytes = utils.IntPtr(v.(int))
	}
	return &mn
}

func flattenManagementNetwork(mn *import1.ManagementNetwork) []map[string]interface{} {
	if mn == nil {
		return nil
	}
	m := map[string]interface{}{
		"ip":      flattenIPAddress(mn.Ip),
		"gateway": flattenIPAddress(mn.Gateway),
	}
	if mn.VlanId != nil {
		m["vlan_id"] = utils.IntValue(mn.VlanId)
	}
	if mn.MtuBytes != nil {
		m["mtu_bytes"] = utils.IntValue(mn.MtuBytes)
	}
	return []map[string]interface{}{m}
}

func expandIPAddressList(in []interface{}) []commonConfig.IPAddress {
	if len(in) == 0 {
		return nil
	}
	out := make([]commonConfig.IPAddress, 0, len(in))
	for _, raw := range in {
		ip := expandSingleIPAddress([]interface{}{raw})
		if ip != nil {
			out = append(out, *ip)
		}
	}
	return out
}

func flattenIPAddressList(in []commonConfig.IPAddress) []map[string]interface{} {
	if len(in) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(in))
	for i := range in {
		flat := flattenIPAddress(&in[i])
		if len(flat) > 0 {
			out = append(out, flat[0])
		}
	}
	return out
}

func expandSingleIPAddress(in []interface{}) *commonConfig.IPAddress {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	ip := commonConfig.IPAddress{}
	if v, ok := val["ipv4"]; ok {
		ip.Ipv4 = expandIPv4(v.([]interface{}))
	}
	if v, ok := val["ipv6"]; ok {
		ip.Ipv6 = expandIPv6(v.([]interface{}))
	}
	if ip.Ipv4 == nil && ip.Ipv6 == nil {
		return nil
	}
	return &ip
}

func expandIPv4(in []interface{}) *commonConfig.IPv4Address {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	addr := commonConfig.IPv4Address{}
	if v, ok := val["value"]; ok && v.(string) != "" {
		addr.Value = utils.StringPtr(v.(string))
	} else {
		return nil
	}
	if v, ok := val["prefix_length"]; ok && v.(int) != 0 {
		addr.PrefixLength = utils.IntPtr(v.(int))
	}
	return &addr
}

func expandIPv6(in []interface{}) *commonConfig.IPv6Address {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	addr := commonConfig.IPv6Address{}
	if v, ok := val["value"]; ok && v.(string) != "" {
		addr.Value = utils.StringPtr(v.(string))
	} else {
		return nil
	}
	if v, ok := val["prefix_length"]; ok && v.(int) != 0 {
		addr.PrefixLength = utils.IntPtr(v.(int))
	}
	return &addr
}

func expandIPAddressOrFQDNList(in []interface{}) []commonConfig.IPAddressOrFQDN {
	if len(in) == 0 {
		return nil
	}
	out := make([]commonConfig.IPAddressOrFQDN, 0, len(in))
	for _, raw := range in {
		if raw == nil {
			continue
		}
		val := raw.(map[string]interface{})
		item := commonConfig.IPAddressOrFQDN{}
		if v, ok := val["ipv4"]; ok {
			item.Ipv4 = expandIPv4(v.([]interface{}))
		}
		if v, ok := val["ipv6"]; ok {
			item.Ipv6 = expandIPv6(v.([]interface{}))
		}
		if v, ok := val["fqdn"]; ok {
			item.Fqdn = expandFQDN(v.([]interface{}))
		}
		if item.Ipv4 == nil && item.Ipv6 == nil && item.Fqdn == nil {
			continue
		}
		out = append(out, item)
	}
	return out
}

func flattenIPAddressOrFQDNList(in []commonConfig.IPAddressOrFQDN) []map[string]interface{} {
	if len(in) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(in))
	for i := range in {
		item := in[i]
		m := map[string]interface{}{}
		if item.Ipv4 != nil {
			m["ipv4"] = []map[string]interface{}{ipv4Map(item.Ipv4)}
		}
		if item.Ipv6 != nil {
			m["ipv6"] = []map[string]interface{}{ipv6Map(item.Ipv6)}
		}
		if item.Fqdn != nil {
			m["fqdn"] = []map[string]interface{}{{"value": utils.StringValue(item.Fqdn.Value)}}
		}
		out = append(out, m)
	}
	return out
}

func expandFQDN(in []interface{}) *commonConfig.FQDN {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	val := in[0].(map[string]interface{})
	if v, ok := val["value"]; ok && v.(string) != "" {
		return &commonConfig.FQDN{Value: utils.StringPtr(v.(string))}
	}
	return nil
}

func ipv4Map(addr *commonConfig.IPv4Address) map[string]interface{} {
	m := map[string]interface{}{"value": utils.StringValue(addr.Value)}
	if addr.PrefixLength != nil {
		m["prefix_length"] = utils.IntValue(addr.PrefixLength)
	}
	return m
}

func ipv6Map(addr *commonConfig.IPv6Address) map[string]interface{} {
	m := map[string]interface{}{"value": utils.StringValue(addr.Value)}
	if addr.PrefixLength != nil {
		m["prefix_length"] = utils.IntValue(addr.PrefixLength)
	}
	return m
}
