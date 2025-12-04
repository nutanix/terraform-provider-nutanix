---
layout: "nutanix"
page_title: "NUTANIX: nutanix_static_routes"
sidebar_current: "docs-nutanix-datasource-static-routes"
description: |-
   This operation retrieves a static routes within VPCs.
---

# nutanix_static_routes

Provides a datasource to retrieve static routes within VPCs given vpc_uuid.

## Example Usage

```hcl
data "nutanix_static_routes" "test1"{
    vpc_reference_uuid = <vpc_reference_uuid>
}

data "nutanix_vpc" "test2"{
    vpc_name = <vpc_name>
}
```

## Argument Reference

The following arguments are supported:

* `vpc_reference_uuid` - vpc UUID

## Attribute Reference

The following attributes are exported:

* `api_version` - API version
* `metadata` -  The vpc_route_table kind metadata
* `spec` - An intentful representation of a vpc_route_table spec
* `status` - An intentful representation of a vpc_route_table status

### spec
* `name` - vpc_route_table Name.
* `resources` - VPC route table resources

### status
* `state` - The state of the vpc_route_table.
* `resources` - VPC route table resources status
* `execution_context` - Execution Context of VPC. 

### resources

* `static_routes_list` - list of static routes
* `default_route_nexthop` - default routes (present in spec resource)
* `default_route` - default route. (present in status resource only )
* `local_routes_list` - list of local routes (present in status resource only )
* `dynamic_routes_list` - list of dynamic routes (present in status resource only)

### static_routes_list

* `nexthop` - Targeted link to use as the nexthop in a route. 
* `destination` - destination ip address with prefix. 
* `priority` - The preference value assigned to this route. A higher value means greater preference. Present in Status Resource.
* `is_active` - Whether this route is currently active. Present in Status Resources. 

### nexthop 

* `external_subnet_reference` - The reference to a subnet
* `direct_connect_virtual_interface_reference` - The reference to a direct_connect_virtual_interface
* `local_subnet_reference` - The reference to a subnet
* `vpn_connection_reference` - The reference to a vpn_connection

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Reference

The  `external_subnet_reference`  attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `uuid`: - the UUID.

See detailed information in [Nutanix Static Route](https://www.nutanix.dev/api_references/prism-central-v3/#/c936631dbba81-get-a-existing-vpc-route-table).