#\VolumegroupApi

##VolumeGroupsListPost
//  Retrieves all volume groups.
var api_instance := nutanix.VolumegroupApi
getEntitiesRequest := nutanix.VolumeGroupListMetadata() // VolumeGroupListMetadata |
volumegroupslistpost_response, api_response, err := api_instance.VolumeGroupsListPost(getEntitiesRequest)

##VolumeGroupsPost
//  Creates a volume group
var api_instance := nutanix.VolumegroupApi
body := nutanix.VolumeGroupIntentInput() // VolumeGroupIntentInput | Volume group object.
volumegroupspost_response, api_response, err := api_instance.VolumeGroupsPost(body)

##VolumeGroupsUuidDelete
//  Deletes a volume group
var api_instance := nutanix.VolumegroupApi
uuid := "uuid_example" // string | The UUID of the entity
volumegroupsuuiddelete_response, api_response, err := api_instance.VolumeGroupsUuidDelete(uuid)

##VolumeGroupsUuidGet
//  Retrieves specified volume group.
var api_instance := nutanix.VolumegroupApi
uuid := "uuid_example" // string | The UUID of the entity
volumegroupsuuidget_response, api_response, err := api_instance.VolumeGroupsUuidGet(uuid)

##VolumeGroupsUuidPut
//  Updates specified volume group
var api_instance := nutanix.VolumegroupApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.VolumeGroupIntentInput() // VolumeGroupIntentInput | Volume group object.
volumegroupsuuidput_response, api_response, err := api_instance.VolumeGroupsUuidPut(uuid, body)
