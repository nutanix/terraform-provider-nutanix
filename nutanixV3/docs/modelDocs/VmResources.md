# VmResources

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BootConfig** | [**VmBootConfig**](vm_boot_config.md) | Indicates which device the VM should boot from. | [optional] [default to null]
**DiskList** | [**[]VmDisk**](vm_disk.md) | Disks attached to the VM. | [optional] [default to null]
**GpuList** | [**[]VmGpu**](vm_gpu.md) | GPUs attached to the VM. | [optional] [default to null]
**GuestCustomization** | [**GuestCustomization**](guest_customization.md) |  | [optional] [default to null]
**MemorySizeMib** | **int64** | Memory size in MiB. | [optional] [default to null]
**NicList** | [**[]VmNic**](vm_nic.md) | NICs attached to the VM. | [optional] [default to null]
**NumSockets** | **int64** | Number of vCPU sockets. | [optional] [default to null]
**NumVcpusPerSocket** | **int64** | Number of vCPUs per socket. | [optional] [default to null]
**ParentReference** | [**Reference**](reference.md) | Reference to an entity that the VM should be cloned from.  | [optional] [default to null]
**PowerState** | **string** | The current or desired power state of the VM. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
