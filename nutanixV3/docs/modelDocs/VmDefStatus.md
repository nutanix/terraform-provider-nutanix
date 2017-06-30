# VmDefStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClusterReference** | [**ClusterReference**](cluster_reference.md) | Reference to the cluster where this VM exists or needs to be migrated to  | [optional] [default to null]
**Description** | **string** | A description or user annotation for the VM. | [optional] [default to null]
**MessageList** | [**[]MessageResource**](message_resource.md) | Any error messages for the VM, if in an error state. | [optional] [default to null]
**Name** | **string** | VM Name. | [default to null]
**Resources** | [**VmDefStatusResources**](vm_def_status_resources.md) |  | [optional] [default to null]
**State** | **string** | The state of the vm entity. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
