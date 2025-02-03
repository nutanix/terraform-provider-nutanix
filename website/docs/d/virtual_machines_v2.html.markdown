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
data "nutanix_virtual_machines_v2" "vms"{}

data "nutanix_virtual_machines_v2" "vms-1"{
    page=0
    limit=2
}

data "nutanix_virtual_machines_v2" "vms-2"{
    filter = "name eq 'test-vm-filter'"
}
```

## Attribute Reference

The following attributes are exported:

- `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions
- `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default
- `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions.
- `vms`: List of all vms

### vms

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
- `memorysizebytes`: Memory size in bytes.
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

See detailed information in [Nutanix Virtual Machines V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0).
