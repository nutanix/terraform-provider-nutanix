package virtualmachine

// Status struct
type Status struct {

ClusterReference ClusterReference `json:"cluster_reference,omitempty"bson:"cluster_reference,omitempty"`
Description string `json:"description,omitempty"bson:"description,omitempty"`
Message string `json:"message,omitempty"bson:"message,omitempty"`
Name string `json:"name,omitempty"bson:"name,omitempty"`
Reason string `json:"reason,omitempty"bson:"reason,omitempty"`
Resources Resources `json:"resources,omitempty"bson:"resources,omitempty"`
State string `json:"state,omitempty"bson:"state,omitempty"`

}