# VolumeGroupResources

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AttachmentList** | [**[]AttachmentReference**](attachment_reference.md) | VMs attached to volume group. | [optional] [default to null]
**DiskList** | [**[]VolumeDiskResource**](volume_disk_resource.md) | Volume group disk specification. | [optional] [default to null]
**FileSystemType** | **string** | File system to be used for volume | [optional] [default to null]
**IscsiTargetName** | **string** | iSCSI target full name | [optional] [default to null]
**IscsiTargetPrefix** | **string** | iSCSI target prefix-name. | [optional] [default to null]
**SharingStatus** | **string** | Whether the volume group can be shared across multiple iSCSI initiators.  | [optional] [default to null]
**SizeMib** | **int64** | The total size of the Volume Group. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


