# VmSnapshotMetadata

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreationTime** | [**time.Time**](time.Time.md) | Time when vm snapshot was created | [optional] [default to null]
**EntityVersion** | **int64** | Monotonically increasing number | [optional] [default to null]
**Kind** | **string** | The kind snapshot name | [optional] [default to null]
**LastUpdateTime** | [**time.Time**](time.Time.md) | Time when vm snapshot was last updated | [optional] [default to null]
**Name** | **string** | vm snapshot name | [optional] [default to null]
**OwnerReference** | [**UserReference**](user_reference.md) |  | [optional] [default to null]
**ParentReference** | [**VmReference**](vm_reference.md) |  | [optional] [default to null]
**Uuid** | **string** | vm snapshot UUID | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
