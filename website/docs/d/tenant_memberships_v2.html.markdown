---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_memberships_v2"
sidebar_current: "docs-nutanix-datasource-tenant-memberships-v2"
description: |-
  Lists all tenant memberships in a flat structure across all domains.
---

# nutanix_tenant_memberships_v2

Lists all tenant memberships in a flat structure across all domains.

## Example

```hcl
data "nutanix_tenant_memberships_v2" "example" {}

data "nutanix_tenant_memberships_v2" "filtered" {
  limit  = 10
  filter = "domainExtId eq '00000000-0000-0000-0000-000000000000'"
}
```

## Argument Reference

The following arguments are supported:

* `page`:- (Optional) A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource.
* `limit`:- (Optional) A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100.
* `filter`:- (Optional) A URL query parameter that allows clients to filter a collection of resources, conforming to the OData V4.01 URL conventions.
* `order_by`:- (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.

## Attribute Reference

The following attributes are exported:

* `memberships`:- List of memberships for the tenant across all domains.
* `total_available_results`:- The total number of entities that are available on the server for this type.

### Memberships

Each entry in the `memberships` list exports the following:

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

See detailed information in [Nutanix List Domain Memberships v4](https://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a1#tag/Membership/operation/listDomainMemberships).
