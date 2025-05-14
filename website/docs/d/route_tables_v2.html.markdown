---
layout: "nutanix"
page_title: "NUTANIX: nutanix_route_tables_v2"
sidebar_current: "docs-nutanix-datasource-route-tables-v2"
description: |-
  List route tables

---

# nutanix_route_tables_v2

Provides Nutanix datasource to List route tables.

## Example Usage

```hcl

data "nutanix_route_tables_v2" "all-tables"{}


data "nutanix_route_tables_v2" "route-tables-with-filter"{
  filter = "vpcReference eq 'f4b4b3b4-4b4b-4b4b-4b4b-4b4b4b4b4b4b'"
}

data "nutanix_route_tables_v2" "route-tables-with-orderby" {
  order_by = "vpcReference"
}

```


## Argument Reference

The following arguments are supported:
* `page`: (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: (Optional) A URL query parameter that allows clients to filter a collection of resources.
    * The filter can be applied to the following fields:
        * `externalRoutingDomainReference`
        * `vpcReference`
* `order_by`: (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default
    * The orderby can be applied to the following fields:
        * `vpcReference`

## Attribute Reference
The following attributes are exported:

* `route_tables`: A list of route tables.

### Route Tables
The `route_tables` object contains the following attributes:

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



See detailed information in [Nutanix Route Tables v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/RouteTables).
