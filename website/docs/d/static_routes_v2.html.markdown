---
layout: "nutanix"
page_title: "NUTANIX: nutanix_route_tables_v2"
sidebar_current: "docs-nutanix-datasource-routes-tables-v2"
description: |-
   This operation retrieves the list route tables.
---

# nutanix_route_table_v2

List route tables.

## Example Usage

```hcl
data "nutanix_route_tables_v2" "test1"{
}

```

## Argument Reference

The following arguments are supported:

* `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.

* `route_tables`: List of static routes.

## route_tables Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
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


See detailed information in [Nutanix Route Tables v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0).