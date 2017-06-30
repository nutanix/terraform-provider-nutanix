#\RoleApi

##RolesListPost
//  List the roles.
var api_instance := nutanix.RoleApi
getEntitiesRequest := nutanix.RoleListMetadata() // RoleListMetadata | 
roleslistpost_response, api_response, err := api_instance.RolesListPost(getEntitiesRequest)

##RolesPost
//  Create a role.
var api_instance := nutanix.RoleApi
body := nutanix.RoleIntentInput() // RoleIntentInput | 
rolespost_response, api_response, err := api_instance.RolesPost(body)

##RolesUuidDelete
//  Delete a role.
var api_instance := nutanix.RoleApi
uuid := "uuid_example" // string | The UUID of the entity
rolesuuiddelete_response, api_response, err := api_instance.RolesUuidDelete(uuid)

##RolesUuidGet
//  Get a role.
var api_instance := nutanix.RoleApi
uuid := "uuid_example" // string | The UUID of the entity
rolesuuidget_response, api_response, err := api_instance.RolesUuidGet(uuid)

##RolesUuidPut
//  Update a role.
var api_instance := nutanix.RoleApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.RoleIntentInput() // RoleIntentInput | 
rolesuuidput_response, api_response, err := api_instance.RolesUuidPut(uuid, body)

