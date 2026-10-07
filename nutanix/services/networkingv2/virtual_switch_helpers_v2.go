package networkingv2

import (
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/config"
	networkingConfig "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func expandVirtualSwitchClusters(clusters []interface{}) []networkingConfig.Cluster {
	if len(clusters) == 0 {
		return nil
	}
	result := make([]networkingConfig.Cluster, len(clusters))
	for i, c := range clusters {
		clusterMap := c.(map[string]interface{})
		cluster := networkingConfig.Cluster{}

		if extID, ok := clusterMap["ext_id"]; ok && extID.(string) != "" {
			cluster.ExtId = utils.StringPtr(extID.(string))
		}
		if gwIP, ok := clusterMap["gateway_ip_address"]; ok {
			gwList := gwIP.([]interface{})
			if len(gwList) > 0 && gwList[0] != nil {
				cluster.GatewayIpAddress = expandIPv4AddressVS(gwList[0].(map[string]interface{}))
			}
		}
		if hosts, ok := clusterMap["hosts"]; ok {
			cluster.Hosts = expandVirtualSwitchHosts(hosts.([]interface{}))
		}
		if vlanID, ok := clusterMap["vlan_identifier"]; ok {
			cluster.VlanIdentifier = utils.IntPtr(vlanID.(int))
		}

		result[i] = cluster
	}
	return result
}

func expandVirtualSwitchHosts(hosts []interface{}) []networkingConfig.Host {
	if len(hosts) == 0 {
		return nil
	}
	result := make([]networkingConfig.Host, len(hosts))
	for i, h := range hosts {
		hostMap := h.(map[string]interface{})
		host := networkingConfig.Host{}

		if extID, ok := hostMap["ext_id"]; ok && extID.(string) != "" {
			host.ExtId = utils.StringPtr(extID.(string))
		}
		if hostNics, ok := hostMap["host_nics"]; ok {
			host.HostNics = expandStringList(hostNics.([]interface{}))
		}
		if bridgeName, ok := hostMap["internal_bridge_name"]; ok && bridgeName.(string) != "" {
			host.InternalBridgeName = utils.StringPtr(bridgeName.(string))
		}
		if ipAddr, ok := hostMap["ip_address"]; ok {
			ipList := ipAddr.([]interface{})
			if len(ipList) > 0 && ipList[0] != nil {
				host.IpAddress = expandIPv4SubnetVS(ipList[0].(map[string]interface{}))
			}
		}
		if activeUplink, ok := hostMap["active_uplink"]; ok && activeUplink.(string) != "" {
			host.ActiveUplink = utils.StringPtr(activeUplink.(string))
		}
		if routeTable, ok := hostMap["route_table"]; ok {
			host.RouteTable = utils.IntPtr(routeTable.(int))
		}

		result[i] = host
	}
	return result
}

func expandIPv4AddressVS(m map[string]interface{}) *commonConfig.IPv4Address {
	addr := &commonConfig.IPv4Address{}
	if value, ok := m["value"]; ok && value.(string) != "" {
		addr.Value = utils.StringPtr(value.(string))
	}
	if prefixLen, ok := m["prefix_length"]; ok {
		addr.PrefixLength = utils.IntPtr(prefixLen.(int))
	}
	return addr
}

func expandIPv4SubnetVS(m map[string]interface{}) *networkingConfig.IPv4Subnet {
	subnet := &networkingConfig.IPv4Subnet{}
	if ip, ok := m["ip"]; ok {
		ipList := ip.([]interface{})
		if len(ipList) > 0 && ipList[0] != nil {
			subnet.Ip = expandIPv4AddressVS(ipList[0].(map[string]interface{}))
		}
	}
	if prefixLen, ok := m["prefix_length"]; ok {
		subnet.PrefixLength = utils.IntPtr(prefixLen.(int))
	}
	return subnet
}

func expandIgmpSpec(igmpSpec []interface{}) *networkingConfig.IgmpSpec {
	if len(igmpSpec) == 0 || igmpSpec[0] == nil {
		return nil
	}
	specMap := igmpSpec[0].(map[string]interface{})
	spec := &networkingConfig.IgmpSpec{}

	if isSnooping, ok := specMap["is_snooping_enabled"]; ok {
		spec.IsSnoopingEnabled = utils.BoolPtr(isSnooping.(bool))
	}
	if snoopingTimeout, ok := specMap["snooping_timeout"]; ok {
		spec.SnoopingTimeout = utils.Int64Ptr(int64(snoopingTimeout.(int)))
	}
	if querierSpec, ok := specMap["querier_spec"]; ok {
		qList := querierSpec.([]interface{})
		if len(qList) > 0 && qList[0] != nil {
			spec.QuerierSpec = expandQuerierSpec(qList[0].(map[string]interface{}))
		}
	}
	return spec
}

func expandQuerierSpec(m map[string]interface{}) *networkingConfig.QuerierSpec {
	spec := &networkingConfig.QuerierSpec{}
	if isEnabled, ok := m["is_querier_enabled"]; ok {
		spec.IsQuerierEnabled = utils.BoolPtr(isEnabled.(bool))
	}
	if vlanIDList, ok := m["vlan_id_list"]; ok {
		spec.VlanIdList = common.ExpandListOfInt(vlanIDList.([]interface{}))
	}
	return spec
}

func expandStringList(list []interface{}) []string {
	if len(list) == 0 {
		return nil
	}
	result := make([]string, 0, len(list))
	for _, v := range list {
		if v != nil {
			result = append(result, v.(string))
		}
	}
	return result
}

func flattenVirtualSwitchClusters(clusters []networkingConfig.Cluster) []map[string]interface{} {
	if len(clusters) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(clusters))
	for i, c := range clusters {
		cluster := map[string]interface{}{
			"ext_id":             utils.StringValue(c.ExtId),
			"gateway_ip_address": flattenIPv4AddressVS(c.GatewayIpAddress),
			"hosts":              flattenVirtualSwitchHosts(c.Hosts),
			"vlan_identifier":    utils.IntValue(c.VlanIdentifier),
			// existing_bridge_name is intentionally NOT emitted here. It is a
			// create-time-only input to the migrate endpoint and the API never
			// returns it on read. Two reasons we don't seed it as "" in this
			// shared map:
			//   1. The data source schemas (data_source_nutanix_virtual_switch_v2
			//      and the list variant) don't declare the key, so writing it
			//      would trip "Invalid address to set" in the SDK.
			//   2. The resource's ReadContext overlays the prior state value
			//      onto this map after flattening, which is the only place the
			//      key has a meaningful value anyway.
		}
		result[i] = cluster
	}
	return result
}

func flattenVirtualSwitchHosts(hosts []networkingConfig.Host) []map[string]interface{} {
	if len(hosts) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(hosts))
	for i, h := range hosts {
		host := map[string]interface{}{
			"ext_id":               utils.StringValue(h.ExtId),
			"host_nics":            h.HostNics,
			"internal_bridge_name": utils.StringValue(h.InternalBridgeName),
			"ip_address":           flattenIPv4SubnetVS(h.IpAddress),
			"active_uplink":        utils.StringValue(h.ActiveUplink),
			"route_table":          utils.IntValue(h.RouteTable),
		}
		result[i] = host
	}
	return result
}

func flattenIPv4AddressVS(addr *commonConfig.IPv4Address) []map[string]interface{} {
	if addr == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"value":         utils.StringValue(addr.Value),
			"prefix_length": utils.IntValue(addr.PrefixLength),
		},
	}
}

func flattenIPv4SubnetVS(subnet *networkingConfig.IPv4Subnet) []map[string]interface{} {
	if subnet == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ip":            flattenIPv4AddressVS(subnet.Ip),
			"prefix_length": utils.IntValue(subnet.PrefixLength),
		},
	}
}

func flattenIgmpSpec(spec *networkingConfig.IgmpSpec) []map[string]interface{} {
	if spec == nil {
		return nil
	}
	result := map[string]interface{}{
		"is_snooping_enabled": utils.BoolValue(spec.IsSnoopingEnabled),
		"snooping_timeout":    utils.Int64Value(spec.SnoopingTimeout),
		"querier_spec":        flattenQuerierSpec(spec.QuerierSpec),
	}
	return []map[string]interface{}{result}
}

func flattenQuerierSpec(spec *networkingConfig.QuerierSpec) []map[string]interface{} {
	if spec == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"is_querier_enabled": utils.BoolValue(spec.IsQuerierEnabled),
			"vlan_id_list":       spec.VlanIdList,
		},
	}
}
