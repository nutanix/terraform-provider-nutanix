package virtualmachine

// BootConfig struct
type BootConfig struct {

DiskAddress DiskAddress `json:"disk_address,omitempty"bson:"disk_address,omitempty"`
MacAddress string `json:"mac_address,omitempty"bson:"mac_address,omitempty"`

}