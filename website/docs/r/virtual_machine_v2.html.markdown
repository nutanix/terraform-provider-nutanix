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
    name= "example-vm-1"
    description =  "vm desc"
    num_cores_per_socket = 1
    num_sockets = 1
    cluster {
        ext_id = "1cefd0f5-6d38-4c9b-a07c-bdd2db004224"
    }
}

resource "nutanix_virtual_machine_v2" "vm-2"{
    name= "example-vm-2"
    description =  "vm desc"
    num_cores_per_socket = 1
    num_sockets = 1
    cluster {
        ext_id = "1cefd0f5-6d38-4c9b-a07c-bdd2db004224"
    }
    project {
    ext_id = "2defe0f5-6e48-4c9b-b07c-bdd2dc004225"
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
                    ext_id = "1cefd0f5-6d38-4c9b-a07c-bdd2db004224"
                }
            }
        }
    }
    boot_config {
        uefi_boot {
            boot_order = ["NETWORK", "DISK", "CDROM", ]
        }
    }
}

resource "nutanix_virtual_machine_v2" "vm-3" {
  name                 = "terraform-example-vm-4-disks"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = "1cefd0f5-6d38-4c9b-a07c-bdd2db004224"
  }
  project {
    ext_id = "2defe0f5-6e48-4c9b-b07c-bdd2dc004225"
  }

  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = "59ec786c-4311-4225-affe-68b65c5ebf10"
            }
          }
        }
        disk_size_bytes = 20 * pow(1024, 3) # 20 GB
      }
    }
  }
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 1
    }
     backing_info {
      vm_disk {
        disk_size_bytes = 10 * pow(1024, 3) # 10 GB
        storage_container {
          ext_id = "5d9b5941-fec3-4996-9d31-f31bed1c7735"
        }
      }
    }
  }

  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 2
    }
    backing_info {
      vm_disk {
        disk_size_bytes = 15 * pow(1024, 3) # 15 GB
        storage_container {
          ext_id = "5d9b5941-fec3-4996-9d31-f31bed1c7735"
        }
      }
    }
  }

  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 3
    }
    backing_info {
      vm_disk {
        disk_size_bytes = 20 * pow(1024, 3) # 20 GB
        storage_container {
          ext_id = "5d9b5941-fec3-4996-9d31-f31bed1c7735"
        }
      }
    }
  }

  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = "7f66e20f-67f4-473f-96bb-c4fcfd487f16"
      }
      vlan_mode = "ACCESS"
    }
  }

  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  power_state = "ON"
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
* `project`: (Optional) Reference to a project.
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


### Source

The `source` attribute supports the following:

* `entity_type`: (Optional) Reference to an entity from which the VM should be cloned or created. Values are:
  - VM_RECOVERY_POINT: Reference to the recovery point entity from which the VM should be cloned or created.
  - VM: Reference to an entity from which the VM should be cloned or created.

### Categories
The `categories` attribute supports the following:

* `ext_id`: A globally unique identifier of a VM category of type UUID.

### Project
The `project` attribute supports the following:

* `ext_id`: The globally unique identifier of an instance of type UUID.


### Ownership Info
The `ownership_info` attribute supports the following:

* `owner`: Reference to the owner.
* `owner.ext_id`: A globally unique identifier of a VM owner type UUID.

### Host
The `host` attribute supports the following:

* `ext_id`: A globally unique identifier of a host of type UUID.

### Cluster
> 💡Cluster automatic selection is supported.

The `cluster` attribute supports the following:

* `ext_id`: The globally unique identifier of a cluster type UUID.

### Availability Zone
The `availability_zone` attribute supports the following:

* `ext_id`: The globally unique identifier of an availability zone type UUID.


### guest_customization
* `config`: (Required) The Nutanix Guest Tools customization settings.

* `config.sysprep`: (Optional) Sysprep config
* `config.cloud_init`: (Optional) CloudInit Config



### Guest Customization
The `guest_customization` attribute supports the following:

* `config`: The Nutanix Guest Tools customization settings.

* `config.sysprep`: Sysprep config
* `config.cloud_init`: CloudInit Config

#### config.sysprep
* `install_type`: (Optional) Indicates whether the guest will be freshly installed using this unattend configuration, or this unattend configuration will be applied to a pre-prepared image. Values allowed is 'PREPARED', 'FRESH'.
* `sysprep_script`: (Optional) Object either UnattendXml or CustomKeyValues
* `sysprep_script.unattend_xml`: (Optional) xml object
* `sysprep_script.unattend_xml.value`: (Optional) base64 encoded sysprep unattended xml
* `sysprep_script.custom_key_values`: (Optional) The list of the individual KeyValuePair elements.


### config.cloud_init
* `datasource_type`: (Optional) Type of datasource. Default: CONFIG_DRIVE_V2
* `metadata`: The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded. Default value is 'CONFIG_DRIVE_V2'.
* `cloud_init_script`: (Optional) The script to use for cloud-init.
* `cloud_init_script.user_data`: (Optional) user data object
* `cloud_init_script.user_data.value`: (Optional) base64 encoded cloud init script as string
* `cloud_init_script.custom_keys`: (Optional) The list of the individual KeyValuePair elements.

#### custom_keys
* `name`: (Optional) The name of the key.
* `value`: (Optional) The value of the key. value can be a:
    - String
    - Integer
    - Boolean
    - Array of strings
    - Object
    - Map of string wrapper
    - Array of integers

### Guest Tools
The `guest_tools` attribute supports the following:

* `is_enabled`: (Optional) Indicates whether Nutanix Guest Tools is enabled or not.
* `capabilities`: (Optional) The list of the application names that are enabled on the guest VM.


### Boot Config
The `boot_config` attribute supports the following:

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



### VTPM Config
The `vtpm_config` attribute supports the following:

* `is_vtpm_enabled`: (Required) Indicates whether the virtual trusted platform module is enabled for the Guest OS or not.


### APC Config
The `apc_config` attribute supports the following:

* `is_apc_enabled`: (Optional) If enabled, the selected CPU model will be retained across live and cold migrations of the VM.
* `cpu_model`: (Optional) CPU model associated with the VM if Advanced Processor Compatibility(APC) is enabled. If APC is enabled and no CPU model is explicitly set, a default baseline CPU model is picked by the system. See the APC documentation for more information
* `cpu_model.name`: (Required) Name of the CPU model associated with the VM.


### Storage Config
The `storage_config` attribute supports the following:

* `is_flash_mode_enabled`: (Optional) Indicates whether the virtual disk is pinned to the hot tier or not.
* `qos_config`: (Optional) QoS parameters to be enforced.
* `qos_config.throttled_iops`: (Optional) Throttled IOPS for the governed entities. The block size for the I/O is 32 kB.


### Disks
The `disks` attribute supports the following:

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


### CD-ROMs
The `cd_roms` attribute supports the following:

* `disk_address`: (Optional) Virtual Machine disk (VM disk).
* `backing_info`: (Optional) Storage provided by Nutanix ADSF
* `iso_type`: Type of ISO image inserted in CD-ROM. Valid values "OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION" .


### NICs
The `nics` attribute supports the following:

* `nic_backing_info`: (Optional) New NIC backing info (v2.4.1+). One of `virtual_ethernet_nic`, `sriov_nic`, `dp_offload_nic`.
* `nic_network_info`: (Optional) New NIC network info (v2.4.1+). One of `virtual_ethernet_nic_network_info`, `sriov_nic_network_info`, `dp_offload_nic_network_info`.
* `backing_info`: (Optional, Deprecated) Use `nic_backing_info.virtual_ethernet_nic` instead.
* `network_info`: (Optional, Deprecated) Use `nic_network_info.virtual_ethernet_nic_network_info` instead.

### nics.nic_backing_info.virtual_ethernet_nic
* `model`: (Optional) Options for the NIC emulation. Valid values "VIRTIO", "E1000".
* `mac_address`: (Optional) MAC address of the emulated NIC.
* `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: (Optional) The number of Tx/Rx queue pairs for this NIC. Default is 1.

### nics.nic_backing_info.sriov_nic
* `sriov_profile_reference`: (Required) SR-IOV profile reference.
* `host_pcie_device_reference`: (Optional) Host PCIe device reference.
* `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
* `mac_address`: (Optional) MAC address of the SR-IOV NIC.

### nics.nic_backing_info.dp_offload_nic
* `dp_offload_profile_reference`: (Required) DP offload profile reference.
* `host_pcie_device_reference`: (Optional) Host PCIe device reference.
* `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
* `mac_address`: (Optional) MAC address of the DP offload NIC.

### nics.nic_network_info.virtual_ethernet_nic_network_info
* `nic_type`: (Optional) NIC type. Valid values "SPAN_DESTINATION_NIC", "NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC".
* `network_function_chain`: (Optional) The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_nic_type`: (Optional) The type of this Network function NIC. Defaults to INGRESS.
* `subnet`: (Optional) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC.
* `subnet.ext_id`: (Optional) The globally unique identifier of a subnet of type UUID.
* `vlan_mode`: (Optional) All the virtual NICs are created in ACCESS mode by default. TRUNKED allows multiple VLANs.
* `trunked_vlans`: (Optional) List of networks to trunk if VLAN mode is TRUNKED.
* `should_allow_unknown_macs`: (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not.
* `ipv4_config`: (Optional) The IP address configurations.

### nics.nic_network_info.sriov_nic_network_info
* `vlan_id`: (Optional) VLAN ID for the SR-IOV NIC.

### nics.nic_network_info.dp_offload_nic_network_info
* `subnet`: (Optional) Network identifier for this adapter.
* `vlan_mode`: (Optional) VLAN mode for DP offload NIC.
* `trunked_vlans`: (Optional) List of networks to trunk if VLAN mode is TRUNKED.
* `should_allow_unknown_macs`: (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not.
* `ipv4_config`: (Optional) The IP address configurations.

### nics.backing_info
* `model`: (Optional) Options for the NIC emulation. Valid values "VIRTIO" , "E1000".
* `mac_address`: (Optional) MAC address of the emulated NIC.
* `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: (Optional) The number of Tx/Rx queue pairs for this NIC. Default is 1.

### nics.network_info
* `nic_type`: (Optional) NIC type. Valid values "SPAN_DESTINATION_NIC",  "NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC" .
* `network_function_chain`: (Optional) The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_nic_type`: (Optional) The type of this Network function NIC. Defaults to INGRESS.
* `subnet`: (Optional) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC
* `subnet.ext_id`: (Optional) The globally unique identifier of a subnet of type UUID.
* `vlan_mode`: (Optional) all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs.
* `trunked_vlans`: (Optional) List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
* `should_allow_unknown_macs`: (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
* `ipv4_config`: (Optional) The IP address configurations.

#### nics.ipv4_config
* `should_assign_ip`: If set to true (default value), an IP address must be assigned to the VM NIC - either the one explicitly specified by the user or allocated automatically by the IPAM service by not specifying the IP address. If false, then no IP assignment is required for this VM NIC.
* `ip_address`: The IP address of the NIC.
* `secondary_ip_address_list`: Secondary IP addresses for the NIC.

##### ip_address, secondary_ip_address_list
* `value`: The IPv4 address of the host.
* `prefix_length`: The prefix length of the IP address.

### gpus
* `mode`: ((Optional)) The mode of this GPU. Valid values "PASSTHROUGH_GRAPHICS", "PASSTHROUGH_COMPUTE", "VIRTUAL" .
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
* `index`: ((Optional)) Index of the serial port.

### protection_policy_state
* `policy`: (Optional) Reference to the policy object in use.
* `policy.ext_id`: (Optional) Reference to the policy object in use.


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

See detailed information in [Nutanix Create Virtual Machine V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/createVm).
