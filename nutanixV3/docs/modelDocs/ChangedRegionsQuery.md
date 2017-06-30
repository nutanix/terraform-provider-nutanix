# ChangedRegionsQuery

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EndOffset** | **int64** | The absolute offset in bytes up to which to query for the changed regions. Note that the interval specified by the start_offset together with the end_offset is right half-open. If the end_offset is not specified, the portion from the start_offset till the end of the file will be included in the query.  | [optional] [default to null]
**ReferenceSnapshotFilePath** | **string** | Absolute path of a file within a snapshot that must be used as the reference in the computation of the changed regions. If this path is not specified, then the changed regions will not be computed. Instead, the sparse and the non-sparse regions of the file specified in snapshot_file_path will be returned.  | [optional] [default to null]
**SnapshotFilePath** | **string** | Absolute path of a file within a snapshot of an entity such as a virtual machine, a volume group, or a protection domain.  | [default to null]
**StartOffset** | **int64** | The absolute offset in bytes from where to query for the changed regions.  | [optional] [default to 0]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


