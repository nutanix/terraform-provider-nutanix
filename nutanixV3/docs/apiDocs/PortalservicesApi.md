#\PortalservicesApi

##PortalServicesSoftwareSoftwareTypeListPost
//  Get all available software on Nutanix Portal
var api_instance := nutanix.PortalservicesApi
softwareType := "softwareType_example" // string | Software type
getEntitiesRequest := nutanix.SoftwareListMetadata() // SoftwareListMetadata |  (optional)
portalservicessoftwaresoftwaretypelistpost_response, api_response, err := api_instance.PortalServicesSoftwareSoftwareTypeListPost(softwareType, getEntitiesRequest=getEntitiesRequest)

##PortalServicesSoftwareSoftwareTypeVersionGet
//  Get specified software information
var api_instance := nutanix.PortalservicesApi
softwareType := "softwareType_example" // string | Software type
version := "version_example" // string | Software version
portalservicessoftwaresoftwaretypeversionget_response, api_response, err := api_instance.PortalServicesSoftwareSoftwareTypeVersionGet(softwareType, version)
