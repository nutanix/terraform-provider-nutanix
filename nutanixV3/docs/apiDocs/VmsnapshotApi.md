#\VmsnapshotApi

##VmSnapshotsListPost
//  Get kind snapshots
var api_instance := nutanix.VmsnapshotApi
getEntitiesRequest := nutanix.VmSnapshotListMetadata() // VmSnapshotListMetadata |
vmsnapshotslistpost_response, api_response, err := api_instance.VmSnapshotsListPost(getEntitiesRequest)

##VmSnapshotsPost
//  Create kind snapshot
var api_instance := nutanix.VmsnapshotApi
body := nutanix.VmSnapshotIntentInput() // VmSnapshotIntentInput |
vmsnapshotspost_response, api_response, err := api_instance.VmSnapshotsPost(body)

##VmSnapshotsUuidDelete
//  Delete kind snapshot
var api_instance := nutanix.VmsnapshotApi
uuid := "uuid_example" // string | The UUID of the entity
vmsnapshotsuuiddelete_response, api_response, err := api_instance.VmSnapshotsUuidDelete(uuid)

##VmSnapshotsUuidGet
//  Get kind snapshot
var api_instance := nutanix.VmsnapshotApi
uuid := "uuid_example" // string | The UUID of the entity
vmsnapshotsuuidget_response, api_response, err := api_instance.VmSnapshotsUuidGet(uuid)

##VmSnapshotsUuidPut
//  Update kind snapshot
var api_instance := nutanix.VmsnapshotApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.VmSnapshotIntentInput() // VmSnapshotIntentInput |
vmsnapshotsuuidput_response, api_response, err := api_instance.VmSnapshotsUuidPut(uuid, body)
