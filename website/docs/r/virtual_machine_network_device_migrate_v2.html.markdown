---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_network_device_migrate_v2"
sidebar_current: "docs-nutanix-resource-vm-network-device-migrate-v2"
description: |-
  Provides a Nutanix Virtual Machine resource to Migrate a nic.
---

# nutanix_vm_network_device_migrate_v2

Provides a Nutanix Virtual Machine resource to Migrate a nic.

## Example Usage

```hcl
data "nutanix_virtual_machines_v2" "vms"{}

data "nutanix_subnets_v2" "subnets"{ }

resource "nutanix_vm_network_device_v2" "nic"{
    vm_ext_id = data.nutanix_virtual_machines_v2.vms.0.data.ext_id
    network_info {
        nic_type = "DIRECT_NIC"
        subnet {
            ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
        }
        ipv4_config {
            should_assign_ip = true
            ip_address {
                value = "10.10.10.10"
                prefix_length = 32
            }
        }
    }
}

resource "nutanix_vm_network_device_migrate_v2" "migrate"{
    vm_ext_id = resource.nutanix_virtual_machine_v4.vms.0.data.ext_id
    ext_id    = resource.nutanix_vm_network_device_v2.nic.ext_id
    subnet {
        ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
    }
    migrate_type = "ASSIGN_IP"
    ip_address {
        value = "10.10.10.11"
        prefix_length = 32
    }
}

```

## Argument Reference

The following arguments are supported:

* `vm_ext_id`: - (Required) The globally unique identifier of a VM. It should be of type UUID.
* `ext_id`: - (Required) The globally unique identifier of a Nic. It should be of type UUID.
* `subnet`: - (Required) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC.
* `migrate_type`: - (Required) The type of IP address management for NIC migration.
  Valid values are:
    - `ASSIGN_IP` The type of NIC is Span-Destination.
    - `RELEASE_IP` The type of NIC is Normal.
* `ip_address`: - (Optional) Ip config settings.

### Subnet

The subnet attribute supports the following:

* `ext_id`: - (Optional) The globally unique identifier of a subnet. It should be of type UUID.

### IP Address

The ip_address attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - Ip address.

See detailed information in [Nutanix Migrate NIC to another Subnet for VM V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0#tag/Vm/operation/migrateNicById).
