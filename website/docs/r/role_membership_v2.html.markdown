---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_membership_v2"
sidebar_current: "docs-nutanix-resource-role-membership-v2"
description: |-
  Creates a role membership in Nutanix.
---

# nutanix_role_membership_v2

Provides a resource to create a role membership in Nutanix.

A role membership assigns a role to an identity (user or group) within a specific scope (project).

## Example Usage

### Basic Role Membership

```hcl
resource "nutanix_role_membership_v2" "example" {
  role_ext_id      = "role-ext-id-here"
  identity_type    = "USER"
  identity_ext_id  = "user-ext-id-here"
  idp_ext_id       = "idp-ext-id-here"
}
```

### Role Membership with Project Scope

```hcl
resource "nutanix_role_membership_v2" "with_project" {
  role_ext_id         = "role-ext-id-here"
  identity_type       = "USER"
  identity_ext_id     = "user-ext-id-here"
  idp_ext_id          = "idp-ext-id-here"
  scope_template_name = "project-scope"
  project_ext_id      = "project-ext-id-here"

  key_value_pairs {
    key   = "projectId"
    value = "project-ext-id-here"
  }

  scope_template_name_values {
    name  = "projectId"
    value = "project-ext-id-here"
  }
}
```

### Role Membership for a Group

```hcl
resource "nutanix_role_membership_v2" "group" {
  role_ext_id      = "role-ext-id-here"
  identity_type    = "GROUP"
  identity_ext_id  = "group-ext-id-here"
  idp_ext_id       = "idp-ext-id-here"
}
```

## Argument Reference

The following arguments are supported:

* `role_ext_id` - (Required) External identifier of the role.
* `identity_type` - (Required) Type of identity. Valid values are `USER`, `GROUP`.
* `identity_ext_id` - (Optional) External identifier of the identity (user or group) associated with the role membership.
* `identity_value` - (Optional) Value of the identity.
* `idp_ext_id` - (Optional) External identifier of the identity provider.
* `scope_template_name` - (Optional) Name of the scope template.
* `scope_template_name_values` - (Optional) Name value pairs to substitute in the scope template variables.
  * `name` - (Optional) The name of the variable.
  * `value` - (Optional) The value to substitute.
* `project_ext_id` - (Optional) External identifier of the project.
* `key_value_pairs` - (Optional) Key-value pairs for the role membership.
  * `key` - (Required) The key.
  * `value` - (Required) The value.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id` - External identifier of the role membership.
* `tenant_id` - Tenant identifier.
* `links` - A HATEOAS style link for the response.
* `authorization_policy_ext_id` - External identifier of the authorization policy.
* `created_by` - User or service name that created the role membership.
* `created_time` - The creation time of the role membership.
* `last_updated_time` - The time when the role membership was last updated.

## Import

Role memberships can be imported using the `ext_id`:

```shell
terraform import nutanix_role_membership_v2.example <ext_id>
```
