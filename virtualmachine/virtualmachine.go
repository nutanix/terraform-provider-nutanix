package virtualmachine

// SubnetReference is SubnetReference data type
type SubnetReference struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

// IPEndpointList is IPEndpointList data type
type IPEndpointList struct {
	Address string `json:"address,omitempty"`
	Type    string `json:"type,omitempty"`
}

// NicList is NicList data type
type NicList struct {
	NicType                       string            `json:"nic_type,omitempty"`
	SubnetReference               *SubnetReference  `json:"subnet_reference,omitempty"`
	NetworkFunctionNicType        string            `json:"network_function_nic_type,omitempty"`
	MacAddress                    string            `json:"mac_address,omitempty"`
	IPEndpointList                []*IPEndpointList `json:"ip_endpoint_list,omitempty"`
	NetworkFunctionChainReference *SubnetReference  `json:"network_function_chain_reference,omitempty"`
}

// NutanixGuestTools is NutanixGuestTools data type
type NutanixGuestTools struct {
	ISOMountState         string    `json:"iso_mount_state,omitempty"`
	State                 string    `json:"state,omitempty"`
	EnabledCapabilityList []*string `json:"enabled_capability_list,omitempty"`
}

// GuestTools is GuestTools data type
type GuestTools struct {
	NutanixGuestTools *NutanixGuestTools `json:"nutanix_guest_tools,omitempty"`
}

// GPUList is GPUList data type
type GPUList struct {
	Vendor   string `json:"vendor,omitempty"`
	Mode     string `json:"mode,omitempty"`
	DeviceID int    `json:"device_id,omitempty"`
}

// CloudInit is CloudInit data type
type CloudInit struct {
	MetaData string `json:"meta_data,omitempty"`
	UserData string `json:"user_data,omitempty"`
}

// Sysprep is Sysprep data type
type Sysprep struct {
	InstallType string `json:"install_type,omitempty"`
	UnattendXML string `json:"unattend_xml,omitempty"`
}

// GuestCustomization is GuestCustomization data type
type GuestCustomization struct {
	CloudInit *CloudInit `json:"cloud_init,omitempty"`
	Sysprep   *Sysprep   `json:"sysprep,omitempty"`
}

// DiskAddress is DiskAddress data type
type DiskAddress struct {
	DeviceIndex int    `json:"device_index,omitempty"`
	AdapterType string `json:"adapter_type,omitempty"`
}

// BootConfig is BootConfig Data type
type BootConfig struct {
	DiskAddress *DiskAddress `json:"disk_address,omitempty"`
	MacAddress  string       `json:"mac_address,omitempty"`
}

// DeviceProperties is DeviceProperties data type
type DeviceProperties struct {
	DiskAddress *DiskAddress `json:"disk_address,omitempty"`
	DeviceType  string       `json:"device_type,omitempty"`
}

// DiskList is DiskList data type
type DiskList struct {
	DataSourceReference *SubnetReference  `json:"data_source_reference,omitempty"`
	DeviceProperties    *DeviceProperties `json:"device_properties,omitempty"`
	UUID                string            `json:"uuid,omitempty"`
	DiskSizeMib         int               `json:"disk_size_mib"`
}

// Resources is Resources data type
type Resources struct {
	NicList               []*NicList          `json:"nic_list,omitempty"`
	GuestOSID             string              `json:"guest_os_id,omitempty"`
	PowerState            string              `json:"power_state,omitempty"`
	GuestTools            *GuestTools         `json:"guest_tools,omitempty"`
	NumVCPUsPerSocket     int                 `json:"num_vcpus_per_socket,omitempty"`
	NumSockets            int                 `json:"num_sockets,omitempty"`
	GPUList               []*GPUList          `json:"gpu_list,omitempty"`
	MemorySizeMb          int                 `json:"memory_size_mb,omitempty"`
	ParentReference       *SubnetReference    `json:"parent_reference,omitempty"`
	HardwareClockTimezone string              `json:"hardware_clock_timezone,omitempty"`
	GuestCustomization    *GuestCustomization `json:"guest_customization,omitempty"`
	BootConfig            *BootConfig         `json:"boot_config,omitempty"`
	DiskList              []*DiskList         `json:"disk_list,omitempty"`
}

// Schedule is Schedule data type
type Schedule struct {
	IntervalMultiple int    `json:"interval_multiple,omitempty"`
	DurationSecs     int    `json:"duration_secs,omitempty"`
	EndTime          string `json:"end_time,omitempty"`
	StartTime        string `json:"start_time,omitempty"`
	IntervalType     string `json:"interval_type,omitempty"`
	IsSuspended      bool   `json:"is_suspended,omitempty"`
}

// ReplicationTarget is ReplicationTarget data type
type ReplicationTarget struct {
	ClusterReference          *SubnetReference `json:"cluster_reference,omitempty"`
	AvailabilityZoneReference *SubnetReference `json:"availability_zone_reference,omitempty"`
}

// SnapshotScheduleList is SnapshotScheduleList data type
type SnapshotScheduleList struct {
	RemoteRetentionQuantity int       `json:"remote_retention_quantity,omitempty"`
	SnapshotType            string    `json:"snapshot_type,omitempty"`
	LocalRetentionQuantity  int       `json:"local_retention_quantity,omitempty"`
	Schedule                *Schedule `json:"schedule,omitempty"`
}

// SnapshotPolicyList is SnapshotPolicyList data type
type SnapshotPolicyList struct {
	ReplicationTarget    *ReplicationTarget      `json:"replication_target,omitempty"`
	SnapshotScheduleList []*SnapshotScheduleList `json:"snapshot_schedule_list,omitempty"`
}

// BackupPolicy is BackupPolicy data type
type BackupPolicy struct {
	DefaultSnapshotType        string                `json:"default_snapshot_type,omitempty"`
	SnapshotPolicyList         []*SnapshotPolicyList `json:"snapshot_policy_list,omitempty"`
	ConsistencyGroupIdentifier string                `json:"consistency_group_identifier,omitempty"`
}

// Spec is Spec data type
type Spec struct {
	Name                      string           `json:"name,omitempty"`
	AvailabilityZoneReference *SubnetReference `json:"availability_zone_reference,omitempty"`
	BackupPolicy              *BackupPolicy    `json:"backup_policy,omitempty"`
	ClusterReference          *SubnetReference `json:"cluster_reference,omitempty"`
	Resources                 *Resources       `json:"resources,omitempty"`
	Description               string           `json:"description,omitempty"`
}

// MetaData is Metadata data type
type MetaData struct {
	LastUpdateTime string                 `json:"last_update_time,omitempty"`
	Kind           string                 `json:"kind,omitempty"`
	UUID           string                 `json:"uuid,omitempty"`
	CreationTime   string                 `json:"creation_time,omitempty"`
	Categories     map[string]interface{} `json:"categories,omitempty"`
	OwnerReference *SubnetReference       `json:"owner_reference,omitempty"`
	SpecVersion    int                    `json:"spec_version,omitempty"`
	EntityVersion  int                    `json:"entity_version,omitempty"`
	Name           string                 `json:"name,omitempty"`
}

// VirtualMachine is struct for the json required for the API call
type VirtualMachine struct {
	Spec       *Spec     `json:"spec,omitempty"`
	APIVersion string    `json:"api_version,omitempty"`
	Metadata   *MetaData `json:"metadata,omitempty"`
}
