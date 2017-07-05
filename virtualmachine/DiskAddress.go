package virtualmachine

// DiskAddress struct
type DiskAddress struct {

AdapterType string `json:"adapter_type,omitempty"bson:"adapter_type,omitempty"`
DeviceIndex int `json:"device_index,omitempty"bson:"device_index,omitempty"`

}