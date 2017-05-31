package setjsonfields

import (
	//	"bufio"
	//	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	st "github.com/ideadevice/terraform-ahv-provider-plugin/jsonstruct"
	//	"os"
)

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

// SetJSONFields function sets fields in struct from ResourceData
func SetJSONFields(d *schema.ResourceData) st.JSONstruct {
	spec := d.Get("spec").(*schema.Set).List()[0].(map[string]interface{})         // spec
	metadata := d.Get("metadata").(*schema.Set).List()[0].(map[string]interface{}) // metadata

	JSON := st.JSONstruct{
		APIVersion: convertToString(d.Get("api_version")), // api_version
		Spec:       SetSpec(spec),
		Metadata:   SetMetadata(metadata),
	}
	return JSON
}

// SetMetadata sets metadata fields in json struct
func SetMetadata(s map[string]interface{}) *(st.MetaDataStruct) {

	var categories map[string]interface{}
	if s["categories"] != nil {
		categories = s["categories"].(map[string]interface{})
	}

	MetadataI := st.MetaDataStruct{
		LastUpdateTime: convertToString(s["last_update_time"]),
		Kind:           convertToString(s["kind"]),
		UUID:           convertToString(s["uuid"]),
		CreationTime:   convertToString(s["creation_time"]),
		Name:           convertToString(s["name"]),
		SpecVersion:    convertToInt(s["spec_version"]),
		EntityVersion:  convertToInt(s["entity_version"]),
		OwnerReference: SetSubnetReference(s["owner_reference"].(*schema.Set).List()),
		Categories:     categories,
	}
	return &MetadataI
}

// SetSubnetReference sets owner_reference fields in json struct
func SetSubnetReference(t []interface{}) *(st.SubnetReferenceStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		SubnetReferenceI := st.SubnetReferenceStruct{
			Kind: convertToString(s["kind"]),
			UUID: convertToString(s["uuid"]),
			Name: convertToString(s["name"]),
		}
		return &SubnetReferenceI
	}
	return nil
}

// SetSpec sets spec fields in json struct
func SetSpec(s map[string]interface{}) *(st.SpecStruct) {

	resources := s["resources"].(*schema.Set).List()[0].(map[string]interface{}) // resources

	SpecI := st.SpecStruct{
		Resources: SetResources(resources),
	}
	return &SpecI
}

// SetResources sets resources fields in json struct
func SetResources(s map[string]interface{}) *(st.ResourcesStruct) {

	var NicListI []*st.NicListStruct
	if s["nic_list"] != nil {
		for i := 0; i < len(s["nic_list"].([]interface{})); i++ {
			elem := SetNicList(s["nic_list"].([]interface{})[i].(map[string]interface{}))
			NicListI = append(NicListI, elem)
		}
	}

	var DiskListI []*st.DiskListStruct
	if s["disk_list"] != nil {
		for i := 0; i < len(s["disk_list"].([]interface{})); i++ {
			elem := SetDiskList(s["disk_list"].([]interface{})[i].(map[string]interface{}))
			DiskListI = append(DiskListI, elem)
		}
	}

	var GPUListI []*st.GPUListStruct
	if s["gpu_list"] != nil {
		for i := 0; i < len(s["gpu_list"].([]interface{})); i++ {
			elem := SetGPUList(s["gpu_list"].([]interface{})[i].(map[string]interface{}))
			GPUListI = append(GPUListI, elem)
		}
	}

	ResourcesI := st.ResourcesStruct{
		NumVCPUsPerSocket:     convertToInt(s["num_vcpus_per_socket"]),                        // num_vcpus_per_socket
		NumSockets:            convertToInt(s["num_sockets"]),                                 // num_sockets
		MemorySizeMb:          convertToInt(s["memory_size_mb"]),                              // memory_size_mb
		PowerState:            convertToString(s["power_state"]),                              // power_state
		GuestOSID:             convertToString(s["guest_os_id"]),                              // guest_os_id
		HardwareClockTimezone: convertToString(s["hardware_clock_timezone"]),                  // hardware_clock_timezone
		NicList:               NicListI,                                                       // nic_list
		DiskList:              DiskListI,                                                      // disk_list
		GPUList:               GPUListI,                                                       // gpu_list
		ParentReference:       SetSubnetReference(s["parent_reference"].(*schema.Set).List()), //parent_reference
		BootConfig:            SetBootConfig(s["boot_config"].(*schema.Set).List()),           // boot_config
		GuestTools:            SetGuestTools(s["guest_tools"].(*schema.Set).List()),           //guest_tools

	}
	return &ResourcesI
}

// SetNicList sets nic_list fields in json struct
func SetNicList(t map[string]interface{}) *(st.NicListStruct) {
	if len(t) > 0 {
		s := t

		var IPEndpointListI []*st.IPEndpointListStruct
		if s["ip_endpoint_list"] != nil {
			for i := 0; i < len(s["ip_endpoint_list"].([]interface{})); i++ {
				elem := SetIPEndpointList(s["ip_endpoint_list"].([]interface{})[i].(map[string]interface{}))
				IPEndpointListI = append(IPEndpointListI, elem)
			}
		}

		NicListI := st.NicListStruct{
			NicType:                       convertToString(s["nic_type"]),
			NetworkFunctionNicType:        convertToString(s["network_function_nic_type"]),
			MacAddress:                    convertToString(s["mac_address"]),
			SubnetReference:               SetSubnetReference(s["subnet_reference"].(*schema.Set).List()),
			NetworkFunctionChainReference: SetSubnetReference(s["network_function_chain_reference"].(*schema.Set).List()),
			IPEndpointList:                IPEndpointListI,
		}
		return &NicListI
	}
	return nil
}

// SetIPEndpointList sets ip_endpoint_list fields in json struct
func SetIPEndpointList(t map[string]interface{}) *(st.IPEndpointListStruct) {
	if len(t) > 0 {
		s := t
		IPEndpointListI := st.IPEndpointListStruct{
			IP:   convertToString(s["ip"]),
			Type: convertToString(s["type"]),
		}
		return &IPEndpointListI
	}
	return nil
}

// SetDiskList sets disk_list fields in json struct
func SetDiskList(t map[string]interface{}) *(st.DiskListStruct) {
	if len(t) > 0 {
		s := t
		DiskListI := st.DiskListStruct{
			UUID:                convertToString(s["uuid"]),
			DiskSizeMib:         convertToInt(s["disk_size_mib"]),
			DataSourceReference: SetSubnetReference(s["data_source_reference"].(*schema.Set).List()),
			DeviceProperties:    SetDeviceProperties(s["device_properties"].(*schema.Set).List()),
		}
		return &DiskListI
	}
	return nil
}

// SetDeviceProperties sets device_properties fields in json struct
func SetDeviceProperties(t []interface{}) *(st.DevicePropertiesStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		DevicePropertiesI := st.DevicePropertiesStruct{
			DeviceType:  convertToString(s["device_type"]),
			DiskAddress: SetDiskAddress(s["disk_address"].(*schema.Set).List()),
		}
		return &DevicePropertiesI
	}
	return nil
}

// SetDiskAddress sets disk_address fields in json struct
func SetDiskAddress(t []interface{}) *(st.DiskAddressStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		DiskAddressI := st.DiskAddressStruct{
			DeviceIndex: convertToInt(s["device_index"]),
			AdapterType: convertToString(s["adapter_type"]),
		}
		return &DiskAddressI
	}
	return nil
}

// SetBootConfig sets boot_config fields in json struct
func SetBootConfig(t []interface{}) *(st.BootConfigStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		BootConfigI := st.BootConfigStruct{
			MacAddress:  convertToString(s["mac_address"]),
			DiskAddress: SetDiskAddress(s["disk_address"].(*schema.Set).List()),
		}
		return &BootConfigI
	}
	return nil
}

// SetGPUList sets gpu_list fields in json struct
func SetGPUList(t map[string]interface{}) *(st.GPUListStruct) {
	if len(t) > 0 {
		s := t
		GPUListI := st.GPUListStruct{
			Vendor:   convertToString(s["vendor"]),
			Mode:     convertToString(s["mode"]),
			DeviceID: convertToInt(s["device_id"]),
		}
		return &GPUListI
	}
	return nil
}

// SetGuestTools sets guest_tools fields in json struct
func SetGuestTools(t []interface{}) *(st.GuestToolsStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})["nutanix_guest_tools"].(*schema.Set).List()[0].(map[string]interface{})
		var str []*string
		if s["enabled_capability_list"] != nil {
			for i := 0; i < len(s["enabled_capability_list"].([]interface{})); i++ {
				elem := s["enabled_capability_list"].([]interface{})[i].(string)
				str = append(str, &elem)
			}
		}

		GuestToolsI := st.GuestToolsStruct{
			NutanixGuestTools: &st.NutanixGuestToolsStruct{
				ISOMountState: convertToString(s["iso_mount_state"]),
				State:         convertToString(s["state"]),
				EnabledCapabilityList: str,
			},
		}
		return &GuestToolsI
	}
	return nil
}

func main() {
}
