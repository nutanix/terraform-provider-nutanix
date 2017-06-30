# BatchRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ActionOnFailure** | **string** | If the specified parameter is CONTINUE, the remaining APIs in the batch continue to be executed.  | [optional] [default to null]
**ApiRequestList** | [**[]ApiRequest**](api_request.md) | A list of API requests in the batch. | [default to null]
**ApiVersion** | **string** | The current API version. | [default to null]
**ExecutionOrder** | **string** | The order of execution of the APIs in the batch. Can be either Sequential (default value) or Parallel.  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
