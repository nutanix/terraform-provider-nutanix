---
layout: "nutanix"
page_title: "NUTANIX: nutanix_clusters_discover_unconfigured_nodes_v2"
sidebar_current: "docs-nutanix-datasource-clusters-discover-unconfigured-nodes-v2"
description: |-
  Get the unconfigured node details such as node UUID, node position, node IP, foundation version and more.
---

# nutanix_clusters_discover_unconfigured_nodes_v2

Get the unconfigured node details such as node UUID, node position, node IP, foundation version and more.

## Example Usage

```hcl
data "nutanix_clusters_discover_unconfigured_nodes_v2" "example"{
  ext_id = "00057b0b-0b0b-0b0b-0b0b-000000000000"
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
        value = "10.0.0.1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) Cluster UUID.
* `address_type`: - (Optional) Address type.
  Valid values are:
    - "IPV4"	IPV4 address type.
    - "IPV6"	IPV6 address type.
* `ip_filter_list`: - (Optional) IP addresses of the unconfigured nodes.
* `uuid_filter_list`: - (Optional) Unconfigured node UUIDs.
* `timeout`: - (Optional) Timeout for the workflow in seconds.
* `interface_filter_list`: - (Optional) Interface name that is used for packet broadcasting.
* `is_manual_discovery`: - (Optional) Indicates if the discovery is manual or not.

### IP Filter List
The ip_filter_list attribute supports the following:

* `ipv4`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv6 format.


#### IPV4, IPV6
The `ipv4`, `ipv6` attributes supports the following:

* `value`: -(Required) The IPv4/IPv6 address of the host.
* `prefix_length`: -(Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.

## Attributes Reference
the following attributes are exported under `unconfigured_nodes`:

* `arch`: Cluster arch type.
* `attributes`: Attributes of a node.
* `cluster_id`: Cluster ID.
* `cpu_type`:  CPU type.
* `current_cvm_vlan_tag`:  Current CVM VLAN tag.
* `current_network_interface`:  Current network interface of a node.
* `cvm_ip`:  CVM IP.
* `foundation_version`:  Foundation version.
* `host_name`:  Host name.
* `host_type`:  Host type.
* `hypervisor_ip`: Hypervisor IP Address.
* `hypervisor_type`: Hypervisor type.
* `hypervisor_version`: Host version of the node.
* `interface_ipv6`: Interface IPV6 address.
* `ipmi_ip`: IPMI IP Address.
* `is_one_node_cluster_supported`: Indicates whether a node can be used to create a single node cluster or not.
* `is_secure_booted`: Secure boot status.
* `is_two_node_cluster_supported`: Indicates whether a node can be used to create a two node cluster or not.
* `node_position`: Position of a node in a rackable unit.
* `node_serial_number`: Node serial number.
* `node_uuid`: UUID of the host.
* `nos_version`: NOS software version of a node.
* `rackable_unit_max_nodes`: Maximum number of nodes in rackable-unit.
* `rackable_unit_model`: Rackable unit model type.
* `rackable_unit_serial`: Rackable unit serial name.

### Attributes
The `attributes` attribute supports the following:

* `default_workload`: Default workload.
* `is_model_supported`: Indicates whether the model is supported or not.
* `is_robo_mixed_hypervisor`: Indicates whether the hypervisor is robo mixed or not.
* `lcm_family`: LCM family name.
* `should_work_with_1g_nic`: Indicates if cvm interface can work with 1 GIG NIC or not.

### CVM IP, Hypervisor IP, IPMI IP
The `cvm_ip`, `hypervisor_ip`, `ipmi_ip` attributes supports the following:

* `ipv4`: An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: An unique address that identifies a device on the internet or a local network in IPv6 format.

#### IPV4, IPV6
The `ipv4`, `ipv6` attributes supports the following:

* `value`: The IPv4/IPv6 address of the host.
* `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.



See detailed information in [Nutanix Discover unconfigured nodes V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/Clusters/operation/discoverUnconfiguredNodes).
