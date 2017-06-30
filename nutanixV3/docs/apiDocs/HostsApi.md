#\HostsApi

##HostsListPost
//  Get a list of Hosts
var api_instance := nutanix.HostsApi
getEntitiesRequest := nutanix.HostListMetadata() // HostListMetadata | 
hostslistpost_response, api_response, err := api_instance.HostsListPost(getEntitiesRequest)

##HostsUuidDelete
//  Delete a Host
var api_instance := nutanix.HostsApi
uuid := "uuid_example" // string | The UUID of the entity
hostsuuiddelete_response, api_response, err := api_instance.HostsUuidDelete(uuid)

##HostsUuidGet
//  Get a Host
var api_instance := nutanix.HostsApi
uuid := "uuid_example" // string | The UUID of the entity
hostsuuidget_response, api_response, err := api_instance.HostsUuidGet(uuid)

##HostsUuidPut
//  Update a Host
var api_instance := nutanix.HostsApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.HostIntentInput() // HostIntentInput | Intent Spec of Host.
hostsuuidput_response, api_response, err := api_instance.HostsUuidPut(uuid, body)

