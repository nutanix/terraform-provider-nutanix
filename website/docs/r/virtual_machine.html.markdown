---
layout: "nutanix"
page_title: "NUTANIX: nutanix_virtual_machine"
sidebar_current: "docs-outscale-resource-virtual-machine"
description: |-
  Provides a Nutanix Virtual Machine resource to Create a virtual machine.
---

# nutanix_virtual_machine

Provides a Nutanix Virtual Machine resource to Create a virtual machine.

## Example Usage

```hcl
resource "nutanix_category_key" "test-category-key"{
    name        = "app-suppport-1"
    description = "App Support Category Key"
}


resource "nutanix_category_value" "test"{
    name        = "${nutanix_category_key.test-category-key.id}"
    description = "Test Category Value"
    value       = "test-value"
}

data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou"

  categories = [{
	  name   = "${nutanix_category_key.test-category-key.id}"
	  value = "${nutanix_category_value.test.id}"
  }]

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  power_state          = "ON"
}
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) The name for the vm.
* `cluster_reference`: - (Required) The reference to a cluster.
* `cluster_name`: - (Optional) The name of the reference to the cluster.
* `categories`: - (Optional) Categories for the vm.
* `project_reference`: - (Optional) The reference to a project.
* `owner_reference`: - (Optional) The reference to a user.
* `availability_zone_reference`: - (Optional) The reference to a availability_zone.
* `api_version`
* `description`: - (Optional) A description for vm.
* `num_vnuma_nodes`: - (Optional) Number of vNUMA nodes. 0 means vNUMA is disabled.
* `nic_list`: - (Optional) NICs attached to the VM.
* `guest_os_id`: - (Optional) Guest OS Identifier. For ESX, refer to VMware documentation [link](https://www.vmware.com/support/developer/converter-sdk/conv43_apireference/vim.vm.GuestOsDescriptor.GuestOsIdentifier.html) for the list of guest OS identifiers.
* `power_state`: - (Optional) The current or desired power state of the VM. (Options : ON , OFF)
* `nutanix_guest_tools`: - (Optional) Information regarding Nutanix Guest Tools.
* `num_vcpus_per_socket`: - (Optional) Number of vCPUs per socket.
* `num_sockets`: - (Optional) Number of vCPU sockets.
* `gpu_list`: - (Optional) GPUs attached to the VM.

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

* `available_version`: - (ReadOnly) Version of Nutanix Guest Tools available on the cluster.
* `iso_mount_state`: - (Optioinal) Desired mount state of Nutanix Guest Tools ISO.
* `state`: - (Optional) Nutanix Guest Tools is enabled or not.
* `version`: - (ReadOnly) Version of Nutanix Guest Tools installed on the VM.
* `guest_os_version`: - (ReadOnly) Version of the operating system on the VM.
* `enabled_capability_list`: - (Optional) Application names that are enabled.
* `vss_snapshot_capable`: - (ReadOnly) Whether the VM is configured to take VSS snapshots through NGT.
* `is_reachable`: - (ReadOnly) Communication from VM to CVM is active or not.
* `vm_mobility_drivers_installed`: - (ReadOnly) Whether VM mobility drivers are installed in the VM.

### NIC List

The nic_list attribute supports the following:

* `nic_type`: - (Optional) The type of this NIC. Defaults to NORMAL_NIC. (Options : NORMAL_NIC , DIRECT_NIC , NETWORK_FUNCTION_NIC)
  .
* `uuid`: - (Optional) The NIC's UUID, which is used to uniquely identify this particular NIC. This UUID may be used to refer to the NIC outside the context of the particular VM it is attached to.

* `floating_ip`: - (ReadOnly)

* `model`: - (Optional) The model of this NIC. (Options : VIRTIO , E1000).
* `network_function_nic_type`: - (Optional) The type of this Network function NIC. Defaults to INGRESS. (Options : INGRESS , EGRESS , TAP).
* `mac_address`: - (Optional) The MAC address for the adapter.
* `ip_endpoint_list`: - (Optional) IP endpoints for the adapter. Currently, IPv4 addresses are supported.
* `network_function_chain_reference`: - (Optional) The reference to a network_function_chain.
* `subnet_reference`: - (Optional) The reference to a subnet.
* `subnet_reference_name`: - (Optional) The name of the subnet reference to.

### ip_endpoint_list

The following attributes are exported:

* `ip`: - Address string.
* `type`: - Address type. It can only be "ASSIGNED" in the spec. If no type is specified in the spec, the default type is set to "ASSIGNED". (Options : ASSIGNED , LEARNED)

## Attributes Reference

The following attributes are exported:

* `metadata`: - The vm kind metadata.
* `state`: -
* `ip_address`: -
* `host_reference`: -
* `hypervisor_type`: -

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when vm was last updated.
* `uuid`: - vm uuid.
* `creation_time`: - UTC date and time in RFC-3339 format when vm was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - vm name.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, `network_function_chain_reference`, `subnet_reference`.

attributes supports the following:

* `kind`: - The kind name (Default value: project)(Required).
* `name`: - the name(Optional).
* `uuid`: - the uuid(Required).

Note: `cluster_reference`, `subnet_reference` does not support the attribute `name`

See detailed information in [Nutanix Virtual Machine](http://developer.nutanix.com/reference/prism_central/v3/#vms).
