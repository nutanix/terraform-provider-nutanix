---
layout: "nutanix"
page_title: "NUTANIX: nutanix_locations_v2"
sidebar_current: "docs-nutanix-datasource-locations-v2"
description: |-
  List locations.
---

# nutanix_locations_v2

List locations.

## Example Usage

```hcl
# List all locations
data "nutanix_locations_v2" "all" {}

# Filter locations by country, sort by name and page the results
data "nutanix_locations_v2" "filtered" {
  limit    = 10
  filter   = "startswith(country, 'United')"
  order_by = "name desc"
}
```

## Argument Reference

The following arguments are supported:

- `page`:- A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results. If both \$page and \$apply query parameters are present, \$page will be applied on entities within the groups.
- `limit`:- A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set. If both \$limit and \$apply query parameters are present, \$limit will be applied on entities within the groups. Default value of limit with \$apply will be 5.
- `filter`:- A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'.
- `order_by`:- A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order.
- `select`:- A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned.

## Attributes Reference

The following attributes are exported:

* `locations`:- List of locations.

## Locations

The `locations` attribute is a list of location objects. Each location supports the following attributes:

* `ext_id`:- A globally unique identifier of an instance that is suitable for external consumption.
* `name`:- Location name.
* `country`:- Location Country.
* `state`:- Location state.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it is immutable from an API consumer perspective.
* `coordinates`:- The location coordinates, that is latitude and longitude.
* `links`:- A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### Coordinates

The `coordinates` attribute supports the following:

* `latitude`:- The latitude part of the location.
* `longitude`:- The longitude part of the location.

### Links

The `links` attribute supports the following:

* `href`:- The URL at which the entity described by the link can be accessed.
* `rel`:- A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

See detailed information in [Nutanix List Locations V4](https://developers.nutanix.com/api-reference?namespace=multidomain&version=v4.4).
