# DiskDefStatusResources

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EnabledFeaturesList** | **[]string** | Disk feature flags - &#39;CanAddAsNewDisk&#39;: Flag to indicate if this disk can be added as    new disk. - &#39;CanAddAsOldDisk&#39;: Flag to indicate if the disk can be added as    old disk. - &#39;BootDisk&#39;: Flag to indicate if its a boot disk. - &#39;OnlyBootDisk&#39;: Flag to indicate if the disk is boot only and    no disk operation to be run on it. - &#39;SelfEncryptingEnabled&#39;: Flag to indicate if the disk has self    encryption enabled. - &#39;PasswordProtected&#39;: Flag to indicate if the disk is password    protected.  | [optional] [default to null]
**FirmwareVersion** | **string** | Firmware version. | [optional] [default to null]
**HostReference** | [**Reference**](reference.md) |  | [optional] [default to null]
**Model** | **string** | Disk model. | [optional] [default to null]
**MountPath** | **string** | Mount path. | [optional] [default to null]
**SerialNumber** | **string** | Disk serial number. | [optional] [default to null]
**SizeMib** | **int64** | Disk size in Mib. | [optional] [default to null]
**SlotNumber** | **int64** | Disk location in a node. | [optional] [default to null]
**StateList** | **[]string** | Array of disk states - DataMigrationInitiated: Data Migration Initiated. - MarkedForRemovalButNotDetachable: Marked for removal, data    migration is in progress. - ReadyToDetach: Flag to indicate the disk is detachable. - DataMigrated: Flag to indicate if data migration is completed for    this disk. - MarkedForRemoval: Flag to indicate if the disk is marked for    removal. - Online: Flag to indicate if the disk is online. - Bad: Flag to indicate if the disk is bad. - Mounted: Flag to indicate if the disk is mounted. - UnderDiagnosis: Flag to indicate if the disk is under diagnosis.  | [optional] [default to null]
**StoragePoolUuid** | **string** | Storage pool uuid. | [optional] [default to null]
**StorageTierType** | **string** | Storage tier type. | [optional] [default to null]
**Vendor** | **string** | Disk vendor. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


