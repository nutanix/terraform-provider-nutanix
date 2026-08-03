---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_membership_v2"
sidebar_current: "docs-nutanix-datasource-role-membership-v2"
description: |-
  Get role membership based on the provided external identifier.
---

# nutanix_role_membership_v2

Displays a role membership based on the provided external identifier.

## Example Usage

```hcl
data "nutanix_role_membership_v2" "example" {
  ext_id = "<role-membership-ext-id>"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) External identifier of the role membership.

## Attributes Reference

The following attributes are exported:

* `tenant_id` - Tenant identifier.
* `links` - A HATEOAS style link for the response.
* `authorization_policy_ext_id` - External identifier of the authorization policy.
* `created_by` - User or Service who created the role membership.
* `created_time` - Creation time of the role membership..
* `last_updated_time` - Last updated time of the role membership.
* `role_ext_id` - External identifier of the role associated with the role membership.
* `identity_type` - Type of identity. Valid values are `USER`, `GROUP`.
* `identity_ext_id` - External identifier of the identity (user or group) associated with the role membership.
* `idp_ext_id` - External Identifier of the identity provider associated with the role membership.
* `scope_template_name` - Name of the scope template.
* `scope_template_name_values` - Name Value pairs to substitute in the scope template variables referenced by the role membership.
  * `name` - The name of the variable.
  * `value` - The value to substitute.
* `project_ext_id` - External identifier of the project associated with the role membership.