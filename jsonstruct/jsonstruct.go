package jsonstruct

// SubnetReferenceStruct is SubnetReference data type
type SubnetReferenceStruct struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

// IPEndpointListStruct is IPEndpointList data type
type IPEndpointListStruct struct {
	IP   string `json:"ip,omitempty"`
	Type string `json:"type,omitempty"`
}

// NicListStruct is NicList data type
type NicListStruct struct {
	NicType                       string                  `json:"nic_type,omitempty"`
	SubnetReference               *SubnetReferenceStruct  `json:"subnet_reference,omitempty"`
	NetworkFunctionNicType        string                  `json:"network_function_nic_type,omitempty"`
	MacAddress                    string                  `json:"mac_address,omitempty"`
	IPEndpointList                []*IPEndpointListStruct `json:"ip_endpoint_list,omitempty"`
	NetworkFunctionChainReference *SubnetReferenceStruct  `json:"network_function_chain_reference,omitempty"`
}

// NutanixGuestToolsStruct is NutanixGuestTools data type
type NutanixGuestToolsStruct struct {
	ISOMountState         string    `json:"iso_mount_state,omitempty"`
	State                 string    `json:"state,omitempty"`
	EnabledCapabilityList []*string `json:"enabled_capability_list,omitempty"`
}

// GuestToolsStruct is GuestTools data type
type GuestToolsStruct struct {
	NutanixGuestTools *NutanixGuestToolsStruct `json:"nutanix_guest_tools,omitempty"`
}

// GPUListStruct is GPUList data type
type GPUListStruct struct {
	Vendor   string `json:"vendor,omitempty"`
	Mode     string `json:"mode,omitempty"`
	DeviceID int    `json:"device_id,omitempty"`
}

// CloudInitStruct is CloudInit data type
type CloudInitStruct struct {
	MetaData        string         `json:"meta_data,omitempty"`
	UserData        string         `json:"user_data,omitempty"`
	CustomKeyValues []*interface{} `json:"custom_key_values,omitempty"`
}

// SysprepStruct is Sysprep data type
type SysprepStruct struct {
	InstallType     string         `json:"install_type,omitempty"`
	UnattendXML     string         `json:"unattend_xml,omitempty"`
	CustomKeyValues []*interface{} `json:"custom_key_values,omitempty"`
}

// GuestCustomizationStruct is GuestCustomization data type
type GuestCustomizationStruct struct {
	CloudInit *CloudInitStruct `json:"cloud_init,omitempty"`
	Sysprep   *SysprepStruct   `json:"sysprep,omitempty"`
}

// DiskAddressStruct is DiskAddress data type
type DiskAddressStruct struct {
	DeviceIndex int    `json:"device_index,omitempty"`
	AdapterType string `json:"adapter_type,omitempty"`
}

// BootConfigStruct is BootConfig Data type
type BootConfigStruct struct {
	DiskAddress *DiskAddressStruct `json:"disk_address,omitempty"`
	MacAddress  string             `json:"mac_address,omitempty"`
}

// DevicePropertiesStruct is DeviceProperties data type
type DevicePropertiesStruct struct {
	DiskAddress *DiskAddressStruct `json:"disk_address,omitempty"`
	DeviceType  string             `json:"device_type,omitempty"`
}

// DiskListStruct is DiskList data type
type DiskListStruct struct {
	DataSourceReference *SubnetReferenceStruct  `json:"data_source_reference,omitempty"`
	DeviceProperties    *DevicePropertiesStruct `json:"device_properties,omitempty"`
	UUID                string                  `json:"uuid,omitempty"`
	DiskSizeMib         int                     `json:"disk_size_mib"`
}

// ResourcesStruct is Resources data type
type ResourcesStruct struct {
	NicList               []*NicListStruct          `json:"nic_list,omitempty"`
	GuestOSID             string                    `json:"guest_os_id,omitempty"`
	PowerState            string                    `json:"power_state,omitempty"`
	GuestTools            *GuestToolsStruct         `json:"guest_tools,omitempty"`
	NumVCPUsPerSocket     int                       `json:"num_vcpus_per_socket,omitempty"`
	NumSockets            int                       `json:"num_sockets,omitempty"`
	GPUList               []*GPUListStruct          `json:"gpu_list,omitempty"`
	MemorySizeMb          int                       `json:"memory_size_mb,omitempty"`
	ParentReference       *SubnetReferenceStruct    `json:"parent_reference,omitempty"`
	HardwareClockTimezone string                    `json:"hardware_clock_timezone,omitempty"`
	GuestCustomization    *GuestCustomizationStruct `json:"guest_customization,omitempty"`
	BootConfig            *BootConfigStruct         `json:"boot_config,omitempty"`
	DiskList              []*DiskListStruct         `json:"disk_list,omitempty"`
}

// ScheduleStruct is Schedule data type
type ScheduleStruct struct {
	IntervalMultiple int    `json:"interval_multiple,omitempty"`
	DurationSecs     int    `json:"duration_secs,omitempty"`
	EndTime          string `json:"end_time,omitempty"`
	StartTime        string `json:"start_time,omitempty"`
	IntervalType     string `json:"interval_type,omitempty"`
	IsSuspended      bool   `json:"is_suspended,omitempty"`
}

// ReplicationTargetStruct is ReplicationTarget data type
type ReplicationTargetStruct struct {
	ClusterReference          *SubnetReferenceStruct `json:"cluster_reference,omitempty"`
	AvailabilityZoneReference *SubnetReferenceStruct `json:"availability_zone_reference,omitempty"`
}

// SnapshotScheduleListStruct is SnapshotScheduleList data type
type SnapshotScheduleListStruct struct {
	RemoteRetentionQuantity int             `json:"remote_retention_quantity,omitempty"`
	SnapshotType            string          `json:"snapshot_type,omitempty"`
	LocalRetentionQuantity  int             `json:"local_retention_quantity,omitempty"`
	Schedule                *ScheduleStruct `json:"schedule,omitempty"`
}

// SnapshotPolicyListStruct is SnapshotPolicyList data type
type SnapshotPolicyListStruct struct {
	ReplicationTarget    *ReplicationTargetStruct      `json:"replication_target,omitempty"`
	SnapshotScheduleList []*SnapshotScheduleListStruct `json:"snapshot_schedule_list,omitempty"`
}

// BackupPolicyStruct is BackupPolicy data type
type BackupPolicyStruct struct {
	DefaultSnapshotType        string                      `json:"default_snapshot_type,omitempty"`
	SnapshotPolicyList         []*SnapshotPolicyListStruct `json:"snapshot_policy_list,omitempty"`
	ConsistencyGroupIdentifier string                      `json:"consistency_group_identifier,omitempty"`
}

// SpecStruct is Spec data type
type SpecStruct struct {
	Name                      string                 `json:"name,omitempty"`
	AvailabilityZoneReference *SubnetReferenceStruct `json:"availability_zone_reference,omitempty"`
	BackupPolicy              *BackupPolicyStruct    `json:"backup_policy,omitempty"`
	ClusterReference          *SubnetReferenceStruct `json:"cluster_reference,omitempty"`
	Resources                 *ResourcesStruct       `json:"resources,omitempty"`
	Description               string                 `json:"description,omitempty"`
}

// MetaDataStruct is Metadata data type
type MetaDataStruct struct {
	LastUpdateTime string                 `json:"last_update_time,omitempty"`
	Kind           string                 `json:"kind,omitempty"`
	UUID           string                 `json:"uuid,omitempty"`
	CreationTime   string                 `json:"creation_time,omitempty"`
	Categories     map[string]interface{} `json:"categories,omitempty"`
	OwnerReference *SubnetReferenceStruct `json:"owner_reference,omitempty"`
	SpecVersion    int                    `json:"spec_version,omitempty"`
	EntityVersion  int                    `json:"entity_version,omitempty"`
	Name           string                 `json:"name,omitempty"`
}

// JSONstruct is struct for the json required for the API call
type JSONstruct struct {
	Spec       *SpecStruct     `json:"spec,omitempty"`
	APIVersion string          `json:"api_version,omitempty"`
	Metadata   *MetaDataStruct `json:"metadata,omitempty"`
}
