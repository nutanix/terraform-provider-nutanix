---
layout: "nutanix"
page_title: "NUTANIX: nutanix_routes_v2"
sidebar_current: "docs-nutanix-datasource-routes-v2"
description: |-
  List Routes request.
---

# nutanix_routes_v2

Provides Nutanix resource to List Routes request.

## Example

```hcl
# List all route tables
data "nutanix_route_tables_v2" "list-route-tables" {}


# List routes by route table id
data "nutanix_routes_v2" "routes"{
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
}

# list all routes with filter
data "nutanix_routes_v2" "filtered-routes" {
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
  filter             = "vpcReference eq '${nutanix_vpc_v2.example.id}' and isActive eq true"
  order_by           = "name desc" # list all routes sorted by name in descending order
  depends_on         = [nutanix_routes_v2.route]
}

# List routes by route table id with limit 3
data "nutanix_routes_v2" "routes"{
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
  limit = 3
}
```

## Argument Reference

The following arguments are supported:

- `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
  - The filter can be applied to the following fields:
    - `destination`
    - `externalRoutingDomainReference`
    - `isActive`
    - `name`
    - `nexthop/nexthopName`
    - `nexthop/nexthopReference`
    - `nexthop/nexthopType`
    - `priority`
    - `routeTableReference`
    - `routeType`
    - `vpcReference`
- `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default
  - The orderby can be applied to the following fields:
    - `isActive`
    - `name`
    - `priority`
    - `routeType`
- `route_table_ext_id`: (Required) Route table UUID

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

### next_hop_ip_address

- `ipv4`: IPv4 Address
- `ipv6`: IPv6 Address

### IPv4/IPv6 Address

- `value`: value of IP address
- `prefix_length`: The prefix length of the network to which this host IPv4/IPv6 address belongs.

See detailed information in [Nutanix List Routes By Route Table Id v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/Routes/operation/listRoutesByRouteTableId).
