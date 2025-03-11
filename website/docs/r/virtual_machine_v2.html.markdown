---
layout: "nutanix"
page_title: "NUTANIX: nutanix_virtual_machine_v2"
sidebar_current: "docs-nutanix-resource-virtual-machine-v2"
description: |-
  Provides a Nutanix Virtual Machine resource to Create a virtual machine.
---

# nutanix_virtual_machine_v2

Creates a Virtual Machine with the provided configuration.

## Example

```hcl

resource "nutanix_virtual_machine_v2" "vm-1"{
    name= "test-vm"
    description =  "vm desc"
    num_cores_per_socket = 1
    num_sockets = 1
    cluster {
        ext_id = "<Cluster uuid>"
    }
}

resource "nutanix_virtual_machine_v2" "vm-2"{
    name= "test-vm"
    description =  "vm desc"
    num_cores_per_socket = 1
    num_sockets = 1
    cluster {
        ext_id = "<Cluster uuid>"
    }
    disks{
        disk_address{
            bus_type = "SCSI"
            index = 0
        }
        backing_info{
            vm_disk{
                disk_size_bytes = "1073741824"
                storage_container{
                    ext_id = "{{ storage ext id}}"
                }
            }
        }
    }
}

resource "nutanix_virtual_machine_v2" "vm-3"{
    name= "test-vm"
    description =  "vm desc"
    num_cores_per_socket = 1
    num_sockets = 1
    cluster {
        ext_id = "<Cluster uuid>"
    }
    disks{
        disk_address{
            bus_type = "SCSI"
            index = 0
        }
        backing_info{
            vm_disk{
                disk_size_bytes = "1073741824"
                storage_container{
                    ext_id = "{{ storage ext id}}"
                }
            }
        }
    }
    nics{
        network_info{
            nic_type = "NORMAL_NIC"
            subnet{
                ext_id = "{{ subnet ext id}}"
            }
            vlan_mode = "ACCESS"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `name`: (Required) VM name.
* `description`: (Optional) VM description
* `source`: (Optional) Reference to an entity that the VM should be cloned or created from. Valid values are "VM", "VM_RECOVERY_POINT".
* `num_sockets`: (Required) Number of vCPU sockets. Value should be at least 1.
* `num_cores_per_socket`: (Optional) Number of cores per socket. Value should be at least 1.
* `num_threads_per_core`: (Optional) Number of threads per core. Value should be at least 1.
* `num_numa_nodes`: (Optional) Number of NUMA nodes. 0 means NUMA is disabled.
* `memory_size_bytes`: (Required) Memory size in bytes.
* `is_vcpu_hard_pinning_enabled`: (Optional) Indicates whether the vCPUs should be hard pinned to specific pCPUs or not.
* `is_cpu_passthrough_enabled`: (Optional) Indicates whether to passthrough the host CPU features to the guest or not. Enabling this will make VM incapable of live migration.
* `enabled_cpu_features`: (Optional) The list of additional CPU features to be enabled. HardwareVirtualization: Indicates whether hardware assisted virtualization should be enabled for the Guest OS or not. Once enabled, the Guest OS can deploy a nested hypervisor. Valid values are "HARDWARE_VIRTUALIZATION".
* `is_memory_overcommit_enabled`: (Optional) Indicates whether the memory overcommit feature should be enabled for the VM or not. If enabled, parts of the VM memory may reside outside of the hypervisor physical memory. Once enabled, it should be expected that the VM may suffer performance degradation.
* `is_gpu_console_enabled`: (Optional) Indicates whether the vGPU console is enabled or not.
* `is_cpu_hotplug_enabled`: (Optional) Indicates whether the VM CPU hotplug is enabled.
* `is_scsi_controller_enabled`: (Optional) Indicates whether the VM SCSI controller is enabled.
* `generation_uuid`: (Optional) Generation UUID of the VM. It should be of type UUID.
* `bios_uuid`: (Optional) BIOS UUID of the VM. It should be of type UUID.
* `categories`: (Optional) Categories for the VM.
* `ownership_info`: Ownership information for the VM.
* `host`: Reference to the host, the VM is running on.
* `cluster`: (Required) Reference to a cluster.
* `guest_customization`: (Optional) Stage a Sysprep or cloud-init configuration file to be used by the guest for the next boot. Note that the Sysprep command must be used to generalize the Windows VMs before triggering this API call.
* `guest_tools`: (Optional) The details about Nutanix Guest Tools for a VM.
* `hardware_clock_timezone`: (Optional) VM hardware clock timezone in IANA TZDB format (America/Los_Angeles).
* `is_branding_enabled`: (Optional) Indicates whether to remove AHV branding from VM firmware tables or not.
* `boot_config`: (Optional) Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order.
* `is_vga_console_enabled`: (Optional) Indicates whether the VGA console should be disabled or not.
* `machine_type`: (Optional) Machine type for the VM. Machine type Q35 is required for secure boot and does not support IDE disks. Valid values are "PSERIES", "Q35", "PC" .
* `vtpm_config`: (Optional) Indicates how the vTPM for the VM should be configured.
* `is_agent_vm`: (Optional) Indicates whether the VM is an agent VM or not. When their host enters maintenance mode, once the normal VMs are evacuated, the agent VMs are powered off. When the host is restored, agent VMs are powered on before the normal VMs are restored. In other words, agent VMs cannot be HA-protected or live migrated.
* `apc_config`: (Optional) Advanced Processor Compatibility configuration for the VM. Enabling this retains the CPU model for the VM across power cycles and migrations.
* `storage_config`: (Optional) Storage configuration for VM.
* `disks`: (Optional) Disks attached to the VM.
* `cd_roms`: (Optional) CD-ROMs attached to the VM.
* `nics`: (Optional) NICs attached to the VM.
* `gpus`: (Optional) GPUs attached to the VM.
* `serial_ports`: (Optional) Serial ports configured on the VM.
* `protection_type`: (Optional) The type of protection applied on a VM. Valid values "PD_PROTECTED", "UNPROTECTED", "RULE_PROTECTED".

### guest_customization
* `config`: (Required) The Nutanix Guest Tools customization settings.

* `config.sysprep`: (Optional) Sysprep config
* `config.cloud_init`: (Optional) CloudInit Config


### config.sysprep
* `install_type`: (Required) Indicates whether the guest will be freshly installed using this unattend configuration, or this unattend configuration will be applied to a pre-prepared image. Values allowed is 'PREPARED', 'FRESH'.

* `sysprep_script`: (Required) Object either UnattendXml or CustomKeyValues
* `sysprep_script.unattend_xml`: (Optional) xml object
* `sysprep_script.custom_key_values`: (Optional) The list of the individual KeyValuePair elements.


### config.cloud_init
* `datasource_type`: (Optional) Type of datasource. Default: CONFIG_DRIVE_V2
* `metadata`: The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded. Default value is 'CONFIG_DRIVE_V2'.
* `cloud_init_script`: (Optional) The script to use for cloud-init.
* `cloud_init_script.user_data`: (Optional) user data object
* `cloud_init_script.custom_keys`: (Optional) The list of the individual KeyValuePair elements.



### guest_tools
* `is_enabled`: (Optional) Indicates whether Nutanix Guest Tools is enabled or not.
* `capabilities`: (Optional) The list of the application names that are enabled on the guest VM.


### boot_config
* `legacy_boot`: (Optional) LegacyBoot config Object
* `uefi_boot`: (Optional) UefiBoot config Object

### boot_config.legacy_boot
* `boot_device`: (Required) Boot Device object
* `boot_device.boot_device_disk`: (Optional) Disk address.
* `boot_device.boot_device_disk.disk_address.bus_type`: (Required) Bus type for the device
* `boot_device.boot_device_disk.disk_address.index`: (Required) Device index on the bus. This field is ignored unless the bus details are specified.

* `boot_device.boot_device_nic`: (Optional) Disk Nic address.
* `boot_device.boot_device_nic.mac_address`: (Required) mac address

* `boot_order`: (Optional) Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order. Valid values are 'CDROM', 'DISK', 'NETWORK'.


### boot_config.uefi_boot
* `is_secure_boot_enabled`: (Optional) Indicate whether to enable secure boot or not
* `nvram_device`: (Optional) Configuration for NVRAM to be presented to the VM.
* `nvram_device.backing_storage_info`: (Required) Storage provided by Nutanix ADSF

### nvram_device.backing_storage_info
* `disk_size_bytes`: (Required) Size of the disk in Bytes
* `storage_container`: (Optional) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: (Optional) Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: (Required) Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: (Optional) A reference to a disk or image that contains the contents of a disk.


### vtpm_config
* `is_vtpm_enabled`: (Required) Indicates whether the virtual trusted platform module is enabled for the Guest OS or not.


### apc_config
* `is_apc_enabled`: (Optional) If enabled, the selected CPU model will be retained across live and cold migrations of the VM.
* `cpu_model`: (Optional) CPU model associated with the VM if Advanced Processor Compatibility(APC) is enabled. If APC is enabled and no CPU model is explicitly set, a default baseline CPU model is picked by the system. See the APC documentation for more information
* `cpu_model.name`: (Required) Name of the CPU model associated with the VM.


### storage_config
* `is_flash_mode_enabled`: (Optional) Indicates whether the virtual disk is pinned to the hot tier or not.
* `qos_config`: (Optional) QoS parameters to be enforced.
* `qos_config.throttled_iops`: (Optional) Throttled IOPS for the governed entities. The block size for the I/O is 32 kB.


### disks
* `disk_address`: (Optional) Disk address.
* `disk_address.bus_type`: (Required) Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
* `disk_address.index`: (Required) Device index on the bus. This field is ignored unless the bus details are specified.
* `backing_info`: (Required) Supporting storage to create virtual disk on.
* `backing_info.vm_disk`:(Optional) backing Info for vmDisk
* `backing_info.adfs_volume_group_reference`: (Required) Volume Group Reference
* `backing_info.adfs_volume_group_reference.volume_group_ext_id`: (Required) The globally unique identifier of an ADSF volume group. It should be of type UUID.

### backing_info.vm_disk
* `disk_size_bytes`: (Required) Size of the disk in Bytes
* `storage_container`: (Required) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: (Optional) Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: (Optional) A reference to a disk or image that contains the contents of a disk.

### backing_info.vm_disk.data_source
* `reference`: (Required) Reference to image or vm disk
* `image_reference`: (Optional) Image Reference
* `image_reference.image_ext_id`: (Required) The globally unique identifier of an image. It should be of type UUID.
* `vm_disk_reference`: (Optional) Vm Disk Reference
* `vm_disk_reference.disk_ext_id`: (Optional) The globally unique identifier of a VM disk. It should be of type UUID.
* `vm_disk_reference.disk_address`: (Optional) Disk address.
* `vm_disk_reference.vm_reference`: (Optional) This is a reference to a VM.


### cd_roms
* `disk_address`: (Optional) Virtual Machine disk (VM disk).
* `backing_info`: (Optional) Storage provided by Nutanix ADSF
* `iso_type`: Type of ISO image inserted in CD-ROM. Valid values "OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION" .


### nics
* `backing_info`: (Optional) Defines a NIC emulated by the hypervisor
* `network_info`: (Optional) Network information for a NIC.

### nics.backing_info
* `model`: (Optional) Options for the NIC emulation. Valid values "VIRTIO" , "E1000".
* `mac_address`: (Optional) MAC address of the emulated NIC.
* `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: (Optional) The number of Tx/Rx queue pairs for this NIC. Default is 1.

### nics.network_info
* `nic_type`: (Optional) NIC type. Valid values "SPAN_DESTINATION_NIC",  "NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC" .
* `network_function_chain`: (Optional) The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_nic_type`: (Optional) The type of this Network function NIC. Defaults to INGRESS.
* `subnet`: (Required) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC
* `vlan_mode`: (Required) all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs.
* `trunked_vlans`: (Optional) List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
* `should_allow_unknown_macs`: (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
* `ipv4_config`: (Optional) The IP address configurations.

### gpus
* `mode`: (Required) The mode of this GPU. Valid values "PASSTHROUGH_GRAPHICS", "PASSTHROUGH_COMPUTE", "VIRTUAL" .
* `device_id`: (Optional) The device Id of the GPU.
* `vendor`: (Optional) The vendor of the GPU. Valid values "NVIDIA", "AMD", "INTEL" .
* `pci_address`: (Optional) The (S)egment:(B)us:(D)evice.(F)unction hardware address.

### gpus.pci_address
* `segment`
* `bus`
* `device`
* `func`

### serial_ports
* `is_connected`: (Optional) Indicates whether the serial port is connected or not.
* `index`: (Required) Index of the serial port.

## Attributes Reference

The following attributes are exported:
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `name`: VM name.
* `description`: VM description
* `create_time`: VM creation time
* `update_time`: VM last updated time.
* `source`: Reference to an entity that the VM should be cloned or created from
* `num_sockets`: Number of vCPU sockets.
* `num_cores_per_socket`: Number of cores per socket.
* `num_threads_per_core`: Number of threads per core
* `num_numa_nodes`: Number of NUMA nodes. 0 means NUMA is disabled.
* `memory_size_bytes`: Memory size in bytes.
* `is_vcpu_hard_pinning_enabled`: Indicates whether the vCPUs should be hard pinned to specific pCPUs or not.
* `is_cpu_passthrough_enabled`: Indicates whether to passthrough the host CPU features to the guest or not. Enabling this will make VM incapable of live migration.
* `enabled_cpu_features`: The list of additional CPU features to be enabled. HardwareVirtualization: Indicates whether hardware assisted virtualization should be enabled for the Guest OS or not. Once enabled, the Guest OS can deploy a nested hypervisor
* `is_memory_overcommit_enabled`: Indicates whether the memory overcommit feature should be enabled for the VM or not. If enabled, parts of the VM memory may reside outside of the hypervisor physical memory. Once enabled, it should be expected that the VM may suffer performance degradation.
* `is_gpu_console_enabled`: Indicates whether the vGPU console is enabled or not.
* `is_cpu_hotplug_enabled`: Indicates whether the VM CPU hotplug is enabled.
* `is_scsi_controller_enabled`: Indicates whether the VM SCSI controller is enabled.
* `generation_uuid`: Generation UUID of the VM. It should be of type UUID.
* `bios_uuid`: BIOS UUID of the VM. It should be of type UUID.
* `categories`: Categories for the VM.
* `ownership_info`: Ownership information for the VM.
* `host`: Reference to the host, the VM is running on.
* `cluster`: Reference to a cluster.
* `guest_tools`: The details about Nutanix Guest Tools for a VM.
* `hardware_clock_timezone`: VM hardware clock timezone in IANA TZDB format (America/Los_Angeles).
* `is_branding_enabled`: Indicates whether to remove AHV branding from VM firmware tables or not.
* `boot_config`: Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order.
* `is_vga_console_enabled`: Indicates whether the VGA console should be disabled or not.
* `machine_type`: Machine type for the VM. Machine type Q35 is required for secure boot and does not support IDE disks.
* `vtpm_config`: Indicates how the vTPM for the VM should be configured.
* `is_agent_vm`: Indicates whether the VM is an agent VM or not. When their host enters maintenance mode, once the normal VMs are evacuated, the agent VMs are powered off. When the host is restored, agent VMs are powered on before the normal VMs are restored. In other words, agent VMs cannot be HA-protected or live migrated.
* `apc_config`: Advanced Processor Compatibility configuration for the VM. Enabling this retains the CPU model for the VM across power cycles and migrations.
* `storage_config`: Storage configuration for VM.
* `disks`: Disks attached to the VM.
* `cd_roms`: CD-ROMs attached to the VM.
* `nics`: NICs attached to the VM.
* `gpus`: GPUs attached to the VM.
* `serial_ports`: Serial ports configured on the VM.
* `protection_type`: The type of protection applied on a VM. PD_PROTECTED indicates a VM is protected using the Prism Element. RULE_PROTECTED indicates a VM protection using the Prism Central.
* `protection_policy_state`: Status of protection policy applied to this VM.

See detailed information in [Nutanix Virtual Machine V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0).


## Import
Nutanix Virtual machine can be imported using the `UUID` eg,

```
terraform import nutanix_virtual_machine_v2.imported_vm 32679dc5-1309-4728-b135-402dc5234907
```


