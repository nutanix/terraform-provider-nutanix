# VmDefStatusResources

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BootConfig** | [**VmBootConfig**](vm_boot_config.md) | Indicates which device the VM should boot from. | [optional] [default to null]
**DiskList** | [**[]VmDisk**](vm_disk.md) | Disks attached to the VM. | [optional] [default to null]
**GpuList** | [**[]VmGpuOutputStatus**](vm_gpu_output_status.md) | GPUs attached to the VM. | [optional] [default to null]
**GuestCustomization** | [**GuestCustomizationStatus**](guest_customization_status.md) |  | [optional] [default to null]
**HostReference** | [**Reference**](reference.md) |  | [optional] [default to null]
**HypervisorType** | **string** | The hypervisor type for the hypervisor the VM is hosted on.  | [optional] [default to null]
**MemorySizeMib** | **int64** | Memory size in MiB. | [optional] [default to null]
**NicList** | [**[]VmNicOutputStatus**](vm_nic_output_status.md) | NICs attached to the VM. | [optional] [default to null]
**NumSockets** | **int64** | Number of vCPU sockets. | [optional] [default to null]
**NumVcpusPerSocket** | **int64** | Number of vCPUs per socket. | [optional] [default to null]
**ParentReference** | [**Reference**](reference.md) | Reference to an entity that the VM cloned from.  | [optional] [default to null]
**PowerState** | **string** | Desired power state of the VM. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
