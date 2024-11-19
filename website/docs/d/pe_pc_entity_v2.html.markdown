---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pe_pc_entity_v2"
sidebar_current: "docs-nutanix-datasource-pe-pc-entity-v2"
description: |-
  Fetches the attributes associated with the domain manager (Prism Central) resource based on the provided external identifier. It includes attributes like config, network, node and other information such as size, environment and resource specifications.


---

# nutanix_pe_pc_entity_v2

Fetches the attributes associated with the domain manager (Prism Central) resource based on the provided external identifier. It includes attributes like config, network, node and other information such as size, environment and resource specifications.



## Example Usage

```hcl
data "nutanix_pe_pc_entity_v2" "pc" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}

```

## Argument Reference

The following arguments are supported:
* `ext_id` : -(Required) The external identifier of the domain manager (Prism Central) resource. 



## Attributes Reference
The following attributes are exported:

* `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `config`: - Domain manager (Prism Central) cluster configuration details.
* `is_registered_with_hosting_cluster`: - Boolean value indicating if the domain manager (Prism Central) is registered with the hosting cluster, that is, Prism Element.
* `network`: - Domain manager (Prism Central) network configuration details.
* `hosting_cluster_ext_id`: - The external identifier of the cluster hosting the domain manager (Prism Central) instance.
* `should_enable_high_availability`: - This configuration enables Prism Central to be deployed in scale-out mode.
* `node_ext_ids`: - Domain manager (Prism Central) nodes external identifier.

### Links
The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Config
The config attribute supports the following:

* `should_enable_lockdown_mode`: - A boolean value indicating whether to enable lockdown mode for a cluster.
* `build_info`: - Currently representing the build information to be used for the cluster creation.
* `name`: - Name of the domain manager (Prism Central).
* `size`: - Domain manager (Prism Central) size is an enumeration of starter, small, large, or extra large starter values. 
   Valid values are:
    - `SMALL`: - Domain manager (Prism Central) of size small.
    - `LARGE`: - Domain manager (Prism Central) of size large.
    - `EXTRALARGE`: - Domain manager (Prism Central) of size extra large.
    - `STARTER`: - Domain manager (Prism Central) of size starter.
* `bootstrap_config`: - Bootstrap configuration details for the domain manager (Prism Central).
* `resource_config`: - This configuration is used to provide the resource-related details like container external identifiers, number of VCPUs, memory size, data disk size of the domain manager (Prism Central). In the case of a multi-node setup, the sum of resources like number of VCPUs, memory size and data disk size are provided.


#### Build Info
The build_info attribute supports the following:

* `version`: - Software version.

#### Bootstrap Config
The bootstrap_config attribute supports the following:

* `environment_info`: - An object denoting the environment information of the PC. It contains the following fields:
  - type: Enums denoting the environment type of the PC.
  - providerType: Enums denoting the provider of the cloud PC.
  - instanceObj: Enums denoting the instance type of the cloud PC.

##### Environment Info
The environment_info attribute supports the following:

* `type`: - Enums denoting the environment type of the PC, that is, on-prem PC or cloud PC.
  Following are the supported entity types:
  - ONPREM: - On-prem environment.
  - NTNX_CLOUD: - Nutanix cloud environment.
* `provider_type`: - Enums denoting the provider of the cloud, in case of environment type a cloud PC.
  The service currently supports the following providers:
  - NTNX: - Nutanix cloud provider.
  - AZURE - Azure cloud provider.
  - AWS - AWS cloud provider.
  - GCP - GCP cloud provider.
  - VSPHERE - Vsphere cloud provider.
  * `provisioning_type`: - Enums denoting the instance type of the cloud PC. It indicates whether the PC is created on bare-metal or on a cloud-provisioned VM. Hence,
     it supports two possible values:
    - NTNX: - Nutanix instance.
    - NATIVE: - Native instance.

#### Resource Config
The resource_config attribute supports the following:

* `num_vcpus`: - This property is used for readOnly purposes to display Prism Central number of VCPUs allocation.
* `memory_size_bytes`: - This property is used for readOnly purposes to display Prism Central RAM allocation at the cluster level.
* `data_disk_size_bytes`: - This property is used for readOnly purposes to display Prism Central data disk size allocation at a cluster level.
* `container_ext_ids`: - The external identifier of the container that will be used to create the domain manager (Prism Central) cluster.



### Network
The network attribute supports the following:

* `external_address`: - An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `name_servers`: - List of name servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
* `ntp_servers`: - List of NTP servers on a cluster. This is part of payload for both cluster create & update operations. For create operation, only ipv4 address / fqdn values are supported currently.
* `fqdn`: - Cluster fully qualified domain name. This is part of payload for cluster update operation only.
* `external_networks`: - This configuration is used to manage Prism Central.


#### external Address
The external_address attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.

#### Name Servers, NTP Servers
The backplane_address attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.
* `fqdn`: - A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

#### External Networks
The external_networks attribute supports the following:

* `default_gateway`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `subnet_mask`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `ip_ranges`: - Range of IPs used for Prism Central network setup.
* `network_ext_id`: - The network external identifier to which Domain Manager (Prism Central) is to be deployed or is already configured.

##### Default Gateway, Subnet Mask
The default_gateway, subnet_mask attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.
* `fqdn`: - A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

##### IP Ranges
The ip_ranges attribute supports the following:

* `begin`: - The beginning of the range of IP addresses.
* `end`: - The end of the range of IP addresses.

##### Begin, End
The begin, end attribute supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.

###### IPV4, IPV6

The ipv4 attribute supports the following:

* `value`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format.
* `prefix_length`: - The prefix length of the network to which this host IPv4/IPv6 address belongs.

###### FQDN

The fqdn attribute supports the following:

* `value`: - The fully qualified domain name of the host.


See detailed information in [Nutanix PC Details](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0.b1#tag/DomainManager/operation/getDomainManagerById).
