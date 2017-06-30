#\DataApi

##DataChangedRegionsPost
//  Query changed regions metadata.
var api_instance := nutanix.DataApi
body := nutanix.ChangedRegionsQuery() // ChangedRegionsQuery |
datachangedregionspost_response, api_response, err := api_instance.DataChangedRegionsPost(body)
