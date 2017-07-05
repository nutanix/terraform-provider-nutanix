package virtualmachine

// VMIntentInput struct
type VMIntentInput struct {

APIVersion string `json:"api_version,omitempty"bson:"api_version,omitempty"`
Metadata Metadata `json:"metadata,omitempty"bson:"metadata,omitempty"`
Spec Spec `json:"spec,omitempty"bson:"spec,omitempty"`

}