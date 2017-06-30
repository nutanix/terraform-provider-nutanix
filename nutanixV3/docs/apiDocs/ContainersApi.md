#\ContainersApi

##ContainersListPost
//  Get a list of Containers
var api_instance := nutanix.ContainersApi
getEntitiesRequest := nutanix.ContainerListMetadata() // ContainerListMetadata |
containerslistpost_response, api_response, err := api_instance.ContainersListPost(getEntitiesRequest)

##ContainersPost
//  Create a Container
var api_instance := nutanix.ContainersApi
body := nutanix.ContainerIntentInput() // ContainerIntentInput |
containerspost_response, api_response, err := api_instance.ContainersPost(body)

##ContainersUuidDelete
//  Delete a Container
var api_instance := nutanix.ContainersApi
uuid := "uuid_example" // string | The UUID of the entity
containersuuiddelete_response, api_response, err := api_instance.ContainersUuidDelete(uuid)

##ContainersUuidGet
//  Get a Container
var api_instance := nutanix.ContainersApi
uuid := "uuid_example" // string | The UUID of the entity
containersuuidget_response, api_response, err := api_instance.ContainersUuidGet(uuid)

##ContainersUuidPut
//  Update a Container
var api_instance := nutanix.ContainersApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.ContainerIntentInput() // ContainerIntentInput | Intent Spec for Container.
containersuuidput_response, api_response, err := api_instance.ContainersUuidPut(uuid, body)
