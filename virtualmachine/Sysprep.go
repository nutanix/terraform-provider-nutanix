package virtualmachine

// Sysprep struct
type Sysprep struct {

InstallType string `json:"install_type,omitempty"bson:"install_type,omitempty"`
UnattendXML string `json:"unattend_xml,omitempty"bson:"unattend_xml,omitempty"`

}