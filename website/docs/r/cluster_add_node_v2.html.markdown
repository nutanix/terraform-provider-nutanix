---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_add_node_v2"
sidebar_current: "docs-nutanix-resource-nutanix-cluster-add-node-v2"
description: |-
  Add node on a cluster identified by {extId}.
---

# nutanix_cluster_add_node_v2

Add node on a cluster identified by {extId}.

-> **Note:** Starting with v2.3.2, users can now perform node add/remove operations directly through the `nutanix_cluster_v2` resource, which offers a more consistent and automated approach to managing cluster scaling operations.



## Example Usage

```hcl
locals {
  # cluster of 3 node uuid that we want to add node
  clusters_ext_id = "00057b8b-0b3b-4b3b-0000-000000000000" # for example
  cvm_ip = "10.xx.xx.xx" # Node Ip address that we want to add
}

## check if the node to add is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-node" {
  ext_id       = local.clusters_ext_id
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = local.cvm_ip
    }
  }

  ## check if the 3 nodes are un configured or not
  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 1
      error_message = "The node ${local.cvm_ip} is configured"
    }
  }
}

## fetch Network info for unconfigured node
resource "nutanix_clusters_unconfigured_node_networks_v2" "node-network-info" {
  ext_id       = local.clusters_ext_id
  request_type = "expand_cluster"
  node_list {
    cvm_ip {
      ipv4 {
        value = local.cvm_ip
      }
    }
    hypervisor_ip {
      ipv4 {
        value = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_ip.0.ipv4.0.value
      }
    }
  }
  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node]
}

## add node to the cluster
resource "nutanix_cluster_add_node_v2" "add-node" {
  cluster_ext_id = local.clusters_ext_id

  should_skip_add_node          = false
  should_skip_pre_expand_checks = false

  node_params {
    should_skip_host_networking = false
    hypervisor_isos {
      type = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_type
    }
    node_list {
      node_uuid                 = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].node_uuid
      model                     = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].rackable_unit_model
      block_id                  = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].rackable_unit_serial
      hypervisor_type           = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_type
      hypervisor_version        = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_version
      node_position             = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].node_position
      nos_version               = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].nos_version
      hypervisor_hostname       = "example"
      current_network_interface = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].name
      # required for adding a node
      hypervisor_ip {
        ipv4 {
          value = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_ip.0.ipv4.0.value
        }
      }
      cvm_ip {
        ipv4 {
          value = local.cvm_ip
        }
      }
      ipmi_ip {
        ipv4 {
          value = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].ipmi_ip.0.ipv4.0.value
        }
      }

      is_robo_mixed_hypervisor = true
      networks {
        name     = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].network_info[0].hci[0].name
        networks = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].network_info[0].hci[0].networks
        uplinks {
          active {
            name  = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].name
            mac   = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].mac
            value = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[0].name
          }
          standby {
            name  = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[1].name
            mac   = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[1].mac
            value = nutanix_clusters_unconfigured_node_networks_v2.node-network-info.nodes_networking_details[0].uplinks[0].uplink_list[1].name
          }
        }
      }
    }

  }

  config_params {
    should_skip_imaging = true
    target_hypervisor   = nutanix_clusters_discover_unconfigured_nodes_v2.cluster-node.unconfigured_nodes[0].hypervisor_type
  }

  remove_node_params {
    extra_params {
      should_skip_upgrade_check = false
      skip_space_check          = false
      should_skip_add_check     = false
    }
    should_skip_remove    = false
    should_skip_prechecks = false
  }

  depends_on = [nutanix_clusters_unconfigured_node_networks_v2.node-network-info]
}

```


## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: -(Required) Cluster UUID.
* `node_params`: -(Required) Parameters of the node to be added.
* `config_params`: -(Optional) Config parameters.
* `should_skip_add_node`: -(Optional) Indicates if node addition can be skipped.
* `should_skip_pre_expand_checks`: -(Optional) Indicates if pre-expand checks can be skipped for node addition.
* `remove_node_params`: -(Optional) configuration for node removal.

### Node Params
The node_params block supports the following:

* `block_list`: -(Optional) Block list of a cluster.
* `node_list`: -(Required) List of nodes in a cluster.
* `computed_node_list`: -(Optional) List of compute only nodes.
* `hypervisor_isos`: -(Optional) Hypervisor type to md5sum map.
* `hyperv_sku`: -(Optional) Hyperv SKU.
* `bundle_info`: -(Optional) Hypervisor bundle information.
* `should_skip_host_networking`: -(Optional) Indicates if the host networking needs to be skipped or not.

#### Block List
The block_list block supports the following:

* `block_id`: -(Required) List of nodes in a block.
* `rack_name`: -(Optional) Indicates if the host networking needs to be skipped or not.

#### Node List
The node_list block supports the following:

* `node_uuid`: -(Optional) Node UUID.
* `block_id`: -(Optional) Block ID.
* `node_position`: -(Optional) Node position.
* `hypervisor_type`: -(Optional) Hypervisor type.
   Valid values are:
    - `XEN`: Xen hypervisor.
    - `HYPERV`: Hyper-V hypervisor.
    - `NATIVEHOST`: NativeHost type where AOS runs natively, without hypervisor.
    - `ESX`: ESX hypervisor.
    - `AHV`: AHV hypervisor.
* `is_robo_mixed_hypervisor`: -(Optional) Is ROBO mixed hypervisor.
* `hypervisor_hostname`: -(Optional) Name of the host.
* `hypervisor_version`: -(Optional) Host version of the node.
* `nos_version`: -(Optional) NOS software version of a node.
* `is_compute_only`: -(Optional) Indicates whether the node is light compute or not.
* `ipmi_ip`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `digital_certificate_map_list`: -(Optional) List of objects containing digital_certificate_base64 and key_management_server_uuid fields for key management server.
* `cvm_ip`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `hypervisor_ip`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `model`: -(Optional) Rackable unit model name.
* `current_network_interface`: -(Optional) Current network interface of a node.
* `networks`: -(Optional) Active and standby uplink information of the target nodes.

##### Networks
The networks block supports the following:

* `name`: -(Optional) Name of the uplink.
* `networks`: -(Optional) List of network types.
* `uplinks`: -(Optional) Active and standby uplink information of the target nodes.

###### Uplinks
The uplinks block supports the following:

* `active`: -(Optional) Active uplink information.
* `standby`: -(Optional) Standby uplink information.

###### Active, Standby
The `active`, `standby` attributes supports the following:

* `mac`: -(Optional) Mac address.
* `name`: -(Optional) Interface name.
* `value`: -(Optional) Interface value.

#### Computed Node List
The computed_node_list block supports the following:

* `node_uuid`: -(Optional) Node UUID.
* `block_id`: -(Optional) Block ID.
* `node_position`: -(Optional) Node position.
* `hypervisor_ip`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `ipmi_ip`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `digital_certificate_map_list`: -(Optional) List of objects containing digital_certificate_base64 and key_management_server_uuid fields for key management server.
* `hypervisor_hostname`: -(Optional) Name of the host.
* `model`: -(Optional) Rackable unit model name.

#### Hypervisor Isos
The hypervisor_isos block supports the following:

* `type`: -(Optional) Hypervisor type.
  Valid values are:
    - `XEN`: Xen hypervisor.
    - `HYPERV`: Hyper-V hypervisor.
    - `NATIVEHOST`: NativeHost type where AOS runs natively, without hypervisor.
    - `ESX`: ESX hypervisor.
    - `AHV`: AHV hypervisor.
* `md5sum`: -(Optional) Md5sum of ISO.

#### Bundle Info
The bundle_info block supports the following:

* `name`: -(Optional) Name of the hypervisor bundle.


#### Digital Certificate Map List
The `digital_certificate_map_list` attribute supports the following:

* `key`: -(Optional) Field containing digital_certificate_base64 and key_management_server_uuid for key management server.
* `name`: -(Optional) Value for the fields digital_certificate_base64 and key_management_server_uuid for key management server.

#### Ip Address Attributes
The `ipmi_ip`, `cvm_ip`, `hypervisor_ip` attributes supports the following:

* `ipv4`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: -(Optional) An unique address that identifies a device on the internet or a local network in IPv6 format.


### Config Params
The config_params block supports the following:

* `should_skip_discovery`: -(Optional) Indicates if node discovery need to be skipped or not.
* `should_skip_imaging`: -(Optional) Indicates if node imaging needs to be skipped or not.
* `should_validate_rack_awareness`: -(Optional) Indicates if rack awareness needs to be validated or not.
* `is_nos_compatible`: -(Optional) Indicates if node is compatible or not.
* `is_compute_only`: -(Optional) Indicates whether the node is compute only or not.
* `is_never_schedulable`: -(Optional) Indicates whether the node is marked to be never schedulable or not.
* `target_hypervisor`: -(Optional) Target hypervisor.
* `hiperv`: -(Optional) HyperV Credentials.

### Hiperv
The hiperv block supports the following:

* `domain_details`: -(Optional) UserName and Password model.
* `failover_cluster_details`: -(Optional) UserName and Password model.

#### Domain Details, Failover Cluster Details
The `domain_details`, `failover_cluster_details` attributes supports the following:

* `username`: -(Optional) Username.
* `password`: -(Optional) Password.
* `cluster_name`: -(Optional) Cluster name. This is part of payload for both cluster create & update operations.

### Remove Node Params
The remove_node_params block supports the following:

* `should_skip_prechecks`: -(Optional) Indicates if prechecks can be skipped for node removal.
* `should_skip_remove`: -(Optional) Indicates if node removal can be skipped.
* `node_uuids`: -(Required) List of node UUIDs to be removed.
* `extra_params`: -(Optional) Extra parameters for node addition.

#### Extra Params
The extra_params block supports the following:

* `should_skip_upgrade_check`: -(Optional) Indicates if upgrade check needs to be skipped or not.
* `skip_space_check`: -(Optional) Indicates if space check needs to be skipped or not.
* `should_skip_add_check`: -(Optional) Indicates if add check needs to be skipped or not.




See detailed information in [Nutanix Cluster - Add Node on a Cluster V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/Clusters/operation/expandCluster).


