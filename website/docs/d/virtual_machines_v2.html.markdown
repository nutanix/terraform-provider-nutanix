---
layout: "nutanix"
page_title: "NUTANIX: nutanix_virtual_machines_v2"
sidebar_current: "docs-nutanix-datasource-virtual-machines-v2"
description: |-
  This operation retrieves a list of all the virtual machines.
---

# nutanix_virtual_machines_v2

Lists the Virtual Machines defined on the system. List of Virtual Machines can be further filtered out using various filtering options.

## Example

```hcl
# List all virtual machines
data "nutanix_virtual_machines_v2" "list-vms"{}

# List virtual machines with  page and limit
data "nutanix_virtual_machines_v2" "paged-vms"{
    page=0
    limit=2
}

# List virtual machines with filter
data "nutanix_virtual_machines_v2" "filtered-vms"{
    filter = "name eq 'example-vm-filter'"
}
```

## Argument Reference
The following arguments are supported:

- `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`:A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - apcConfig/cpuModel/extId
  - apcConfig/cpuModel/name
  - apcConfig/isApcEnabled
  - cluster/extId
  - extId
  - gpus
  - hardwareClockTimezone
  - host/extId
  - isAgentVm
  - isCpuHotplugEnabled
  - isCpuPassthroughEnabled
  - isGpuConsoleEnabled
  - isLiveMigrateCapable
  - isMemoryOvercommitEnabled
  - isScsiControllerEnabled
  - isVcpuHardPinningEnabled
  - isVgaConsoleEnabled
  - machineType
  - memorySizeBytes
  - name
  - nics
  - numCoresPerSocket
  - numNumaNodes
  - numSockets
  - numThreadsPerCore
  - powerState
- `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - memorySizeBytes
  - name
  - numCoresPerSocket
  - numNumaNodes
  - numSockets
  - numThreadsPerCore
- `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - apcConfig/cpuModel/extId
  - apcConfig/cpuModel/name
  - apcConfig/isApcEnabled
  - availabilityZone
  - biosUuid
  - bootConfig
  - categories
  - cdRoms
  - cluster/extId
  - createTime
  - description
  - disks
  - enabledCpuFeatures
  - extId
  - generationUuid
  - gpus
  - guestCustomization
  - guestTools
  - hardwareClockTimezone
  - host/extId
  - isAgentVm
  - isBrandingEnabled
  - isCpuHotplugEnabled
  - isCpuPassthroughEnabled
  - isCrossClusterMigrationInProgress
  - isGpuConsoleEnabled
  - isLiveMigrateCapable
  - isMemoryOvercommitEnabled
  - isScsiControllerEnabled
  - isVcpuHardPinningEnabled
  - isVgaConsoleEnabled
  - links
  - machineType
  - memorySizeBytes
  - name
  - nics
  - numCoresPerSocket
  - numNumaNodes
  - numSockets
  - numThreadsPerCore
  - ownershipInfo
  - pcieDevices
  - powerState
  - protectionPolicyState
  - protectionType
  - serialPorts
  - source
  - storageConfig
  - tenantId
  - updateTime
  - vtpmConfig

## Attribute Reference

The following attributes are exported:

- `vms`: List of all vms

## VMs
The `vms` is a list of all virtual machines. Each virtual machine has the following attributes:

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `name`: VM name.
- `description`: VM description
- `create_time`: VM creation time
- `update_time`: VM last updated time.
- `source`: Reference to an entity that the VM should be cloned or created from
- `num_sockets`: Number of vCPU sockets.
- `num_cores_per_socket`: Number of cores per socket.
- `num_threads_per_core`: Number of threads per core
- `num_numa_nodes`: Number of NUMA nodes. 0 means NUMA is disabled.
- `memory_size_bytes`: Memory size in bytes.
- `is_vcpu_hard_pinning_enabled`: Indicates whether the vCPUs should be hard pinned to specific pCPUs or not.
- `is_cpu_passthrough_enabled`: Indicates whether to passthrough the host CPU features to the guest or not. Enabling this will make VM incapable of live migration.
- `enabled_cpu_features`: The list of additional CPU features to be enabled. HardwareVirtualization: Indicates whether hardware assisted virtualization should be enabled for the Guest OS or not. Once enabled, the Guest OS can deploy a nested hypervisor
- `is_memory_overcommit_enabled`: Indicates whether the memory overcommit feature should be enabled for the VM or not. If enabled, parts of the VM memory may reside outside of the hypervisor physical memory. Once enabled, it should be expected that the VM may suffer performance degradation.
- `is_gpu_console_enabled`: Indicates whether the vGPU console is enabled or not.
- `is_cpu_hotplug_enabled`: Indicates whether the VM CPU hotplug is enabled.
- `is_scsi_controller_enabled`: Indicates whether the VM SCSI controller is enabled.
- `generation_uuid`: Generation UUID of the VM. It should be of type UUID.
- `bios_uuid`: BIOS UUID of the VM. It should be of type UUID.
- `categories`: Categories for the VM.
* `project`: Reference to a project.
- `ownership_info`: Ownership information for the VM.
- `host`: Reference to the host, the VM is running on.
- `cluster`: Reference to a cluster.
- `guest_customization`: Stage a Sysprep or cloud-init configuration file to be used by the guest for the next boot. Note that the Sysprep command must be used to generalize the Windows VMs before triggering this API call.
- `guest_tools`: The details about Nutanix Guest Tools for a VM.
- `hardware_clock_timezone`: VM hardware clock timezone in IANA TZDB format (America/Los_Angeles).
- `is_branding_enabled`: Indicates whether to remove AHV branding from VM firmware tables or not.
- `boot_config`: Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order.
- `is_vga_console_enabled`: Indicates whether the VGA console should be disabled or not.
- `machine_type`: Machine type for the VM. Machine type Q35 is required for secure boot and does not support IDE disks.
- `vtpm_config`: Indicates how the vTPM for the VM should be configured.
- `is_agent_vm`: Indicates whether the VM is an agent VM or not. When their host enters maintenance mode, once the normal VMs are evacuated, the agent VMs are powered off. When the host is restored, agent VMs are powered on before the normal VMs are restored. In other words, agent VMs cannot be HA-protected or live migrated.
- `apc_config`: Advanced Processor Compatibility configuration for the VM. Enabling this retains the CPU model for the VM across power cycles and migrations.
- `storage_config`: Storage configuration for VM.
- `disks`: Disks attached to the VM.
- `cd_roms`: CD-ROMs attached to the VM.
- `nics`: NICs attached to the VM.
- `gpus`: GPUs attached to the VM.
- `serial_ports`: Serial ports configured on the VM.
- `protection_type`: The type of protection applied on a VM. PD_PROTECTED indicates a VM is protected using the Prism Element. RULE_PROTECTED indicates a VM protection using the Prism Central.
- `protection_policy_state`: Status of protection policy applied to this VM.

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Source

The `source` attribute supports the following:

* `entity_type`: Reference to an entity from which the VM should be cloned or created. Values are:
  - VM_RECOVERY_POINT: Reference to the recovery point entity from which the VM should be cloned or created.
  - VM: Reference to an entity from which the VM should be cloned or created.
* `ext_id`: A globally unique identifier of a VM of type UUID.

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
The `cluster` attribute supports the following:

* `ext_id`: The globally unique identifier of a cluster type UUID.

### Availability Zone
The `availability_zone` attribute supports the following:

* `ext_id`: The globally unique identifier of an availability zone type UUID.

### Guest Customization
The `guest_customization` attribute supports the following:

* `config`: The Nutanix Guest Tools customization settings.

* `config.sysprep`: Sysprep config
* `config.cloud_init`: CloudInit Config



#### config.sysprep
* `install_type`: Indicates whether the guest will be freshly installed using this unattend configuration, or this unattend configuration will be applied to a pre-prepared image. Default is 'PREPARED'.
* `sysprep_script`: Object either UnattendXml or CustomKeyValues
* `sysprep_script.unattend_xml`: xml object
* `sysprep_script.custom_key_values`: The list of the individual KeyValuePair elements.


#### config.cloud_init
* `datasource_type`: Type of datasource. Default: CONFIG_DRIVE_V2
* `metadata`: The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded.
* `cloud_init_script`: The script to use for cloud-init.
* `cloud_init_script.user_data`: user data object
* `cloud_init_script.custom_keys`: The list of the individual KeyValuePair elements.


### Guest Tools
The `guest_tools` attribute supports the following:

* `version`: Version of Nutanix Guest Tools installed on the VM.
* `is_installed`: Indicates whether Nutanix Guest Tools is installed on the VM or not.
* `is_iso_inserted`: Indicates whether Nutanix Guest Tools ISO is inserted or not.
* `available_version`: Version of Nutanix Guest Tools available on the cluster.
* `guest_os_version`: Version of the operating system on the VM
* `is_reachable`: Indicates whether the communication from VM to CVM is active or not.
* `is_vss_snapshot_capable`: Indicates whether the VM is configured to take VSS snapshots through NGT or not.
* `is_vm_mobility_drivers_installed`: Indicates whether the VM mobility drivers are installed on the VM or not.
* `is_enabled`: Indicates whether Nutanix Guest Tools is enabled or not.
* `capabilities`: The list of the application names that are enabled on the guest VM.

### Boot Config
The `boot_config` attribute supports the following:

* `legacy_boot`: LegacyBoot config Object
* `uefi_boot`: UefiBoot config Object

#### boot_config.legacy_boot
* `boot_device`: Boot Device object
* `boot_device.boot_device_disk`: Disk address.
* `boot_device.boot_device_disk.disk_address.bus_type`: Bus type for the device
* `boot_device.boot_device_disk.disk_address.index`: Device index on the bus. This field is ignored unless the bus details are specified.

* `boot_device.boot_device_nic`: Disk Nic address.
* `boot_device.boot_device_nic.mac_address`: mac address

* `boot_order`: Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order.


#### boot_config.uefi_boot
* `is_secure_boot_enabled`: Indicate whether to enable secure boot or not
* `nvram_device`: Configuration for NVRAM to be presented to the VM.
* `nvram_device.backing_storage_info`: Storage provided by Nutanix ADSF

##### nvram_device.backing_storage_info
* `disk_ext_id`: The globally unique identifier of a VM disk. It should be of type UUID.
* `disk_size_bytes`: Size of the disk in Bytes
* `storage_container`: This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_config`: Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: A reference to a disk or image that contains the contents of a disk.
* `is_migration_in_progress`: Indicates if the disk is undergoing migration to another container.

### VTPM Config
The `vtpm_config` attribute supports the following:

* `is_vtpm_enabled`: Indicates whether the virtual trusted platform module is enabled for the Guest OS or not.
* `version`: Virtual trusted platform module version.

### APC Config
The `apc_config` attribute supports the following:

* `is_apc_enabled`: If enabled, the selected CPU model will be retained across live and cold migrations of the VM.
* `cpu_model`: CPU model associated with the VM if Advanced Processor Compatibility(APC) is enabled. If APC is enabled and no CPU model is explicitly set, a default baseline CPU model is picked by the system. See the APC documentation for more information
* `cpu_model.ext_id`: The globally unique identifier of the CPU model associated with the VM.
* `cpu_model.name`: Name of the CPU model associated with the VM.


### Storage Config
The `storage_config` attribute supports the following:

* `is_flash_mode_enabled`: Indicates whether the virtual disk is pinned to the hot tier or not.
* `qos_config`: QoS parameters to be enforced.
* `qos_config.throttled_iops`: Throttled IOPS for the governed entities. The block size for the I/O is 32 kB.

### Disks
The `disks` attribute supports the following:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `disk_address`: Disk address.
* `disk_address.bus_type`: Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
* `disk_address.index`: Device index on the bus. This field is ignored unless the bus details are specified.
* `backing_info`: Supporting storage to create virtual disk on.
* `backing_info.vm_disk`: backing Info for vmDisk
* `backing_info.adfs_volume_group_reference`: Volume Group Reference
* `backing_info.adfs_volume_group_reference.volume_group_ext_id`: The globally unique identifier of an ADSF volume group. It should be of type UUID.


#### backing_info.vm_disk
* `disk_ext_id`: The globally unique identifier of a VM disk. It should be of type UUID.
* `disk_size_bytes`: Size of the disk in Bytes
* `storage_container`: This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
* `storage_container.ext_id`: A globally unique identifier of a VM disk container. It should be of type UUID.
* `storage_config`: Storage configuration for VM disks
* `storage_config.is_flash_mode_enabled`: Indicates whether the virtual disk is pinned to the hot tier or not.
* `data_source`: A reference to a disk or image that contains the contents of a disk.
* `is_migration_in_progress`: Indicates if the disk is undergoing migration to another container.

#### backing_info.vm_disk.data_source
* `reference`: Reference to image or vm disk
* `reference.image_reference`: Image Reference
* `reference.image_reference.image_ext_id`: The globally unique identifier of an image. It should be of type UUID.
* `reference.vm_disk_reference`: Vm Disk Reference
* `reference.vm_disk_reference.disk_ext_id`:  The globally unique identifier of a VM disk. It should be of type UUID.
* `reference.vm_disk_reference.disk_address`: Disk address.
* `reference.vm_disk_reference.disk_address.bus_type`: Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
* `reference.vm_disk_reference.disk_address.index`: Device index on the bus. This field is ignored unless the bus details are specified.
* `reference.vm_disk_reference.vm_reference`: This is a reference to a VM.
* `reference.vm_disk_reference.vm_reference.ext_id`: A globally unique identifier of a VM of type UUID.




### CD-ROMs
The `cd_roms` attribute supports the following:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `disk_address`: Virtual Machine disk (VM disk).
* `backing_info`: Storage provided by Nutanix ADSF
* `iso_type`: Type of ISO image inserted in CD-ROM

### NICs
The `nics` attribute supports the following:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption
* `backing_info`: Defines a NIC emulated by the hypervisor
* `network_info`: Network information for a NIC.

### nics.backing_info
* `model`: Options for the NIC emulation.
* `mac_address`: MAC address of the emulated NIC.
* `is_connected`: Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: The number of Tx/Rx queue pairs for this NIC

### nics.network_info
* `nic_type`: NIC type. Defaults to NORMAL_NIC. The acceptable values are: SPAN_DESTINATION_NIC, NORMAL_NIC, DIRECT_NIC, NETWORK_FUNCTION_NIC.
* `network_function_chain`: The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_chain.ext_id`: The globally unique identifier of a network function chain. It should be of type UUID.
* `network_function_nic_type`: The type of this Network function NIC. Defaults to INGRESS.  values are: TAP, EGRESS, INGRESS.
* `subnet`: Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC.
* `subnet.ext_id`: The globally unique identifier of a subnet of type UUID.
* `vlan_mode`: all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs. values are: ACCESS, TRUNKED.
* `trunked_vlans`: List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
* `should_allow_unknown_macs`: Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
* `ipv4_config`: The IP address configurations.
* `ipv4_info`: The runtime IP address information of the NIC.

#### nics.ipv4_config
* `should_assign_ip`: If set to true (default value), an IP address must be assigned to the VM NIC - either the one explicitly specified by the user or allocated automatically by the IPAM service by not specifying the IP address. If false, then no IP assignment is required for this VM NIC.
* `ip_address`: The IP address of the NIC.
* `secondary_ip_address_list`: Secondary IP addresses for the NIC.

##### ip_address, secondary_ip_address_list
* `value`: The IPv4 address of the host.
* `prefix_length`: The prefix length of the IP address.

#### nics.ipv4_info
* `learned_ip_addresses`: The list of IP addresses learned by the NIC.

##### learned_ip_addresses
* `value`: The IPv4 address of the host.
* `prefix_length`: The prefix length of the IP address.


### gpus
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `mode`: The mode of this GPU.
* `device_id`: The device Id of the GPU.
* `vendor`: The vendor of the GPU.
* `pci_address`: The (S)egment:(B)us:(D)evice.(F)unction hardware address. See
* `guest_driver_version`: Last determined guest driver version.
* `name`: Name of the GPU resource.
* `frame_buffer_size_bytes`: GPU frame buffer size in bytes.
* `num_virtual_display_heads`: Number of supported virtual display heads.
* `fraction`: Fraction of the physical GPU assigned.

### gpus.pci_address
* `segment`
* `bus`
* `device`
* `func`

### serial_ports
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `is_connected`: Indicates whether the serial port is connected or not.
* `index`: Index of the serial port.

### protection_policy_state
* `policy`: Reference to the policy object in use.




See detailed information in [Nutanix List Virtual Machines V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/listVms).
