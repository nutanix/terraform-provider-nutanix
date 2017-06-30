#\VmApi

##VmsListPost
//  Get a list of VMs
var api_instance := nutanix.VmApi
getEntitiesRequest := nutanix.VmListMetadata() // VmListMetadata | 
vmslistpost_response, api_response, err := api_instance.VmsListPost(getEntitiesRequest)

##VmsPost
//  Create a VM
var api_instance := nutanix.VmApi
body := nutanix.VmIntentInput() // VmIntentInput | 
vmspost_response, api_response, err := api_instance.VmsPost(body)

##VmsUuidDelete
//  Delete a VM
var api_instance := nutanix.VmApi
uuid := "uuid_example" // string | The UUID of the entity
vmsuuiddelete_response, api_response, err := api_instance.VmsUuidDelete(uuid)

##VmsUuidGet
//  Get a VM
var api_instance := nutanix.VmApi
uuid := "uuid_example" // string | The UUID of the entity
vmsuuidget_response, api_response, err := api_instance.VmsUuidGet(uuid)

##VmsUuidPut
//  Update a VM
var api_instance := nutanix.VmApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.VmIntentInput() // VmIntentInput | 
vmsuuidput_response, api_response, err := api_instance.VmsUuidPut(uuid, body)

