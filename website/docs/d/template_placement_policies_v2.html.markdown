---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_placement_policies_v2"
sidebar_current: "docs-nutanix-datasource-template-placement-policies-v2"
description: |-
  Lists the template placement policies created on Prism Central, including details such as name, description and more.
---

# nutanix_template_placement_policies_v2

Lists the template placement policies created on Prism Central, including details such as name, description and more. The API supports operations like filtering, sorting, selection and pagination.

## Example

```hcl
data "nutanix_template_placement_policies_v2" "example" {}

data "nutanix_template_placement_policies_v2" "filtered" {
  filter = "name eq 'my-policy'"
  limit  = 10
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) A URL query parameter that specifies the page number of the result set.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attribute Reference

The following attributes are exported:

* `template_placement_policies` - List of template placement policies.

### template_placement_policies

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `name` - Name of the placement policy.
* `description` - Description of the placement policy.
* `placement_type` - The placement type. Valid values: `SOFT`.
* `cluster_filter` - Category filter for clusters.
* `content_filter` - Category filter for content.
* `create_time` - The time when the placement policy was created.
* `created_by` - External identifier of the user who created the placement policy.
* `update_time` - The time when the placement policy was last updated.
* `updated_by` - External identifier of the user who updated the placement policy.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

### cluster_filter, content_filter

* `category_ext_ids` - Filter to match entities based on the provided categories.
* `type` - The match type for categories. Valid values: `CATEGORIES_MATCH_ALL`, `CATEGORIES_MATCH_ANY`.
