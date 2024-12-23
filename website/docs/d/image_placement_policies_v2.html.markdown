---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_placements_v2"
sidebar_current: "docs-nutanix-datasource-image-placements-v2"
description: |-
 Describes a Image placement policies
---

# nutanix_image_v2

List image placement policies details.

## Example

```hcl
data "nutanix_image_placement_policies_v2" "ipp"{
    page=0
    limit=10
    filter="startswith(name,'<name-prefix>')"
}

```


## Argument Reference

The following arguments are supported:

* `page`: A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit`: A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter`: A URL query parameter that allows clients to filter a collection of resources. The expression specified with $filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the $filter must conform to the OData V4.01 URL conventions
* `order_by`: A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default
* `select`: A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions.


## Attribute Reference

* `placement_policies`: List of all image placement policies

### placement_policies

The following arguments are supported:
* `ext_id`: The external identifier of an image placement policy.
* `name`: (Required) Name of the image placement policy.
* `description`: (Optional) Description of the image placement policy.
* `placement_type`: (Required) Type of the image placement policy. Valid values "HARD", "SOFT"
* `image_entity_filter`: (Required) Category-based entity filter.
* `cluster_entity_filter`: (Required) Category-based entity filter.
* `enforcement_state`: (Optional) Enforcement status of the image placement policy. Valid values "ACTIVE", "SUSPENDED"

### image_entity_filter
* `type`: (Required) Filter matching type. Valid values "CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"
* `category_ext_ids`: Array of strings

### cluster_entity_filter
* `type`: (Required) Filter matching type. Valid values "CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"
* `category_ext_ids`: Array of strings

See detailed information in [Nutanix Image placement policy V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0)