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

data "nutanix_clusters_v2" "clusters" {}
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# create a subnet
resource "nutanix_subnet_v2" "ext-subnet" {
  name              = "subnet_for_route"
  description       = "subnet to test create route"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = "198"
  is_external       = true
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "10.44.3.192"
        }
        prefix_length = "27"
      }
      default_gateway_ip {
        value = "10.44.3.193"
      }
      pool_list {
        start_ip {
          value = "10.44.3.198"
        }
        end_ip {
          value = "10.44.3.207"
        }
      }
    }
  }
}

# crete a vpc
resource "nutanix_vpc_v2" "vpc" {
  name        = "terraform_example_vpc_1"
  description = "terraform example vpc 1 to test create route"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.ext-subnet.id
  }
  depends_on = [nutanix_subnet_v2.ext-subnet]
}

# get route table
data "nutanix_route_tables_v2" "list-route-tables" {
  filter     = "vpcReference eq '${nutanix_vpc_v2.vpc.id}'"
  depends_on = [nutanix_vpc_v2.vpc]
}

# create a project
resource "nutanix_project" "example-project" {
  name        = "tf-example-project"
  description = "terraform example project"
  default_subnet_reference {
    kind = "subnet"
    uuid = nutanix_subnet_v2.ext-subnet.id
  }
  lifecycle {
    ignore_changes = [default_subnet_reference]
  }
}
# create a route
resource "nutanix_routes_v2" "route" {
  name               = "terraform_example_route"
  description        = "terraform example route to example create route"
  vpc_reference      = nutanix_vpc_v2.vpc.id
  route_table_ext_id = data.nutanix_route_tables_v2.list-route-tables.route_tables[0].ext_id
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
    next_hop_reference = nutanix_subnet_v2.ext-subnet.id
  }
  metadata {
    owner_reference_id   = nutanix_vpc_v2.vpc.id
    project_reference_id = nutanix_project.example-project.metadata.uuid
  }
  route_type = "STATIC"
}

```

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



See detailed information in [Nutanix Routes v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/Routes/operation/createRouteForRouteTable).
