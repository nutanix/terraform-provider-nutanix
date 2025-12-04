---
layout: "nutanix"
page_title: "NUTANIX: nutanix_routes_v2"
sidebar_current: "docs-nutanix-resource-routes-v2"
description: |-
  Create Route.
---

# nutanix_routes_v2

Provides Nutanix resource to Create Route.

## Example

```hcl

# create a route
resource "nutanix_routes_v2" "route" {
  name               = "terraform_example_route"
  description        = "terraform example route to example create route"
  vpc_reference      = "8a938cc5-282b-48c4-81be-de22de145d07"
  route_table_ext_id = "c2c249b0-98a0-43fa-9ff6-dcde578d3936"
  destination {
    ipv4 {
      ip {
        value = "10.0.0.2"
      }
      prefix_length = 32
    }
  }
  next_hop {
    next_hop_type      = "EXTERNAL_SUBNET"
    next_hop_reference = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  }
  metadata {
    owner_reference_id   = "a8fe48c4-f0d3-49c7-a017-efc30dd8fb2b"
    project_reference_id = "ab520e1d-4950-1db1-917f-a9e2ea35b8e3"
  }
  route_type = "STATIC"
}

```

## Argument Reference

The following arguments are supported:

- `route_table_ext_id`: (Required) Route table UUID
- `metadata`: (Optional) Metadata associated with this resource.
- `name`: (Optional) Route name.
- `description`: (Optional) BGP session description.
- `destination`: (Optional) Destination IP Subnet Configuration.
- `next_hop` : (Optional) Route nexthop.
- `route_table_reference`: (Optional) Route table reference.
- `vpc_reference`: (Optional) VPC reference.
- `external_routing_domain_reference`: (Optional) External routing domain associated with this route table.
- `route_type`: (Optional) Route type. Acceptable values are "STATIC", "LOCAL", "DYNAMIC"

### metadata

- `owner_reference_id` : (Optional) A globally unique identifier that represents the owner of this resource.
- `owner_user_name` : (Optional) The userName of the owner of this resource.
- `project_reference_id` : (Optional) A globally unique identifier that represents the project this resource belongs to.
- `project_name` : (Optional) The name of the project this resource belongs to.
- `category_ids` : (Optional) A list of globally unique identifiers that represent all the categories the resource is associated with.

### destination

- `ipv4`: (Optional) IPv4 Subnet Object
- `ipv4.ip`: (Required) An unique address that identifies a device on the internet or a local network in IPv4 format.
- `ipv4.ip.value`: (Required) The IPv4 address of the host.
- `ipv4.ip.prefix_length`: (Optional) The prefix length of the network to which this host IPv4 address belongs.
- `ipv4.prefix_length`: (Required) The prefix length of the network to which this host IPv4 address belongs.

- `ipv6`: (Optional) IPv6 Subnet Object
- `ipv6.ip`: (Required) IP address format
- `ipv6.ip.value`: (Required) The IPv6 address of the host.
- `ipv6.ip.prefix_length`: (Optional) The prefix length of the network to which this host IPv6 address belongs.
- `ipv6.prefix_length`: (Required) The prefix length of the network to which this host IPv6 address belongs.

### next_hop_ip_address

- `ipv4`: (Optional) IPv4 Address
- `ipv6`: (Optional) IPv6 Address

### IPv4/IPv6 Address

- `value`: (Optional) value of IP address
- `prefix_length`: (Optional) The prefix length of the network to which this host IPv4/IPv6 address belongs.

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

### Links

The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

## Import

This helps to manage existing entities which are not created through terraform. Route can be imported using the route table and route uuid `routeTableUUID/routeUUID` (ext_id in v4 terms). eg,

**Note**:To import Route, you need to have the Route Table Ext ID and Route Ext ID, and provide them in the format mentioned above while importing.

```hcl
// create its configuration in the root module. For example:
resource "nutanix_routes_v2" "import_route"{}

// execute the below command. UUID can be fetched using datasource. Example:
data "nutanix_routes_v2" "fetch_templates"{
  route_table_ext_id = "c2c249b0-98a0-43fa-9ff6-dcde578d3936"
}

terraform import nutanix_routes_v2.import_route <routeTableUUID>/<routeUUID>
```

See detailed information in [Nutanix Routes v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/Routes/operation/createRouteForRouteTable).
