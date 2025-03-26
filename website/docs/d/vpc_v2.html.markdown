---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vpc_v2"
sidebar_current: "docs-nutanix-datasource-vpc-v2"
description: |-
   This operation retrieves a vpc based on the input parameters.
---

# nutanix_vpc_v2

Provides a datasource to retrieve VPC with vpc_uuid .

## Example Usage

```hcl
data "nutanix_vpc_v2" "get-vpc"{
  ext_id = "8a938cc5-282b-48c4-81be-de22de145d07"
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) vpc UUID

## Attribute Reference

The following attributes are exported:

* `name`: Name of the VPC.
* `description`: Description of the VPC.
* `common_dhcp_options`: List of DHCP options to be configured.
* `vpc_type`: Type of VPC.
* `snat_ips`: List of IP Addresses used for SNAT.
* `external_subnets`: List of external subnets that the VPC is attached to.
* `external_routing_domain_reference`: External routing domain associated with this route table
* `externally_routable_prefixes`: CIDR blocks from the VPC which can talk externally without performing NAT. This is applicable when connecting to external subnets which have disabled NAT.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `metadata`: Metadata associated with this resource.


### common_dhcp_options

* `domain_name_servers`: List of Domain Name Server addresses
* `domain_name_servers.ipv4`: Reference to address configuration
* `domain_name_servers.ipv6`: Reference to address configuration


### external_subnets

* `subnet_reference`: External subnet reference.
* `external_ips`: List of IP Addresses used for SNAT, if NAT is enabled on the external subnet. If NAT is not enabled, this specifies the IP address of the VPC port connected to the external gateway.
* `gateway_nodes`: List of gateway nodes that can be used for external connectivity.
* `active_gateway_node`: Reference of gateway nodes
* `active_gateway_count`: Maximum number of active gateway nodes for the VPC external subnet association.


### snat_ips, external_ips

* `ipv4`: Reference to address configuration
* `ipv6`: Reference to address configuration


### externally_routable_prefixes
* `ipv4`: IP V4 Configuration
* `ipv4.ip`: Reference to address configuration
* `ipv4.prefix_length`: The prefix length of the network

* `ipv6`: IP V6 Configuration
* `ipv6.ip`: Reference to address configuration
* `ipv6.prefix_length`: The prefix length of the network


### ipv4, ipv6 (Reference to address configuration)

* `value`: value of address
* `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.



See detailed information in [Nutanix Get VPC v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/Vpcs/operation/getVpcById).
