---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_node_network_details"
sidebar_current: "docs-nutanix-datasource-foundation-node-network-details"
description: |-
 Gets hypervisor, CVM & IPMI info of the discovered nodes.
---

# nutanix_foundation_node_network_details

Gets hypervisor, CVM & IPMI info of the discovered nodes using their ipv6 address.

## Example Usage

```hcl
data "nutanix_foundation_node_network_details" "network_details" {
    ipv6_addresses = [
        "<ipv6-address-1>", "<ipv6-address-2>"
    ]
    timeout = "30"
}
```

## Argument Reference

The following arguments are supported:

* `ipv6_addresses`: (Required) list of ipv6 addresses
* `timeout`: timeout in seconds

## Attribute Reference

The following attributes are exported:

* `nodes`: nodes array.

### nodes

* `cvm_gateway`:  Gateway of CVM.
* `ipmi_netmask`: IPMI netmask.
* `ipv6_address`: IPV6 address of the CVM.
* `cvm_vlan_id`: CVM vlan tag.
* `hypervisor_hostname`: Hypervisor hostname.
* `hypervisor_netmask`: Netmask of the hypervisor.
* `cvm_netmask`: Netmask of CVM.
* `ipmi_ip`: IPMI IP address.
* `hypervisor_gateway`: Gateway of the hypervisor.
* `error`: Only exists when failed to fetch node_info, with the reason of failure. all other fields will be empty.
* `cvm_ip`: CVM IP address.
* `hypervisor_ip`: Hypervisor IP address.
* `ipmi_gateway`: IPMI gateway.
* `node_serial`: Node serial.

See detailed information in [Nutanix Foundation Node Network Details](https://www.nutanix.dev/api_references/foundation/#/b3A6MjIyMjMzOTU-get-the-network-configuration-details-of-the-node).
