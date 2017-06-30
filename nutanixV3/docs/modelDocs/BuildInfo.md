# BuildInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BuildType** | **string** | Build type, one of {dbg, opt, release}. | [default to null]
**CommitDate** | [**time.Time**](time.Time.md) | Date/time of the last commit. | [default to null]
**CommitId** | **string** | Last Git commit id which the build is based on. | [default to null]
**ShortCommitId** | **string** | First 6 characters of the last Git commit id. | [default to null]
**Version** | **string** | Version string in format &lt;code_name&gt;-&lt;version_numbers&gt;-&lt;branch_type&gt;, i.e master, danube-4.5.0.2-stable  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
