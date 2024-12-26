---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_network_device_v2"
sidebar_current: "docs-nutanix-resource-vm-network-device-v2"
description: |-
  Provides a Nutanix Virtual Machine resource to Create a virtual machine nic.
---

# nutanix_vm_network_device_v2

Provides a Nutanix Virtual Machine resource to Create a virtual machine nic.

## Example Usage

```hcl
data "nutanix_virtual_machines_v2" "vms" {}

data "nutanix_subnets_v2" "subnets" { }

resource "nutanix_vm_network_device_v2" "test" {
    vm_ext_id = data.nutanix_virtual_machines_v2.vms.0.data.ext_id
    network_info {
        nic_type = "DIRECT_NIC"
        subnet {
        ext_id = data.nutanix_subnets_v2.subnets.subnets.0.ext_id
        }
    }
}

```

## Argument Reference

The following arguments are supported:

* `vm_ext_id`: - (Required) The globally unique identifier of a VM. It should be of type UUID.
* `ext_id`: - (Required) The globally unique identifier of a Nic. It should be of type UUID.
* `backing_info`: - (Optional) Defines a NIC emulated by the hypervisor
* `network_info`: - (Optional) Network information for a NIC.

### Backing Info

The backing_info attribute supports the following:

* `model`: - (Optional) Options for the NIC emulation.
  Valid values are:
    - `VIRTIO` The NIC emulation model is Virtio.
    - `E1000` The NIC emulation model is E1000.
* `mac_address`: - (Optional) MAC address of the emulated NIC.
* `is_connected`: - (Optional) Indicates whether the NIC is connected or not. Default is True.
* `num_queues`: - (Optional) The number of Tx/Rx queue pairs for this NIC.

### Network Info

The network_info attribute supports the following:

* `nic_type`: - (Optional) NIC type.
  Defaults to NORMAL_NIC.
  Valid values are:
    - `SPAN_DESTINATION_NIC` The type of NIC is Span-Destination.
    - `NORMAL_NIC` The type of NIC is Normal.
    - `DIRECT_NIC` The type of NIC is Direct.
    - `NETWORK_FUNCTION_NIC` The type of NIC is Network-Function.
* `network_function_chain`: - (Optional)The network function chain associates with the NIC. Only valid if nic_type is NORMAL_NIC.
* `network_function_nic_type`: - (Optional) The type of this Network function NIC.
  Defaults to INGRESS.
  Valid values are:
    - `TAP` The type of Network-Function NIC is Tap.
    - `EGRESS` The type of Network-Function NIC is Egress.
    - `INGRESS` The type of Network-Function NIC is Ingress.
* `subnet`: - (Optional) Network identifier for this adapter. Only valid if nic_type is NORMAL_NIC or DIRECT_NIC.
* `vlan_mode`: - (Optional) By default, all the virtual NICs are created in ACCESS mode, which permits only one VLAN per virtual network. TRUNKED mode allows multiple VLANs on a single VM NIC for network-aware user VMs.
  Valid values are:
    - `TRUNK` The virtual NIC is created in TRUNKED mode.
    - `ACCESS` The virtual NIC is created in ACCESS mode.
* `trunked_vlans`: - (Optional) List of networks to trunk if VLAN mode is marked as TRUNKED. If empty and VLAN mode is set to TRUNKED, all the VLANs are trunked.
* `should_allow_unknown_macs`: - (Optional) Indicates whether an unknown unicast traffic is forwarded to this NIC or not. This is applicable only for the NICs on the overlay subnets.
* `ipv4_config`: - (Optional) The IP address configurations.

### Network Function Chain

The network_function_chain attribute supports the following:

* `ext_id`: - (Optional) The globally unique identifier of a network function chain. It should be of type UUID.

### Subnet

The subnet attribute supports the following:

* `ext_id`: - (Optional) The globally unique identifier of a subnet. It should be of type UUID.

### IPV4 Config

The ipv4_config attribute supports the following:

* `should_assign_ip`: - (Optional) If set to true (default value), an IP address must be assigned to the VM NIC - either the one explicitly specified by the user or allocated automatically by the IPAM service by not specifying the IP address. If false, then no IP assignment is required for this VM NIC.
  `ip_address`: - (Optional) Ip config settings.
  `secondary_ip_address_list`: - (Optional) Secondary IP addresses for the NIC.

### IP Address

The ip_address attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - Ip address.

### Secondary IP Address List

The secondary_ip_address_list attribute supports the following:

* `prefix_length`: - (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `value`: - Ip address.

See detailed information in [Nutanix Virtual Machine](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0).

## Import
Nutanix Virtual machines can be imported using the `UUID` eg,

`
terraform import nutanix_vm_network_device_v2.nic01 0F75E6A7-55FB-44D9-A50D-14AD72E2CF7C
`