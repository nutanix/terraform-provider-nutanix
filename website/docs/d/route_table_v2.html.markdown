---
layout: "nutanix"
page_title: "NUTANIX: nutanix_route_table_v2"
sidebar_current: "docs-nutanix-datasource-route-table-v2"
description: |-
  Get the route table for the specified extId.

---

# nutanix_route_table_v2

Provides Nutanix datasource Get the route table for the specified extId.


## Example Usage

```hcl

data "nutanix_route_table_v2" "get-route-table"{
  ext_id = "ad7692ff-50f8-4a47-b3c9-fd5a3462d455"
}

```


## Argument Reference

The following arguments are supported:
* `ext_id`: (Required) Route UUID

## Attribute Reference
The following attributes are exported:
* `ext_id`: Route UUID
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `metadata`: Metadata associated with this resource.
* `vpc_reference`:  VPC reference.
* `external_routing_domain_reference`:  External routing domain associated with this route table.

### metadata
* `owner_reference_id` :  A globally unique identifier that represents the owner of this resource.
* `owner_user_name` :  The userName of the owner of this resource.
* `project_reference_id` :  A globally unique identifier that represents the project this resource belongs to.
* `project_name` :  The name of the project this resource belongs to.
* `category_ids` :  A list of globally unique identifiers that represent all the categories the resource is associated with.



See detailed information in [Nutanix Get Route Table v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/RouteTables/operation/getRouteTableById).
