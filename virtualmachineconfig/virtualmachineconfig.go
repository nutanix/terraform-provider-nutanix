package virtualmachineconfig

import (
	"github.com/hashicorp/terraform/helper/schema"
	vmdefn "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachine"
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

// SetMachineConfig function sets fields in struct from ResourceData
func SetMachineConfig(d *schema.ResourceData) vmdefn.VirtualMachine {
	spec := d.Get("spec").(*schema.Set).List()[0].(map[string]interface{})         // spec
	metadata := d.Get("metadata").(*schema.Set).List()[0].(map[string]interface{}) // metadata

	JSON := vmdefn.VirtualMachine{
		APIVersion: convertToString(d.Get("api_version")), // api_version
		Spec:       SetSpec(spec),
		Metadata:   SetMetadata(metadata),
	}
	return JSON
}

// SetMetadata sets metadata fields in json struct
func SetMetadata(s map[string]interface{}) *(vmdefn.MetaDataStruct) {

	var categories map[string]interface{}
	if s["categories"] != nil {
		categories = s["categories"].(map[string]interface{})
	}

	MetadataI := vmdefn.MetaDataStruct{
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
func SetSubnetReference(t []interface{}) *(vmdefn.SubnetReferenceStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		SubnetReferenceI := vmdefn.SubnetReferenceStruct{
			Kind: convertToString(s["kind"]),
			UUID: convertToString(s["uuid"]),
			Name: convertToString(s["name"]),
		}
		return &SubnetReferenceI
	}
	return nil
}

// SetSpec sets spec fields in json struct
func SetSpec(s map[string]interface{}) *(vmdefn.SpecStruct) {
	SpecI := vmdefn.SpecStruct{
		Resources:                 SetResources(s["resources"].(*schema.Set).List()[0].(map[string]interface{})), //resources
		Name:                      convertToString(s["name"]),                                                    //name
		Description:               convertToString(s["description"]),                                             //description
		ClusterReference:          SetSubnetReference(s["cluster_reference"].(*schema.Set).List()),               // cluster_description
		AvailabilityZoneReference: SetSubnetReference(s["availability_zone_reference"].(*schema.Set).List()),     //availability_zone_reference
		BackupPolicy:              SetBackupPolicy(s["backup_policy"].(*schema.Set).List()),                      // backup_policy
	}
	return &SpecI
}

// SetResources sets resources fields in json struct
func SetResources(s map[string]interface{}) *(vmdefn.ResourcesStruct) {

	var NicListI []*vmdefn.NicListStruct
	if s["nic_list"] != nil {
		for i := 0; i < len(s["nic_list"].([]interface{})); i++ {
			elem := SetNicList(s["nic_list"].([]interface{})[i].(map[string]interface{}))
			NicListI = append(NicListI, elem)
		}
	}

	var DiskListI []*vmdefn.DiskListStruct
	if s["disk_list"] != nil {
		for i := 0; i < len(s["disk_list"].([]interface{})); i++ {
			elem := SetDiskList(s["disk_list"].([]interface{})[i].(map[string]interface{}))
			DiskListI = append(DiskListI, elem)
		}
	}

	var GPUListI []*vmdefn.GPUListStruct
	if s["gpu_list"] != nil {
		for i := 0; i < len(s["gpu_list"].([]interface{})); i++ {
			elem := SetGPUList(s["gpu_list"].([]interface{})[i].(map[string]interface{}))
			GPUListI = append(GPUListI, elem)
		}
	}

	ResourcesI := vmdefn.ResourcesStruct{
		NumVCPUsPerSocket:     convertToInt(s["num_vcpus_per_socket"]),                              // num_vcpus_per_socket
		NumSockets:            convertToInt(s["num_sockets"]),                                       // num_sockets
		MemorySizeMb:          convertToInt(s["memory_size_mb"]),                                    // memory_size_mb
		PowerState:            convertToString(s["power_state"]),                                    // power_state
		GuestOSID:             convertToString(s["guest_os_id"]),                                    // guest_os_id
		HardwareClockTimezone: convertToString(s["hardware_clock_timezone"]),                        // hardware_clock_timezone
		NicList:               NicListI,                                                             // nic_list
		DiskList:              DiskListI,                                                            // disk_list
		GPUList:               GPUListI,                                                             // gpu_list
		ParentReference:       SetSubnetReference(s["parent_reference"].(*schema.Set).List()),       //parent_reference
		BootConfig:            SetBootConfig(s["boot_config"].(*schema.Set).List()),                 // boot_config
		GuestTools:            SetGuestTools(s["guest_tools"].(*schema.Set).List()),                 //guest_tools
		GuestCustomization:    SetGuestCustomization(s["guest_customization"].(*schema.Set).List()), //guest_customization
	}
	return &ResourcesI
}

// SetNicList sets nic_list fields in json struct
func SetNicList(t map[string]interface{}) *(vmdefn.NicListStruct) {
	if len(t) > 0 {
		s := t

		var IPEndpointListI []*vmdefn.IPEndpointListStruct
		if s["ip_endpoint_list"] != nil {
			for i := 0; i < len(s["ip_endpoint_list"].([]interface{})); i++ {
				elem := SetIPEndpointList(s["ip_endpoint_list"].([]interface{})[i].(map[string]interface{}))
				IPEndpointListI = append(IPEndpointListI, elem)
			}
		}

		NicListI := vmdefn.NicListStruct{
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
func SetIPEndpointList(t map[string]interface{}) *(vmdefn.IPEndpointListStruct) {
	if len(t) > 0 {
		s := t
		IPEndpointListI := vmdefn.IPEndpointListStruct{
			Address: convertToString(s["address"]),
			Type:    convertToString(s["type"]),
		}
		return &IPEndpointListI
	}
	return nil
}

// SetDiskList sets disk_list fields in json struct
func SetDiskList(t map[string]interface{}) *(vmdefn.DiskListStruct) {
	if len(t) > 0 {
		s := t
		DiskListI := vmdefn.DiskListStruct{
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
func SetDeviceProperties(t []interface{}) *(vmdefn.DevicePropertiesStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		DevicePropertiesI := vmdefn.DevicePropertiesStruct{
			DeviceType:  convertToString(s["device_type"]),
			DiskAddress: SetDiskAddress(s["disk_address"].(*schema.Set).List()),
		}
		return &DevicePropertiesI
	}
	return nil
}

// SetBackupPolicy sets backup-policy fields in json struct
func SetBackupPolicy(t []interface{}) *(vmdefn.BackupPolicyStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		var SnapshotPolicyListI []*vmdefn.SnapshotPolicyListStruct
		if s["snapshot_policy_list"] != nil {
			for i := 0; i < len(s["snapshot_policy_list"].([]interface{})); i++ {
				elem := SetSnapshotPolicyList(s["snapshot_policy_list"].([]interface{})[i].(map[string]interface{}))
				SnapshotPolicyListI = append(SnapshotPolicyListI, elem)
			}
		}
		BackupPolicyI := vmdefn.BackupPolicyStruct{
			DefaultSnapshotType:        convertToString(s["default_snapshot_type"]),
			ConsistencyGroupIdentifier: convertToString(s["consistency_group_identifier"]),
			SnapshotPolicyList:         SnapshotPolicyListI,
		}
		return &BackupPolicyI
	}
	return nil
}

// SetSnapshotPolicyList sets snapshot_policy_list fields in json struct
func SetSnapshotPolicyList(t map[string]interface{}) *(vmdefn.SnapshotPolicyListStruct) {
	if len(t) > 0 {
		s := t
		var SnapshotScheduleListI []*vmdefn.SnapshotScheduleListStruct
		if s["snapshot_schedule_list"] != nil {
			for i := 0; i < len(s["snapshot_schedule_list"].([]interface{})); i++ {
				elem := SetSnapshotScheduleList(s["snapshot_schedule_list"].([]interface{})[i].(map[string]interface{}))
				SnapshotScheduleListI = append(SnapshotScheduleListI, elem)
			}
		}

		SnapshotPolicyListI := vmdefn.SnapshotPolicyListStruct{
			ReplicationTarget:    SetReplicationTarget(s["replication_target"].(*schema.Set).List()),
			SnapshotScheduleList: SnapshotScheduleListI,
		}
		return &SnapshotPolicyListI
	}
	return nil
}

// SetSnapshotScheduleList sets snapshot_schedule_list fields in json struct
func SetSnapshotScheduleList(t map[string]interface{}) *(vmdefn.SnapshotScheduleListStruct) {
	if len(t) > 0 {
		s := t
		SnapshotScheduleListI := vmdefn.SnapshotScheduleListStruct{
			Schedule:                SetSchedule(s["schedule"].(*schema.Set).List()),
			SnapshotType:            convertToString(s["snapshot_type"]),
			LocalRetentionQuantity:  convertToInt(s["local_retention_quantity"]),
			RemoteRetentionQuantity: convertToInt(s["remote_retention_quantity"]),
		}
		return &SnapshotScheduleListI
	}
	return nil
}

// SetSchedule sets schedule fields in json struct
func SetSchedule(t []interface{}) *(vmdefn.ScheduleStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		ScheduleI := vmdefn.ScheduleStruct{
			IntervalMultiple: convertToInt(s["interval_multiple"]),
			DurationSecs:     convertToInt(s["duration_secs"]),
			EndTime:          convertToString(s["end_time"]),
			StartTime:        convertToString(s["start_time"]),
			IntervalType:     convertToString(s["interval_type"]),
			IsSuspended:      convertToBool(s["is_suspended"]),
		}
		return &ScheduleI
	}
	return nil
}

// SetReplicationTarget sets replication_target fields in json struct
func SetReplicationTarget(t []interface{}) *(vmdefn.ReplicationTargetStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		ReplicationTargetI := vmdefn.ReplicationTargetStruct{
			ClusterReference:          SetSubnetReference(s["cluster_reference"].(*schema.Set).List()),
			AvailabilityZoneReference: SetSubnetReference(s["availability_zone_reference"].(*schema.Set).List()),
		}
		return &ReplicationTargetI
	}
	return nil
}

// SetDiskAddress sets disk_address fields in json struct
func SetDiskAddress(t []interface{}) *(vmdefn.DiskAddressStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		DiskAddressI := vmdefn.DiskAddressStruct{
			DeviceIndex: convertToInt(s["device_index"]),
			AdapterType: convertToString(s["adapter_type"]),
		}
		return &DiskAddressI
	}
	return nil
}

// SetBootConfig sets boot_config fields in json struct
func SetBootConfig(t []interface{}) *(vmdefn.BootConfigStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		BootConfigI := vmdefn.BootConfigStruct{
			MacAddress:  convertToString(s["mac_address"]),
			DiskAddress: SetDiskAddress(s["disk_address"].(*schema.Set).List()),
		}
		return &BootConfigI
	}
	return nil
}

// SetGuestCustomization sets guest_customization fields in json struct
func SetGuestCustomization(t []interface{}) *(vmdefn.GuestCustomizationStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		GuestCustomizationI := vmdefn.GuestCustomizationStruct{
			CloudInit: SetCloudInit(s["cloud_init"].(*schema.Set).List()),
			Sysprep:   SetSysprep(s["sysprep"].(*schema.Set).List()),
		}
		return &GuestCustomizationI
	}
	return nil
}

// SetCloudInit sets cloud_init fields in json struct
func SetCloudInit(t []interface{}) *(vmdefn.CloudInitStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		CloudInitI := vmdefn.CloudInitStruct{
			MetaData: convertToString(s["meta_data"]),
			UserData: convertToString(s["user_data"]),
		}
		return &CloudInitI
	}
	return nil
}

// SetSysprep sets sys_prep fields in json struct
func SetSysprep(t []interface{}) *(vmdefn.SysprepStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})
		SysprepI := vmdefn.SysprepStruct{
			InstallType: convertToString(s["install_type"]),
			UnattendXML: convertToString(s["unattend_xml"]),
		}
		return &SysprepI
	}
	return nil
}

// SetGPUList sets gpu_list fields in json struct
func SetGPUList(t map[string]interface{}) *(vmdefn.GPUListStruct) {
	if len(t) > 0 {
		s := t
		GPUListI := vmdefn.GPUListStruct{
			Vendor:   convertToString(s["vendor"]),
			Mode:     convertToString(s["mode"]),
			DeviceID: convertToInt(s["device_id"]),
		}
		return &GPUListI
	}
	return nil
}

// SetGuestTools sets guest_tools fields in json struct
func SetGuestTools(t []interface{}) *(vmdefn.GuestToolsStruct) {
	if len(t) > 0 {
		s := t[0].(map[string]interface{})["nutanix_guest_tools"].(*schema.Set).List()[0].(map[string]interface{})
		var str []*string
		if s["enabled_capability_list"] != nil {
			for i := 0; i < len(s["enabled_capability_list"].([]interface{})); i++ {
				elem := s["enabled_capability_list"].([]interface{})[i].(string)
				str = append(str, &elem)
			}
		}

		GuestToolsI := vmdefn.GuestToolsStruct{
			NutanixGuestTools: &vmdefn.NutanixGuestToolsStruct{
				ISOMountState: convertToString(s["iso_mount_state"]),
				State:         convertToString(s["state"]),
				EnabledCapabilityList: str,
			},
		}
		return &GuestToolsI
	}
	return nil
}
