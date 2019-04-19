package v3

import "time"

// Reference ...
type Reference struct {
	Kind *string `json:"kind"`
	Name *string `json:"name,omitempty"`
	UUID *string `json:"uuid"`
}

// VMVnumaConfig Indicates how VM vNUMA should be configured
type VMVnumaConfig struct {

	// Number of vNUMA nodes. 0 means vNUMA is disabled.
	NumVnumaNodes *int64 `json:"num_vnuma_nodes,omitempty"`
}

type VMSerialPort struct {
	Index       *int64 `json:"index,omitempty"`
	IsConnected *bool  `json:"is_connected,omitempty"`
}

// IPAddress An IP address.
type IPAddress struct {

	// Address *string.
	IP *string `json:"ip,omitempty"`

	// Address type. It can only be \"ASSIGNED\" in the spec. If no type is specified in the spec, the default type is
	// set to \"ASSIGNED\".
	Type *string `json:"type,omitempty"`
}

// VMNic Virtual Machine NIC.
type VMNic struct {

	// IP endpoints for the adapter. Currently, IPv4 addresses are supported.
	IPEndpointList []*IPAddress `json:"ip_endpoint_list,omitempty"`

	// The MAC address for the adapter.
	MacAddress *string `json:"mac_address,omitempty"`

	// The model of this NIC.
	Model *string `json:"model,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty"`

	// The type of this Network function NIC. Defaults to INGRESS.
	NetworkFunctionNicType *string `json:"network_function_nic_type,omitempty"`

	// The type of this NIC. Defaults to NORMAL_NIC.
	NicType *string `json:"nic_type,omitempty"`

	SubnetReference *Reference `json:"subnet_reference,omitempty"`

	// The NIC's UUID, which is used to uniquely identify this particular NIC. This UUID may be used to refer to the NIC
	// outside the context of the particular VM it is attached to.
	UUID *string `json:"uuid,omitempty"`

	IsConnected *bool `json:"is_connected,omitempty"`
}

// DiskAddress Disk Address.
type DiskAddress struct {
	AdapterType *string `json:"adapter_type,omitempty"`
	DeviceIndex *int64  `json:"device_index,omitempty"`
}

// VMBootDevice Indicates which device a VM should boot from. One of disk_address or mac_address should be provided.
type VMBootDevice struct {

	// Address of disk to boot from.
	DiskAddress *DiskAddress `json:"disk_address,omitempty"`

	// MAC address of nic to boot from.
	MacAddress *string `json:"mac_address,omitempty"`
}

// VMBootConfig Indicates which device a VM should boot from.
type VMBootConfig struct {

	// Indicates which device a VM should boot from. Boot device takes precdence over boot device order. If both are
	// given then specified boot device will be primary boot device and remaining devices will be assigned boot order
	// according to boot device order field.
	BootDevice *VMBootDevice `json:"boot_device,omitempty"`

	// Indicates the order of device types in which VM should try to boot from. If boot device order is not provided the
	// system will decide appropriate boot device order.
	BootDeviceOrderList []*string `json:"boot_device_order_list,omitempty"`
}

// NutanixGuestToolsSpec Information regarding Nutanix Guest Tools.
type NutanixGuestToolsSpec struct {
	State                 *string           `json:"state,omitempty"`                   // Nutanix Guest Tools is enabled or not.
	Version               *string           `json:"version,omitempty"`                 // Version of Nutanix Guest Tools installed on the VM.
	NgtState              *string           `json:"ngt_state,omitempty"`               //Nutanix Guest Tools installed or not.
	Credentials           map[string]string `json:"credentials,omitempty"`             //Credentials to login server
	IsoMountState         *string           `json:"iso_mount_state,omitempty"`         // Desired mount state of Nutanix Guest Tools ISO.
	EnabledCapabilityList []*string         `json:"enabled_capability_list,omitempty"` // Application names that are enabled.
}

// GuestToolsSpec Information regarding guest tools.
type GuestToolsSpec struct {

	// Nutanix Guest Tools information
	NutanixGuestTools *NutanixGuestToolsSpec `json:"nutanix_guest_tools,omitempty"`
}

// VMGpu Graphics resource information for the Virtual Machine.
type VMGpu struct {

	// The device ID of the GPU.
	DeviceID *int64 `json:"device_id,omitempty"`

	// The mode of this GPU.
	Mode *string `json:"mode,omitempty"`

	// The vendor of the GPU.
	Vendor *string `json:"vendor,omitempty"`
}

// GuestCustomizationCloudInit If this field is set, the guest will be customized using cloud-init. Either user_data or
// custom_key_values should be provided. If custom_key_ves are provided then the user data will be generated using these
// key-value pairs.
type GuestCustomizationCloudInit struct {

	// Generic key value pair used for custom attributes
	CustomKeyValues map[string]string `json:"custom_key_values,omitempty"`

	// The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must
	// be base64 encoded.
	MetaData *string `json:"meta_data,omitempty"`

	// The contents of the user_data configuration for cloud-init. This can be formatted as YAML, JSON, or could be a
	// shell script. The value must be base64 encoded.
	UserData *string `json:"user_data,omitempty"`
}

// GuestCustomizationSysprep If this field is set, the guest will be customized using Sysprep. Either unattend_xml or
// custom_key_values should be provided. If custom_key_values are provided then the unattended answer file will be
// generated using these key-value pairs.
type GuestCustomizationSysprep struct {

	// Generic key value pair used for custom attributes
	CustomKeyValues map[string]string `json:"custom_key_values,omitempty"`

	// Whether the guest will be freshly installed using this unattend configuration, or whether this unattend
	// configuration will be applied to a pre-prepared image. Default is \"PREPARED\".
	InstallType *string `json:"install_type,omitempty"`

	// This field contains a Sysprep unattend xml definition, as a *string. The value must be base64 encoded.
	UnattendXML *string `json:"unattend_xml,omitempty"`
}

// GuestCustomization VM guests may be customized at boot time using one of several different methods. Currently,
// cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or
// cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting
// to change it after creation will result in an error. Additional properties can be specified. For example - in the
// context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own
// custom script.
type GuestCustomization struct {
	CloudInit *GuestCustomizationCloudInit `json:"cloud_init,omitempty"`

	// Flag to allow override of customization by deployer.
	IsOverridable *bool `json:"is_overridable,omitempty"`

	Sysprep *GuestCustomizationSysprep `json:"sysprep,omitempty"`
}

// VMGuestPowerStateTransitionConfig Extra configs related to power state transition.
type VMGuestPowerStateTransitionConfig struct {

	// Indicates whether to execute set script before ngt shutdown/reboot.
	EnableScriptExec *bool `json:"enable_script_exec,omitempty"`

	// Indicates whether to abort ngt shutdown/reboot if script fails.
	ShouldFailOnScriptFailure *bool `json:"should_fail_on_script_failure,omitempty"`
}

// VMPowerStateMechanism Indicates the mechanism guiding the VM power state transition. Currently used for the transition
// to \"OFF\" state.
type VMPowerStateMechanism struct {
	GuestTransitionConfig *VMGuestPowerStateTransitionConfig `json:"guest_transition_config,omitempty"`

	// Power state mechanism (ACPI/GUEST/HARD).
	Mechanism *string `json:"mechanism,omitempty"`
}

// VMDiskDeviceProperties ...
type VMDiskDeviceProperties struct {
	DeviceType  *string      `json:"device_type,omitempty"`
	DiskAddress *DiskAddress `json:"disk_address,omitempty"`
}

// VMDisk VirtualMachine Disk (VM Disk).
type VMDisk struct {
	DataSourceReference *Reference `json:"data_source_reference,omitempty"`

	DeviceProperties *VMDiskDeviceProperties `json:"device_properties,omitempty"`

	// Size of the disk in Bytes.
	DiskSizeBytes *int64 `json:"disk_size_bytes,omitempty"`

	// Size of the disk in MiB. Must match the size specified in 'disk_size_bytes' - rounded up to the nearest MiB -
	// when that field is present.
	DiskSizeMib *int64 `json:"disk_size_mib,omitempty"`

	// The device ID which is used to uniquely identify this particular disk.
	UUID *string `json:"uuid,omitempty"`

	VolumeGroupReference *Reference `json:"volume_group_reference,omitempty"`
}

// VMResources VM Resources Definition.
type VMResources struct {

	// Indicates which device the VM should boot from.
	BootConfig *VMBootConfig `json:"boot_config,omitempty"`

	// Disks attached to the VM.
	DiskList []*VMDisk `json:"disk_list,omitempty"`

	// GPUs attached to the VM.
	GpuList []*VMGpu `json:"gpu_list,omitempty"`

	GuestCustomization *GuestCustomization `json:"guest_customization,omitempty"`

	// Guest OS Identifier. For ESX, refer to VMware documentation link
	// https://www.vmware.com/support/orchestrator/doc/vro-vsphere65-api/html/VcVirtualMachineGuestOsIdentifier.html
	// for the list of guest OS identifiers.
	GuestOsID *string `json:"guest_os_id,omitempty"`

	// Information regarding guest tools.
	GuestTools *GuestToolsSpec `json:"guest_tools,omitempty"`

	// VM's hardware clock timezone in IANA TZDB format (America/Los_Angeles).
	HardwareClockTimezone *string `json:"hardware_clock_timezone,omitempty"`

	// Memory size in MiB.
	MemorySizeMib *int64 `json:"memory_size_mib,omitempty"`

	// NICs attached to the VM.
	NicList []*VMNic `json:"nic_list,omitempty"`

	// Number of vCPU sockets.
	NumSockets *int64 `json:"num_sockets,omitempty"`

	// Number of vCPUs per socket.
	NumVcpusPerSocket *int64 `json:"num_vcpus_per_socket,omitempty"`

	// *Reference to an entity that the VM should be cloned from.
	ParentReference *Reference `json:"parent_reference,omitempty"`

	// The current or desired power state of the VM.
	PowerState *string `json:"power_state,omitempty"`

	PowerStateMechanism *VMPowerStateMechanism `json:"power_state_mechanism,omitempty"`

	// Indicates whether VGA console should be enabled or not.
	VgaConsoleEnabled *bool `json:"vga_console_enabled,omitempty"`

	// Information regarding vNUMA configuration.
	VMVnumaConfig *VMVnumaConfig `json:"vnuma_config,omitempty"`

	SerialPortList []*VMSerialPort `json:"serial_port_list,omitempty"`
}

// VM An intentful representation of a vm spec
type VM struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty"`

	// A description for vm.
	Description *string `json:"description,omitempty"`

	// vm Name.
	Name *string `json:"name"`

	Resources *VMResources `json:"resources,omitempty"`
}

// VMIntentInput ...
type VMIntentInput struct {
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata"`

	Spec *VM `json:"spec"`
}

// MessageResource ...
type MessageResource struct {

	// Custom key-value details relevant to the status.
	Details map[string]string `json:"details,omitempty"`

	// If state is ERROR, a message describing the error.
	Message *string `json:"message"`

	// If state is ERROR, a machine-readable snake-cased *string.
	Reason *string `json:"reason"`
}

// VMStatus The status of a REST API call. Only used when there is a failure to report.
type VMStatus struct {
	APIVersion *string `json:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty"`

	State *string `json:"state,omitempty"`
}

// VMNicOutputStatus Virtual Machine NIC Status.
type VMNicOutputStatus struct {

	// The Floating IP associated with the vnic.
	FloatingIP *string `json:"floating_ip,omitempty"`

	// IP endpoints for the adapter. Currently, IPv4 addresses are supported.
	IPEndpointList []*IPAddress `json:"ip_endpoint_list,omitempty"`

	// The MAC address for the adapter.
	MacAddress *string `json:"mac_address,omitempty"`

	// The model of this NIC.
	Model *string `json:"model,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty"`

	// The type of this Network function NIC. Defaults to INGRESS.
	NetworkFunctionNicType *string `json:"network_function_nic_type,omitempty"`

	// The type of this NIC. Defaults to NORMAL_NIC.
	NicType *string `json:"nic_type,omitempty"`

	SubnetReference *Reference `json:"subnet_reference,omitempty"`

	// The NIC's UUID, which is used to uniquely identify this particular NIC. This UUID may be used to refer to the NIC
	// outside the context of the particular VM it is attached to.
	UUID *string `json:"uuid,omitempty"`

	IsConnected *bool `json:"is_connected,omitempty"`
}

// NutanixGuestToolsStatus Information regarding Nutanix Guest Tools.
type NutanixGuestToolsStatus struct {
	// Version of Nutanix Guest Tools available on the cluster.
	AvailableVersion *string `json:"available_version,omitempty"`
	//Nutanix Guest Tools installed or not.
	NgtState *string `json:"ngt_state,omitempty"`
	// Desired mount state of Nutanix Guest Tools ISO.
	IsoMountState *string `json:"iso_mount_state,omitempty"`
	// Nutanix Guest Tools is enabled or not.
	State *string `json:"state,omitempty"`
	// Version of Nutanix Guest Tools installed on the VM.
	Version *string `json:"version,omitempty"`
	// Application names that are enabled.
	EnabledCapabilityList []*string `json:"enabled_capability_list,omitempty"`
	//Credentials to login server
	Credentials map[string]string `json:"credentials,omitempty"`
	// Version of the operating system on the VM.
	GuestOsVersion *string `json:"guest_os_version,omitempty"`
	// Whether the VM is configured to take VSS snapshots through NGT.
	VSSSnapshotCapable *bool `json:"vss_snapshot_capable,omitempty"`
	// Communication from VM to CVM is active or not.
	IsReachable *bool `json:"is_reachable,omitempty"`
	// Whether VM mobility drivers are installed in the VM.
	VMMobilityDriversInstalled *bool `json:"vm_mobility_drivers_installed,omitempty"`
}

// GuestToolsStatus Information regarding guest tools.
type GuestToolsStatus struct {

	// Nutanix Guest Tools information
	NutanixGuestTools *NutanixGuestToolsStatus `json:"nutanix_guest_tools,omitempty"`
}

// VMGpuOutputStatus Graphics resource status information for the Virtual Machine.
type VMGpuOutputStatus struct {

	// The device ID of the GPU.
	DeviceID *int64 `json:"device_id,omitempty"`

	// Fraction of the physical GPU assigned.
	Fraction *int64 `json:"fraction,omitempty"`

	// GPU frame buffer size in MiB.
	FrameBufferSizeMib *int64 `json:"frame_buffer_size_mib,omitempty"`

	// Last determined guest driver version.
	GuestDriverVersion *string `json:"guest_driver_version,omitempty"`

	// The mode of this GPU
	Mode *string `json:"mode,omitempty"`

	// Name of the GPU resource.
	Name *string `json:"name,omitempty"`

	// Number of supported virtual display heads.
	NumVirtualDisplayHeads *int64 `json:"num_virtual_display_heads,omitempty"`

	// GPU {segment:bus:device:function} (sbdf) address if assigned.
	PCIAddress *string `json:"pci_address,omitempty"`

	// UUID of the GPU.
	UUID *string `json:"uuid,omitempty"`

	// The vendor of the GPU.
	Vendor *string `json:"vendor,omitempty"`
}

// GuestCustomizationStatus VM guests may be customized at boot time using one of several different methods. Currently,
// cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or
// cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting
// to change it after creation will result in an error. Additional properties can be specified. For example - in the
// context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own
// custom script.
type GuestCustomizationStatus struct {
	CloudInit *GuestCustomizationCloudInit `json:"cloud_init,omitempty"`

	// Flag to allow override of customization by deployer.
	IsOverridable *bool `json:"is_overridable,omitempty"`

	Sysprep *GuestCustomizationSysprep `json:"sysprep,omitempty"`
}

// VMResourcesDefStatus VM Resources Status Definition.
type VMResourcesDefStatus struct {

	// Indicates which device the VM should boot from.
	BootConfig *VMBootConfig `json:"boot_config,omitempty"`

	// Disks attached to the VM.
	DiskList []*VMDisk `json:"disk_list,omitempty"`

	// GPUs attached to the VM.
	GpuList []*VMGpuOutputStatus `json:"gpu_list,omitempty"`

	GuestCustomization *GuestCustomizationStatus `json:"guest_customization,omitempty"`

	// Guest OS Identifier. For ESX, refer to VMware documentation link
	// https://www.vmware.com/support/orchestrator/doc/vro-vsphere65-api/html/VcVirtualMachineGuestOsIdentifier.html
	// for the list of guest OS identifiers.
	GuestOsID *string `json:"guest_os_id,omitempty"`

	// Information regarding guest tools.
	GuestTools *GuestToolsStatus `json:"guest_tools,omitempty"`

	// VM's hardware clock timezone in IANA TZDB format (America/Los_Angeles).
	HardwareClockTimezone *string `json:"hardware_clock_timezone,omitempty"`

	HostReference *Reference `json:"host_reference,omitempty"`

	// The hypervisor type for the hypervisor the VM is hosted on.
	HypervisorType *string `json:"hypervisor_type,omitempty"`

	// Memory size in MiB.
	MemorySizeMib *int64 `json:"memory_size_mib,omitempty"`

	// NICs attached to the VM.
	NicList []*VMNicOutputStatus `json:"nic_list,omitempty"`

	// Number of vCPU sockets.
	NumSockets *int64 `json:"num_sockets,omitempty"`

	// Number of vCPUs per socket.
	NumVcpusPerSocket *int64 `json:"num_vcpus_per_socket,omitempty"`

	// *Reference to an entity that the VM cloned from.
	ParentReference *Reference `json:"parent_reference,omitempty"`

	// Current power state of the VM.
	PowerState *string `json:"power_state,omitempty"`

	PowerStateMechanism *VMPowerStateMechanism `json:"power_state_mechanism,omitempty"`

	// Indicates whether VGA console has been enabled or not.
	VgaConsoleEnabled *bool `json:"vga_console_enabled,omitempty"`

	// Information regarding vNUMA configuration.
	VnumaConfig *VMVnumaConfig `json:"vnuma_config,omitempty"`

	SerialPortList []*VMSerialPort `json:"serial_port_list,omitempty"`
}

// VMDefStatus An intentful representation of a vm status
type VMDefStatus struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty"`

	// A description for vm.
	Description *string `json:"description,omitempty"`

	// Any error messages for the vm, if in an error state.
	MessageList []*MessageResource `json:"message_list,omitempty"`

	// vm Name.
	Name *string `json:"name,omitempty"`

	Resources *VMResourcesDefStatus `json:"resources,omitempty"`

	// The state of the vm.
	State *string `json:"state,omitempty"`

	ExecutionContext *ExecutionContext `json:"execution_context,omitempty"`
}

//ExecutionContext ...
type ExecutionContext struct {
	TaskUUID interface{} `json:"task_uuid,omitempty"`
}

// VMIntentResponse Response object for intentful operations on a vm
type VMIntentResponse struct {
	APIVersion *string `json:"api_version"`

	Metadata *Metadata `json:"metadata,omitempty"`

	Spec *VM `json:"spec,omitempty"`

	Status *VMDefStatus `json:"status,omitempty"`
}

// DSMetadata All api calls that return a list will have this metadata block as input
type DSMetadata struct {

	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty"`
}

// VMIntentResource Response object for intentful operations on a vm
type VMIntentResource struct {
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata"`

	Spec *VM `json:"spec,omitempty"`

	Status *VMDefStatus `json:"status,omitempty"`
}

// VMListIntentResponse Response object for intentful operation of vms
type VMListIntentResponse struct {
	APIVersion *string `json:"api_version"`

	Entities []*VMIntentResource `json:"entities,omitempty"`

	Metadata *ListMetadataOutput `json:"metadata"`
}

// SubnetMetadata The subnet kind metadata
type SubnetMetadata struct {

	// Categories for the subnet
	Categories map[string]string `json:"categories,omitempty"`

	// UTC date and time in RFC-3339 format when subnet was created
	CreationTime *time.Time `json:"creation_time,omitempty"`

	// The kind name
	Kind *string `json:"kind"`

	// UTC date and time in RFC-3339 format when subnet was last updated
	LastUpdateTime *time.Time `json:"last_update_time,omitempty"`

	// subnet name
	Name *string `json:"name,omitempty"`

	OwnerReference *Reference `json:"owner_reference,omitempty"`

	// project reference
	ProjectReference *Reference `json:"project_reference,omitempty"`

	// Hash of the spec. This will be returned from server.
	SpecHash *string `json:"spec_hash,omitempty"`

	// Version number of the latest spec.
	SpecVersion *int64 `json:"spec_version,omitempty"`

	// subnet uuid
	UUID *string `json:"uuid,omitempty"`
}

// Address represents the Host address.
type Address struct {

	// Fully qualified domain name.
	FQDN *string `json:"fqdn,omitempty"`

	// IPV4 address.
	IP *string `json:"ip,omitempty"`

	// IPV6 address.
	IPV6 *string `json:"ipv6,omitempty"`

	// Port Number
	Port *int64 `json:"port,omitempty"`
}

// IPPool represents IP pool.
type IPPool struct {

	// Range of IPs (example: 10.0.0.9 10.0.0.19).
	Range *string `json:"range,omitempty"`
}

// DHCPOptions Spec for defining DHCP options.
type DHCPOptions struct {
	BootFileName *string `json:"boot_file_name,omitempty"`

	DomainName *string `json:"domain_name,omitempty"`

	DomainNameServerList []*string `json:"domain_name_server_list,omitempty"`

	DomainSearchList []*string `json:"domain_search_list,omitempty"`

	TFTPServerName *string `json:"tftp_server_name,omitempty"`
}

// IPConfig represents the configurtion of IP.
type IPConfig struct {

	// Default gateway IP address.
	DefaultGatewayIP *string `json:"default_gateway_ip,omitempty"`

	DHCPOptions *DHCPOptions `json:"dhcp_options,omitempty"`

	DHCPServerAddress *Address `json:"dhcp_server_address,omitempty"`

	PoolList []*IPPool `json:"pool_list,omitempty"`

	PrefixLength *int64 `json:"prefix_length,omitempty"`

	// Subnet IP address.
	SubnetIP *string `json:"subnet_ip,omitempty"`
}

// SubnetResources represents Subnet creation/modification spec.
type SubnetResources struct {
	IPConfig *IPConfig `json:"ip_config,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty"`

	SubnetType *string `json:"subnet_type"`

	VlanID *int64 `json:"vlan_id,omitempty"`

	VswitchName *string `json:"vswitch_name,omitempty"`
}

// Subnet An intentful representation of a subnet spec
type Subnet struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty"`

	// A description for subnet.
	Description *string `json:"description,omitempty"`

	// subnet Name.
	Name *string `json:"name"`

	Resources *SubnetResources `json:"resources,omitempty"`
}

// SubnetIntentInput An intentful representation of a subnet
type SubnetIntentInput struct {
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata"`

	Spec *Subnet `json:"spec"`
}

// SubnetStatus represents The status of a REST API call. Only used when there is a failure to report.
type SubnetStatus struct {
	APIVersion *string `json:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty"`

	State *string `json:"state,omitempty"`
}

// SubnetResourcesDefStatus represents a Subnet creation/modification status.
type SubnetResourcesDefStatus struct {
	IPConfig *IPConfig `json:"ip_config,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty"`

	SubnetType *string `json:"subnet_type"`

	VlanID *int64 `json:"vlan_id,omitempty"`

	VswitchName *string `json:"vswitch_name,omitempty"`
}

// SubnetDefStatus An intentful representation of a subnet status
type SubnetDefStatus struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty"`

	// A description for subnet.
	Description *string `json:"description"`

	// Any error messages for the subnet, if in an error state.
	MessageList []*MessageResource `json:"message_list,omitempty"`

	// subnet Name.
	Name *string `json:"name"`

	Resources *SubnetResourcesDefStatus `json:"resources,omitempty"`

	// The state of the subnet.
	State *string `json:"state,omitempty"`

	ExecutionContext *ExecutionContext `json:"execution_context,omitempty"`
}

// SubnetIntentResponse represents the response object for intentful operations on a subnet
type SubnetIntentResponse struct {
	APIVersion *string `json:"api_version"`

	Metadata *Metadata `json:"metadata,omitempty"`

	Spec *Subnet `json:"spec,omitempty"`

	Status *SubnetDefStatus `json:"status,omitempty"`
}

// SubnetIntentResource represents Response object for intentful operations on a subnet
type SubnetIntentResource struct {
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata"`

	Spec *Subnet `json:"spec,omitempty"`

	Status *SubnetDefStatus `json:"status,omitempty"`
}

// SubnetListIntentResponse represents the response object for intentful operation of subnets
type SubnetListIntentResponse struct {
	APIVersion *string `json:"api_version"`

	Entities []*SubnetIntentResponse `json:"entities,omitempty"`

	Metadata *ListMetadataOutput `json:"metadata"`
}

// SubnetListMetadata ...
type SubnetListMetadata struct {

	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty"`
}

// Checksum represents the image checksum
type Checksum struct {
	ChecksumAlgorithm *string `json:"checksum_algorithm"`
	ChecksumValue     *string `json:"checksum_value"`
}

// ImageVersionResources The image version, which is composed of a product name and product version.
type ImageVersionResources struct {

	// Name of the producer/distribution of the image. For example windows or red hat.
	ProductName *string `json:"product_name"`

	// Version *string for the disk image.
	ProductVersion *string `json:"product_version"`
}

// ImageResources describes the image spec resources object.
type ImageResources struct {

	// The supported CPU architecture for a disk image.
	Architecture *string `json:"architecture,omitempty"`

	// Checksum of the image. The checksum is used for image validation if the image has a source specified. For images
	// that do not have their source specified the checksum is generated by the image service.
	Checksum *Checksum `json:"checksum,omitempty"`

	// The type of image.
	ImageType *string `json:"image_type,omitempty"`

	// The source URI points at the location of a the source image which is used to create/update image.
	SourceURI *string `json:"source_uri,omitempty"`

	// The image version
	Version *ImageVersionResources `json:"version,omitempty"`
}

// Image An intentful representation of a image spec
type Image struct {

	// A description for image.
	Description *string `json:"description,omitempty"`

	// image Name.
	Name *string `json:"name,omitempty"`

	Resources *ImageResources `json:"resources"`
}

// ImageMetadata Metadata The image kind metadata
type ImageMetadata struct {

	// Categories for the image
	Categories map[string]string `json:"categories,omitempty"`

	// UTC date and time in RFC-3339 format when vm was created
	CreationTime *time.Time `json:"creation_time,omitempty"`

	// The kind name
	Kind *string `json:"kind"`

	// UTC date and time in RFC-3339 format when image was last updated
	LastUpdateTime *time.Time `json:"last_update_time,omitempty"`

	// image name
	Name *string `json:"name,omitempty"`

	// project reference
	ProjectReference *Reference `json:"project_reference,omitempty"`

	OwnerReference *Reference `json:"owner_reference,omitempty"`

	// Hash of the spec. This will be returned from server.
	SpecHash *string `json:"spec_hash,omitempty"`

	// Version number of the latest spec.
	SpecVersion *int64 `json:"spec_version,omitempty"`

	// image uuid
	UUID *string `json:"uuid,omitempty"`
}

// ImageIntentInput An intentful representation of a image
type ImageIntentInput struct {
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty"`

	Spec *Image `json:"spec,omitempty"`
}

// ImageStatus represents the status of a REST API call. Only used when there is a failure to report.
type ImageStatus struct {
	APIVersion *string `json:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty"`

	State *string `json:"state,omitempty"`
}

// ImageVersionStatus represents the image version, which is composed of a product name and product version.
type ImageVersionStatus struct {

	// Name of the producer/distribution of the image. For example windows or red hat.
	ProductName *string `json:"product_name"`

	// Version *string for the disk image.
	ProductVersion *string `json:"product_version"`
}

// ImageResourcesDefStatus describes the image status resources object.
type ImageResourcesDefStatus struct {

	// The supported CPU architecture for a disk image.
	Architecture *string `json:"architecture,omitempty"`

	// Checksum of the image. The checksum is used for image validation if the image has a source specified. For images
	// that do not have their source specified the checksum is generated by the image service.
	Checksum *Checksum `json:"checksum,omitempty"`

	// The type of image.
	ImageType *string `json:"image_type,omitempty"`

	// List of URIs where the raw image data can be accessed.
	RetrievalURIList []*string `json:"retrieval_uri_list,omitempty"`

	// The size of the image in bytes.
	SizeBytes *int64 `json:"size_bytes,omitempty"`

	// The source URI points at the location of a the source image which is used to create/update image.
	SourceURI *string `json:"source_uri,omitempty"`

	// The image version
	Version *ImageVersionStatus `json:"version,omitempty"`
}

// ImageDefStatus represents an intentful representation of a image status
type ImageDefStatus struct {
	AvailabilityZoneReference *Reference `json:"availability_zone_reference,omitempty"`

	ClusterReference *Reference `json:"cluster_reference,omitempty"`

	// A description for image.
	Description *string `json:"description,omitempty"`

	// Any error messages for the image, if in an error state.
	MessageList []*MessageResource `json:"message_list,omitempty"`

	// image Name.
	Name *string `json:"name"`

	Resources ImageResourcesDefStatus `json:"resources"`

	// The state of the image.
	State *string `json:"state,omitempty"`

	ExecutionContext *ExecutionContext `json:"execution_context,omitempty"`
}

// ImageIntentResponse represents the response object for intentful operations on a image
type ImageIntentResponse struct {
	APIVersion *string `json:"api_version"`

	Metadata *Metadata `json:"metadata"`

	Spec *Image `json:"spec,omitempty"`

	Status *ImageDefStatus `json:"status,omitempty"`
}

// ImageListMetadata represents metadata input
type ImageListMetadata struct {

	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty"`
}

// ImageIntentResource represents the response object for intentful operations on a image
type ImageIntentResource struct {
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata"`

	Spec *Image `json:"spec,omitempty"`

	Status *ImageDefStatus `json:"status,omitempty"`
}

// ImageListIntentResponse represents the response object for intentful operation of images
type ImageListIntentResponse struct {
	APIVersion *string `json:"api_version"`

	Entities []*ImageIntentResponse `json:"entities,omitempty"`

	Metadata *ListMetadataOutput `json:"metadata"`
}

// ClusterListIntentResponse ...
type ClusterListIntentResponse struct {
	APIVersion *string                  `json:"api_version"`
	Entities   []*ClusterIntentResource `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput      `json:"metadata"`
}

// ClusterIntentResource ...
type ClusterIntentResource struct {
	APIVersion *string           `json:"api_version,omitempty"`
	Metadata   *Metadata         `json:"metadata"`
	Spec       *Cluster          `json:"spec,omitempty"`
	Status     *ClusterDefStatus `json:"status,omitempty"`
}

// ClusterIntentResponse ...
type ClusterIntentResponse struct {
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *Metadata `json:"metadata"`

	Spec *Cluster `json:"spec,omitempty"`

	Status *ClusterDefStatus `json:"status,omitempty"`
}

// Cluster ...
type Cluster struct {
	Name      *string          `json:"name,omitempty"`
	Resources *ClusterResource `json:"resources,omitempty"`
}

// ClusterDefStatus ...
type ClusterDefStatus struct {
	State       *string            `json:"state,omitempty"`
	MessageList []*MessageResource `json:"message_list,omitempty"`
	Name        *string            `json:"name,omitempty"`
	Resources   *ClusterObj        `json:"resources,omitempty"`
}

// ClusterObj ...
type ClusterObj struct {
	Nodes             *ClusterNodes    `json:"nodes,omitempty"`
	Config            *ClusterConfig   `json:"config,omitempty"`
	Network           *ClusterNetwork  `json:"network,omitempty"`
	Analysis          *ClusterAnalysis `json:"analysis,omitempty"`
	RuntimeStatusList []*string        `json:"runtime_status_list,omitempty"`
}

// ClusterNodes ...
type ClusterNodes struct {
	HypervisorServerList []*HypervisorServer `json:"hypervisor_server_list,omitempty"`
}

// SoftwareMapValues ...
type SoftwareMapValues struct {
	SoftwareType *string `json:"software_type,omitempty"`
	Status       *string `json:"status,omitempty"`
	Version      *string `json:"version,omitempty"`
}

// SoftwareMap ...
type SoftwareMap struct {
	NCC *SoftwareMapValues `json:"ncc,omitempty"`
	NOS *SoftwareMapValues `json:"nos,omitempty"`
}

// ClusterConfig ...
type ClusterConfig struct {
	GpuDriverVersion              *string                    `json:"gpu_driver_version,omitempty"`
	ClientAuth                    *ClientAuth                `json:"client_auth,omitempty"`
	AuthorizedPublicKeyList       []*PublicKey               `json:"authorized_public_key_list,omitempty"`
	SoftwareMap                   *SoftwareMap               `json:"software_map,omitempty"`
	EncryptionStatus              *string                    `json:"encryption_status,omitempty"`
	SslKey                        *SslKey                    `json:"ssl_key,omitempty"`
	ServiceList                   []*string                  `json:"service_list,omitempty"`
	SupportedInformationVerbosity *string                    `json:"supported_information_verbosity,omitempty"`
	CertificationSigningInfo      *CertificationSigningInfo  `json:"certification_signing_info,omitempty"`
	RedundancyFactor              *int64                     `json:"redundancy_factor,omitempty"`
	ExternalConfigurations        *ExternalConfigurations    `json:"external_configurations,omitempty"`
	OperationMode                 *string                    `json:"operation_mode,omitempty"`
	CaCertificateList             []*CaCert                  `json:"ca_certificate_list,omitempty"`
	EnabledFeatureList            []*string                  `json:"enabled_feature_list,omitempty"`
	IsAvailable                   *bool                      `json:"is_available,omitempty"`
	Build                         *BuildInfo                 `json:"build,omitempty"`
	Timezone                      *string                    `json:"timezone,omitempty"`
	ClusterArch                   *string                    `json:"cluster_arch,omitempty"`
	ManagementServerList          []*ClusterManagementServer `json:"management_server_list,omitempty"`
}

// ClusterManagementServer ...
type ClusterManagementServer struct {
	IP         *string   `json:"ip,omitempty"`
	DrsEnabled *bool     `json:"drs_enabled,omitempty"`
	StatusList []*string `json:"status_list,omitempty"`
	Type       *string   `json:"type,omitempty"`
}

// BuildInfo ...
type BuildInfo struct {
	CommitID      *string `json:"commit_id,omitempty"`
	FullVersion   *string `json:"full_version,omitempty"`
	CommitDate    *string `json:"commit_date,omitempty"`
	Version       *string `json:"version,omitempty"`
	ShortCommitID *string `json:"short_commit_id,omitempty"`
	BuildType     *string `json:"build_type,omitempty"`
}

// CaCert ...
type CaCert struct {
	CaName      *string `json:"ca_name,omitempty"`
	Certificate *string `json:"certificate,omitempty"`
}

// ExternalConfigurations ...
type ExternalConfigurations struct {
	CitrixConnectorConfig *CitrixConnectorConfigDetails `json:"citrix_connector_config,omitempty"`
}

// CitrixConnectorConfigDetails ...
type CitrixConnectorConfigDetails struct {
	CitrixVMReferenceList *[]Reference            `json:"citrix_vm_reference_list,omitempty"`
	ClientSecret          *string                 `json:"client_secret,omitempty"`
	CustomerID            *string                 `json:"customer_id,omitempty"`
	ClientID              *string                 `json:"client_id,omitempty"`
	ResourceLocation      *CitrixResourceLocation `json:"resource_location,omitempty"`
}

// CitrixResourceLocation ...
type CitrixResourceLocation struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// SslKey ...
type SslKey struct {
	KeyType        *string                   `json:"key_type,omitempty"`
	KeyName        *string                   `json:"key_name,omitempty"`
	SigningInfo    *CertificationSigningInfo `json:"signing_info,omitempty"`
	ExpireDatetime *string                   `json:"expire_datetime,omitempty"`
}

// CertificationSigningInfo ...
type CertificationSigningInfo struct {
	City             *string `json:"city,omitempty"`
	CommonNameSuffix *string `json:"common_name_suffix,omitempty"`
	State            *string `json:"state,omitempty"`
	CountryCode      *string `json:"country_code,omitempty"`
	CommonName       *string `json:"common_name,omitempty"`
	Organization     *string `json:"organization,omitempty"`
	EmailAddress     *string `json:"email_address,omitempty"`
}

// PublicKey ...
type PublicKey struct {
	Key  *string `json:"key,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ClientAuth ...
type ClientAuth struct {
	Status  *string `json:"status,omitempty"`
	CaChain *string `json:"ca_chain,omitempty"`
	Name    *string `json:"name,omitempty"`
}

// HypervisorServer ...
type HypervisorServer struct {
	IP      *string `json:"ip,omitempty"`
	Version *string `json:"version,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// ClusterResource ...
type ClusterResource struct {
	Config            *ConfigClusterSpec `json:"config,omitempty"`
	Network           *ClusterNetwork    `json:"network,omitempty"`
	RunTimeStatusList []*string          `json:"runtime_status_list,omitempty"`
}

// ConfigClusterSpec ...
type ConfigClusterSpec struct {
	GpuDriverVersion              *string                     `json:"gpu_driver_version,omitempty"`
	ClientAuth                    *ClientAuth                 `json:"client_auth,omitempty"`
	AuthorizedPublicKeyList       []*PublicKey                `json:"authorized_public_key_list,omitempty"`
	SoftwareMap                   map[string]interface{}      `json:"software_map,omitempty"`
	EncryptionStatus              string                      `json:"encryption_status,omitempty"`
	RedundancyFactor              *int64                      `json:"redundancy_factor,omitempty"`
	CertificationSigningInfo      *CertificationSigningInfo   `json:"certification_signing_info,omitempty"`
	SupportedInformationVerbosity *string                     `json:"supported_information_verbosity,omitempty"`
	ExternalConfigurations        *ExternalConfigurationsSpec `json:"external_configurations,omitempty"`
	EnabledFeatureList            []*string                   `json:"enabled_feature_list,omitempty"`
	Timezone                      *string                     `json:"timezone,omitempty"`
	OperationMode                 *string                     `json:"operation_mode,omitempty"`
}

// ExternalConfigurationsSpec ...
type ExternalConfigurationsSpec struct {
	CitrixConnectorConfig *CitrixConnectorConfigDetailsSpec `json:"citrix_connector_config,omitempty"`
}

// CitrixConnectorConfigDetailsSpec ...
type CitrixConnectorConfigDetailsSpec struct {
	CitrixVMReferenceList []*Reference                `json:"citrix_connector_config,omitempty"`
	ClientSecret          *string                     `json:"client_secret,omitempty"`
	CustomerID            *string                     `json:"customer_id,omitempty"`
	ClientID              *string                     `json:"client_id,omitempty"`
	ResourceLocation      *CitrixResourceLocationSpec `json:"resource_location,omitempty"`
}

// CitrixResourceLocationSpec ...
type CitrixResourceLocationSpec struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ClusterNetwork ...
type ClusterNetwork struct {
	MasqueradingPort       *int64                  `json:"masquerading_port,omitempty"`
	MasqueradingIP         *string                 `json:"masquerading_ip,omitempty"`
	ExternalIP             *string                 `json:"external_ip,omitempty"`
	HTTPProxyList          []*ClusterNetworkEntity `json:"http_proxy_list,omitempty"`
	SMTPServer             *SMTPServer             `json:"smtp_server,omitempty"`
	NTPServerIPList        []*string               `json:"ntp_server_ip_list,omitempty"`
	ExternalSubnet         *string                 `json:"external_subnet,omitempty"`
	NFSSubnetWhitelist     []*string               `json:"nfs_subnet_whitelist,omitempty"`
	ExternalDataServicesIP *string                 `json:"external_data_services_ip,omitempty"`
	DomainServer           *ClusterDomainServer    `json:"domain_server,omitempty"`
	NameServerIPList       []*string               `json:"name_server_ip_list,omitempty"`
	HTTPProxyWhitelist     []*HTTPProxyWhitelist   `json:"http_proxy_whitelist,omitempty"`
	InternalSubnet         *string                 `json:"internal_subnet,omitempty"`
}

// HTTPProxyWhitelist ...
type HTTPProxyWhitelist struct {
	Target     *string `json:"target,omitempty"`
	TargetType *string `json:"target_type,omitempty"`
}

// ClusterDomainServer ...
type ClusterDomainServer struct {
	Nameserver        *string      `json:"nameserver,omitempty"`
	Name              *string      `json:"name,omitempty"`
	DomainCredentials *Credentials `json:"external_data_services_ip,omitempty"`
}

// SMTPServer ...
type SMTPServer struct {
	Type         *string               `json:"type,omitempty"`
	EmailAddress *string               `json:"email_address,omitempty"`
	Server       *ClusterNetworkEntity `json:"server,omitempty"`
}

// ClusterNetworkEntity ...
type ClusterNetworkEntity struct {
	Credentials   *Credentials `json:"credentials,omitempty"`
	ProxyTypeList []*string    `json:"proxy_type_list,omitempty"`
	Address       *Address     `json:"address,omitempty"`
}

// Credentials ...
type Credentials struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

// VMEfficiencyMap ...
type VMEfficiencyMap struct {
	BullyVMNum           *string `json:"bully_vm_num,omitempty"`
	ConstrainedVMNum     *string `json:"constrained_vm_num,omitempty"`
	DeadVMNum            *string `json:"dead_vm_num,omitempty"`
	InefficientVMNum     *string `json:"inefficient_vm_num,omitempty"`
	OverprovisionedVMNum *string `json:"overprovisioned_vm_num,omitempty"`
}

// ClusterAnalysis ...
type ClusterAnalysis struct {
	VMEfficiencyMap *VMEfficiencyMap `json:"vm_efficiency_map,omitempty"`
}

// CategoryListMetadata All api calls that return a list will have this metadata block as input
type CategoryListMetadata struct {

	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty"`
}

// CategoryKeyStatus represents Category Key Definition.
type CategoryKeyStatus struct {

	// API version.
	APIVersion *string `json:"api_version,omitempty"`

	// Description of the category.
	Description *string `json:"description,omitempty"`

	// Name of the category.
	Name *string `json:"name"`

	// Specifying whether its a system defined category.
	SystemDefined *bool `json:"system_defined,omitempty"`
}

// CategoryKeyListResponse represents the category key list response.
type CategoryKeyListResponse struct {

	// API Version.
	APIVersion *string `json:"api_version,omitempty"`

	Entities []*CategoryKeyStatus `json:"entities,omitempty"`

	Metadata *CategoryListMetadata `json:"metadata,omitempty"`
}

// CategoryKey represents category key definition.
type CategoryKey struct {

	// API version.
	APIVersion *string `json:"api_version,omitempty"`

	// Description of the category.
	Description *string `json:"description,omitempty"`

	// Name of the category.
	Name *string `json:"name"`
}

// CategoryStatus represents The status of a REST API call. Only used when there is a failure to report.
type CategoryStatus struct {
	APIVersion *string `json:"api_version,omitempty"`

	// The HTTP error code.
	Code *int64 `json:"code,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty"`

	MessageList []*MessageResource `json:"message_list,omitempty"`

	State *string `json:"state,omitempty"`
}

// CategoryValueListResponse represents Category Value list response.
type CategoryValueListResponse struct {
	APIVersion *string `json:"api_version,omitempty"`

	Entities []*CategoryValueStatus `json:"entities,omitempty"`

	Metadata *CategoryListMetadata `json:"metadata,omitempty"`
}

// CategoryValueStatus represents Category value definition.
type CategoryValueStatus struct {

	// API version.
	APIVersion *string `json:"api_version,omitempty"`

	// Description of the category value.
	Description *string `json:"description,omitempty"`

	// The name of the category.
	Name *string `json:"name,omitempty"`

	// Specifying whether its a system defined category.
	SystemDefined *bool `json:"system_defined,omitempty"`

	// The value of the category.
	Value *string `json:"value,omitempty"`
}

// CategoryFilter represents A category filter.
type CategoryFilter struct {

	// List of kinds associated with this filter.
	KindList []*string `json:"kind_list,omitempty"`

	// A list of category key and list of values.
	Params map[string][]string `json:"params,omitempty"`

	// The type of the filter being used.
	Type *string `json:"type,omitempty"`
}

// CategoryQueryInput represents Categories query input object.
type CategoryQueryInput struct {

	// API version.
	APIVersion *string `json:"api_version,omitempty"`

	CategoryFilter *CategoryFilter `json:"category_filter,omitempty"`

	// The maximum number of members to return per group.
	GroupMemberCount *int64 `json:"group_member_count,omitempty"`

	// The offset into the total member set to return per group.
	GroupMemberOffset *int64 `json:"group_member_offset,omitempty"`

	// TBD: USED_IN - to get policies in which specified categories are used. APPLIED_TO - to get entities attached to
	// specified categories.
	UsageType *string `json:"usage_type,omitempty"`
}

// CategoryQueryResponseMetadata represents Response metadata.
type CategoryQueryResponseMetadata struct {

	// The maximum number of records to return per group.
	GroupMemberCount *int64 `json:"group_member_count,omitempty"`

	// The offset into the total records set to return per group.
	GroupMemberOffset *int64 `json:"group_member_offset,omitempty"`

	// Total number of matched results.
	TotalMatches *int64 `json:"total_matches,omitempty"`

	// TBD: USED_IN - to get policies in which specified categories are used. APPLIED_TO - to get entities attached to specified categories.
	UsageType *string `json:"usage_type,omitempty"`
}

// EntityReference Reference to an entity.
type EntityReference struct {

	// Categories for the entity.
	Categories map[string]string `json:"categories,omitempty"`

	// Kind of the reference.
	Kind *string `json:"kind,omitempty"`

	// Name of the entity.
	Name *string `json:"name,omitempty"`

	// The type of filter being used. (Options : CATEGORIES_MATCH_ALL , CATEGORIES_MATCH_ANY)
	Type *string `json:"type,omitempty"`

	// UUID of the entity.
	UUID *string `json:"uuid,omitempty"`
}

// CategoryQueryResponseResults ...
type CategoryQueryResponseResults struct {

	// List of entity references.
	EntityAnyReferenceList []*EntityReference `json:"entity_any_reference_list,omitempty"`

	// Total number of filtered results.
	FilteredEntityCount *int64 `json:"filtered_entity_count,omitempty"`

	// The entity kind.
	Kind *string `json:"kind,omitempty"`

	// Total number of the matched results.
	TotalEntityCount *int64 `json:"total_entity_count,omitempty"`
}

// CategoryQueryResponse represents Categories query response object.
type CategoryQueryResponse struct {

	// API version.
	APIVersion *string `json:"api_version,omitempty"`

	Metadata *CategoryQueryResponseMetadata `json:"metadata,omitempty"`

	Results []*CategoryQueryResponseResults `json:"results,omitempty"`
}

// CategoryValue represents Category value definition.
type CategoryValue struct {

	// API version.
	APIVersion *string `json:"api_version,omitempty"`

	// Description of the category value.
	Description *string `json:"description,omitempty" `

	// Value for the category.
	Value *string `json:"value,omitempty"`
}

// PortRange represents Range of TCP/UDP ports.
type PortRange struct {
	EndPort *int64 `json:"end_port,omitempty"`

	StartPort *int64 `json:"start_port,omitempty"`
}

// IPSubnet IP subnet provided as an address and prefix length.
type IPSubnet struct {

	// IPV4 address.
	IP *string `json:"ip,omitempty"`

	PrefixLength *int64 `json:"prefix_length,omitempty"`
}

// NetworkRuleIcmpTypeCodeList ..
type NetworkRuleIcmpTypeCodeList struct {
	Code *int64 `json:"code,omitempty"`

	Type *int64 `json:"type,omitempty"`
}

// NetworkRule ...
type NetworkRule struct {

	// Timestamp of expiration time.
	ExpirationTime *string `json:"expiration_time,omitempty"`

	// The set of categories that matching VMs need to have.
	Filter *CategoryFilter `json:"filter,omitempty"`

	// List of ICMP types and codes allowed by this rule.
	IcmpTypeCodeList []*NetworkRuleIcmpTypeCodeList `json:"icmp_type_code_list,omitempty"`

	IPSubnet *IPSubnet `json:"ip_subnet,omitempty"`

	NetworkFunctionChainReference *Reference `json:"network_function_chain_reference,omitempty"`

	// The set of categories that matching VMs need to have.
	PeerSpecificationType *string `json:"peer_specification_type,omitempty"`

	// Select a protocol to allow.  Multiple protocols can be allowed by repeating network_rule object.  If a protocol
	// is not configured in the network_rule object then it is allowed.
	Protocol *string `json:"protocol,omitempty"`

	// List of TCP ports that are allowed by this rule.
	TCPPortRangeList []*PortRange `json:"tcp_port_range_list,omitempty"`

	// List of UDP ports that are allowed by this rule.
	UDPPortRangeList []*PortRange `json:"udp_port_range_list,omitempty"`
}

// TargetGroup ...
type TargetGroup struct {

	// Default policy for communication within target group.
	DefaultInternalPolicy *string `json:"default_internal_policy,omitempty"`

	// The set of categories that matching VMs need to have.
	Filter *CategoryFilter `json:"filter,omitempty"`

	// Way to identify the object for which rule is applied.
	PeerSpecificationType *string `json:"peer_specification_type,omitempty"`
}

// NetworkSecurityRuleResourcesRule These rules are used for quarantining suspected VMs. Target group is a required
// attribute.  Empty inbound_allow_list will not allow anything into target group. Empty outbound_allow_list will allow
// everything from target group.
type NetworkSecurityRuleResourcesRule struct {
	Action            *string        `json:"action,omitempty"`             // Type of action.
	InboundAllowList  []*NetworkRule `json:"inbound_allow_list,omitempty"` //
	OutboundAllowList []*NetworkRule `json:"outbound_allow_list,omitempty"`
	TargetGroup       *TargetGroup   `json:"target_group,omitempty"`
}

// NetworkSecurityRuleIsolationRule These rules are used for environmental isolation.
type NetworkSecurityRuleIsolationRule struct {
	Action             *string         `json:"action,omitempty"`               // Type of action.
	FirstEntityFilter  *CategoryFilter `json:"first_entity_filter,omitempty"`  // The set of categories that matching VMs need to have.
	SecondEntityFilter *CategoryFilter `json:"second_entity_filter,omitempty"` // The set of categories that matching VMs need to have.
}

// NetworkSecurityRuleResources ...
type NetworkSecurityRuleResources struct {
	AppRule        *NetworkSecurityRuleResourcesRule `json:"app_rule,omitempty"`
	IsolationRule  *NetworkSecurityRuleIsolationRule `json:"isolation_rule,omitempty"`
	QuarantineRule *NetworkSecurityRuleResourcesRule `json:"quarantine_rule,omitempty"`
}

// NetworkSecurityRule ...
type NetworkSecurityRule struct {
	Description *string                       `json:"description"`
	Name        *string                       `json:"name,omitempty"`
	Resources   *NetworkSecurityRuleResources `json:"resources,omitempty" `
}

// Metadata Metadata The kind metadata
type Metadata struct {
	LastUpdateTime   *time.Time        `json:"last_update_time,omitempty"`  //
	Kind             *string           `json:"kind"`                        //
	UUID             *string           `json:"uuid,omitempty"`              //
	ProjectReference *Reference        `json:"project_reference,omitempty"` // project reference
	CreationTime     *time.Time        `json:"creation_time,omitempty"`
	SpecVersion      *int64            `json:"spec_version,omitempty"`
	SpecHash         *string           `json:"spec_hash,omitempty"`
	OwnerReference   *Reference        `json:"owner_reference,omitempty"`
	Categories       map[string]string `json:"categories,omitempty"`
	Name             *string           `json:"name,omitempty"`
}

// NetworkSecurityRuleIntentInput An intentful representation of a network_security_rule
type NetworkSecurityRuleIntentInput struct {
	APIVersion *string              `json:"api_version,omitempty"`
	Metadata   *Metadata            `json:"metadata"`
	Spec       *NetworkSecurityRule `json:"spec"`
}

// NetworkSecurityRuleDefStatus ... Network security rule status
type NetworkSecurityRuleDefStatus struct {
	AppRule          *NetworkSecurityRuleResourcesRule `json:"app_rule,omitempty"`
	IsolationRule    *NetworkSecurityRuleIsolationRule `json:"isolation_rule,omitempty"`
	QuarantineRule   *NetworkSecurityRuleResourcesRule `json:"quarantine_rule,omitempty"`
	State            *string                           `json:"state,omitmepty"`
	ExecutionContext *ExecutionContext                 `json:"execution_context,omitempty"`
}

// NetworkSecurityRuleIntentResponse Response object for intentful operations on a network_security_rule
type NetworkSecurityRuleIntentResponse struct {
	APIVersion *string                      `json:"api_version,omitempty"`
	Metadata   *Metadata                    `json:"metadata"`
	Spec       *NetworkSecurityRule         `json:"spec,omitempty"`
	Status     NetworkSecurityRuleDefStatus `json:"status,omitempty" bson:"status,omitempty"`
}

// NetworkSecurityRuleStatus The status of a REST API call. Only used when there is a failure to report.
type NetworkSecurityRuleStatus struct {
	APIVersion  *string            `json:"api_version,omitempty"` //
	Code        *int64             `json:"code,omitempty"`        // The HTTP error code.
	Kind        *string            `json:"kind,omitempty"`        // The kind name
	MessageList []*MessageResource `json:"message_list,omitempty"`
	State       *string            `json:"state,omitempty"`
}

// ListMetadata All api calls that return a list will have this metadata block as input
type ListMetadata struct {
	Filter        *string `json:"filter,omitempty"`         // The filter in FIQL syntax used for the results.
	Kind          *string `json:"kind,omitempty"`           // The kind name
	Length        *int64  `json:"length,omitempty"`         // The number of records to retrieve relative to the offset
	Offset        *int64  `json:"offset,omitempty"`         // Offset from the start of the entity list
	SortAttribute *string `json:"sort_attribute,omitempty"` // The attribute to perform sort on
	SortOrder     *string `json:"sort_order,omitempty"`     // The sort order in which results are returned
}

// ListMetadataOutput All api calls that return a list will have this metadata block
type ListMetadataOutput struct {
	Filter        *string `json:"filter,omitempty"`         // The filter used for the results
	Kind          *string `json:"kind,omitempty"`           // The kind name
	Length        *int64  `json:"length,omitempty"`         // The number of records retrieved relative to the offset
	Offset        *int64  `json:"offset,omitempty"`         // Offset from the start of the entity list
	SortAttribute *string `json:"sort_attribute,omitempty"` // The attribute to perform sort on
	SortOrder     *string `json:"sort_order,omitempty"`     // The sort order in which results are returned
	TotalMatches  *int64  `json:"total_matches,omitempty"`  // Total matches found
}

// NetworkSecurityRuleIntentResource ... Response object for intentful operations on a network_security_rule
type NetworkSecurityRuleIntentResource struct {
	APIVersion *string                       `json:"api_version,omitempty"`
	Metadata   *Metadata                     `json:"metadata,omitempty"`
	Spec       *NetworkSecurityRule          `json:"spec,omitempty"`
	Status     *NetworkSecurityRuleDefStatus `json:"status,omitempty"`
}

// NetworkSecurityRuleListIntentResponse Response object for intentful operation of network_security_rules
type NetworkSecurityRuleListIntentResponse struct {
	APIVersion string                               `json:"api_version"`
	Entities   []*NetworkSecurityRuleIntentResource `json:"entities,omitempty" bson:"entities,omitempty"`
	Metadata   *ListMetadataOutput                  `json:"metadata"`
}

// VolumeGroupInput Represents the request body for create volume_grop request
type VolumeGroupInput struct {
	APIVersion *string      `json:"api_version,omitempty"` // default 3.1.0
	Metadata   *Metadata    `json:"metadata,omitempty"`    // The volume_group kind metadata.
	Spec       *VolumeGroup `json:"spec,omitempty"`        // Volume group input spec.
}

// VolumeGroup Represents volume group input spec.
type VolumeGroup struct {
	Name        *string               `json:"name"`                  // Volume Group name (required)
	Description *string               `json:"description,omitempty"` // Volume Group description.
	Resources   *VolumeGroupResources `json:"resources"`             // Volume Group resources.
}

// VolumeGroupResources Represents the volume group resources
type VolumeGroupResources struct {
	FlashMode         *string         `json:"flash_mode,omitempty"`          // Flash Mode, if enabled all disks of the VG are pinned to SSD
	FileSystemType    *string         `json:"file_system_type,omitempty"`    // File system to be used for volume
	SharingStatus     *string         `json:"sharing_status,omitempty"`      // Whether the VG can be shared across multiple iSCSI initiators
	AttachmentList    []*VMAttachment `json:"attachment_list,omitempty"`     // VMs attached to volume group.
	DiskList          []*VGDisk       `json:"disk_list,omitempty"`           // VGDisk Volume group disk specification.
	IscsiTargetPrefix *string         `json:"iscsi_target_prefix,omitempty"` // iSCSI target prefix-name.
}

// VMAttachment VMs attached to volume group.
type VMAttachment struct {
	VMReference        *Reference `json:"vm_reference"`         // Reference to a kind
	IscsiInitiatorName *string    `json:"iscsi_initiator_name"` // Name of the iSCSI initiator of the workload outside Nutanix cluster.
}

// VGDisk Volume group disk specification.
type VGDisk struct {
	VmdiskUUID           *string    `json:"vmdisk_uuid"`            // The UUID of this volume disk
	Index                *int64     `json:"index"`                  // Index of the volume disk in the group.
	DataSourceReference  *Reference `json:"data_source_reference"`  // Reference to a kind
	DiskSizeMib          *int64     `json:"disk_size_mib"`          // Size of the disk in MiB.
	StorageContainerUUID *string    `json:"storage_container_uuid"` // Container UUID on which to create the disk.
}

// VolumeGroupResponse Response object for intentful operations on a volume_group
type VolumeGroupResponse struct {
	APIVersion *string               `json:"api_version"`      //
	Metadata   *Metadata             `json:"metadata"`         // The volume_group kind metadata
	Spec       *VolumeGroup          `json:"spec,omitempty"`   // Volume group input spec.
	Status     *VolumeGroupDefStatus `json:"status,omitempty"` // Volume group configuration.
}

// VolumeGroupDefStatus  Volume group configuration.
type VolumeGroupDefStatus struct {
	State       *string               `json:"state"`        // The state of the volume group entity.
	MessageList []*MessageResource    `json:"message_list"` // Volume group message list.
	Name        *string               `json:"name"`         // Volume group name.
	Resources   *VolumeGroupResources `json:"resources"`    // Volume group resources.
	Description *string               `json:"description"`  // Volume group description.
}

// VolumeGroupListResponse Response object for intentful operation of volume_groups
type VolumeGroupListResponse struct {
	APIVersion *string                `json:"api_version"`
	Entities   []*VolumeGroupResponse `json:"entities,omitempty"`
	Metadata   *ListMetadataOutput    `json:"metadata"`
}

// TasksResponse ...
type TasksResponse struct {
	Status               *string      `json:"status,omitempty"`
	LastUpdateTime       *time.Time   `json:"last_update_time,omitempty"`
	LogicalTimestamp     *int64       `json:"logical_timestamp,omitempty"`
	EntityReferenceList  []*Reference `json:"entity_reference_list,omitempty"`
	StartTime            *time.Time   `json:"start_time,omitempty"`
	CreationTime         *time.Time   `json:"creation_time,omitempty"`
	ClusterReference     *Reference   `json:"cluster_reference,omitempty"`
	SubtaskReferenceList []*Reference `json:"subtask_reference_list,omitempty"`
	CompletionTime       *time.Time   `json:"completion_timev"`
	ProgressMessage      *string      `json:"progress_message,omitempty"`
	OperationType        *string      `json:"operation_type,omitempty"`
	PercentageComplete   *int64       `json:"percentage_complete,omitempty"`
	APIVersion           *string      `json:"api_version,omitempty"`
	UUID                 *string      `json:"uuid,omitempty"`
	ErrorDetail          *string      `json:"error_detail,omitempty"`
}

// DeleteResponse ...
type DeleteResponse struct {
	Status     *DeleteStatus `json:"status"`
	Spec       string        `json:"spec"`
	APIVersion string        `json:"api_version"`
	Metadata   *Metadata     `json:"metadata"`
}

// DeleteStatus ...
type DeleteStatus struct {
	State            string            `json:"state"`
	ExecutionContext *ExecutionContext `json:"execution_context"`
}
