package virtualmachine

// Metadata struct
type Metadata struct {

Categories map[string]string `json:"categories,omitempty"bson:"categories,omitempty"`
CreationTime string `json:"creation_time,omitempty"bson:"creation_time,omitempty"`
EntityVersion int `json:"entity_version,omitempty"bson:"entity_version,omitempty"`
Kind string `json:"kind,omitempty"bson:"kind,omitempty"`
LastUpdateTime string `json:"last_update_time,omitempty"bson:"last_update_time,omitempty"`
Name string `json:"name,omitempty"bson:"name,omitempty"`
OwnerReference OwnerReference `json:"owner_reference,omitempty"bson:"owner_reference,omitempty"`
ParentReference string `json:"parent_reference,omitempty"bson:"parent_reference,omitempty"`
UUID string `json:"uuid,omitempty"bson:"uuid,omitempty"`

}