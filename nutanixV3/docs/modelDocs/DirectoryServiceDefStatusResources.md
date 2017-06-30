# DirectoryServiceDefStatusResources

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AdminGroupList** | [**[]UserGroup**](user_group.md) | List of distinguished names of the admin group in the directory service.  | [optional] [default to null]
**AdminUserReferenceList** | [**[]UserReference**](user_reference.md) | The list of admin users available in the directory service.  | [optional] [default to null]
**ConfiguredFeatureList** | **[]string** | List of features configured with directory service. | [optional] [default to null]
**DirectoryType** | **string** | Type of the directory service. | [default to null]
**DomainName** | **string** | The domain name of the directory service. | [default to null]
**ServiceAccount** | [**ServiceAccount**](service_account.md) | Validates and connects to the directory service with the given credentials.  | [optional] [default to null]
**Url** | **string** | URL of the directory. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


