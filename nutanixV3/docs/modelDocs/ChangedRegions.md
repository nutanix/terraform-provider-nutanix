# ChangedRegions

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FileSize** | **int64** | Size of the file specified by snapshot_file_path | [optional] [default to null]
**NextOffset** | **int64** | The offset from where the client must continue the request. This field will not be set when there are no more changed regions to be returned. Note that the next_offset can be outside the endOffset specified by the client in the request. This helps clients reach the next changed offset faster.  | [optional] [default to null]
**RegionList** | [**[]Region**](region.md) | List of regions describing the change for the interval [start_offset, next_offset].  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
