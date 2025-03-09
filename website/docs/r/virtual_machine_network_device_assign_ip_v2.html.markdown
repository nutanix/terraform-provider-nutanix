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
data "nutanix_virtual_machines_v2" "vms"{}

data "nutanix_subnets_v2" "subnets"{}

resource "nutanix_vm_network_device_v2" "nic"{
    vm_ext_id = data.nutanix_virtual_machines_v2.vms.0.data.ext_id
    network_info {
        nic_type = "DIRECT_NIC"
        subnet {
        ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
        }
    }
}

resource "nutanix_vm_network_device_assign_ip_v2" "nic_assign_ip"{
    vm_ext_id = resource.nutanix_virtual_machine_v4.vms.0.ext_id
    ext_id    = resource.nutanix_vm_network_device_v2.nic.ext_id
    ip_address {
        value = "10.10.10.10"
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
