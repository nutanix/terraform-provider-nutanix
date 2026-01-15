---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_clone_v2"
sidebar_current: "docs-nutanix-resource-vm-clone-v2"
description: |-
  Provides a Nutanix Virtual Machine resource to Create a virtual machine clone.
---

# nutanix_vm_clone_v2

Provides a Nutanix Virtual Machine resource to Create a virtual machine clone.

## Example Usage

```hcl
data "nutanix_virtual_machines_v2" "vm-list"{}

resource "nutanix_vm_clone_v2" "vm1"{
  vm_ext_id = data.nutanix_virtual_machines_v2.vm-list.vms.0.data.ext_id
  name = "test-dou"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
}
```

## Argument Reference

The following arguments are supported:

* `vm_ext_id`: - (Required) The globally unique identifier of a VM. It should be of type UUID.
* `name`: - (Optional) The name for the vm.
* `num_sockets`: - (Optional) Number of vCPU sockets.
* `num_cores_per_socket`: - (Optional) Number of cores per socket.
* `memory_size_mib`: - (Optional) Memory size in MiB.
* `num_threads_per_core`: - (Optional) Number of threads per core.
* `guest_customization`: - (Optional) Stage a Sysprep or cloud-init configuration file to be used by the guest for the next boot. Note that the Sysprep command must be used to generalize the Windows VMs before triggering this API call.
* `boot_config`: - (Optional) Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order.
* `nics`: - (Optional) NICs attached to the VM.

### Nics

The nics attribute supports the following:

* `ext_id`: - (Optional) A globally unique identifier of an instance that is suitable for external consumption.
* `nic_backing_info`: (Optional) New NIC backing info (v2.4.1+). One of `virtual_ethernet_nic`, `sriov_nic`, `dp_offload_nic`.
* `nic_network_info`: (Optional) New NIC network info (v2.4.1+). One of `virtual_ethernet_nic_network_info`, `sriov_nic_network_info`, `dp_offload_nic_network_info`.
* `backing_info`: - (Optional, Deprecated) Use `nic_backing_info.virtual_ethernet_nic` instead.
* `network_info`: - (Optional, Deprecated) Use `nic_network_info.virtual_ethernet_nic_network_info` instead.

### Nic Backing Info (new)

* `nic_backing_info.virtual_ethernet_nic`: Virtual Ethernet NIC backing info.
* `nic_backing_info.sriov_nic`: SR-IOV NIC backing info.
* `nic_backing_info.dp_offload_nic`: DP offload NIC backing info.

### Nic Network Info (new)

* `nic_network_info.virtual_ethernet_nic_network_info`: Virtual Ethernet NIC network info.
* `nic_network_info.sriov_nic_network_info`: SR-IOV NIC network info.
* `nic_network_info.dp_offload_nic_network_info`: DP offload NIC network info.

### Backing Info

The backing_info attribute supports the following:

* `model`: - (Optional) Options for the NIC emulation.
    Valid values are:
     - `VIRTIO` The NIC emulation model is Virtio.
     - `E1000` The NIC emulation model is E1000.
* `mac_address`: - (Optional) MAC address of the emulated NIC.
* `is_connected`: - (Optional) Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: - (Optional) The number of Tx/Rx queue pairs for this NIC.

### Network Info

The network_info attribute supports the following:

* `nic_type`: - (Optional) NIC type.
Defaults to NORMAL_NIC.
    Valid values are:
     - `SPAN_DESTINATION_NIC` The type of NIC is Span-Destination.
     - `NORMAL_NIC` The type of NIC is Normal.
     - `DIRECT_NIC` The type of NIC is Direct.
     - `NETWORK_FUNCTION_NIC` The type of NIC is Network-Function.
* `network_function_chain`: - (Optional)The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_nic_type`: - (Optional) The type of this Network function NIC.
Defaults to INGRESS.
    Valid values are:
     - `TAP` The type of Network-Function NIC is Tap.
     - `EGRESS` The type of Network-Function NIC is Egress.
     - `INGRESS` The type of Network-Function NIC is Ingress.
* `subnet`: - (Optional) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC.
* `vlan_mode`: - (Optional) By default, all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs.
    Valid values are:
     - `TRUNK` The virtual NIC is created in TRUNKED mode.
     - `ACCESS` The virtual NIC is created in ACCESS mode.
* `trunked_vlans`: - (Optional) List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
* `should_allow_unknown_macs`: - (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
* `ipv4_config`: - (Optional) The IP address configurations.

### Network Function Chain

The network_function_chain attribute supports the following:

* `ext_id`: - (Optional) The globally unique identifier of a network function chain. It should be of type UUID.

### Subnet

The subnet attribute supports the following:

* `ext_id`: - (Optional) The globally unique identifier of a subnet. It should be of type UUID.

### IPV4 Config

The ipv4_config attribute supports the following:

* `should_assign_ip`: - (Optional) If set to true (default value), an IP address must be assigned to the VM NIC - either the one explicitly specified by the user or allocated automatically by the IPAM service by not specifying the IP address. If false, then no IP assignment is required for this VM NIC.
`ip_address`: - (Optional) Ip config settings.
`secondary_ip_address_list`: - (Optional) Secondary IP addresses for the NIC.

### IP Address

The ip_address attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - Ip address.

### Secondary IP Address List

The secondary_ip_address_list attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - Ip address.

### Boot Config

The boot_config attribute supports the following:

* `legacy_config`: - (Optional) The Nutanix Legacy Boot Config settings.
* `uefi_config`: - (Optional) The Nutanix Uefi Boot Config settings.

### Legacy Boot Config

The legacy_boot attribute supports the following:

* `boot_device`: - (Optional) The Boot Device settings.
* `boot_order`: - (Optional) Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order.

### Boot Device

The boot_device attribute supports the following:

* `boot_device_disk`: - (Optional) The Boot Device Disk settings.
* `boot_device_nic`: - (Optional) The Boot Device Nic settings.

### Boot Device Disk

The boot_device_disk attribute supports the following:

* `disk_address`: - (Optional) Address of disk to boot from.


### Disk Address

 The disk_address attribute supports the following:

* `index`: - (Optional) Device index on the bus. This field is ignored unless the bus details are specified.
* `bus_type`: - (Optional) Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
    Valid values are:
    - `SCSI` The type of disk bus is SCSI.
    - `SPAPR` The type of disk bus is SPAPR.
    - `PCI` The type of disk bus is PCI.
    - `PCI` The type of disk bus is PCI.
    - `SATA` The type of disk bus is SATA.

### Boot Device Nic

The boot_device_nic attribute supports the following:

* `mac_address`: - (Optional) MAC address of nic to boot from.

### Uefi Boot Config

The uefi_boot attribute supports the following:

* `is_secure_boot_enabled`: - (Optional) Indicate whether to enable secure boot or not.
* `nvram_device`: - (Optional) Configuration for NVRAM to be presented to the VM.

### Nvram Device

The nvram_device attribute supports the following:

* `backing_storage_info`: - (Optional) Storage provided by Nutanix ADSF.

### Backing Storage Info

The backing_storage_info attribute supports the following:

* `disk_size_bytes`: - (Optional) Size of the disk in Bytes.
* `storage_container`: - (Optional) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: - (Optional) Storage configuration for VM disks.
* `data_source`: - (Optional) A reference to a disk or image that contains the contents of a disk.

### Storage Container

The storage_container attribute supports the following:

* `ext_id`: - (Optional) The globally unique identifier of a VM disk container. It should be of type UUID.

### Storage Config

The storage_config attribute supports the following:

* `is_flash_mode_enabled`: - (Optional) Indicates whether the virtual disk is pinned to the hot tier or not.

### Data Source

The data_source attribute supports the following:

* `reference`: - (Optional) Data Source Reference settings.

### Data Source Reference

The reference attribute supports the following:

* `image_reference`: - (Optional) Data Source Image Reference settings.
* `vm_disk_reference`: - (Optional) Data Source VM Disk Reference settings.

### Image Reference

The image_reference attribute supports the following:

* `image_ext_id`: - (Optional) The globally unique identifier of an image. It should be of type UUID.

### VM Disk Reference

The vm_disk_reference attribute supports the following:

* `disk_ext_id`: - (Optional) The globally unique identifier of a VM disk. It should be of type UUID.
* `disk_address`: - (Optional) Address of disk.
* `vm_reference`: - (Optional) Reference to a VM.

### VM Reference

The vm_reference attribute supports the following:

* `ext_id`: - (Optional) The globally unique identifier of a VM. It should be of type UUID.

### Guest Customization

The guest_customization attribute supports the following:

* `config`: - (Optional) The Nutanix Guest Tools customization settings.

### Config

The config attribute supports the following:

* `sysprep`: - (Optional) VM guests may be customized at boot time using one of several different methods. Currently, cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting to change it after creation will result in an error. Additional properties can be specified. For example - in the context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own custom script.

* `cloud_init`: - (Optional) VM guests may be customized at boot time using one of several different methods. Currently, cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting to change it after creation will result in an error. Additional properties can be specified. For example - in the context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own custom script.

### Sysprep

The sysprep attribute supports the following:

* `install_type`: - (Optional) Whether the guest will be freshly installed using this unattend configuration, or whether this unattend configuration will be applied to a pre-prepared image. Default is `PREPARED`.
    Valid values are:
    - `PREPARED` is done when sysprep is used to finalize Windows installation from an installed Windows and file name it is searching `unattend.xml` for `unattend_xml` parameter
    - `FRESH` is done when sysprep is used to install Windows from ISO and file name it is searching `autounattend.xml` for `unattend_xml` parameter
* `unattend_xml`: - (Optional) Generic key value pair used for custom attributes.

### Cloud Init

The cloud_init attribute supports the following:

* `datasource_type`: - (Optional) Type of datasource.
Default: CONFIG_DRIVE_V2Default is `CONFIG_DRIVE_V2`.
    Valid values are:
    - `CONFIG_DRIVE_V2` The type of datasource for cloud-init is Config Drive V2.
* `metadata` - (Optional) The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded.
* `cloud_init_script`: - (Optional) The script to use for cloud-init.

### Cloud Init Script

The cloud_init_script attribute supports the following:

* `user_data`: - (Optional) The contents of the user_data configuration for cloud-init. This can be formatted as YAML, JSON, or could be a shell script. The value must be base64 encoded.
* `custom_key_values`: - (Optional) Generic key value pair used for custom attributes in cloud init.

### User Data

The user_data attribute supports the following:

* `value`: - (Optional) The value for the cloud-init user_data.

### Custom Key Values

The custom_key_values attribute supports the following:

* `key_value_pairs`: - (Optional) The list of the individual KeyValuePair elements.

### Key Value Pairs

The key_value_pairs attribute supports the following:

* `name`: - (Optional) The key of this key-value pair
* `value`: - (Optional) The value associated with the key for this key-value pair.

See detailed information in [Nutanix Clone Virtual Machine V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/cloneVm).
