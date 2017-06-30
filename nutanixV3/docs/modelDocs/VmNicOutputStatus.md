# VmNicOutputStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IpEndpointList** | [**[]IpAddress**](ip_address.md) | IP endpoints for the adapter. Currently, IPv4 addresses are supported.  | [optional] [default to null]
**MacAddress** | **string** | The MAC address for the adapter. | [optional] [default to null]
**NetworkFunctionChainReference** | [**NetworkFunctionChainReference**](network_function_chain_reference.md) |  | [optional] [default to null]
**NetworkFunctionNicType** | **string** | The type of this Network function NIC. Defaults to INGRESS.  | [optional] [default to null]
**NicType** | **string** | The type of this NIC. Defaults to NORMAL_NIC. | [optional] [default to null]
**SubnetReference** | [**Reference**](reference.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
