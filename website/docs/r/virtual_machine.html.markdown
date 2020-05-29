---
layout: "nutanix"
page_title: "NUTANIX: nutanix_virtual_machine"
sidebar_current: "docs-nutanix-resource-virtual-machine"
description: |-
  Provides a Nutanix Virtual Machine resource to Create a virtual machine.
---

# nutanix_virtual_machine

Provides a Nutanix Virtual Machine resource to Create a virtual machine.

## Example Usage

```hcl
data "nutanix_clusters" "clusters" {}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou"
  cluster_uuid = data.nutanix_clusters.clusters.entities.0.metadata.uuid

  categories {
		name   = "Environment"
    value  = "Staging"
	}


  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) The name for the vm.
* `cluster_uuid`: - (Required) The UUID of the cluster.
* `categories`: - (Optional) Categories for the vm.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.
* `availability_zone_reference`: - (Optional) The reference to a availability_zone.
* `description`: - (Optional) A description for vm.
* `num_vnuma_nodes`: - (Optional) Number of vNUMA nodes. 0 means vNUMA is disabled.
* `nic_list`: - (Optional) Spec NICs attached to the VM.
* `serial_port_list`: - (Optional) Serial Ports configured on the VM.
* `guest_os_id`: - (Optional) Guest OS Identifier. For ESX, refer to VMware documentation [link](https://www.vmware.com/support/developer/converter-sdk/conv43_apireference/vim.vm.GuestOsDescriptor.GuestOsIdentifier.html) for the list of guest OS identifiers.
* `power_state`: - (Optional) The current or desired power state of the VM. (Options : ON , OFF)
* `nutanix_guest_tools`: - (Optional) Information regarding Nutanix Guest Tools.
* `ngt_credentials`: - (Ooptional) Credentials to login server.
* `ngt_enabled_capability_list` - (Optional) Application names that are enabled.
* `num_vcpus_per_socket`: - (Optional) Number of vCPUs per socket.
* `num_sockets`: - (Optional) Number of vCPU sockets.
* `gpu_list`: - (Optional) GPUs attached to the VM.
* `parent_referece`: - (Optional) Reference to an entity that the VM cloned from.
* `memory_size_mib`: - (Optional) Memory size in MiB.
* `boot_device_order_list`: - (Optional) Indicates the order of device types in which VM should try to boot from. If boot device order is not provided the system will decide appropriate boot device order.
* `boot_device_disk_address`: - (Optional) Address of disk to boot from.
* `boot_device_mac_address`: - (Optional) MAC address of nic to boot from.
* `hardware_clock_timezone`: - (Optional) VM's hardware clock timezone in IANA TZDB format (America/Los_Angeles).
* `guest_customization_cloud_init_user_data`: - (Optional) The contents of the user_data configuration for cloud-init. This can be formatted as YAML, JSON, or could be a shell script. The value must be base64 encoded.
* `guest_customization_cloud_init_meta_data` - (Optional) The contents of the meta_data configuration for cloud-init. This can be formatted as YAML or JSON. The value must be base64 encoded.
* `guest_customization_cloud_init_custom_key_values`: - (Optional) Generic key value pair used for custom attributes in cloud init.
* `guest_customization_is_overridable`: - (Optional) Flag to allow override of customization by deployer.
* `guest_customization_sysprep`: - (Optional) VM guests may be customized at boot time using one of several different methods. Currently, cloud-init w/ ConfigDriveV2 (for Linux VMs) and Sysprep (for Windows VMs) are supported. Only ONE OF sysprep or cloud_init should be provided. Note that guest customization can currently only be set during VM creation. Attempting to change it after creation will result in an error. Additional properties can be specified. For example - in the context of VM template creation if \"override_script\" is set to \"True\" then the deployer can upload their own custom script.
* `guest_customization_sysrep_custom_key_values`: - (Optional) Generic key value pair used for custom attributes in sysrep.
* `should_fail_on_script_failure`: - (Optional)  Extra configs related to power state transition. Indicates whether to abort ngt shutdown/reboot if script fails.
* `enable_script_exec`: - (Optional) Extra configs related to power state transition. Indicates whether to execute set script before ngt shutdown/reboot.
* `power_state_mechanism`: - (Optional) Indicates the mechanism guiding the VM power state transition. Currently used for the transition to \"OFF\" state. Power state mechanism (ACPI/GUEST/HARD).
* `vga_console_enabled`: - (Optional) Indicates whether VGA console should be enabled or not.
* `disk_list` Disks attached to the VM.
* `use_hot_add`: - (Optional) Use Hot Add when modifying VM resources. Passing value false will result in VM reboots. Default value is true.

### Disk List

The disk_list attribute supports the following:

* `uuid`: - (Optional) The device ID which is used to uniquely identify this particular disk.
* `disk_size_bytes` - (Optional) Size of the disk in Bytes.
* `disk_size_mib` - Size of the disk in MiB. Must match the size specified in 'disk_size_bytes' - rounded up to the nearest MiB - when that field is present.
* `device_properties` - Properties to a device.
* `data_source_reference` - Reference to a data source.
* `volume_group_reference` - Reference to a volume group.

The disk_size (the disk size_mib and the disk_size_bytes attributes) is only honored by creating an empty disk. When you are creating from an image, the size is ignored and the disk becomes the size of the image from which it was cloned. In VM creation, you can't set either disk size_mib or disk_size_bytes when you set data_source_reference but, you can update the disk_size after creation (second apply).

### Device Properties

The device_properties attribute supports the following.

* `device_type`: - A Disk type (default: DISK).
* `disk_address`: - Address of disk to boot from.

### Sysprep

The guest_customization_sysprep attribute supports the following:

* `install_type`: - (Optional) Whether the guest will be freshly installed using this unattend configuration, or whether this unattend configuration will be applied to a pre-prepared image. Default is \"PREPARED\".
* `unattend_xml`: - (Optional) Generic key value pair used for custom attributes.

### Disk Address

 The boot_device_disk_address attribute supports the following:

* `device_index`: - (Optional) The index of the disk address.
* `adapter_type`: - (Optional) The adapter type of the disk address.

### GPU List

The gpu_list attribute supports the following:

* `frame_buffer_size_mib`: - (ReadOnly) GPU frame buffer size in MiB.
* `vendor`: - (Optional) The vendor of the GPU.
* `uuid`: - (ReadOnly) UUID of the GPU.
* `name`: - (ReadOnly) Name of the GPU resource.
* `pci_address` - (ReadOnly) GPU {segment:bus:device:function} (sbdf) address if assigned.
* `fraction` - (ReadOnly) Fraction of the physical GPU assigned.
* `mode`: - (Optional) The mode of this GPU.
* `num_virtual_display_heads`: - (ReadOnly) Number of supported virtual display heads.
* `guest_driver_version`: - (ReadOnly) Last determined guest driver version.
* `device_id`: - (Computed) The device ID of the GPU.

### Nutanix Guest Tools

The nutanix_guest_tools attribute supports the following:

* `state`: - (Optional) Nutanix Guest Tools is enabled or not.
* `ngt_state`: - (Optional) Nutanix Guest Tools is enabled or not.
* `iso_mount_state`: - (Optional) Desired mount state of Nutanix Guest Tools ISO.
* `version`: - (ReadOnly) Version of Nutanix Guest Tools installed on the VM.
* `available_version`: - (ReadOnly) Version of Nutanix Guest Tools available on the cluster.
* `guest_os_version`: - (ReadOnly) Version of the operating system on the VM.
* `vss_snapshot_capable`: - (ReadOnly) Whether the VM is configured to take VSS snapshots through NGT.
* `is_reachable`: - (ReadOnly) Communication from VM to CVM is active or not.
* `vm_mobility_drivers_installed`: - (ReadOnly) Whether VM mobility drivers are installed in the VM.

### NIC List

The nic_list attribute supports the following:

* `nic_type`: - The type of this NIC. Defaults to NORMAL_NIC. (Options : NORMAL_NIC , DIRECT_NIC , NETWORK_FUNCTION_NIC).
* `uuid`: - The NIC's UUID, which is used to uniquely identify this particular NIC. This UUID may be used to refer to the NIC outside the context of the particular VM it is attached to.
* `model`: - The model of this NIC. (Options : VIRTIO , E1000).
* `network_function_nic_type`: - The type of this Network function NIC. Defaults to INGRESS. (Options : INGRESS , EGRESS , TAP).
* `mac_address`: - The MAC address for the adapter.
* `ip_endpoint_list`: - IP endpoints for the adapter. Currently, IPv4 addresses are supported.
* `network_function_chain_reference`: - The reference to a network_function_chain.
* `subnet_uuid`: - The reference to a subnet.
* `subnet_name`: - The name of the subnet reference to.
* `floating_ip`: -  The Floating IP associated with the vnic. (Only in `nic_list_status`)

### Serial Port List

The `serial_port_list` attribute supports the following:

* `index`: - Index of the serial port (int).
* `is_connected`: - Indicates whether the serial port connection is connected or not (`true` or `false`).

### ip_endpoint_list

The following attributes are exported:

* `ip`: - Address string.
* `type`: - Address type. It can only be "ASSIGNED" in the spec. If no type is specified in the spec, the default type is set to "ASSIGNED". (Options : ASSIGNED , LEARNED)

## Attributes Reference

The following attributes are exported:

* `metadata`: - The vm kind metadata.
* `api_version` - The version of the API.
* `state`: - The state of the vm.
* `cluster_name`: - The name of the cluster.
* `host_reference`: - Reference to a host.
* `hypervisor_type`: - The hypervisor type for the hypervisor the VM is hosted on.
* `nic_list_status`: - Status NICs attached to the VM.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid`: - vm UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when vm was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - vm name.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `network_function_chain_reference`, `data_source_reference`, `volume_group_reference` attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `uuid`: - the UUID(Required).

See detailed information in [Nutanix Virtual Machine](http://developer.nutanix.com/reference/prism_central/v3/#vms).
