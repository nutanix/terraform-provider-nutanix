---
layout: "nutanix"
page_title: "NUTANIX: nutanix_static_routes"
sidebar_current: "docs-nutanix-resource-static-routes"
description: |-
  Create Static Routes within VPCs .
---

# nutanix_static_routes

Provides Nutanix resource to create Static Routes within VPCs.

## Example Usage

## create one static route for vpc uuid with external subnet

```hcl
resource "nutanix_static_routes" "scn" {
  vpc_uuid = "{{vpc_uuid}}"

  static_routes_list{
    destination= "10.x.x.x/x"
    external_subnet_reference_uuid = "{{ext_subnet_uuid}}" 
  }
}
```


## create one static route with default route for vpc name with external subnet

```hcl
resource "nutanix_static_routes" "scn" {
  vpc_name = "{{vpc_name}}"

  static_routes_list{
    destination= "10.x.x.x/x"
    external_subnet_reference_uuid = "{{ext_subnet_uuid}}" 
  }
  default_route_nexthop{
	  external_subnet_reference_uuid = "{{ext_subnet_uuid}}"
  }
}
```

#### Note: destination with 0.0.0.0/0 will be default route. 

## Argument Reference

The following arguments are supported:

* `vpc_uuid` - (Required) Reference to a VPC UUID. Should not be used with vpc_name.
* `vpc_name` - (Required) vpc Name. Should not be used with vpc_uuid. 
* `static_routes_list` - (Optional) Static Routes. 
* `default_route_nexthop`- (Optional) Default Route

### static_routes_list

* `destination` - (Required) Destination ip with prefix. 
* `external_subnet_reference_uuid` - (Optional) Reference to a subnet. Supported with 2022.x . 
* `vpn_connection_reference_uuid` - (Optional) Reference to a vpn connection.


### default_route_nexthop
* `external_subnet_reference_uuid` - (Required) Reference to a subnet.

## Attributes Reference

The following attributes are exported:

* `metadata` - The vpc_route_table kind metadata.
* `api_version` - The version of the API.

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when subnet was last updated.
* `UUID`: - subnet UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when subnet was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - subnet name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### lifecylce

Static Route can be managed but there is no destroy for resource. To delete the existing route you can remove the `static_route_list` from the `nutanix_static_routes` resource. Therefore, your existing routes created by resource will be deleted.  Refer example below

```hcl
resource "nutanix_static_routes" "scn" {
  vpc_uuid = "{{vpc_uuid}}"
}
```

See detailed information in [Nutanix Static Routes](https://www.nutanix.dev/api_references/prism-central-v3/#/56796ae9af040-update-a-existing-vpc-route-table).