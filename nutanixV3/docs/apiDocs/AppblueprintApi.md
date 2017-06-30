#\AppblueprintApi

##AppBlueprintsListPost
//  List the App Blueprints
var api_instance := nutanix.AppblueprintApi
getEntitiesRequest := nutanix.AppBlueprintListMetadata() // AppBlueprintListMetadata | 
appblueprintslistpost_response, api_response, err := api_instance.AppBlueprintsListPost(getEntitiesRequest)

##AppBlueprintsPost
//  Create an App Blueprint
var api_instance := nutanix.AppblueprintApi
body := nutanix.AppBlueprintIntentInput() // AppBlueprintIntentInput | 
appblueprintspost_response, api_response, err := api_instance.AppBlueprintsPost(body)

##AppBlueprintsRenderPost
//  Render and Create an AppBlueprint from the given input
var api_instance := nutanix.AppblueprintApi
body := nutanix.AppBlueprintRenderInput() // AppBlueprintRenderInput | 
appblueprintsrenderpost_response, api_response, err := api_instance.AppBlueprintsRenderPost(body)

##AppBlueprintsUuidDelete
//  Delete App Blueprint
var api_instance := nutanix.AppblueprintApi
uuid := "uuid_example" // string | The UUID of the entity
appblueprintsuuiddelete_response, api_response, err := api_instance.AppBlueprintsUuidDelete(uuid)

##AppBlueprintsUuidGet
//  Get App Blueprint
var api_instance := nutanix.AppblueprintApi
uuid := "uuid_example" // string | The UUID of the entity
appblueprintsuuidget_response, api_response, err := api_instance.AppBlueprintsUuidGet(uuid)

##AppBlueprintsUuidPut
//  Update the App Blueprint
var api_instance := nutanix.AppblueprintApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.AppBlueprintIntentInput() // AppBlueprintIntentInput | 
appblueprintsuuidput_response, api_response, err := api_instance.AppBlueprintsUuidPut(uuid, body)

