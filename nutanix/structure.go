package nutanix

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
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

			if v.IsConnected != nil {
				nic["is_connected"] = strconv.FormatBool(utils.BoolValue(v.IsConnected))
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
				if utils.StringValue(v1.Type) != "LEARNED" {
					ipEndpoint := make(map[string]interface{})
					ipEndpoint["ip"] = utils.StringValue(v1.IP)
					ipEndpoint["type"] = utils.StringValue(v1.Type)
					ipEndpointList = append(ipEndpointList, ipEndpoint)
				}
			}
			nic["ip_endpoint_list"] = ipEndpointList
			nic["network_function_chain_reference"] = flattenReferenceValues(v.NetworkFunctionChainReference)

			if v.SubnetReference != nil {
				nic["subnet_uuid"] = utils.StringValue(v.SubnetReference.UUID)
				nic["subnet_name"] = utils.StringValue(v.SubnetReference.Name)
			}

			if v.IsConnected != nil {
				nic["is_connected"] = strconv.FormatBool(utils.BoolValue(v.IsConnected))
			}

			nicLists[k] = nic
		}
	}

	return nicLists
}

func flattenDiskList(disks []*v3.VMDisk, diskFromState []*v3.VMDisk) []map[string]interface{} {
	diskList := make([]map[string]interface{}, 0)
	sortedDisks := sortDiskList(disks, diskFromState)
	for _, v := range sortedDisks {
		if v == nil {
			continue
		}
		var deviceProps []map[string]interface{}
		var storageConfig []map[string]interface{}
		if v.DeviceProperties != nil {
			deviceProps = make([]map[string]interface{}, 1)
			index := fmt.Sprintf("%d", utils.Int64Value(v.DeviceProperties.DiskAddress.DeviceIndex))
			adapter := *v.DeviceProperties.DiskAddress.AdapterType
			deviceType := *v.DeviceProperties.DeviceType

			if deviceType == "CDROM" && index == "3" && adapter == IDE {
				continue
			}

			deviceProps[0] = map[string]interface{}{
				"device_type": v.DeviceProperties.DeviceType,
				"disk_address": map[string]interface{}{
					"device_index": index,
					"adapter_type": adapter,
				},
			}
		}

		if v.StorageConfig != nil {
			storageConfig = append(storageConfig, map[string]interface{}{
				"flash_mode": cast.ToString(v.StorageConfig.FlashMode),
				"storage_container_reference": []map[string]interface{}{
					{
						"url":  cast.ToString(v.StorageConfig.StorageContainerReference.URL),
						"kind": cast.ToString(v.StorageConfig.StorageContainerReference.Kind),
						"name": cast.ToString(v.StorageConfig.StorageContainerReference.Name),
						"uuid": cast.ToString(v.StorageConfig.StorageContainerReference.UUID),
					},
				},
			})
		}

		diskList = append(diskList, map[string]interface{}{
			"uuid":                   utils.StringValue(v.UUID),
			"disk_size_bytes":        utils.Int64Value(v.DiskSizeBytes),
			"disk_size_mib":          utils.Int64Value(v.DiskSizeMib),
			"device_properties":      deviceProps,
			"storage_config":         storageConfig,
			"data_source_reference":  flattenReferenceValues(v.DataSourceReference),
			"volume_group_reference": flattenReferenceValues(v.VolumeGroupReference),
		})
	}
	return diskList
}

func sortDiskList(disks []*v3.VMDisk, disksFromState []*v3.VMDisk) []*v3.VMDisk {
	resList := make([]*v3.VMDisk, len(disks)) //make sure the length is according the API VM diskList
	//offset := 0
	notFoundList := make([]*v3.VMDisk, 0)
	for _, disk := range disks {
		if disk.DeviceProperties != nil {
			deviceType := disk.DeviceProperties.DeviceType
			if disk.DeviceProperties.DiskAddress != nil {
				if *deviceType == "CDROM" && *disk.DeviceProperties.DiskAddress.DeviceIndex == int64(CDROOMIndex) && *disk.DeviceProperties.DiskAddress.AdapterType == IDE {
					continue
				}
			}
		}
		index := findIndexDiskList(disk, disksFromState)

		if index > -1 && index < len(disks) {
			resList[index] = disk
		}
		if index == -1 {
			notFoundList = append(notFoundList, disk)
		}
	}
	resList = append(resList, notFoundList...)
	return resList
}

func findIndexDiskList(disk *v3.VMDisk, disks []*v3.VMDisk) int {
	return Search(len(disks), func(i int) bool {
		return reflect.DeepEqual(disks[i].DeviceProperties, disk.DeviceProperties)
	})
}

func Search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}

func flattenSerialPortList(serialPorts []*v3.VMSerialPort) []map[string]interface{} {
	serialPortList := make([]map[string]interface{}, 0)
	if serialPorts != nil {
		serialPortList = make([]map[string]interface{}, len(serialPorts))
		for k, v := range serialPorts {
			serialPort := make(map[string]interface{})
			serialPort["index"] = utils.Int64Value(v.Index)
			serialPort["is_connected"] = utils.BoolValue(v.IsConnected)
			serialPortList[k] = serialPort
		}
	}
	return serialPortList
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

func flattenNutanixGuestTools(d *schema.ResourceData, guest *v3.GuestToolsStatus) error {
	nutanixGuestTools := make(map[string]interface{})
	ngtCredentials := make(map[string]string)
	ngtEnabledCapabilityList := make([]string, 0)

	if guest != nil && guest.NutanixGuestTools != nil {
		tools := guest.NutanixGuestTools
		ngtCredentials = tools.Credentials
		ngtEnabledCapabilityList = utils.StringValueSlice(tools.EnabledCapabilityList)

		nutanixGuestTools["available_version"] = utils.StringValue(tools.AvailableVersion)
		nutanixGuestTools["iso_mount_state"] = utils.StringValue(tools.IsoMountState)
		nutanixGuestTools["ngt_state"] = utils.StringValue(tools.NgtState)
		nutanixGuestTools["state"] = utils.StringValue(tools.State)
		nutanixGuestTools["version"] = utils.StringValue(tools.Version)
		nutanixGuestTools["guest_os_version"] = utils.StringValue(tools.GuestOsVersion)
		nutanixGuestTools["vss_snapshot_capable"] = strconv.FormatBool(utils.BoolValue(tools.VSSSnapshotCapable))
		nutanixGuestTools["is_reachable"] = strconv.FormatBool(utils.BoolValue(tools.IsReachable))
		nutanixGuestTools["vm_mobility_drivers_installed"] = strconv.FormatBool(utils.BoolValue(tools.VMMobilityDriversInstalled))
	}

	if err := d.Set("ngt_enabled_capability_list", ngtEnabledCapabilityList); err != nil {
		return err
	}

	if err := d.Set("ngt_credentials", ngtCredentials); err != nil {
		return err
	}

	if err := d.Set("nutanix_guest_tools", nutanixGuestTools); err != nil {
		return err
	}
	return nil
}
