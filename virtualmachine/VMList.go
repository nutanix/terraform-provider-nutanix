package virtualmachine

// VMList struct
type VMList struct {

APIVersion string `json:"api_version,omitempty"bson:"api_version,omitempty"`
Entities []Entities `json:"entities,omitempty"bson:"entities,omitempty"`
Metadata Metadata `json:"metadata,omitempty"bson:"metadata,omitempty"`

}