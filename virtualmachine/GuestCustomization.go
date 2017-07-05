package virtualmachine

// GuestCustomization struct
type GuestCustomization struct {

CloudInit CloudInit `json:"cloud_init,omitempty"bson:"cloud_init,omitempty"`
Sysprep Sysprep `json:"sysprep,omitempty"bson:"sysprep,omitempty"`

}