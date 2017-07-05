package virtualmachine

// GPUList struct
type GPUList struct {

DeviceID int `json:"device_id,omitempty"bson:"device_id,omitempty"`
Mode string `json:"mode,omitempty"bson:"mode,omitempty"`
Vendor string `json:"vendor,omitempty"bson:"vendor,omitempty"`

}