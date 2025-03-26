---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_network_device_assign_ip_v2"
sidebar_current: "docs-nutanix-resource-vm-network-device-assign-ip-v2"
description: |-
  Provides a Nutanix Virtual Machine resource to Assign IP.
---

# nutanix_vm_network_device_assign_ip_v2

Provides a Nutanix Virtual Machine resource to Assign IP.

## Example Usage

```hcl

resource "nutanix_vm_network_device_assign_ip_v2" "nic_assign_ip"{
    vm_ext_id = "246f6e8a-ff05-4057-af6b-b1fd23a46d7d" # VM UUID
    ext_id    = "eb0157e7-4a87-4ba6-ac8f-62cfe6251b8b" # Vm NIC Device UUID
    ip_address {
        value = "10.10.10.10" # IP Address
        prefix_length = 32
    }
}

```

## Argument Reference

The following arguments are supported:

* `vm_ext_id`: - (Required) The globally unique identifier of a VM. It should be of type UUID.
* `ext_id`: - (Required) The globally unique identifier of a Nic. It should be of type UUID.
* `ip_address`: - (Optional) Ip config settings.

### IP Address

The ip_address attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - Ip address.

See detailed information in [Nutanix Assign an IP address to the VM V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/assignIpById).
