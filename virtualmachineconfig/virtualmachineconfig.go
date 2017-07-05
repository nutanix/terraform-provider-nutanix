package virtualmachineconfig

import (
	"github.com/hashicorp/terraform/helper/schema"
	vm "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachine"
)

func convertToBool(a interface{}) bool {
	if a != nil {
		return a.(bool)
	}
	return false
}
func convertToInt(a interface{}) int {
	if a != nil {
		return a.(int)
	}
	return 0
}
func convertToString(a interface{}) string {
	if a != nil {
		return a.(string)
	}
	return ""
}


// SetClusterReference sets cluster_reference fields in  json struct
func SetClusterReference(t []interface{}, i int) vm.ClusterReference {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		ClusterReference := vm.ClusterReference{
			Kind:		convertToString(s["kind"]),
			Name:		convertToString(s["name"]),
			UUID:		convertToString(s["uuid"]),
		}
		return ClusterReference
	}
	return vm.ClusterReference{}
}

// SetGPUList sets gpu_list fields in  json struct
func SetGPUList(t []interface{}, i int) vm.GPUList {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		GPUList := vm.GPUList{
			Vendor:		convertToString(s["vendor"]),
			Mode:		convertToString(s["mode"]),
			DeviceID:		convertToInt(s["device_id"]),
		}
		return GPUList
	}
	return vm.GPUList{}
}

// SetCloudInit sets cloud_init fields in  json struct
func SetCloudInit(t []interface{}, i int) vm.CloudInit {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		CloudInit := vm.CloudInit{
			MetaData:		convertToString(s["meta_data"]),
			UserData:		convertToString(s["user_data"]),
		}
		return CloudInit
	}
	return vm.CloudInit{}
}

// SetSysprep sets sysprep fields in  json struct
func SetSysprep(t []interface{}, i int) vm.Sysprep {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		Sysprep := vm.Sysprep{
			UnattendXML:		convertToString(s["unattend_xml"]),
			InstallType:		convertToString(s["install_type"]),
		}
		return Sysprep
	}
	return vm.Sysprep{}
}

// SetGuestCustomization sets guest_customization fields in  json struct
func SetGuestCustomization(t []interface{}, i int) vm.GuestCustomization {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		GuestCustomization := vm.GuestCustomization{
			CloudInit:		SetCloudInit(s["cloud_init"].(*schema.Set).List(), 0),
			Sysprep:		SetSysprep(s["sysprep"].(*schema.Set).List(), 0),
		}
		return GuestCustomization
	}
	return vm.GuestCustomization{}
}

// SetDiskAddress sets disk_address fields in  json struct
func SetDiskAddress(t []interface{}, i int) vm.DiskAddress {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		DiskAddress := vm.DiskAddress{
			DeviceIndex:		convertToInt(s["device_index"]),
			AdapterType:		convertToString(s["adapter_type"]),
		}
		return DiskAddress
	}
	return vm.DiskAddress{}
}

// SetBootConfig sets boot_config fields in  json struct
func SetBootConfig(t []interface{}, i int) vm.BootConfig {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		BootConfig := vm.BootConfig{
			DiskAddress:		SetDiskAddress(s["disk_address"].(*schema.Set).List(), 0),
			MacAddress:		convertToString(s["mac_address"]),
		}
		return BootConfig
	}
	return vm.BootConfig{}
}

// SetSubnetReference sets subnet_reference fields in  json struct
func SetSubnetReference(t []interface{}, i int) vm.SubnetReference {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		SubnetReference := vm.SubnetReference{
			UUID:		convertToString(s["uuid"]),
			Kind:		convertToString(s["kind"]),
			Name:		convertToString(s["name"]),
		}
		return SubnetReference
	}
	return vm.SubnetReference{}
}

// SetIPEndpointList sets ip_endpoint_list fields in  json struct
func SetIPEndpointList(t []interface{}, i int) vm.IPEndpointList {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		IPEndpointList := vm.IPEndpointList{
			Type:		convertToString(s["type"]),
			Address:		convertToString(s["address"]),
		}
		return IPEndpointList
	}
	return vm.IPEndpointList{}
}

// SetNetworkFunctionChainReference sets network_function_chain_reference fields in  json struct
func SetNetworkFunctionChainReference(t []interface{}, i int) vm.NetworkFunctionChainReference {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		NetworkFunctionChainReference := vm.NetworkFunctionChainReference{
			Kind:		convertToString(s["kind"]),
			Name:		convertToString(s["name"]),
			UUID:		convertToString(s["uuid"]),
		}
		return NetworkFunctionChainReference
	}
	return vm.NetworkFunctionChainReference{}
}

// SetNicList sets nic_list fields in  json struct
func SetNicList(t []interface{}, i int) vm.NicList {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		var IPEndpointList []vm.IPEndpointList
		if s["ip_endpoint_list"] != nil {
			for i := 0; i< len(s["ip_endpoint_list"].([]interface{})); i++ {
				elem := SetIPEndpointList(s["ip_endpoint_list"].([]interface{}),	i)
				IPEndpointList = append(IPEndpointList, elem)
			}
		}


		NicList := vm.NicList{
			NicType:		convertToString(s["nic_type"]),
			SubnetReference:		SetSubnetReference(s["subnet_reference"].(*schema.Set).List(), 0),
			NetworkFunctionNicType:		convertToString(s["network_function_nic_type"]),
			MacAddress:		convertToString(s["mac_address"]),
			IPEndpointList:		IPEndpointList,
			NetworkFunctionChainReference:		SetNetworkFunctionChainReference(s["network_function_chain_reference"].(*schema.Set).List(), 0),
		}
		return NicList
	}
	return vm.NicList{}
}

// SetParentReference sets parent_reference fields in  json struct
func SetParentReference(t []interface{}, i int) vm.ParentReference {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		ParentReference := vm.ParentReference{
			Kind:		convertToString(s["kind"]),
			UUID:		convertToString(s["uuid"]),
			Name:		convertToString(s["name"]),
		}
		return ParentReference
	}
	return vm.ParentReference{}
}

// SetDataSourceReference sets data_source_reference fields in  json struct
func SetDataSourceReference(t []interface{}, i int) vm.DataSourceReference {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		DataSourceReference := vm.DataSourceReference{
			Kind:		convertToString(s["kind"]),
			UUID:		convertToString(s["uuid"]),
			Name:		convertToString(s["name"]),
		}
		return DataSourceReference
	}
	return vm.DataSourceReference{}
}

// SetDeviceProperties sets device_properties fields in  json struct
func SetDeviceProperties(t []interface{}, i int) vm.DeviceProperties {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		DeviceProperties := vm.DeviceProperties{
			DeviceType:		convertToString(s["device_type"]),
			DiskAddress:		SetDiskAddress(s["disk_address"].(*schema.Set).List(), 0),
		}
		return DeviceProperties
	}
	return vm.DeviceProperties{}
}

// SetDiskList sets disk_list fields in  json struct
func SetDiskList(t []interface{}, i int) vm.DiskList {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		DiskList := vm.DiskList{
			DataSourceReference:		SetDataSourceReference(s["data_source_reference"].(*schema.Set).List(), 0),
			DeviceProperties:		SetDeviceProperties(s["device_properties"].(*schema.Set).List(), 0),
			UUID:		convertToString(s["uuid"]),
			DiskSizeMib:		convertToInt(s["disk_size_mib"]),
		}
		return DiskList
	}
	return vm.DiskList{}
}

// SetResources sets resources fields in  json struct
func SetResources(t []interface{}, i int) vm.Resources {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		var GPUList []vm.GPUList
		if s["gpu_list"] != nil {
			for i := 0; i< len(s["gpu_list"].([]interface{})); i++ {
				elem := SetGPUList(s["gpu_list"].([]interface{}),	i)
				GPUList = append(GPUList, elem)
			}
		}


		var NicList []vm.NicList
		if s["nic_list"] != nil {
			for i := 0; i< len(s["nic_list"].([]interface{})); i++ {
				elem := SetNicList(s["nic_list"].([]interface{}),	i)
				NicList = append(NicList, elem)
			}
		}


		var DiskList []vm.DiskList
		if s["disk_list"] != nil {
			for i := 0; i< len(s["disk_list"].([]interface{})); i++ {
				elem := SetDiskList(s["disk_list"].([]interface{}),	i)
				DiskList = append(DiskList, elem)
			}
		}


		Resources := vm.Resources{
			PowerState:		convertToString(s["power_state"]),
			NumSockets:		convertToInt(s["num_sockets"]),
			MemorySizeMb:		convertToInt(s["memory_size_mb"]),
			GPUList:		GPUList,
			GuestCustomization:		SetGuestCustomization(s["guest_customization"].(*schema.Set).List(), 0),
			BootConfig:		SetBootConfig(s["boot_config"].(*schema.Set).List(), 0),
			NicList:		NicList,
			NumVcpusPerSocket:		convertToInt(s["num_vcpus_per_socket"]),
			ParentReference:		SetParentReference(s["parent_reference"].(*schema.Set).List(), 0),
			DiskList:		DiskList,
		}
		return Resources
	}
	return vm.Resources{}
}

// SetSpec sets spec fields in  json struct
func SetSpec(t []interface{}, i int) vm.Spec {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		Spec := vm.Spec{
			ClusterReference:		SetClusterReference(s["cluster_reference"].(*schema.Set).List(), 0),
			Resources:		SetResources(s["resources"].(*schema.Set).List(), 0),
			Name:		convertToString(s["name"]),
			Description:		convertToString(s["description"]),
		}
		return Spec
	}
	return vm.Spec{}
}

// SetCategories sets categories fields in  json struct
func SetCategories(s map[string]interface{}) map[string]string {
	var CategoriesI map[string]interface{}
	if s["categories"] != nil{
		CategoriesI = s["categories"].(map[string]interface{})
	}
	Categories := make(map[string]string)
	for key, value := range CategoriesI {
		 switch value := value.(type) {
		case string:
			Categories[key] = value
		}
	}
	return Categories
}


// SetOwnerReference sets owner_reference fields in  json struct
func SetOwnerReference(t []interface{}, i int) vm.OwnerReference {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		OwnerReference := vm.OwnerReference{
			Kind:		convertToString(s["kind"]),
			Name:		convertToString(s["name"]),
			UUID:		convertToString(s["uuid"]),
		}
		return OwnerReference
	}
	return vm.OwnerReference{}
}

// SetMetadata sets metadata fields in  json struct
func SetMetadata(t []interface{}, i int) vm.Metadata {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		Metadata := vm.Metadata{
			Kind:		convertToString(s["kind"]),
			ParentReference:		convertToString(s["parent_reference"]),
			Name:		convertToString(s["name"]),
			LastUpdateTime:		convertToString(s["last_update_time"]),
			UUID:		convertToString(s["uuid"]),
			CreationTime:		convertToString(s["creation_time"]),
			Categories:		SetCategories(s),
			OwnerReference:		SetOwnerReference(s["owner_reference"].(*schema.Set).List(), 0),
			EntityVersion:		convertToInt(s["entity_version"]),
		}
		return Metadata
	}
	return vm.Metadata{}
}

// SetVMIntentInput sets VmIntentInput fields in  json struct
func SetVMIntentInput(t []interface{}, i int) vm.VMIntentInput {
	if len(t) > 0 {
		s := t[i].(map[string]interface{})

		VMIntentInput := vm.VMIntentInput{
			Spec:		SetSpec(s["spec"].(*schema.Set).List(), 0),
			APIVersion:		convertToString(s["api_version"]),
			Metadata:		SetMetadata(s["metadata"].(*schema.Set).List(), 0),
		}
		return VMIntentInput
	}
	return vm.VMIntentInput{}
}
