#\WebhookApi

##WebhooksListPost
//  Get a list of Webhooks
var api_instance := nutanix.WebhookApi
getEntitiesRequest := nutanix.WebhookListMetadata() // WebhookListMetadata |
webhookslistpost_response, api_response, err := api_instance.WebhooksListPost(getEntitiesRequest)

##WebhooksPost
//  Create a Webhook
var api_instance := nutanix.WebhookApi
body := nutanix.WebhookIntentInput() // WebhookIntentInput |
webhookspost_response, api_response, err := api_instance.WebhooksPost(body)

##WebhooksUuidDelete
//  Delete a Webhook
var api_instance := nutanix.WebhookApi
uuid := "uuid_example" // string | The UUID of the entity
webhooksuuiddelete_response, api_response, err := api_instance.WebhooksUuidDelete(uuid)

##WebhooksUuidGet
//  Get a Webhook
var api_instance := nutanix.WebhookApi
uuid := "uuid_example" // string | The UUID of the entity
webhooksuuidget_response, api_response, err := api_instance.WebhooksUuidGet(uuid)

##WebhooksUuidPut
//  Update a Webhook
var api_instance := nutanix.WebhookApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.WebhookIntentInput() // WebhookIntentInput |
webhooksuuidput_response, api_response, err := api_instance.WebhooksUuidPut(uuid, body)
