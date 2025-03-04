---
layout: "nutanix"
page_title: "NUTANIX: nutanix_route_v2"
sidebar_current: "docs-nutanix-datasource-route-v2"
description: |-
  Get Route for the specified {extId}.
---

# nutanix_routes_v2

Provides Nutanix datasource Get Route for the specified {extId}.

## Example

```hcl

data "nutanix_route_v2" "route-by-id" {
  route_table_ext_id = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
  ext_id             = "7f66e20f-67f4-473f-96bb-c4fcfd487f16"
}

```

## Argument Reference

The following arguments are supported:

- `route_table_ext_id`: (Required) Route table UUID
- `ext_id`: (Required) Route UUID.

## Attribute Reference

The following attributes are exported:

- `ext_id`: Route UUID
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `metadata`: Metadata associated with this resource.
- `name`: Route name.
- `description`: BGP session description.
- `destination`: Destination IP Subnet Configuration.
- `next_hop` : Route nexthop.
- `route_table_reference`: Route table reference.
- `vpc_reference`: VPC reference.
- `external_routing_domain_reference`: External routing domain associated with this route table.
- `route_type`: Route type. Acceptable values are "STATIC", "LOCAL", "DYNAMIC"
- `is_active`: Indicates whether the route is active in the forwarding plane.
- `priority`: Route priority. A higher value implies greater preference is assigned to the route.

### metadata

- `owner_reference_id` : A globally unique identifier that represents the owner of this resource.
- `owner_user_name` : The userName of the owner of this resource.
- `project_reference_id` : A globally unique identifier that represents the project this resource belongs to.
- `project_name` : The name of the project this resource belongs to.
- `category_ids` : A list of globally unique identifiers that represent all the categories the resource is associated with.

### destination

- `ipv4`: IPv4 Subnet Object
- `ipv4.ip`: An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv4.ip.value`: The IPv4 address of the host.
- `ipv4.ip.prefix_length`: The prefix length of the network to which this host IPv4 address belongs.
- `ipv4.prefix_length`: The prefix length of the network to which this host IPv4 address belongs.

- `ipv6`: IPv6 Subnet Object
- `ipv6.ip`: IP address format
- `ipv6.ip.value`: The IPv6 address of the host.
- `ipv6.ip.prefix_length`: The prefix length of the network to which this host IPv6 address belongs.
- `ipv6.prefix_length`: The prefix length of the network to which this host IPv6 address belongs.

### next_hop

- `next_hop_type`: Nexthop type.
  - supported values:
    - `LOCAL_SUBNET`: - Next hop is an internal subnet.
    - `DIRECT_CONNECT_VIF`: - Next hop is a direct connect VIF.
    - `VPN_CONNECTION`: - Next hop is a VPN connection.
    - `IP_ADDRESS`: - Next hop is an IP address.
    - `EXTERNAL_SUBNET`: - Next hop is an external subnet.
- `next_hop_reference`: The reference to a link, such as a VPN connection or a subnet.
- `next_hop_ip_address`: An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.
- `next_hope_name`: Name of the nexthop, where the nexthop is either an IP address, a VPN connection, or a subnet.

### next_hop_ip_address

- `ipv4`: IPv4 Address
- `ipv6`: IPv6 Address

### IPv4/IPv6 Address

- `value`: value of IP address
- `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.

See detailed information in [Nutanix Get Route For Route Table v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/Routes/operation/getRouteForRouteTableById).
