---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_membership_v2"
sidebar_current: "docs-nutanix-datasource-tenant-membership-v2"
description: |-
  Retrieves a specific tenant membership by its identifier.
---

# nutanix_tenant_membership_v2

Retrieves a specific tenant membership by its identifier.

## Example

```hcl
data "nutanix_tenant_membership_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`:- (Required) External identifier of the membership.

## Attribute Reference

The following attributes are exported:

* `ext_id`:- A globally unique identifier of an instance that is suitable for external consumption.
* `domain_ext_id`:- External identifier of the domain.
* `role_ext_id`:- The external identifier of the role assigned to the principal.
* `identity_ext_id`:- The external identifier of an existing principal entity.
* `identity_details`:- Details used to resolve an external principal.
  * `identity_provider_ext_id`:- The external identifier of the identity provider.
  * `principal_value`:- The unique lookup value for the user, such as a UPN, username, or email address.
  * `type`:- The identity type. Valid values: `"USER"`, `"GROUP"`.
* `created_by`:- The username of the user who created the membership.
* `created_time`:- The date and time when the membership was created.
* `principal`:- The resolved principal name associated with this membership.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `links`:- A HATEOAS style link for the response.
  * `href`:- The URL at which the entity described by the link can be accessed.
  * `rel`:- A name that identifies the relationship of the link to the object that is returned by the URL.
* `state_info`:- The current lifecycle state of the membership.
  * `state`:- The membership state. Valid values: `"PENDING"`, `"ACTIVE"`, `"ERROR"`, `"DELETING"`.

See detailed information in [Nutanix Get Domain Membership v4](https://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a1#tag/Membership/operation/getDomainMembershipById).
