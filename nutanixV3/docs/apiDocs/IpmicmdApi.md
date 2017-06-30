#\IpmicmdApi

##HostsUuidRunIpmiCmdPost
//  Run IPMI command on the given host
var api_instance := nutanix.IpmicmdApi
uuid := "uuid_example" // string | The UUID of the entity
ipmiArgs := nutanix.HostIpmiArgs() // HostIpmiArgs | The arguments for the IPMI tool as a single string
hostsuuidrunipmicmdpost_response, api_response, err := api_instance.HostsUuidRunIpmiCmdPost(uuid, ipmiArgs)

