package virtualmachine

// HostReference struct
type HostReference struct {

Kind string `json:"kind,omitempty"bson:"kind,omitempty"`
Name string `json:"name,omitempty"bson:"name,omitempty"`
UUID string `json:"uuid,omitempty"bson:"uuid,omitempty"`

}