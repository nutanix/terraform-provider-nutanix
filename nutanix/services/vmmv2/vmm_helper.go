package vmmv2

import (
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Expand helper functions for VMM v2
// expandVirtualEthernetNic expands either:
// - nic_backing_info.virtual_ethernet_nic (VirtualEthernetNic), OR
// - nic_network_info.virtual_ethernet_nic_network_info (VirtualEthernetNicNetworkInfo).
//
// It returns the concrete (non-pointer) SDK object to satisfy OneOf setters, or nil if empty.
func expandVirtualEthernetNic(pr interface{}) interface{} {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}

	// Backing object uses model/mac/is_connected/num_queues keys.
	if _, hasModel := val["model"]; hasModel || val["mac_address"] != nil || val["num_queues"] != nil || val["is_connected"] != nil {
		ven := config.NewVirtualEthernetNic()

		if model, ok := val["model"]; ok && model != nil && model.(string) != "" {
			ven.Model = common.ExpandEnum[config.VirtualEthernetNicModel](model.(string))
		}
		if mac, ok := val["mac_address"]; ok && mac != nil && mac.(string) != "" {
			ven.MacAddress = utils.StringPtr(mac.(string))
		}
		if isConn, ok := val["is_connected"]; ok && isConn != nil {
			ven.IsConnected = utils.BoolPtr(isConn.(bool))
		}
		if nq, ok := val["num_queues"]; ok && nq != nil {
			ven.NumQueues = utils.IntPtr(nq.(int))
		}
		return *ven
	}

	// Otherwise treat as network info object.
	venNI := config.NewVirtualEthernetNicNetworkInfo()

	// nic_type (optional; defaults to NORMAL_NIC server-side, but we infer to avoid "missing NIC" diffs)
	if nicType, ok := val["nic_type"]; ok && nicType != nil && nicType.(string) != "" {
		venNI.NicType = common.ExpandEnum[config.NicType](nicType.(string))
	}
	if ntwkFunc, ok := val["network_function_chain"]; ok && ntwkFunc != nil {
		venNI.NetworkFunctionChain = expandNetworkFunctionChainReference(ntwkFunc)
	}
	if ntwkFuncNicType, ok := val["network_function_nic_type"]; ok && ntwkFuncNicType != nil && ntwkFuncNicType.(string) != "" {
		venNI.NetworkFunctionNicType = common.ExpandEnum[config.NetworkFunctionNicType](ntwkFuncNicType.(string))
	}
	if subnet, ok := val["subnet"]; ok && subnet != nil {
		venNI.Subnet = expandSubnetReference(subnet)
	}
	if vlanMode, ok := val["vlan_mode"]; ok && vlanMode != nil && vlanMode.(string) != "" {
		venNI.VlanMode = common.ExpandEnum[config.VlanMode](vlanMode.(string))
	}
	if trunkedVlans, ok := val["trunked_vlans"]; ok && trunkedVlans != nil {
		vlansList := trunkedVlans.([]interface{})
		trunkInts := make([]int, len(vlansList))
		for i, v := range vlansList {
			trunkInts[i] = v.(int)
		}
		venNI.TrunkedVlans = trunkInts
	}
	if unknownMacs, ok := val["should_allow_unknown_macs"]; ok && unknownMacs != nil {
		venNI.ShouldAllowUnknownMacs = utils.BoolPtr(unknownMacs.(bool))
	}
	if ipv4, ok := val["ipv4_config"]; ok && ipv4 != nil {
		venNI.Ipv4Config = expandIpv4Config(ipv4)
	}
	if ipv4Info, ok := val["ipv4_info"]; ok && ipv4Info != nil {
		venNI.Ipv4Info = expandIPv4Info(ipv4Info)
	}

	// Infer nic_type if absent (helps avoid perpetual "nics { ... }" diffs when config sets subnet/ipv4 only).
	if venNI.NicType == nil {
		if venNI.NetworkFunctionNicType != nil || venNI.NetworkFunctionChain != nil {
			p := config.NICTYPE_NETWORK_FUNCTION_NIC
			venNI.NicType = p.Ref()
		} else if venNI.Subnet != nil || venNI.Ipv4Config != nil || venNI.VlanMode != nil || len(venNI.TrunkedVlans) > 0 || venNI.ShouldAllowUnknownMacs != nil {
			p := config.NICTYPE_NORMAL_NIC
			venNI.NicType = p.Ref()
		}
	}

	return *venNI
}

// expandSriovNic expands either:
// - nic_backing_info.sriov_nic (SriovNic), OR
// - nic_network_info.sriov_nic_network_info (SriovNicNetworkInfo).
func expandSriovNic(pr interface{}) interface{} {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}

	// Network info is only vlan_id.
	if _, hasVlanID := val["vlan_id"]; hasVlanID {
		ni := config.NewSriovNicNetworkInfo()
		if vlanID, ok := val["vlan_id"]; ok && vlanID != nil {
			ni.VlanId = utils.IntPtr(vlanID.(int))
		}
		return *ni
	}

	// Backing info.
	sriov := config.NewSriovNic()
	if mac, ok := val["mac_address"]; ok && mac != nil && mac.(string) != "" {
		sriov.MacAddress = utils.StringPtr(mac.(string))
	}
	if isConn, ok := val["is_connected"]; ok && isConn != nil {
		sriov.IsConnected = utils.BoolPtr(isConn.(bool))
	}
	if hostRef, ok := val["host_pcie_device_reference"]; ok && hostRef != nil {
		sriov.HostPcieDeviceReference = expandHostPcieDeviceReference(hostRef)
	}
	if prof, ok := val["sriov_profile_reference"]; ok && prof != nil {
		sriov.SriovProfileReference = expandNicProfileReference(prof)
	}
	return *sriov
}

// expandDpOffloadNic expands either:
// - nic_backing_info.dp_offload_nic (DpOffloadNic), OR
// - nic_network_info.dp_offload_nic_network_info (DpOffloadNicNetworkInfo).
func expandDpOffloadNic(pr interface{}) interface{} {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}

	// Backing info contains dp_offload_profile_reference / host_pcie_device_reference / mac/is_connected.
	if _, hasProfile := val["dp_offload_profile_reference"]; hasProfile || val["host_pcie_device_reference"] != nil || val["mac_address"] != nil || val["is_connected"] != nil {
		dp := config.NewDpOffloadNic()
		if mac, ok := val["mac_address"]; ok && mac != nil && mac.(string) != "" {
			dp.MacAddress = utils.StringPtr(mac.(string))
		}
		if isConn, ok := val["is_connected"]; ok && isConn != nil {
			dp.IsConnected = utils.BoolPtr(isConn.(bool))
		}
		if hostRef, ok := val["host_pcie_device_reference"]; ok && hostRef != nil {
			dp.HostPcieDeviceReference = expandHostPcieDeviceReference(hostRef)
		}
		if prof, ok := val["dp_offload_profile_reference"]; ok && prof != nil {
			dp.DpOffloadProfileReference = expandNicProfileReference(prof)
		}
		return *dp
	}

	// Otherwise network info.
	ni := config.NewDpOffloadNicNetworkInfo()
	if subnet, ok := val["subnet"]; ok && subnet != nil {
		ni.Subnet = expandSubnetReference(subnet)
	}
	if vlanMode, ok := val["vlan_mode"]; ok && vlanMode != nil && vlanMode.(string) != "" {
		ni.VlanMode = common.ExpandEnum[config.VlanMode](vlanMode.(string))
	}
	if trunkedVlans, ok := val["trunked_vlans"]; ok && trunkedVlans != nil {
		vlansList := trunkedVlans.([]interface{})
		trunkInts := make([]int, len(vlansList))
		for i, v := range vlansList {
			trunkInts[i] = v.(int)
		}
		ni.TrunkedVlans = trunkInts
	}
	if unknownMacs, ok := val["should_allow_unknown_macs"]; ok && unknownMacs != nil {
		ni.ShouldAllowUnknownMacs = utils.BoolPtr(unknownMacs.(bool))
	}
	if ipv4, ok := val["ipv4_config"]; ok && ipv4 != nil {
		ni.Ipv4Config = expandIpv4Config(ipv4)
	}
	if ipv4Info, ok := val["ipv4_info"]; ok && ipv4Info != nil {
		ni.Ipv4Info = expandIPv4Info(ipv4Info)
	}
	return *ni
}

// Flatten helper functions for VMM v2

func flattenVirtualEthernetNicAsBackingInfo(pr *config.VirtualEthernetNic) []map[string]interface{} {
	if pr == nil {
		return nil
	}

	nic := make(map[string]interface{})

	if pr.Model != nil {
		nic["model"] = flattenVirtualEthernetNicModel(pr.Model)
	}
	if pr.MacAddress != nil {
		nic["mac_address"] = pr.MacAddress
	}
	if pr.IsConnected != nil {
		nic["is_connected"] = pr.IsConnected
	}
	if pr.NumQueues != nil {
		nic["num_queues"] = pr.NumQueues
	}

	return []map[string]interface{}{nic}
}

func flattenVirtualEthernetNicModel(pr *config.VirtualEthernetNicModel) string {
	if pr != nil {
		if *pr == config.VIRTUALETHERNETNICMODEL_VIRTIO {
			return "VIRTIO"
		}
		if *pr == config.VIRTUALETHERNETNICMODEL_E1000 {
			return "E1000"
		}
	}
	return "UNKNOWN"
}

func flattenVirtualEthernetNicNetworkInfo(pr *config.VirtualEthernetNicNetworkInfo) []map[string]interface{} {
	if pr == nil {
		return nil
	}

	nic := make(map[string]interface{})

	if pr.NicType != nil {
		nic["nic_type"] = flattenNicType(pr.NicType)
	}
	if pr.NetworkFunctionChain != nil {
		nic["network_function_chain"] = flattenNetworkFunctionChainReference(pr.NetworkFunctionChain)
	}
	if pr.NetworkFunctionNicType != nil {
		nic["network_function_nic_type"] = flattenNetworkFunctionNicType(pr.NetworkFunctionNicType)
	}
	if pr.Subnet != nil {
		nic["subnet"] = flattenSubnetReference(pr.Subnet)
	}
	if pr.VlanMode != nil {
		nic["vlan_mode"] = flattenVlanMode(pr.VlanMode)
	}
	if pr.TrunkedVlans != nil {
		nic["trunked_vlans"] = pr.TrunkedVlans
	}
	if pr.ShouldAllowUnknownMacs != nil {
		nic["should_allow_unknown_macs"] = pr.ShouldAllowUnknownMacs
	}
	if pr.Ipv4Config != nil {
		nic["ipv4_config"] = flattenIpv4Config(pr.Ipv4Config)
	}
	if pr.Ipv4Info != nil {
		nic["ipv4_info"] = flattenIpv4Info(pr.Ipv4Info)
	}

	return []map[string]interface{}{nic}
}

func flattenSriovNicAsBackingInfo(pr *config.SriovNic) []map[string]interface{} {
	if pr == nil {
		return nil
	}

	nic := make(map[string]interface{})

	if pr.SriovProfileReference != nil {
		nic["sriov_profile_reference"] = flattenNicProfileReference(pr.SriovProfileReference)
	}
	if pr.HostPcieDeviceReference != nil {
		nic["host_pcie_device_reference"] = flattenHostPcieDeviceReference(pr.HostPcieDeviceReference)
	}
	if pr.IsConnected != nil {
		nic["is_connected"] = pr.IsConnected
	}
	if pr.MacAddress != nil {
		nic["mac_address"] = pr.MacAddress
	}

	return []map[string]interface{}{nic}
}

func flattenDpOffloadNicAsBackingInfo(pr *config.DpOffloadNic) []map[string]interface{} {
	if pr == nil {
		return nil
	}

	nic := make(map[string]interface{})

	if pr.DpOffloadProfileReference != nil {
		nic["dp_offload_profile_reference"] = flattenNicProfileReference(pr.DpOffloadProfileReference)
	}
	if pr.HostPcieDeviceReference != nil {
		nic["host_pcie_device_reference"] = flattenHostPcieDeviceReference(pr.HostPcieDeviceReference)
	}
	if pr.IsConnected != nil {
		nic["is_connected"] = pr.IsConnected
	}
	if pr.MacAddress != nil {
		nic["mac_address"] = pr.MacAddress
	}

	return []map[string]interface{}{nic}
}

func flattenSriovNicNetworkInfo(pr *config.SriovNicNetworkInfo) []map[string]interface{} {
	if pr == nil {
		return nil
	}
	nic := make(map[string]interface{})
	if pr.VlanId != nil {
		nic["vlan_id"] = pr.VlanId
	}
	return []map[string]interface{}{nic}
}

func flattenDpOffloadNicNetworkInfo(pr *config.DpOffloadNicNetworkInfo) []map[string]interface{} {
	if pr == nil {
		return nil
	}

	nic := make(map[string]interface{})

	if pr.Subnet != nil {
		nic["subnet"] = flattenSubnetReference(pr.Subnet)
	}
	if pr.VlanMode != nil {
		nic["vlan_mode"] = flattenVlanMode(pr.VlanMode)
	}
	if pr.TrunkedVlans != nil {
		nic["trunked_vlans"] = pr.TrunkedVlans
	}
	if pr.ShouldAllowUnknownMacs != nil {
		nic["should_allow_unknown_macs"] = pr.ShouldAllowUnknownMacs
	}
	if pr.Ipv4Config != nil {
		nic["ipv4_config"] = flattenIpv4Config(pr.Ipv4Config)
	}
	if pr.Ipv4Info != nil {
		nic["ipv4_info"] = flattenIpv4Info(pr.Ipv4Info)
	}

	return []map[string]interface{}{nic}
}

func flattenHostPcieDeviceReference(ref *config.HostPcieDeviceReference) []map[string]interface{} {
	if ref == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ext_id": utils.StringValue(ref.ExtId),
		},
	}
}

func flattenNicProfileReference(ref *config.NicProfileReference) []map[string]interface{} {
	if ref == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ext_id": utils.StringValue(ref.ExtId),
		},
	}
}
