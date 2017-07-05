package virtualmachine

// Resources struct
type Resources struct {

BootConfig BootConfig `json:"boot_config,omitempty"bson:"boot_config,omitempty"`
DiskList []DiskList `json:"disk_list,omitempty"bson:"disk_list,omitempty"`
GPUList []GPUList `json:"gpu_list,omitempty"bson:"gpu_list,omitempty"`
GuestCustomization GuestCustomization `json:"guest_customization,omitempty"bson:"guest_customization,omitempty"`
MemorySizeMb int `json:"memory_size_mb,omitempty"bson:"memory_size_mb,omitempty"`
NicList []NicList `json:"nic_list,omitempty"bson:"nic_list,omitempty"`
NumSockets int `json:"num_sockets,omitempty"bson:"num_sockets,omitempty"`
NumVcpusPerSocket int `json:"num_vcpus_per_socket,omitempty"bson:"num_vcpus_per_socket,omitempty"`
ParentReference ParentReference `json:"parent_reference,omitempty"bson:"parent_reference,omitempty"`
PowerState string `json:"power_state,omitempty"bson:"power_state,omitempty"`

}