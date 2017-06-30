#\DisksApi

##DisksListPost
//  Get a list of Disks
var api_instance := nutanix.DisksApi
getEntitiesRequest := nutanix.DiskListMetadata() // DiskListMetadata |
diskslistpost_response, api_response, err := api_instance.DisksListPost(getEntitiesRequest)

##DisksUuidDelete
//  Delete a Disk
var api_instance := nutanix.DisksApi
uuid := "uuid_example" // string | The UUID of the entity
disksuuiddelete_response, api_response, err := api_instance.DisksUuidDelete(uuid)

##DisksUuidGet
//  Get a Disk
var api_instance := nutanix.DisksApi
uuid := "uuid_example" // string | The UUID of the entity
disksuuidget_response, api_response, err := api_instance.DisksUuidGet(uuid)
