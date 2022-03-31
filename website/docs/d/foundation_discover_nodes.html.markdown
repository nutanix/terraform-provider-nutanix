---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_discover_nodes"
sidebar_current: "docs-nutanix-datasource-foundation-discover-nodes"
description: |-
 Discovers and lists Nutanix-imaged nodes within an IPv6 network.
---

# nutanix_foundation_discover_nodes

Discovers and lists Nutanix-imaged nodes within an IPv6 network.

## Example Usage

```hcl
data "nutanix_foundation_discover_nodes" "discovered_nodes" {}
```

## Argument Reference

No arguments are supported

## Attribute Reference

The following attributes are exported:

* `entities`: List of Nutanix-imaged nodes within an IPv6 network

### entities
* `model`: Model name of the block.
* `nodes`: Node level properties.
* `block_id`: Chassis serial number.
* `chassis_n`: ID number of the block.

### nodes
* `foundation_version`: Version of foundation.
* `ipv6_address`: IPV6 address of the node.
* `node_uuid`: UUID of the node.
* `current_network_interface`: Current network interface of the node.
* `node_position`: Position of the node in the block.
* `hypervisor`: Type of hypervisor installed on the node.
* `configured`: Whether the node is configured.
* `nos_version`: Version of NOS installed on the node.
* `cluster_id`: ID of the cluster the node is part of.
* `current_cvm_vlan_tag`: vlan tag of cvm.
* `hypervisor_version`: Version of hypervisor installed.
* `svm_ip`: IP address of CVM.
* `model`: Model name of the node.
* `node_serial`: Node serial of the node.

See detailed information in [Nutanix Foundation Discover Nodes](https://www.nutanix.dev/api_references/foundation/#/b3A6MjIyMjM0MDM-i-pv6-discovery).
