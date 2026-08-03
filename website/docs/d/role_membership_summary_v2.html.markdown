---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_membership_summary_v2"
sidebar_current: "docs-nutanix-datasource-role-membership-summary-v2"
description: |-
  Lists role membership summaries.
---

# nutanix_role_membership_summary_v2

Lists role membership summaries. Each record represents a project and returns the count of identities (users and groups) and roles for that project. Use the $filter query parameter to filter by extId to get the summary for a specific project.

## Example Usage

```hcl
data "nutanix_role_membership_summary_v2" "example" {
  ext_id = "project_uuid"
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) A URL query parameter that specifies the page number of the result set.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
    * extId
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
    * extId
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.
    * extId
    * groupsCount
    * rolesCount
    * totalIdentitiesCount
    * usersCount
## Attributes Reference

The following attributes are exported:

* `summaries` - List of role membership summaries.

### Summary

Each summary in `summaries` exports the following:

* `ext_id` - External identifier of the role membership summary. This attribute is used to fetch the summary of a specific project.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links` - A HATEOAS style link for the response.
* `users_count` - Count of distinct users in the project.
* `groups_count` - Count of distinct groups in the project.
* `roles_count` - Count of distinct roles assigned in the project.
* `total_identities_count` - Total count of distinct identities (users and groups) in the project.
