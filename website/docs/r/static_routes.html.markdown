---
layout: "nutanix"
page_title: "NUTANIX: nutanix_static_routes"
sidebar_current: "docs-nutanix-resource-static-routes"
description: |-
  Create Static Routes within VPCs .
---

# nutanix_static_routes

Provides Nutanix resource to create Floating IPs. 

## Example Usage

## create one static route for vpc with external subnet

```hcl
resource "nutanix_static_routes" "scn" {
  vpc_uuid = "{{vpc_uuid}}"

  static_routes_list{
    destination= "10.x.x.x/x"
    external_subnet_reference_uuid = "{{ext_subnet_uuid}}" 
  }
}
```

## Argument Reference

The following arguments are supported:

*`vpc_uuid` - (Required) Reference to a VPC .
*`static_routes_list` - (Required) Static Routes. 

## static_routes_list

*`destination` - (Required) Destination ip with prefix. 
*`external_subnet_reference_uuid` - (Required) Reference to a subnet.


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