#\SubnetApi

##SubnetsListPost
//  Get a list of subnets
var api_instance := nutanix.SubnetApi
getEntitiesRequest := nutanix.SubnetListMetadata() // SubnetListMetadata |
subnetslistpost_response, api_response, err := api_instance.SubnetsListPost(getEntitiesRequest)

##SubnetsPost
//  Create a subnet
var api_instance := nutanix.SubnetApi
body := nutanix.SubnetIntentInput() // SubnetIntentInput |
subnetspost_response, api_response, err := api_instance.SubnetsPost(body)

##SubnetsUuidDelete
//  Delete a subnet
var api_instance := nutanix.SubnetApi
uuid := "uuid_example" // string | The UUID of the entity
subnetsuuiddelete_response, api_response, err := api_instance.SubnetsUuidDelete(uuid)

##SubnetsUuidGet
//  Get a subnet
var api_instance := nutanix.SubnetApi
uuid := "uuid_example" // string | The UUID of the entity
subnetsuuidget_response, api_response, err := api_instance.SubnetsUuidGet(uuid)

##SubnetsUuidPut
//  Update a subnet
var api_instance := nutanix.SubnetApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.SubnetIntentInput() // SubnetIntentInput |
subnetsuuidput_response, api_response, err := api_instance.SubnetsUuidPut(uuid, body)
