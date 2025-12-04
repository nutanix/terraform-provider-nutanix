---
layout: "nutanix"
page_title: "NUTANIX: nutanix_clusters_unconfigured_node_networks_v2"
sidebar_current: "docs-nutanix-datasource-clusters-unconfigured-node-networks-v2"
description: |-
  Get a dictionary of cluster networks and available uplinks on the given nodes. This API is not supported for XEN hypervisor type.
---

# nutanix_clusters_unconfigured_node_networks_v2

Get a dictionary of cluster networks and available uplinks on the given nodes. This API is not supported for XEN hypervisor type.

## Example Usage

```hcl
# ## fetch Network info for unconfigured node
resource "nutanix_clusters_unconfigured_node_networks_v2" "node-network-info" {
  ext_id       = "0005b6b0-0b0b-0000-0000-000000000000"
  request_type = "expand_cluster"
  node_list {
    cvm_ip {
      ipv4 {
        value = "10.73.23.55"
      }
    }
    hypervisor_ip {
      ipv4 {
        value = "10.33.44.12"
      }
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) Cluster UUID.
* `node_list`: -(Required) List of nodes for which the network information is required.
* `request_type`: -(Optional) Request type

### Node List
The node_list attribute supports the following:

* `node_uuid`: -(Optional) Node UUID.
* `block_id`: -(Optional) Block ID.
* `node_position`: -(Optional) Node position.
* `hypervisor_type`: -(Optional) Hypervisor type.
* `is_robo_mixed_hypervisor`: -(Optional) Is ROBO mixed hypervisor.
* `hypervisor_version`: -(Optional) Hypervisor version.
* `nos_version`: -(Optional) NOS version.
* `is_compute_only`: -(Optional) Is compute only.
* `ipmi_ip`: -(Optional) IPMI IP.
* `digital_certificate_map_list`: -(Optional) Digital certificate map list.
* `cvm_ip`: -(Optional) CVM IP.
* `hypervisor_ip`: -(Optional) Hypervisor IP.
* `model`: -(Optional) Model name.
* `current_network_interface`: -(Optional) Current network interface.

#### Digital Certificate Map List
The `digital_certificate_map_list` attribute supports the following:

* `key`: -(Optional) Field containing digital_certificate_base64 and key_management_server_uuid for key management server.
* `name`: -(Optional) Value for the fields digital_certificate_base64 and key_management_server_uuid for key management server.

#### Ip Address Attributes
The `ipmi_ip`, `cvm_ip`, `hypervisor_ip` attributes supports the following:

* `ipv4`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv6 format.


#### IPV4, IPV6
The `ipv4`, `ipv6` attributes supports the following:

* `value`: -(Required) The IPv4/IPv6 address of the host.
* `prefix_length`: -(Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.

## Attributes Reference
The following attributes are exported under `nodes_networking_details`:

* `network_info`: - Network information for the given nodes.
* `uplinks`: - List of uplinks information for each CVM IP.
* `warnings`: - List of warning messages.


### Network Info
The `network_info` attribute supports the following:

* `hci`: - Network information of HCI nodes.
* `so`: - Network information of SO nodes.

#### HCI and SO
The `hci` and `so` attribute supports the following:

* `hypervisor_type`: - Hypervisor type.
* `name`: - Interface name.
* `networks`: - List of networks for interface.

### Uplinks
The `uplinks` attribute supports the following:

* `cvm_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `uplink_list`: - Uplink details for a controller VM.

#### Uplink List
The `uplink_list` attribute supports the following:

* `name`: - Interface name.
* `mac`: - MAC address.

See detailed information in [Nutanix Network Information of Unconfigured Nodes V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/Clusters/operation/fetchNodeNetworkingDetails).
