---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_add_node_v2"
sidebar_current: "docs-nutanix-resource-nutanix-cluster-add-node-v2"
description: |-
  Add node on a cluster identified by {extId}.
---

# nutanix_cluster_add_node_v2

Add node on a cluster identified by {extId}.

## Example Usage

```hcl
resource "nutanix_cluster_add_node_v2" "cluster_node"{
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  node_params {
    block_list{
      node_list{
        node_uuid = "00000000-0000-0000-0000-000000000000"
        block_uuid = "00000000-0000-0000-0000-000000000000"
        node_position = "<node_position>"
        hypervisor_type = "XEN"
        is_robo_mixed_hypervisor = true
        hypervisor_hostname = "<hypervisor_hostname>"
        hypervisor_version = "9.9.99"
        nos_version = "9.9.99"
        ipmi_ip {
          ipv4{
            value = "10.0.0.1"
          }
        }
        digital_certificate_map_list{
          key = "key"
          value = "value"   
        }
        model = "<model>"
      }
      should_skip_host_networking = false
    }
    node_list{
        node_uuid = "00000000-0000-0000-0000-000000000000"
        block_uuid = "00000000-0000-0000-0000-000000000000"
        node_position = "<node_position>"
        hypervisor_type = "XEN"
        is_robo_mixed_hypervisor = true
        hypervisor_hostname = "<hypervisor_hostname>"
        hypervisor_version = "9.9.99"
        nos_version = "9.9.99"
        ipmi_ip {
          ipv4 {
            value = "10.0.0.1"
          }
        }
      
    }
    bundle_info{
      name = "<name>"
    }
  }
  config_params{
    should_skip_discovery = false
    should_skip_imaging = true
    is_nos_compatible = true
    target_hypervisor = "<target_hypervisor>"
  }
  should_skip_add_node = false
  should_skip_pre_expand_checks = true

  remove_node_params {
    extra_params {
      should_skip_upgrade_check = false
      skip_space_check          = false
      should_skip_add_check     = false
    }
    should_skip_remove    = false
    should_skip_prechecks = false
  }
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




See detailed information in [Nutanix Cluster - Add Node on a Cluster V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0.b2#tag/Clusters/operation/expandCluster).


