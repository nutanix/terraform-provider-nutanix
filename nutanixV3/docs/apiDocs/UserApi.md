#\UserApi

##LogoutGet
//  Logs out the current user
var api_instance := nutanix.UserApi
logoutget_response, api_response, err := api_instance.LogoutGet()

##UsersMeGet
//  Retrieves currently logged in user.
var api_instance := nutanix.UserApi
usersmeget_response, api_response, err := api_instance.UsersMeGet()

##UsersUuidGet
//  Retrieves specified user.
var api_instance := nutanix.UserApi
uuid := "uuid_example" // string | The UUID of the entity
usersuuidget_response, api_response, err := api_instance.UsersUuidGet(uuid)

##UsersUuidProjectUsageSummaryGet
//  Retrieves specified user resource domain information.
var api_instance := nutanix.UserApi
uuid := "uuid_example" // string | The UUID of the entity
usersuuidprojectusagesummaryget_response, api_response, err := api_instance.UsersUuidProjectUsageSummaryGet(uuid)

