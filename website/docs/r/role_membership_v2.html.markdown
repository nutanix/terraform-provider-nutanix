---
layout: "nutanix"
page_title: "NUTANIX: nutanix_role_membership_v2"
sidebar_current: "docs-nutanix-resource-role-membership-v2"
description: |-
  Creates a role membership.
---

# nutanix_role_membership_v2

Provides a resource to create a role membership.

A role membership assigns a role to an identity (user or group) within a specific scope (project).

## Example Usage

### Role Membership with Project Scope

```hcl
resource "nutanix_role_membership_v2" "example" {
  role_ext_id         = "ca386756-e45f-5555-8625-5b68ae17393b" # Role uuid
  identity_type       = "USER"
  identity_ext_id     = "8a49f561-6bd7-5d26-b53e-661d63e7bdb8" # User uuid
  idp_ext_id          = "d711f713-cdf5-5ee9-9936-4e67373eb842" # Identity provider uuid
  scope_template_name = "ProjectsScopeTemplate"
  project_ext_id      = "78198f9c-063d-590e-9e9a-939b51829a39" # Project uuid

  scope_template_name_values {
    name  = "projectExtId"
    value = "78198f9c-063d-590e-9e9a-939b51829a39" # Project uuid
  }
}
```

### Create Role Membership for a Group

```hcl
resource "nutanix_role_membership_v2" "group" {
  role_ext_id         = "ca386756-e45f-5555-8625-5b68ae17393b" # Role uuid
  identity_type       = "GROUP"
  identity_ext_id     = "8a49f561-6bd7-5d26-b53e-661d63e7bdb8" # UserGroup uuid
  idp_ext_id          = "d711f713-cdf5-5ee9-9936-4e67373eb842" # Identity provider uuid
  scope_template_name = "ProjectsScopeTemplate"
  project_ext_id      = "78198f9c-063d-590e-9e9a-939b51829a39" # Project uuid

  scope_template_name_values {
    name  = "projectExtId"
    value = "78198f9c-063d-590e-9e9a-939b51829a39" # Project uuid
  }
}
```

## Argument Reference

The following arguments are supported:

* `role_ext_id` - (Required) External identifier of the role associated with the role membership.
* `identity_type` - (Required) Type of identity. Valid values are `USER`, `GROUP`.
* `identity_ext_id` - (Required) External identifier of the identity (user or group) associated with the role membership.
* `idp_ext_id` - (Required) External Identifier of the identity provider associated with the role membership.
* `scope_template_name` - (Required) Name of the scope template.
* `scope_template_name_values` - (Optional) Name Value pairs to substitute in the scope template variables referenced by the role membership.
  * `name` - (Optional) The name of the variable.
  * `value` - (Optional) The value to substitute.
* `project_ext_id` - (Optional) External identifier of the project associated with the role membership.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id` - External identifier of the role membership.
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

## Import

Role memberships can be imported using the `ext_id`:

```hcl
// create its configuration in the root module. For example:
resource "nutanix_role_membership_v2" "import_role_membership" {}

terraform import nutanix_role_membership_v2.import_role_membership <ext_id>
```

See detailed information in [Nutanix Role Membership V4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.1.b3#tag/RoleMembership/operation/createRoleMembership).