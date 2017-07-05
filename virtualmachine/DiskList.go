package virtualmachine

// DiskList struct
type DiskList struct {

DataSourceReference DataSourceReference `json:"data_source_reference,omitempty"bson:"data_source_reference,omitempty"`
DeviceProperties DeviceProperties `json:"device_properties,omitempty"bson:"device_properties,omitempty"`
DiskSizeMib int `json:"disk_size_mib,omitempty"bson:"disk_size_mib,omitempty"`
UUID string `json:"uuid,omitempty"bson:"uuid,omitempty"`

}