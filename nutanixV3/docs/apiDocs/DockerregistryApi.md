#\DockerregistryApi

##DockerRegistriesListPost
//  List all docker registries
var api_instance := nutanix.DockerregistryApi
getEntitiesRequest := nutanix.DockerRegistryListMetadata() // DockerRegistryListMetadata | 
dockerregistrieslistpost_response, api_response, err := api_instance.DockerRegistriesListPost(getEntitiesRequest)

##DockerRegistriesPost
//  Create a docker registry
var api_instance := nutanix.DockerregistryApi
body := nutanix.DockerRegistryIntentInput() // DockerRegistryIntentInput | Docker registry spec
dockerregistriespost_response, api_response, err := api_instance.DockerRegistriesPost(body)

##DockerRegistriesUuidDelete
//  Deletes a Docker registry
var api_instance := nutanix.DockerregistryApi
uuid := "uuid_example" // string | The UUID of the entity
dockerregistriesuuiddelete_response, api_response, err := api_instance.DockerRegistriesUuidDelete(uuid)

##DockerRegistriesUuidGet
//  Retrieve a Docker registry
var api_instance := nutanix.DockerregistryApi
uuid := "uuid_example" // string | The UUID of the entity
dockerregistriesuuidget_response, api_response, err := api_instance.DockerRegistriesUuidGet(uuid)

##DockerRegistriesUuidPut
//  Update a docker registry
var api_instance := nutanix.DockerregistryApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.DockerRegistryIntentInput() // DockerRegistryIntentInput | Docker registry spec
dockerregistriesuuidput_response, api_response, err := api_instance.DockerRegistriesUuidPut(uuid, body)

##DockerRegistriesUuidSearchListPost
//  Searches docker containers for specified registry
var api_instance := nutanix.DockerregistryApi
uuid := "uuid_example" // string | The UUID of the entity
getEntitiesRequest := nutanix.DockerRegistryListMetadata() // DockerRegistryListMetadata | 
dockerregistriesuuidsearchlistpost_response, api_response, err := api_instance.DockerRegistriesUuidSearchListPost(uuid, getEntitiesRequest)

