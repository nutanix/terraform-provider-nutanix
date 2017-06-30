# ProjectResources

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DefaultSubnetReference** | [**SubnetReference**](subnet_reference.md) | Optional default subnet if one is specified | [optional] [default to null]
**ExternalUserGroupList** | **[]string** | List of directory service group&#39;s distinguished name. Those groups are not managed by Nutanix.  | [optional] [default to null]
**ResourceDomain** | [**ResourceDomainSpec**](resource_domain_spec.md) |  | [optional] [default to null]
**RoleReference** | [**RoleReference**](role_reference.md) | The role assigned to project users | [default to null]
**SubnetReferenceList** | [**[]SubnetReference**](subnet_reference.md) | List of subnets for the project. | [optional] [default to null]
**UserReferenceList** | [**[]UserReference**](user_reference.md) | List of users in the project. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
