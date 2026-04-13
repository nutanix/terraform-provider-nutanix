---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_membership_v2"
sidebar_current: "docs-nutanix-datasource-role-membership-v2"
description: |-
  Describes a role membership in Nutanix.
---

# nutanix_role_membership_v2

Describes a role membership in Nutanix. A role membership assigns a role to an identity (user or group).

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
