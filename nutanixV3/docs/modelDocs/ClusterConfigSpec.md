# ClusterConfigSpec

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthorizedPublicKeyList** | [**[]PublicKey**](public_key.md) | List of valid ssh keys for the cluster. | [optional] [default to null]
**CertificationSigningInfo** | [**CertificationSigningInfo**](certification_signing_info.md) |  | [optional] [default to null]
**ClientAuth** | [**ClientAuth**](client_auth.md) |  | [optional] [default to null]
**EnabledFeatureList** | **[]string** | Array of enabled features. | [optional] [default to null]
**EncryptionStatus** | **string** | Cluster encryption status. | [optional] [default to null]
**RedundancyFactor** | **int64** | Cluster supported redundancy factor. Default is 2. | [optional] [default to null]
**SoftwareMap** | [**map[string]ClusterSoftware**](cluster_software.md) | Map of software on the cluster with software type as the key.  | [optional] [default to null]
**SupportedInformationVerbosity** | **string** | Verbosity level settings for populating support information. - &#39;Nothing&#39;: Send nothing - &#39;Basic&#39;: Send basic information - skip core dump and hypervisor            stats information - &#39;BasicPlusCoreDump&#39;: Send basic and core dump information - &#39;All&#39;: Send all information  | [optional] [default to null]
**Timezone** | **string** | Zone name used in value of TZ environment variable. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
