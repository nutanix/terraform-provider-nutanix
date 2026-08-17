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
    * authorizationPolicyExtId
    * createdBy
    * extId
    * identityExtId
    * identityType
    * idpExtId
    * projectExtId
    * roleExtId
    * scopeTemplateName
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
    * authorizationPolicyExtId
    * createdTime
    * extId
    * identityExtId
    * identityType
    * idpExtId
    * lastUpdatedTime
    * projectExtId
    * roleExtId
* `expand` - (Optional) A URL query parameter that allows clients to request related resources.
    * roleExtId
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.
    * authorizationPolicyExtId
    * createdBy
    * createdTime
    * extId
    * identityExtId
    * identityType
    * idpExtId
    * lastUpdatedTime
    * projectExtId
    * roleExtId
    * scopeTemplateName

## Attributes Reference

The following attributes are exported:

* `role_memberships` - List of role memberships.

### Role Membership

Each role membership in `role_memberships` exports the following:

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

See detailed information in [Nutanix List Role Memberships V2](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.1.b3#tag/RoleMembership/operation/listRoleMemberships).