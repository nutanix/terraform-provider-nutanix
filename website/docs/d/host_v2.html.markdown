---
layout: "nutanix"
page_title: "NUTANIX: nutanix_host_v2"
sidebar_current: "docs-nutanix-datasource-host-v2"
description: |-
 Describes the statistics data of the host identified by {hostExtId} belonging to the cluster identified by {clusterExtId}.
---

# nutanix_host_v2

Describes the statistics data of the host identified by {hostExtId} belonging to the cluster identified by {clusterExtId}.

## Example Usage

```hcl
data "nutanix_host_v2" "host"{
   cluster_ext_id = "021151dc-3ed1-4fec-a81d-39606451750c"
   ext_id = "919c9488-0b50-4fc8-9159-923e56a3abca"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: - (Required) host uuid
* `cluster_ext_id`: - (Required) cluster uuid

## Attributes Reference
The following attributes are exported:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: - image uuid.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `host_name`: - Name of the host.
* `host_type`: - Type of the host.
   * `HYPER_CONVERGED`: Hyper-converged node.
   * `COMPUTE_ONLY`: Compute only node.
   * `STORAGE_ONLY`: Storage only node.
* `hypervisor`: - Hypervisor details.
* `cluster`: - Cluster reference for an entity.
* `controller_vm`: - Host entity with its attributes.
* `disk`: - Disks attached to host.
* `is_degraded`: - Node degraded status.
* `is_secure_booted`: - Secure boot status.
* `is_hardware_virtualized`: - Indicates whether the hardware is virtualized or not.
* `has_csr`: - Certificate signing request status.
* `key_management_device_to_cert_status`: - Mapping of key management device to certificate status list.
* `number_of_cpu_cores`: - Number of CPU cores.
* `number_of_cpu_threads`: - Number of CPU threads.
* `number_of_cpu_sockets`: - Number of CPU sockets.
* `cpu_capacity_hz`: - CPU capacity in Hz.
* `cpu_frequency_hz`: - CPU frequency in Hz.
* `cpu_model`: - CPU model name.
* `gpu_driver_version`: - GPU driver version.
* `gpu_list`: - GPU attached list.
* `default_vhd_location`: - Default VHD location.
* `default_vhd_container_uuid`: - Default VHD container UUID.
* `default_vm_location`: - Default VM location.
* `default_vm_container_uuid`: - Default VM container UUID.
* `reboot_pending`: - Reboot pending status.
* `failover_cluster_fqdn`: - Failover cluster FQDN.
* `failover_cluster_node_status`: - Failover cluster node status.
* `boot_time_usecs`: - Boot time in secs.
* `memory_size_bytes`: - Memory size in bytes.
* `block_serial`: - Rackable unit serial name.
* `block_model`: - Rackable unit model name.
* `maintenance_state`: - Host Maintenance State.
* `node_status`: - Node status.
   * `TO_BE_PREPROTECTED`: Node to be preprotected.
   * `TO_BE_REMOVED`: Node to be removed.
   * `PREPROTECTED`: Node is preprotected.
   * `OK_TO_BE_REMOVED`: Indicates whether removing the node from the cluster is adequate.
   * `NORMAL`: Normal node.
   * `NEW_NODE`: New node.

### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Hypervisor
The hypervisor attribute supports the following:

* `external_address`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `user_name`: - Hypervisor user name.
* `full_name`: - Hypervisor full name.
* `type`: - Hypervisor type.
   * `XEN`: Xen hypervisor.
   * `HYPERV`: HyperV hypervisor.
   * `ESX`: ESX hypervisor.
   * `AHV`: AHV hypervisor.
* `number_of_vms`: - Number of VMs.
* `state`: - Hypervisor state.
   * `HA_HEALING_TARGET`: Hypervisor in HA healing target state.
   * `ENTERING_MAINTENANCE_MODE`: Hypervisor entering maintenance mode.
   * `RESERVED_FOR_HA_FAILOVER`: Hypervisor reserved for HA failover.
   * `HA_HEALING_SOURCE`: Hypervisor in HA healing source state.
   * `RESERVING_FOR_HA_FAILOVER`: Hypervisor that is planned to be reserved for HA failover.
   * `HA_FAILOVER_SOURCE`: Hypervisor in HA failover source state.
   * `ACROPOLIS_NORMAL`: Hypervisor in Acropolis normal state.
   * `ENTERED_MAINTENANCE_MODE`: Hypervisor entered maintenance mode.
   * `ENTERING_MAINTENANCE_MODE_FROM_HA_FAILOVER`: Hypervisor entering maintenance mode from HA failover.
   * `HA_FAILOVER_TARGET`: Hypervisor in HA failover target state.
* `acropolis_connection_state`: - Status of Acropolis connection to hypervisor.
   * `DISCONNECTED`: Acropolis disconnected.
   * `CONNECTED`: Acropolis connected.


### Cluster
The cluster attribute supports the following:

* `uuid`: - Cluster UUID.
* `name`: - Cluster name. This is part of payload for both cluster create & update operations.


### Controller VM
The controller_vm attribute supports the following:

* `id`: - Controller VM Id.
* `external_address`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `backplane_address`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `rdma_backplane_address`: - RDMA backplane address.
* `ipmi`: - IPMI reference.
* `nat_ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `nat_port`: - NAT port.
* `maintenance_mode`: - Maintenance mode status.
* `rackable_unit_uuid`: - Rackable unit UUID.


### Disk
The disk attribute supports the following:

* `uuid`: - Disk UUID.
* `mount_path`: - Disk mount path.
* `size_in_bytes`: - Disk size.
* `serial_id`: - Disk serial Id.
* `storage_tier`: - Disk storage Tier type.
   * `HDD`: HDD storage tier.
   * `PCIE_SSD`: PCIE SSD storage tier.
   * `SATA_SSD`: SATA SSD storage tier.

### key Management Device To Cert Status
The key_management_device_to_cert_status attribute supports the following:

* `key_management_server_name`: - Key management server name.
* `status`: - Certificate status.


#### external Address
The external_address attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.

#### Backplane Address
The backplane_address attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.

#### Rdma Backplane Address
The rdma_backplane_address attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.


#### ipmi
The ipmi attribute supports the following:

* `ip`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `username`: - IPMI username.


#### Nat Ip
The nat_ip attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.


##### ip

The ip attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.


###### IPV4

The ipv4 attribute supports the following:

* `value`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `prefix_length`: - The prefix length of the network to which this host IPv4 address belongs.

###### IPV6

The ipv6 attribute supports the following:

* `value`: - An unique address that identifies a device on the internet or a local network in IPv6 format.
* `prefix_length`: - The prefix length of the network to which this host IPv6 address belongs.



See detailed information in [Nutanix Get Host V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0#tag/Clusters/operation/getHostById).
