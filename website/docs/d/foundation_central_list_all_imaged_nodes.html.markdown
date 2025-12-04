---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_imaged_nodes_list"
sidebar_current: "docs-nutanix-datasource-foundation-central-imaged-nodes"
description: |-
 List all the nodes registered with Foundation Central
---

# nutanix_foundation_central_imaged_nodes_list

List all the nodes registered with Foundation Central

## Example Usage

```hcl
data "nutanix_foundation_central_imaged_nodes_list" "nodes_list" {}
```

## Argument Reference

No arguments are supported

## Attribute Reference

The following attributes are exported:

* `metadata`: List metadata output for all list apis.

### metadata
* `total_matches`: Total matches found.
* `length`: The number of records retrieved.
* `offset`: Offset from the start of the object list.

### imaged_nodes
* `cvm_vlan_id`: Vlan tag of the cvm, if the cvm is on a vlan.
* `node_type`: Specifies the type of node - on-prem, AWS, GCP etc.
* `created_timestamp`: Time when the node was discovered in Foundation Central.
* `ipv6_interface`: Name of the cvm interface having ipv6 address.
* `api_key_uuid`: API key used to register the node.
* `foundation_version`: Foundation version installed on the node.
* `current_time`: Current time of Foundation Central.
* `node_position`: Position of the node in the block.
* `cvm_netmask`: netmask of the cvm.
* `ipmi_ip`: IP address of the ipmi.
* `cvm_uuid`: Node UUID from the node's cvm.
* `cvm_ipv6`: IPv6 address of the cvm.
* `imaged_cluster_uuid`: UUID of the cluster to which the node belongs, if any.
* `cvm_up`: Denotes whether the CVM is up or not on this node.
* `available`: Specifies whether the node is available for cluster creation.
* `object_version`: Version of the node used for CAS.
* `ipmi_netmask`: netmask of the ipmi.
* `hypervisor_hostname`: Name of the hypervisor host.
* `node_state`: Specifies whether the node is discovering, available or unavailable for cluster creation.
* `hypervisor_version`: Version of the hypervisor currently installed on the node.
* `hypervisor_ip`: IP address of the hypervisor.
* `model`: Model of the node.
* `ipmi_gateway`: gateway of the ipmi.
* `hardware_attributes`: Hardware attributes json of the node.
* `cvm_gateway`: gateway of the cvm.
* `node_serial`: Serial number of the node.
* `imaged_node_uuid`: UUID of the node.
* `block_serial`: Serial number of the block to which the node belongs.
* `hypervisor_type`: Hypervisor type currently installed on the node. Must be one of {kvm, esx, hyperv}.
* `latest_hb_ts_list`: List of timestamps when the node has sent heartbeats to Foundation Central.
* `hypervisor_netmask`: netmask of the hypervisor.
* `hypervisor_gateway`: gateway of the hypervisor.
* `cvm_ip`: IP address of the cvm.
* `aos_version`: AOS version currently installed on the node.

See detailed information in [Nutanix Foundation Central List all the nodes](https://www.nutanix.dev/api_references/foundation-central/#/26192129ae504-list-all-the-nodes).