---
layout: "nutanix"
page_title: "NUTANIX: nutanix_entity_groups_v2"
sidebar_current: "docs-nutanix-datasource-entity-groups-v2"
description: |-
  Retrieves a list of entity groups for microsegmentation.
---

# nutanix_entity_groups_v2

Retrieves a list of entity groups. Use this data source when you need to list or filter entity groups rather than fetch a single group by `ext_id` (use the `nutanix_entity_group_v2` data source for that).

## Example Usage

```hcl
# List all entity groups
data "nutanix_entity_groups_v2" "list" {
}

# List entity groups with filter
data "nutanix_entity_groups_v2" "filtered" {
  filter = "name eq 'my-entity-group'"
}

# List entity groups with limit
data "nutanix_entity_groups_v2" "with_limit" {
  limit = 10
}

# List entity groups with filter and limit
data "nutanix_entity_groups_v2" "filtered_limit" {
  filter = "name eq 'my-entity-group'"
  limit  = 1
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources. The expression must conform to OData V4.01 URL conventions. The filter can be applied to the following fields:
  - `name`
  - `extId`
  - `description`
  - `creationTime`
  - `lastUpdateTime`
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. The orderby can be applied to the following fields:
  - `name`
  - `extId`
  - `description`
  - `creationTime`
  - `lastUpdateTime`
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the $select must conform to the OData V4.01 URL conventions.

## Attribute Reference

The following attributes are exported:

* `entity_groups` - List of entity groups.
* `id` - A synthetic identifier used internally by Terraform.

### Entity Group (list element)

Each element in `entity_groups` contains:

* `ext_id` - A globally unique identifier (UUID) of the entity group.
* `name` - A short identifier of the Entity Group.
* `description` - A user defined annotation for an Entity Group.
* `allowed_config` - Configuration of the allowed entities in the Entity Group.
* `except_config` - Configuration of except entities in the Entity Group.
* `policy_ext_ids` - List of policy external identifiers associated with the entity group.
* `creation_time` - The timestamp when the Entity Group was created.
* `last_update_time` - The timestamp when the Entity Group was last updated.
* `links` - A HATEOAS style link for the response.
* `owner_ext_id` - The external identifier of the user who created the Entity Group.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

### allowed_config (list element)

* `entities` - List of allowed entities. Each entity may contain `type`, `selected_by`, `addresses`, `ip_ranges`, `kube_entities`, `reference_ext_ids`.

### except_config (list element)

* `entities` - List of except entities. Each entity may contain `addresses`, `ip_ranges`, `reference_ext_ids`.

See detailed information in [Nutanix List Entity Groups V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.2#tag/EntityGroups/operation/listEntityGroups).
