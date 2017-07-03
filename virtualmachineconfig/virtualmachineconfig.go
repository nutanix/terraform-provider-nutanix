package virtualmachineconfig

import (
	"log"
	"time"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	nutanixV3 "nutanixV3"
)

func convertToBool(a interface{}) bool {
	if a != nil {
		return a.(bool)
	}
	return false
}

func convertToInt(a interface{}) int64 {
	if a != nil {
		i :=  a.(int)
		return int64(i)
	}
	return 0
}

func convertToString(a interface{}) string {
	if a != nil {
		return a.(string)
	}
	return ""
}

// SetMachineConfig function sets fields in struct from ResourceData
func SetMachineConfig(d *schema.ResourceData) nutanixV3.VmIntentInput {
	spec := d.Get("spec").(*schema.Set).List()[0].(map[string]interface{}) // spec
	metadata := d.Get("metadata").(*schema.Set).List()                     // metadata
	JSON := nutanixV3.VmIntentInput{
		ApiVersion: "3.0", // api_version
		Spec:       SetSpec(spec), 	//Spec
		Metadata:   SetMetadata(metadata), 	//Metadata
	}
	return JSON
}

// SetMetadata sets metadata fields in json struct
func SetMetadata(t []interface{}) (nutanixV3.VmMetadata) {
	if len(t) == 0 {
		MetadataI := nutanixV3.VmMetadata{
			Kind:        "vm",
			Name:        "",
		}
		return MetadataI
	}
	s := t[0].(map[string]interface{})
	var categoriesI map[string]interface{}
	if s["categories"] != nil {
		categoriesI = s["categories"].(map[string]interface{})
	}
	categories := make(map[string]string)
	for key, value := range categoriesI {        
		switch value := value.(type) {
	    case string:
		     categories[key] = value
	    }
	}


	var temp string
	var err error
	var lastUpdateTime time.Time
	var creationTime time.Time
	temp = convertToString(s["last_update_time"])
	if temp != "" {
		lastUpdateTime, err = time.Parse(temp, temp)
	}	
	if err != nil {
		log.Fatal(err)
	}
	temp = convertToString(s["creation_time"])
	if temp != "" {
		creationTime, err = time.Parse(temp, temp)
	}	
	if err != nil {
		log.Fatal(err)
	}
	MetadataI := nutanixV3.VmMetadata{
		LastUpdateTime: lastUpdateTime,
		Kind:           "vm",
		Uuid:           convertToString(s["uuid"]),
		CreationTime:   creationTime,
		Name:           convertToString(s["name"]),
		EntityVersion:  convertToInt(s["entity_version"]),
		OwnerReference: SetSubnetReference(s["owner_reference"].(*schema.Set).List()),
		Categories:     categories,
	}
	return MetadataI
}

// SetSubnetReference sets owner_reference fields in json struct
func SetSubnetReference(t []interface{}) (nutanixV3.UserReference) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		UserReferenceI := nutanixV3.UserReference{
			Kind: convertToString(s["kind"]),
			Uuid: convertToString(s["uuid"]),
			Name: convertToString(s["name"]),
		}
		return UserReferenceI
	}
	return nutanixV3.UserReference{}
}

// SetSpec sets spec fields in json struct
func SetSpec(s map[string]interface{}) (nutanixV3.Vm) {
	SpecI := nutanixV3.Vm{
		Resources:                 SetResources(s["resources"].(*schema.Set).List()[0].(map[string]interface{})), //resources
		Name:                      convertToString(s["name"]),                                                    //name
		Description:               convertToString(s["description"]),                                             //description
		ClusterReference:          nutanixV3.ClusterReference(SetSubnetReference(s["cluster_reference"].(*schema.Set).List())),               // cluster_description
	}
	return SpecI
}

// SetResources sets resources fields in json struct
func SetResources(s map[string]interface{}) (nutanixV3.VmResources) {

	var NicListI []nutanixV3.VmNic
	if s["nic_list"] != nil {
		for i := 0; i < len(s["nic_list"].([]interface{})); i++ {
			elem := SetNicList(s["nic_list"].([]interface{})[i].(map[string]interface{}))
			NicListI = append(NicListI, elem)
		}
	}

	var DiskListI []nutanixV3.VmDisk
	if s["disk_list"] != nil {
		for i := 0; i < len(s["disk_list"].([]interface{})); i++ {
			elem := SetDiskList(s["disk_list"].([]interface{})[i].(map[string]interface{}))
			DiskListI = append(DiskListI, elem)
		}
	}

	var GPUListI []nutanixV3.VmGpu
	if s["gpu_list"] != nil {
		for i := 0; i < len(s["gpu_list"].([]interface{})); i++ {
			elem := SetGPUList(s["gpu_list"].([]interface{})[i].(map[string]interface{}))
			GPUListI = append(GPUListI, elem)
		}
	}
	powerState := "POWERED_OFF"
	if strings.ToUpper(convertToString(s["power_state"])) == "ON" {
		powerState = "POWERED_ON"
	}

	ResourcesI := nutanixV3.VmResources{
		NumVcpusPerSocket:     convertToInt(s["num_vcpus_per_socket"]),                              // num_vcpus_per_socket
		NumSockets:            convertToInt(s["num_sockets"]),                                       // num_sockets
		MemorySizeMib:          convertToInt(s["memory_size_mb"]),                                    // memory_size_mb
		PowerState:            powerState,                                                           // power_state
		NicList:               NicListI,                                                             // nic_list
		DiskList:              DiskListI,                                                            // disk_list
		GpuList:               GPUListI,                                                             // gpu_list
		ParentReference:       nutanixV3.Reference(SetSubnetReference(s["parent_reference"].(*schema.Set).List())),       //parent_reference
		BootConfig:            SetBootConfig(s["boot_config"].(*schema.Set).List()),                 // boot_config
		GuestCustomization:    SetGuestCustomization(s["guest_customization"].(*schema.Set).List()), //guest_customization
	}
	return ResourcesI
}

// SetNicList sets nic_list fields in json struct
func SetNicList(t map[string]interface{}) (nutanixV3.VmNic) {
	if len(t) > 0 {
		s := t

		var IPEndpointListI []nutanixV3.IpAddress
		if s["ip_endpoint_list"] != nil {
			for i := 0; i < len(s["ip_endpoint_list"].([]interface{})); i++ {
				elem := SetIPEndpointList(s["ip_endpoint_list"].([]interface{})[i].(map[string]interface{}))
				IPEndpointListI = append(IPEndpointListI, elem)
			}
		}

		NicListI := nutanixV3.VmNic{
			NicType:                       convertToString(s["nic_type"]),
			NetworkFunctionNicType:        convertToString(s["network_function_nic_type"]),
			MacAddress:                    convertToString(s["mac_address"]),
			SubnetReference:               nutanixV3.SubnetReference(SetSubnetReference(s["subnet_reference"].(*schema.Set).List())),
			NetworkFunctionChainReference: nutanixV3.NetworkFunctionChainReference(SetSubnetReference(s["network_function_chain_reference"].(*schema.Set).List())),
			IpEndpointList:                IPEndpointListI,
		}
		return NicListI
	}
	return nutanixV3.VmNic{}
}

// SetIPEndpointList sets ip_endpoint_list fields in json struct
func SetIPEndpointList(t map[string]interface{}) (nutanixV3.IpAddress) {
	if len(t) > 0 {
		s := t
		IPEndpointListI := nutanixV3.IpAddress{
			Address: convertToString(s["address"]),
			Type_:    convertToString(s["type"]),
		}
		return IPEndpointListI
	}
	return nutanixV3.IpAddress{}
}

// SetDiskList sets disk_list fields in json struct
func SetDiskList(t map[string]interface{}) (nutanixV3.VmDisk) {
	if len(t) > 0 {
		s := t
		DiskListI := nutanixV3.VmDisk{
			DiskSizeMib:         convertToInt(s["disk_size_mib"]),
			DataSourceReference: nutanixV3.Reference(SetSubnetReference(s["data_source_reference"].(*schema.Set).List())),
			DeviceProperties:    SetDeviceProperties(s["device_properties"].(*schema.Set).List()),
		}
		return DiskListI
	}
	return nutanixV3.VmDisk{}
}

// SetDeviceProperties sets device_properties fields in json struct
func SetDeviceProperties(t []interface{}) (nutanixV3.VmDiskDeviceProperties) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		DevicePropertiesI := nutanixV3.VmDiskDeviceProperties{
			DeviceType:  convertToString(s["device_type"]),
			DiskAddress: SetDiskAddress(s["disk_address"].(*schema.Set).List()),
		}
		return DevicePropertiesI
	}
	return nutanixV3.VmDiskDeviceProperties{}
}


// SetDiskAddress sets disk_address fields in json struct
func SetDiskAddress(t []interface{}) (nutanixV3.DiskAddress) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		DiskAddressI := nutanixV3.DiskAddress{
			DeviceIndex: convertToInt(s["device_index"]),
			AdapterType: convertToString(s["adapter_type"]),
		}
		return DiskAddressI
	}
	return nutanixV3.DiskAddress{}
}

// SetBootConfig sets boot_config fields in json struct
func SetBootConfig(t []interface{}) (nutanixV3.VmBootConfig) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		BootConfigI := nutanixV3.VmBootConfig{
			MacAddress:  convertToString(s["mac_address"]),
			DiskAddress: SetDiskAddress(s["disk_address"].(*schema.Set).List()),
		}
		return BootConfigI
	}
	return nutanixV3.VmBootConfig{}
}

// SetGuestCustomization sets guest_customization fields in json struct
func SetGuestCustomization(t []interface{}) (nutanixV3.GuestCustomization) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		GuestCustomizationI := nutanixV3.GuestCustomization{
			CloudInit: SetCloudInit(s["cloud_init"].(*schema.Set).List()),
			Sysprep:   SetSysprep(s["sysprep"].(*schema.Set).List()),
		}
		return GuestCustomizationI
	}
	return nutanixV3.GuestCustomization{}
}

// SetCloudInit sets cloud_init fields in json struct
func SetCloudInit(t []interface{}) (nutanixV3.GuestCustomizationCloudInit) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		CloudInitI := nutanixV3.GuestCustomizationCloudInit{
			MetaData: convertToString(s["meta_data"]),
			UserData: convertToString(s["user_data"]),
		}
		return CloudInitI
	}
	return nutanixV3.GuestCustomizationCloudInit{}
}

// SetSysprep sets sys_prep fields in json struct
func SetSysprep(t []interface{}) (nutanixV3.GuestCustomizationSysprep) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		SysprepI := nutanixV3.GuestCustomizationSysprep{
			InstallType: convertToString(s["install_type"]),
			UnattendXml: convertToString(s["unattend_xml"]),
		}
		return SysprepI
	}
	return nutanixV3.GuestCustomizationSysprep{}
}

// SetGPUList sets gpu_list fields in json struct
func SetGPUList(t map[string]interface{}) (nutanixV3.VmGpu) {
	if len(t) > 0 {
		s := t
		GPUListI := nutanixV3.VmGpu{
			Vendor:   convertToString(s["vendor"]),
			Mode:     convertToString(s["mode"]),
			DeviceId: convertToInt(s["device_id"]),
		}
		return GPUListI
	}
	return nutanixV3.VmGpu{}
}
