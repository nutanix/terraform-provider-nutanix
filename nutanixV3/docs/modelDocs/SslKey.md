# SslKey

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExpireDatetime** | [**time.Time**](time.Time.md) | UTC date and time in RFC-3339 format when the key expires | [optional] [default to null]
**KeyName** | **string** |  | [optional] [default to null]
**KeyType** | **string** | SSL key type. Key types with RSA_2048, ECDSA_256 and ECDSA_384 are supported for key generation and importing.  | [default to null]
**SigningInfo** | [**CertificationSigningInfo**](certification_signing_info.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


