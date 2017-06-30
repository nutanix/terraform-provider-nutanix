#\BatchApi

##BatchPost
//  Submit a list of one or more intentful REST APIs to be processed
var api_instance := nutanix.BatchApi
intentList := nutanix.BatchRequest() // BatchRequest | List of intent APIs
batchpost_response, api_response, err := api_instance.BatchPost(intentList)

