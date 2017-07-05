package virtualmachine

// Spec struct
type Spec struct {

ClusterReference ClusterReference `json:"cluster_reference,omitempty"bson:"cluster_reference,omitempty"`
Description string `json:"description,omitempty"bson:"description,omitempty"`
Name string `json:"name,omitempty"bson:"name,omitempty"`
Resources Resources `json:"resources,omitempty"bson:"resources,omitempty"`

}