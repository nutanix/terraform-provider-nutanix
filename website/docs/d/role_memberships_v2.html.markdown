---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_memberships_v2"
sidebar_current: "docs-nutanix-datasource-role-memberships-v2"
description: |-
  Lists role memberships in Nutanix.
---

# nutanix_role_memberships_v2

Lists role memberships in Nutanix.

## Example Usage

```hcl
data "nutanix_role_memberships_v2" "example" {}
```

### With Filters

```hcl
data "nutanix_role_memberships_v2" "filtered" {
  page   = 0
  limit  = 10
  filter = "identityType eq 'USER'"
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) A URL query parameter that specifies the page number of the result set.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `expand` - (Optional) A URL query parameter that allows clients to request related resources.
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attributes Reference

The following attributes are exported:

* `role_memberships` - List of role memberships.

### Role Membership

Each role membership in `role_memberships` exports the following:

* `ext_id` - External identifier of the role membership.
* `tenant_id` - Tenant identifier.
* `links` - A HATEOAS style link for the response.
* `authorization_policy_ext_id` - External identifier of the authorization policy.
* `role_ext_id` - External identifier of the role.
* `identity_ext_id` - External identifier of the identity (user or group).
* `identity_type` - Type of identity. Valid values are `USER`, `GROUP`.
* `identity_value` - Value of the identity.
* `idp_ext_id` - External identifier of the identity provider.
* `scope_template_name` - Name of the scope template.
* `scope_template_name_values` - Name value pairs for the scope template.
  * `name` - The name of the variable.
  * `value` - The value.
* `project_ext_id` - External identifier of the project.
* `key_value_pairs` - Key-value pairs for the role membership.
  * `key` - The key.
  * `value` - The value.
* `created_by` - User or service name that created the role membership.
* `created_time` - The creation time of the role membership.
* `last_updated_time` - The time when the role membership was last updated.
