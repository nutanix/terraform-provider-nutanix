package virtualmachine

// DeviceProperties struct
type DeviceProperties struct {

DeviceType string `json:"device_type,omitempty"bson:"device_type,omitempty"`
DiskAddress DiskAddress `json:"disk_address,omitempty"bson:"disk_address,omitempty"`

}