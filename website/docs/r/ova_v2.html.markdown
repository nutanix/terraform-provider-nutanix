---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ova_v2 "
sidebar_current: "docs-nutanix-resource-ova-v2"
description: |-
  Creates an OVA using the provided request body. The name and source are mandatory fields to create an OVA.
---

# nutanix_ova_v2

Creates an OVA using the provided request body. The name and source are mandatory fields to create an OVA.

## Example Usage

```hcl

// Create a new OVA using the ova_vm_source
resource "nutanix_ova_v2 "ova-vm" {
  name = "tf-example-ova"
  source {
    ova_vm_source {
      vm_ext_id        = "7034016e-f3d4-472a-8c7b-cd13658b7315"
      disk_file_format = "QCOW2"
    }
  }
}

// Create a new OVA using the ova_url_source
resource "nutanix_ova_v2" "ova-url"{
  name = "tf-example-ova-url"
  source {
    ova_url_source {
      url                       = "https://example.com/path/to/ova/file.ova"
      should_allow_insecure_url = true
    }
  }
}

// Create a new OVA using the object_lite_source
resource "nutanix_ova_v2" "ova-object-lite"{
  name = "tf-example-ova-object-lite"
  source {
    object_lite_source {
      key = "object-key"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

- `name`: -(Required) Name of the OVA.
- `checksum`: -(Optional) The checksum of an OVA.
- `source`: -(Required) Source of the created OVA file. The source can either be a VM, URL, or a local upload.
- `created_by`: -(Optional) Information of the user.
- `cluster_location_ext_ids`: -(Optional) List of cluster identifiers where the OVA is located. This field is required when creating an OVA from URL or Objects lite upload. its `mandatory` when creating an OVA from URL or object lite source .
- `vm_config`: -(Optional) VM configuration.
- `disk_format`: -(Optional) Disk format of an OVA.
  |ENUM |Description |
  |---|---|
  | VMDK | The VMDK disk format of an OVA. |
  | QCOW2 | The QCOW2 disk format of an OVA. |

### checksum

The checksum argument supports the following :

- `ova_sha1_checksum`: -(Optional) The SHA1 checksum of the OVA file.
- `ova_sha256_checksum`: -(Optional) The SHA256 checksum of the OVA file.

#### ova_sha1_checksum, ova_sha256_checksum

The `ova_sha1_checksum` and `ova_sha256_checksum` arguments support the following:

- `hex_digest`: -(Required) The hexadecimal representation of the checksum.

### source

The `source` argument should be one of the following:

- `ova_url_source`: -(Optional) The source of the OVA file when it is being created from a URL.
- `ova_vm_source`: -(Optional) The source of the OVA file when it is being created from a VM.
- `object_lite_source`: -(Optional) The source of the OVA file when it is being created from an object lite upload.

#### Ova Url Source

The `ova_url_source` argument supports the following:

- `url`: -(Required) The URL from which the OVA file can be downloaded.
- `should_allow_insecure_url`: -(Optional) Flag to allow insecure URLs.
- `basic_auth`: -(Optional) Basic authentication credentials for accessing the OVA file.

##### Basic Auth

The `basic_auth` argument supports the following:

- `username`: -(Required) The username for basic authentication.
- `password`: -(Required) The password for basic authentication.

#### Ova VM Source

The `ova_vm_source` argument supports the following:

- `vm_ext_id`: -(Required) The external identifier of the VM from which the OVA file is being created.
- `disk_file_format`: -(Required) The disk file format of the VM.

#### Object Lite Source

The `object_lite_source` argument supports the following:

- `key`: -(Required) The identifier of the object from which the OVA file is being created.

### created_by

The `created_by` argument supports the following:

- `username`: -(Required) Identifier for the User in the form an email address.
- `user_type`: -(Required) Type of the User.
- `idp_id`: -(Optional) Identifier of the IDP for the User.
- `display_name`: -(Optional) Display name for the User.
- `first_name`: -(Optional) First name for the User.
- `middle_initial`: -(Optional) Middle name for the User.
- `last_name`: -(Optional) Last name for the User.
- `email_id`: -(Optional) Email Id for the User.
- `locale`: -(Optional) Default locale for the User.
- `region`: -(Optional) Default Region for the User.
- `password`: -(Optional) Password of the user.
- `is_force_reset_password_enabled`: -(Optional) Flag to force the User to reset password.
- `additional_attributes`: -(Optional) Any additional attribute for the User.
- `status`: -(Optional) Status of the User.
- `description`: -(Optional) Description of the User.
- `creation_type`: -(Optional) Creation type of the User.
  |ENUM |Description |
  |---|---|
  | PREDEFINED | Predefined creator workflow type is for entity created by the system. |
  | SERVICEDEFINED | Service defined creator workflow type is for entity created by the service. |
  | USERDEFINED | User defined creator workflow type is for entity created by the users. |

#### Additional Attributes

The `additional_attributes` argument supports the following:

- `name`: -(Optional) The URL at which the entity described by the link can be accessed.
- `value`: -(Optional) A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### VM Config

The `vm_config` argument supports the following:

- `name`: -(Optional) VM name.
- `description`: -(Optional) VM description
- `create_time`: -(Optional) VM creation time
- `update_time`: -(Optional) VM last updated time.
- `source`: -(Optional) Reference to an entity that the VM should be cloned or created from
- `num_sockets`: -(Optional) Number of vCPU sockets.
- `num_cores_per_socket`: -(Optional) Number of cores per socket.
- `num_threads_per_core`: -(Optional) Number of threads per core
- `num_numa_nodes`: -(Optional) Number of NUMA nodes. 0 means NUMA is disabled.
- `memory_size_bytes`: -(Optional) Memory size in bytes.
- `is_vcpu_hard_pinning_enabled`: -(Optional) Indicates whether the vCPUs should be hard pinned to specific pCPUs or not.
- `is_cpu_passthrough_enabled`: -(Optional) Indicates whether to passthrough the host CPU features to the guest or not. Enabling this will make VM incapable of live migration.
- `enabled_cpu_features`: -(Optional) The list of additional CPU features to be enabled. HardwareVirtualization: Indicates whether hardware assisted virtualization should be enabled for the Guest OS or not. Once enabled, the Guest OS can deploy a nested hypervisor
- `is_memory_overcommit_enabled`: -(Optional) Indicates whether the memory overcommit feature should be enabled for the VM or not. If enabled, parts of the VM memory may reside outside of the hypervisor physical memory. Once enabled, it should be expected that the VM may suffer performance degradation.
- `is_gpu_console_enabled`: -(Optional) Indicates whether the vGPU console is enabled or not.
- `is_cpu_hotplug_enabled`: -(Optional) Indicates whether the VM CPU hotplug is enabled.
- `is_scsi_controller_enabled`: -(Optional) Indicates whether the VM SCSI controller is enabled.
- `generation_uuid`: -(Optional) Generation UUID of the VM. It should be of type UUID.
- `bios_uuid`: -(Optional) BIOS UUID of the VM. It should be of type UUID.
- `categories`: -(Optional) Categories for the VM.
- `project`: -(Optional) Reference to a project.
- `ownership_info`: -(Optional) Ownership information for the VM.
- `host`: -(Optional) Reference to the host, the VM is running on.
- `cluster`: -(Optional) Reference to a cluster.
- `availability_zone`: -(Optional) Reference to an availability zone.
- `guest_customization`: -(Optional) Stage a Sysprep or cloud-init configuration file to be used by the guest for the next boot. Note that the Sysprep command must be used to generalize the Windows VMs before triggering this API call.
- `guest_tools`: -(Optional) The details about Nutanix Guest Tools for a VM.
- `hardware_clock_timezone`: -(Optional) VM hardware clock timezone in IANA TZDB format (America/Los_Angeles).
- `is_branding_enabled`: -(Optional) Indicates whether to remove AHV branding from VM firmware tables or not.
- `boot_config`: -(Optional) Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order.
- `is_vga_console_enabled`: -(Optional) Indicates whether the VGA console should be disabled or not.
- `machine_type`: -(Optional) Machine type for the VM. Machine type Q35 is required for secure boot and does not support IDE disks.
- `vtpm_config`: -(Optional) Indicates how the vTPM for the VM should be configured.
- `is_agent_vm`: -(Optional) Indicates whether the VM is an agent VM or not. When their host enters maintenance mode, once the normal VMs are evacuated, the agent VMs are powered off. When the host is restored, agent VMs are powered on before the normal VMs are restored. In other words, agent VMs cannot be HA-protected or live migrated.
- `apc_config`: -(Optional) Advanced Processor Compatibility configuration for the VM. Enabling this retains the CPU model for the VM across power cycles and migrations.
- `storage_config`: -(Optional) Storage configuration for VM.
- `disks`: -(Optional) Disks attached to the VM.
- `cd_roms`: -(Optional) CD-ROMs attached to the VM.
- `nics`: -(Optional) NICs attached to the VM.
- `gpus`: -(Optional) GPUs attached to the VM.
- `serial_ports`: -(Optional) Serial ports configured on the VM.
- `protection_type`: -(Optional) The type of protection applied on a VM. PD_PROTECTED indicates a VM is protected using the Prism Element. RULE_PROTECTED indicates a VM protection using the Prism Central.
- `protection_policy_state`: -(Optional) Status of protection policy applied to this VM.

#### Source

The `source` attribute supports the following:

- `entity_type`: -(Optional) Reference to an entity from which the VM should be cloned or created. Values are:
  - VM_RECOVERY_POINT: Reference to the recovery point entity from which the VM should be cloned or created.
  - VM: Reference to an entity from which the VM should be cloned or created.

#### Categories

The `categories` attribute supports the following:

- `ext_id`: -(Optional) A globally unique identifier of a VM category of type UUID.

#### Ownership Info

The `ownership_info` attribute supports the following:

- `owner`: Reference to the owner.
- `owner.ext_id`: -(Optional) A globally unique identifier of a VM owner type UUID.

#### Host

The `host` attribute supports the following:

- `ext_id`: -(Optional) A globally unique identifier of a host of type UUID.

#### Cluster

The `cluster` attribute supports the following:

- `ext_id`: -(Optional) The globally unique identifier of a cluster type UUID.

#### Availability Zone

The `availability_zone` attribute supports the following:

- `ext_id`: -(Optional) The globally unique identifier of an availability zone type UUID.

#### Guest Customization

The `guest_customization` attribute supports the following:

- `config`: -(Optional) The Nutanix Guest Tools customization settings.

- `config.sysprep`: -(Optional) Sysprep config
- `config.cloud_init`: -(Optional) CloudInit Config

##### config.sysprep

- `install_type`: -(Optional) Indicates whether the guest will be freshly installed using this unattend configuration, or this unattend configuration will be applied to a pre-prepared image. Default is 'PREPARED'.
- `sysprep_script`: -(Optional) Object either UnattendXml or CustomKeyValues
- `sysprep_script.unattend_xml`: -(Optional) xml object
- `sysprep_script.custom_key_values`: -(Optional) The list of the individual KeyValuePair elements.

##### config.cloud_init

- `datasource_type`: -(Optional) Type of datasource. Default: CONFIG_DRIVE_V2
- `metadata`: -(Optional) The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded.
- `cloud_init_script`: -(Optional) The script to use for cloud-init.
- `cloud_init_script.user_data`: -(Optional) user data object
- `cloud_init_script.custom_keys`: -(Optional) The list of the individual KeyValuePair elements.

#### Guest Tools

The `guest_tools` attribute supports the following:

- `is_enabled`: -(Optional) Indicates whether Nutanix Guest Tools is enabled or not.
- `capabilities`: -(Optional) The list of the application names that are enabled on the guest VM.

#### Boot Config

The `boot_config` attribute supports the following:

- `legacy_boot`: (Optional) LegacyBoot config Object
- `uefi_boot`: (Optional) UefiBoot config Object

##### boot_config.legacy_boot

- `boot_device`: (Required) Boot Device object
- `boot_device.boot_device_disk`: (Optional) Disk address.
- `boot_device.boot_device_disk.disk_address.bus_type`: (Required) Bus type for the device
- `boot_device.boot_device_disk.disk_address.index`: (Required) Device index on the bus. This field is ignored unless the bus details are specified.

- `boot_device.boot_device_nic`: (Optional) Disk Nic address.
- `boot_device.boot_device_nic.mac_address`: (Required) mac address

- `boot_order`: (Optional) Indicates the order of device types in which the VM should try to boot from. If the boot device order is not provided the system will decide an appropriate boot device order. Valid values are 'CDROM', 'DISK', 'NETWORK'.

##### boot_config.uefi_boot

- `is_secure_boot_enabled`: (Optional) Indicate whether to enable secure boot or not
- `nvram_device`: (Optional) Configuration for NVRAM to be presented to the VM.
- `nvram_device.backing_storage_info`: (Required) Storage provided by Nutanix ADSF

###### nvram_device.backing_storage_info

- `disk_size_bytes`: (Required) Size of the disk in Bytes
- `storage_container`: (Optional) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
- `storage_config`: (Optional) Storage configuration for VM disks
- `storage_config.is_flash_mode_enabled`: (Required) Indicates whether the virtual disk is pinned to the hot tier or not.
- `data_source`: (Optional) A reference to a disk or image that contains the contents of a disk.

#### VTPM Config

The `vtpm_config` attribute supports the following:

- `is_vtpm_enabled`: (Required) Indicates whether the virtual trusted platform module is enabled for the Guest OS or not.

#### APC Config

The `apc_config` attribute supports the following:

- `is_apc_enabled`: (Optional) If enabled, the selected CPU model will be retained across live and cold migrations of the VM.
- `cpu_model`: (Optional) CPU model associated with the VM if Advanced Processor Compatibility(APC) is enabled. If APC is enabled and no CPU model is explicitly set, a default baseline CPU model is picked by the system. See the APC documentation for more information
- `cpu_model.name`: (Required) Name of the CPU model associated with the VM.

#### Storage Config

The `storage_config` attribute supports the following:

- `is_flash_mode_enabled`: (Optional) Indicates whether the virtual disk is pinned to the hot tier or not.
- `qos_config`: (Optional) QoS parameters to be enforced.
- `qos_config.throttled_iops`: (Optional) Throttled IOPS for the governed entities. The block size for the I/O is 32 kB.

#### Disks

The `disks` attribute supports the following:

- `disk_address`: (Optional) Disk address.
- `disk_address.bus_type`: (Required) Bus type for the device. The acceptable values are: SCSI, IDE, PCI, SATA, SPAPR (only PPC).
- `disk_address.index`: (Required) Device index on the bus. This field is ignored unless the bus details are specified.
- `backing_info`: (Required) Supporting storage to create virtual disk on.
- `backing_info.vm_disk`:(Optional) backing Info for vmDisk
- `backing_info.adfs_volume_group_reference`: (Required) Volume Group Reference
- `backing_info.adfs_volume_group_reference.volume_group_ext_id`: (Required) The globally unique identifier of an ADSF volume group. It should be of type UUID.

##### backing_info.vm_disk

- `disk_size_bytes`: (Required) Size of the disk in Bytes
- `storage_container`: (Required) This reference is for disk level storage container preference. This preference specifies the storage container to which this disk belongs.
- `storage_config`: (Optional) Storage configuration for VM disks
- `storage_config.is_flash_mode_enabled`: Indicates whether the virtual disk is pinned to the hot tier or not.
- `data_source`: (Optional) A reference to a disk or image that contains the contents of a disk.
  container.

##### backing_info.vm_disk.data_source

- `reference`: (Required) Reference to image or vm disk
- `image_reference`: (Optional) Image Reference
- `image_reference.image_ext_id`: (Required) The globally unique identifier of an image. It should be of type UUID.
- `vm_disk_reference`: (Optional) Vm Disk Reference
- `vm_disk_reference.disk_ext_id`: (Optional) The globally unique identifier of a VM disk. It should be of type UUID.
- `vm_disk_reference.disk_address`: (Optional) Disk address.
- `vm_disk_reference.vm_reference`: (Optional) This is a reference to a VM.

#### CD-ROMs

The `cd_roms` attribute supports the following:

- `disk_address`: (Optional) Virtual Machine disk (VM disk).
- `backing_info`: (Optional) Storage provided by Nutanix ADSF
- `iso_type`: Type of ISO image inserted in CD-ROM. Valid values "OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION" .

#### NICs

The `nics` attribute supports the following:

- `nic_backing_info`: (Optional) New NIC backing info (v2.4.1+). One of `virtual_ethernet_nic`, `sriov_nic`, `dp_offload_nic`.
- `nic_network_info`: (Optional) New NIC network info (v2.4.1+). One of `virtual_ethernet_nic_network_info`, `sriov_nic_network_info`, `dp_offload_nic_network_info`.
- `backing_info`: (Optional, Deprecated) Use `nic_backing_info.virtual_ethernet_nic` instead.
- `network_info`: (Optional, Deprecated) Use `nic_network_info.virtual_ethernet_nic_network_info` instead.

#### nics.backing_info

- `model`: (Optional) Options for the NIC emulation. Valid values "VIRTIO" , "E1000".
- `mac_address`: (Optional) MAC address of the emulated NIC.
- `is_connected`: (Optional) Indicates whether the NIC is connected or not. Default is True.
- `num_queues`: (Optional) The number of Tx/Rx queue pairs for this NIC. Default is 1.

#### nics.network_info

- `nic_type`: (Optional) NIC type. Valid values "SPAN_DESTINATION_NIC", "NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC" .
- `network_function_chain`: (Optional) The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
- `network_function_nic_type`: (Optional) The type of this Network function NIC. Defaults to INGRESS.
- `subnet`: (Optional) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC
- `subnet.ext_id`: (Optional) The globally unique identifier of a subnet of type UUID.
- `vlan_mode`: (Optional) all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs.
- `trunked_vlans`: (Optional) List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
- `should_allow_unknown_macs`: (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
- `ipv4_config`: (Optional) The IP address configurations.

##### nics.ipv4_config

- `should_assign_ip`: If set to true (default value), an IP address must be assigned to the VM NIC - either the one explicitly specified by the user or allocated automatically by the IPAM service by not specifying the IP address. If false, then no IP assignment is required for this VM NIC.
- `ip_address`: The IP address of the NIC.
- `secondary_ip_address_list`: Secondary IP addresses for the NIC.

###### ip_address, secondary_ip_address_list

- `value`: The IPv4 address of the host.
- `prefix_length`: The prefix length of the IP address.

#### gpus

- `mode`: ((Optional)) The mode of this GPU. Valid values "PASSTHROUGH_GRAPHICS", "PASSTHROUGH_COMPUTE", "VIRTUAL" .
- `device_id`: (Optional) The device Id of the GPU.
- `vendor`: (Optional) The vendor of the GPU. Valid values "NVIDIA", "AMD", "INTEL" .
- `pci_address`: (Optional) The (S)egment:(B)us:(D)evice.(F)unction hardware address.

#### gpus.pci_address

- `segment`
- `bus`
- `device`
- `func`

#### serial_ports

- `is_connected`: -(Optional) Indicates whether the serial port is connected or not.
- `index`: -(Optional) Index of the serial port.

#### protection_policy_state

- `policy`: (Optional) Reference to the policy object in use.
- `policy.ext_id`: (Optional) Reference to the policy object in use.

## Import

This helps to manage existing entities which are not created through terraform. OVAs can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_ova_v2" "import_ova" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_ovas_v2" "fetch_ovas"{}
terraform import nutanix_ova_v2.import_ova <UUID>
```

See detailed information in [Nutanix Get Ova Details V4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/getDomainManagerById).
