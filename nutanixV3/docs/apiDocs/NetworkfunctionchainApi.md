#\NetworkfunctionchainApi

##NetworkFunctionChainsListPost
//  Get a list of Network Function Chains
var api_instance := nutanix.NetworkfunctionchainApi
getEntitiesRequest := nutanix.NetworkFunctionChainListMetadata() // NetworkFunctionChainListMetadata | 
networkfunctionchainslistpost_response, api_response, err := api_instance.NetworkFunctionChainsListPost(getEntitiesRequest)

##NetworkFunctionChainsPost
//  Create a Network Function Chain
var api_instance := nutanix.NetworkfunctionchainApi
body := nutanix.NetworkFunctionChainIntentInput() // NetworkFunctionChainIntentInput | 
networkfunctionchainspost_response, api_response, err := api_instance.NetworkFunctionChainsPost(body)

##NetworkFunctionChainsUuidDelete
//  Delete a Network Function Chain
var api_instance := nutanix.NetworkfunctionchainApi
uuid := "uuid_example" // string | The UUID of the entity
networkfunctionchainsuuiddelete_response, api_response, err := api_instance.NetworkFunctionChainsUuidDelete(uuid)

##NetworkFunctionChainsUuidGet
//  Get a Network Function Chain
var api_instance := nutanix.NetworkfunctionchainApi
uuid := "uuid_example" // string | The UUID of the entity
networkfunctionchainsuuidget_response, api_response, err := api_instance.NetworkFunctionChainsUuidGet(uuid)

##NetworkFunctionChainsUuidPut
//  Update a Network Function Chain
var api_instance := nutanix.NetworkfunctionchainApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.NetworkFunctionChainIntentInput() // NetworkFunctionChainIntentInput | 
networkfunctionchainsuuidput_response, api_response, err := api_instance.NetworkFunctionChainsUuidPut(uuid, body)

