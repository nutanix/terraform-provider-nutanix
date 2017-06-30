# VmSnapshotDefStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MessageList** | [**[]MessageResource**](message_resource.md) | Any error messages for the vm, if in an error state. | [optional] [default to null]
**Name** | **string** | Name of the snapshot. | [optional] [default to null]
**Resources** | [**VmSnapshotResources**](vm_snapshot_resources.md) |  | [optional] [default to null]
**SnapshotFileList** | [**[]VmSnapshotDefStatusSnapshotFileList**](vm_snapshot_def_status_snapshot_file_list.md) | Describes the files that are included in the snapshot.  | [default to null]
**SnapshotType** | **string** | The consistency level desired while creating the snapshot. | [optional] [default to null]
**State** | **string** | The state of the vm entity. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


