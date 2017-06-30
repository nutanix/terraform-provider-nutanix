# PortalSoftware

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CompatiblePeVersionList** | **[]string** | List of Prism Element compatible versions | [optional] [default to null]
**CompatibleVersionList** | **[]string** | List of software versions that this version can be upgraded from  | [optional] [default to null]
**Md5sum** | **string** | MD5 checksum of the software file | [optional] [default to null]
**ReleaseDate** | [**time.Time**](time.Time.md) | Release date of this software in RFC3339 format.  | [optional] [default to null]
**ReleaseNoteUrl** | **string** | URL to point to the support portal release note of this software. Currently only set and used for NOS releases  | [optional] [default to null]
**SizeInMib** | **int64** | Total size of the software file in mebibytes | [optional] [default to null]
**SoftwareType** | **string** | Software type | [optional] [default to null]
**UpgradeNotification** | [**UpgradeNotification**](upgrade_notification.md) |  | [optional] [default to null]
**Version** | **string** | Software version string | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


