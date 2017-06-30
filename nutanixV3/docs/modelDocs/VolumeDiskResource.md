# VolumeDiskResource

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Index** | **int64** | Volume group index of the disk. | [optional] [default to null]
**NfsPath** | **string** | NFS path of the image. For creating/updating operations. If NFS path is specified, then the operation would be clone and vmdisk_uuid should not be set.  | [optional] [default to null]
**SizeMib** | **int64** | Size of the disk in MiB. | [optional] [default to null]
**StorageContainerUuid** | **string** | Container UUID on which to create the disk. | [optional] [default to null]
**VmdiskUuid** | **string** | UUID of the vmdisk. For creating/updating operation, this will be the UUID of the disk to clone from. If UUID is specified, then the operation would be clone and nfs_path should not be set.  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
