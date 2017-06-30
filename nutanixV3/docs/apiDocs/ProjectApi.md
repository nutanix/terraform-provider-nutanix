#\ProjectApi

##ProjectsListPost
//  Retrieves all Projects.
var api_instance := nutanix.ProjectApi
getEntitiesRequest := nutanix.ProjectListMetadata() // ProjectListMetadata | 
projectslistpost_response, api_response, err := api_instance.ProjectsListPost(getEntitiesRequest)

##ProjectsPost
//  Creates a project.
var api_instance := nutanix.ProjectApi
body := nutanix.ProjectIntentInput() // ProjectIntentInput | Project object.
projectspost_response, api_response, err := api_instance.ProjectsPost(body)

##ProjectsUuidDelete
//  Deletes a project.
var api_instance := nutanix.ProjectApi
uuid := "uuid_example" // string | The UUID of the entity
projectsuuiddelete_response, api_response, err := api_instance.ProjectsUuidDelete(uuid)

##ProjectsUuidGet
//  Retrieves specified Project.
var api_instance := nutanix.ProjectApi
uuid := "uuid_example" // string | The UUID of the entity
projectsuuidget_response, api_response, err := api_instance.ProjectsUuidGet(uuid)

##ProjectsUuidPut
//  Updates a project.
var api_instance := nutanix.ProjectApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.ProjectIntentInput() // ProjectIntentInput | Project object.
projectsuuidput_response, api_response, err := api_instance.ProjectsUuidPut(uuid, body)

