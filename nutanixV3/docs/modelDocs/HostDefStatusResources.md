# HostDefStatusResources

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Block** | [**Block**](block.md) |  | [optional] [default to null]
**ControllerVm** | [**ControllerVm**](controller_vm.md) |  | [optional] [default to null]
**CpuCapacityHz** | **int64** | Host CPU capacity. | [optional] [default to null]
**CpuModel** | **string** | Host CPU model. | [optional] [default to null]
**FailoverCluster** | [**FailoverCluster**](failover_cluster.md) |  | [optional] [default to null]
**GpuList** | [**[]HostGpu**](host_gpu.md) | List of GPUs on the host. | [optional] [default to null]
**HostDisksReferenceList** | [**[]DiskReference**](disk_reference.md) |  | [optional] [default to null]
**HostNicsIdList** | **[]string** | Host NICs. | [optional] [default to null]
**Hypervisor** | [**Hypervisor**](hypervisor.md) |  | [optional] [default to null]
**Ipmi** | [**Ipmi**](ipmi.md) |  | [optional] [default to null]
**MemoryCapacityMib** | **int64** | Host memory capacity in MiB. | [optional] [default to null]
**MonitoringState** | **string** | Host monitoring status. | [optional] [default to null]
**NumCpuCores** | **int64** | Number of CPU cores on Host. | [optional] [default to null]
**NumCpuSockets** | **int64** | Number of CPU sockets. | [optional] [default to null]
**SerialNumber** | **string** | Node serial number. | [optional] [default to null]
**WindowsDomain** | [**WindowsDomain**](windows_domain.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
