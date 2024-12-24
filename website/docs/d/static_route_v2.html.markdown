---
layout: "nutanix"
page_title: "NUTANIX: nutanix_route_table_v2"
sidebar_current: "docs-nutanix-datasource-route-table-v2"
description: |-
   This operation retrieves the route table for the specified extId.
---

# nutanix_route_table_v2

Provides a datasource to retrieve a route table.

## Example Usage

```hcl
data "nutanix_route_table_v2" "test1"{
    ext_id = {{ route table uuid }}
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: Route table UUID

## Attribute Reference

The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. 
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `metadata`: Metadata associated with this resource.
* `vpc_reference`: VPC
* `external_routing_domain_reference`: External routing domain associated with this route table
* `static_routes`: Static routes
* `dynamic_routes`: Dynamic routes
* `local_routes`: Routes to local subnets


### static_routes, dynamic_routes, local_routes
* `is_active`: Indicates whether the route is active or inactive.
* `priority`: Route priority. A higher value implies greater preference is assigned to the route.
* `destination`: Destination IPv4/IPv6 Object. 
* `next_hop_type`: Next hop type.
* `next_hop_reference`: The reference to a link, such as a VPN connection or a subnet. 
* `next_hop_ip_address`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
* `next_hop_name`: Name of the next hop, where the next hop is either a VPN connection, direct connect virtual interface, or a subnet.
* `source`: The source of a dynamic route is either a VPN connection, direct connect virtual interface, or a BGP session


### destination
* `ipv4`: IPv4 Subnet Object
* `ipv4.ip`: IP address format
* `ipv4.prefix_length`: The prefix length of the network to which this host IPv4 address belongs.

* `ipv6`: IPv6 Subnet Object
* `ipv6.ip`: IP address format
* `ipv6.prefix_length`: The prefix length of the network to which this host IPv6 address belongs.


### next_hop_ip_address
* `ipv4`: IPv4 Address
* `ipv6`: IPv6 Address


### IPv4/IPv6 Address
* `value`: value of IP address
* `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.


See detailed information in [Nutanix Route Table v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0).