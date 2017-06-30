#\DirectoryserviceApi

##DirectoryServicesListPost
//  Get a list of directory service configurations
var api_instance := nutanix.DirectoryserviceApi
getEntitiesRequest := nutanix.DirectoryServiceListMetadata() // DirectoryServiceListMetadata | 
directoryserviceslistpost_response, api_response, err := api_instance.DirectoryServicesListPost(getEntitiesRequest)

##DirectoryServicesPost
//  Create a directory service configuration
var api_instance := nutanix.DirectoryserviceApi
body := nutanix.DirectoryServiceIntentInput() // DirectoryServiceIntentInput | 
directoryservicespost_response, api_response, err := api_instance.DirectoryServicesPost(body)

##DirectoryServicesUuidDelete
//  Delete a directory service configuration
var api_instance := nutanix.DirectoryserviceApi
uuid := "uuid_example" // string | The UUID of the entity
directoryservicesuuiddelete_response, api_response, err := api_instance.DirectoryServicesUuidDelete(uuid)

##DirectoryServicesUuidGet
//  Get a directory service configuration
var api_instance := nutanix.DirectoryserviceApi
uuid := "uuid_example" // string | The UUID of the entity
directoryservicesuuidget_response, api_response, err := api_instance.DirectoryServicesUuidGet(uuid)

##DirectoryServicesUuidPut
//  Update a directory service configuration
var api_instance := nutanix.DirectoryserviceApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.DirectoryServiceIntentInput() // DirectoryServiceIntentInput | 
directoryservicesuuidput_response, api_response, err := api_instance.DirectoryServicesUuidPut(uuid, body)

##DirectoryServicesUuidSearchPost
//  Search user or group in the directory service.
var api_instance := nutanix.DirectoryserviceApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.DirectoryServiceSearchMetadata() // DirectoryServiceSearchMetadata | 
directoryservicesuuidsearchpost_response, api_response, err := api_instance.DirectoryServicesUuidSearchPost(uuid, body)

