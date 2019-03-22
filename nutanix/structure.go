package nutanix

import (
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	// CDROM ...
	CDROM = "CDROM"
)

func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, utils.StringPtr(v.(string)))
		}
	}
	return vs
}

func flattenNicListStatus(nics []*v3.VMNicOutputStatus) []map[string]interface{} {
	nicLists := make([]map[string]interface{}, 0)
	if nics != nil {
		nicLists = make([]map[string]interface{}, len(nics))
		for k, v := range nics {
			nic := make(map[string]interface{})
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["floating_ip"] = utils.StringValue(v.FloatingIP)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)
			var ipEndpointList []map[string]interface{}
			for _, v1 := range v.IPEndpointList {
				ipEndpoint := make(map[string]interface{})
				ipEndpoint["ip"] = utils.StringValue(v1.IP)
				ipEndpoint["type"] = utils.StringValue(v1.Type)
				ipEndpointList = append(ipEndpointList, ipEndpoint)
			}
			nic["ip_endpoint_list"] = ipEndpointList
			nic["network_function_chain_reference"] = flattenReferenceValues(v.NetworkFunctionChainReference)

			if v.SubnetReference != nil {
				nic["subnet_uuid"] = utils.StringValue(v.SubnetReference.UUID)
				nic["subnet_name"] = utils.StringValue(v.SubnetReference.Name)

			}

			nicLists[k] = nic
		}
	}

	return nicLists
}

func flattenNicList(nics []*v3.VMNic) []map[string]interface{} {
	nicLists := make([]map[string]interface{}, 0)
	if nics != nil {
		nicLists = make([]map[string]interface{}, len(nics))
		for k, v := range nics {
			nic := make(map[string]interface{})
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)
			var ipEndpointList []map[string]interface{}
			for _, v1 := range v.IPEndpointList {
				ipEndpoint := make(map[string]interface{})
				ipEndpoint["ip"] = utils.StringValue(v1.IP)
				ipEndpoint["type"] = utils.StringValue(v1.Type)
				ipEndpointList = append(ipEndpointList, ipEndpoint)
			}
			nic["ip_endpoint_list"] = ipEndpointList
			nic["network_function_chain_reference"] = flattenReferenceValues(v.NetworkFunctionChainReference)

			if v.SubnetReference != nil {
				nic["subnet_uuid"] = utils.StringValue(v.SubnetReference.UUID)
				nic["subnet_name"] = utils.StringValue(v.SubnetReference.Name)

			}

			nicLists[k] = nic
		}
	}

	return nicLists
}

func flattenGPUList(gpu []*v3.VMGpuOutputStatus) []map[string]interface{} {
	gpuList := make([]map[string]interface{}, 0)
	if gpu != nil {
		gpuList = make([]map[string]interface{}, len(gpu))
		for k, v := range gpu {
			gpu := make(map[string]interface{})
			gpu["frame_buffer_size_mib"] = utils.Int64Value(v.FrameBufferSizeMib)
			gpu["vendor"] = utils.StringValue(v.Vendor)
			gpu["uuid"] = utils.StringValue(v.UUID)
			gpu["name"] = utils.StringValue(v.Name)
			gpu["pci_address"] = utils.StringValue(v.PCIAddress)
			gpu["fraction"] = utils.Int64Value(v.Fraction)
			gpu["mode"] = utils.StringValue(v.Mode)
			gpu["num_virtual_display_heads"] = utils.Int64Value(v.NumVirtualDisplayHeads)
			gpu["guest_driver_version"] = utils.StringValue(v.GuestDriverVersion)
			gpu["device_id"] = utils.Int64Value(v.DeviceID)
			gpuList[k] = gpu
		}
	}
	return gpuList
}

func setDiskList(disk []*v3.VMDisk, hasCloudInit *v3.GuestCustomizationStatus) []map[string]interface{} {
	var diskList []map[string]interface{}
	if len(disk) > 0 {
		for _, v1 := range disk {

			if hasCloudInit != nil {
				if hasCloudInit.CloudInit != nil && utils.StringValue(v1.DeviceProperties.DeviceType) == CDROM {
					continue
				}
			}

			disk := make(map[string]interface{})
			disk["uuid"] = utils.StringValue(v1.UUID)
			disk["disk_size_bytes"] = utils.Int64Value(v1.DiskSizeBytes)
			disk["disk_size_mib"] = utils.Int64Value(v1.DiskSizeMib)
			if v1.DataSourceReference != nil {
				disk["data_source_reference"] = flattenReferenceValues(v1.DataSourceReference)
			}

			if v1.VolumeGroupReference != nil {
				disk["volume_group_reference"] = flattenReferenceValues(v1.VolumeGroupReference)
			}

			dp := make([]map[string]interface{}, 1)
			deviceProps := make(map[string]interface{})
			deviceProps["device_type"] = utils.StringValue(v1.DeviceProperties.DeviceType)
			dp[0] = deviceProps

			da := make([]map[string]interface{}, 1)
			diskAddress := make(map[string]interface{})
			if v1.DeviceProperties.DiskAddress != nil {
				diskAddress["device_index"] = utils.Int64Value(v1.DeviceProperties.DiskAddress.DeviceIndex)
				diskAddress["adapter_type"] = utils.StringValue(v1.DeviceProperties.DiskAddress.AdapterType)
			}
			da[0] = diskAddress
			deviceProps["disk_address"] = da

			disk["device_properties"] = dp

			diskList = append(diskList, disk)
		}
	}

	if diskList == nil {
		return make([]map[string]interface{}, 0)
	}

	return diskList
}

func setNutanixGuestTools(guest *v3.GuestToolsStatus) map[string]interface{} {
	nutanixGuestTools := make(map[string]interface{})
	if guest != nil {
		tools := guest.NutanixGuestTools
		nutanixGuestTools["available_version"] = utils.StringValue(tools.AvailableVersion)
		nutanixGuestTools["iso_mount_state"] = utils.StringValue(tools.IsoMountState)
		nutanixGuestTools["state"] = utils.StringValue(tools.State)
		nutanixGuestTools["version"] = utils.StringValue(tools.Version)
		nutanixGuestTools["guest_os_version"] = utils.StringValue(tools.GuestOsVersion)
		nutanixGuestTools["enabled_capability_list"] = utils.StringValueSlice(tools.EnabledCapabilityList)
		nutanixGuestTools["vss_snapshot_capable"] = utils.BoolValue(tools.VSSSnapshotCapable)
		nutanixGuestTools["is_reachable"] = utils.BoolValue(tools.IsReachable)
		nutanixGuestTools["vm_mobility_drivers_installed"] = utils.BoolValue(tools.VMMobilityDriversInstalled)
	}
	return nutanixGuestTools
}
