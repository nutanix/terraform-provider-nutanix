#\IdempotenceidentifiersApi

##IdempotenceIdentifiersClientIdentifierDelete
//  Deletes an idempotence identifier object.
var api_instance := nutanix.IdempotenceidentifiersApi
clientIdentifier := "clientIdentifier_example" // string |
idempotenceidentifiersclientidentifierdelete_response, api_response, err := api_instance.IdempotenceIdentifiersClientIdentifierDelete(clientIdentifier)

##IdempotenceIdentifiersClientIdentifierGet
//  Get an idempotence identifier object.
var api_instance := nutanix.IdempotenceidentifiersApi
clientIdentifier := "clientIdentifier_example" // string |
idempotenceidentifiersclientidentifierget_response, api_response, err := api_instance.IdempotenceIdentifiersClientIdentifierGet(clientIdentifier)

##IdempotenceIdentifiersPost
//  Creates an idempotence identifier
var api_instance := nutanix.IdempotenceidentifiersApi
body := nutanix.IdempotenceIdentifiersInput() // IdempotenceIdentifiersInput | An idempotence identifier object. (optional)
idempotenceidentifierspost_response, api_response, err := api_instance.IdempotenceIdentifiersPost(body=body)
