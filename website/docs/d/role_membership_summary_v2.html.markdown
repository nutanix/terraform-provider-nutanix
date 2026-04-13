---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_membership_summary_v2"
sidebar_current: "docs-nutanix-datasource-role-membership-summary-v2"
description: |-
  Lists role membership summaries in Nutanix.
---

# nutanix_role_membership_summary_v2

Lists role membership summaries in Nutanix. Provides aggregated counts of users, groups, roles, and total identities.

## Example Usage

```hcl
data "nutanix_role_membership_summary_v2" "example" {}
```

### With Filters

```hcl
data "nutanix_role_membership_summary_v2" "filtered" {
  page  = 0
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) A URL query parameter that specifies the page number of the result set.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

The following attributes are exported:

* `summaries` - List of role membership summaries.

### Summary

Each summary in `summaries` exports the following:

* `ext_id` - External identifier of the role membership summary.
* `tenant_id` - Tenant identifier.
* `links` - A HATEOAS style link for the response.
* `users_count` - Count of distinct users.
* `groups_count` - Count of distinct groups.
* `roles_count` - Count of distinct roles.
* `total_identities_count` - Total count of identities.
