#\PermissionApi

##PermissionsListPost
//  List the permissions.
var api_instance := nutanix.PermissionApi
getEntitiesRequest := nutanix.PermissionListMetadata() // PermissionListMetadata | 
permissionslistpost_response, api_response, err := api_instance.PermissionsListPost(getEntitiesRequest)

##PermissionsPost
//  Create a permission.
var api_instance := nutanix.PermissionApi
body := nutanix.PermissionIntentInput() // PermissionIntentInput | 
permissionspost_response, api_response, err := api_instance.PermissionsPost(body)

##PermissionsUuidDelete
//  Delete a permission.
var api_instance := nutanix.PermissionApi
uuid := "uuid_example" // string | The UUID of the entity
permissionsuuiddelete_response, api_response, err := api_instance.PermissionsUuidDelete(uuid)

##PermissionsUuidGet
//  Get a permission.
var api_instance := nutanix.PermissionApi
uuid := "uuid_example" // string | The UUID of the entity
permissionsuuidget_response, api_response, err := api_instance.PermissionsUuidGet(uuid)

##PermissionsUuidPut
//  Update a permission.
var api_instance := nutanix.PermissionApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.PermissionIntentInput() // PermissionIntentInput | 
permissionsuuidput_response, api_response, err := api_instance.PermissionsUuidPut(uuid, body)

