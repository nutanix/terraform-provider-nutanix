---
layout: "nutanix"
page_title: "NUTANIX: nutanix_routes_v2"
sidebar_current: "docs-nutanix-resource-routes-v2"
description: |-
  Create Route.
---

# nutanix_routes_v2

Provides Nutanix resource to Create Route.


## Argument Reference

The following arguments are supported:

* `route_table_ext_id`: (Required) Route table UUID
* `metadata`: (Optional) Metadata associated with this resource.
* `name`: (Optional) Route name.
* `description`: (Optional) BGP session description.
* `destination`: (Optional) Destination IP Subnet Configuration.
* `next_hop` : (Optional) Route nexthop.
* `route_table_reference`: (Optional) Route table reference.
* `vpc_reference`: (Optional) VPC reference.
* `external_routing_domain_reference`: (Optional) External routing domain associated with this route table.
* `route_type`: (Optional) Route type. Acceptable values are "STATIC", "LOCAL", "DYNAMIC"

### metadata
* `owner_reference_id` : (Optional) A globally unique identifier that represents the owner of this resource.
* `owner_user_name` : (Optional) The userName of the owner of this resource.
* `project_reference_id` : (Optional) A globally unique identifier that represents the project this resource belongs to.
* `project_name` : (Optional) The name of the project this resource belongs to.
* `category_ids` : (Optional) A list of globally unique identifiers that represent all the categories the resource is associated with.


### destination
* `ipv4`: (Optional) IPv4 Subnet Object
* `ipv4.ip`: (Required) An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv4.ip.value`: (Required) The IPv4 address of the host.
* `ipv4.ip.prefix_length`: (Optional) The prefix length of the network to which this host IPv4 address belongs.
* `ipv4.prefix_length`: (Required) The prefix length of the network to which this host IPv4 address belongs.

* `ipv6`: (Optional) IPv6 Subnet Object
* `ipv6.ip`: (Required) IP address format
* `ipv6.ip.value`: (Required) The IPv6 address of the host.
* `ipv6.ip.prefix_length`: (Optional) The prefix length of the network to which this host IPv6 address belongs.
* `ipv6.prefix_length`: (Required) The prefix length of the network to which this host IPv6 address belongs.


### next_hop_ip_address
* `ipv4`: (Optional) IPv4 Address
* `ipv6`: (Optional) IPv6 Address


### IPv4/IPv6 Address
* `value`: (Optional) value of IP address
* `prefix_length`: (Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.

## Attribute Reference
The following attributes are exported:
* `ext_id`: Route UUID
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `metadata`: Metadata associated with this resource.
* `name`:  Route name.
* `description`:  BGP session description.
* `destination`:  Destination IP Subnet Configuration.
* `next_hop` :  Route nexthop.
* `route_table_reference`:  Route table reference.
* `vpc_reference`:  VPC reference.
* `external_routing_domain_reference`:  External routing domain associated with this route table.
* `route_type`: Route type. Acceptable values are "STATIC", "LOCAL", "DYNAMIC"
* `is_active`:  Indicates whether the route is active in the forwarding plane.
* `priority`:  Route priority. A higher value implies greater preference is assigned to the route.

### Links
The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.



See detailed information in [Nutanix Routes v2](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0.b1).