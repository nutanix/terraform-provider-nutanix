#\PrismcentralApi

##DeployPrismCentralPost
//  Create a Prism Central
var api_instance := nutanix.PrismcentralApi
body := nutanix.PrismCentral() // PrismCentral | 
deployprismcentralpost_response, api_response, err := api_instance.DeployPrismCentralPost(body)

