package virtualmachine

// CloudInit struct
type CloudInit struct {

MetaData string `json:"meta_data,omitempty"bson:"meta_data,omitempty"`
UserData string `json:"user_data,omitempty"bson:"user_data,omitempty"`

}