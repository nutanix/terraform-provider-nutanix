#\GroupsApi

##GroupsPost
//  Get Entities.
var api_instance := nutanix.GroupsApi
getEntitiesRequest := nutanix.GroupsGetEntitiesRequest() // GroupsGetEntitiesRequest | 
groupspost_response, api_response, err := api_instance.GroupsPost(getEntitiesRequest)

