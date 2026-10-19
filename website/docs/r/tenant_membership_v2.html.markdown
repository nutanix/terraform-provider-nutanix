---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_membership_v2"
sidebar_current: "docs-nutanix-resource-tenant-membership-v2"
description: |-
  Create a tenant membership for a user within a single domain.
---

# nutanix_tenant_membership_v2

Creates a tenant membership for a user within a single domain. This operation is asynchronous and returns a membership reference.

## Example

### Membership by existing identity ext_id

```hcl
resource "nutanix_tenant_membership_v2" "example" {
  domain_ext_id   = "00000000-0000-0000-0000-000000000000"
  role_ext_id     = "00000000-0000-0000-0000-000000000001"
  identity_ext_id = "00000000-0000-0000-0000-000000000002"
}
```

### Membership by external user details

```hcl
resource "nutanix_tenant_membership_v2" "example" {
  domain_ext_id = "00000000-0000-0000-0000-000000000000"
  role_ext_id   = "00000000-0000-0000-0000-000000000001"

  identity_details {
    identity_provider_ext_id = "00000000-0000-0000-0000-000000000003"
    principal_value          = "user@example.com"
    type                     = "USER"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_ext_id`:- (Required) External identifier of the domain.
* `role_ext_id`:- (Required) The external identifier of the role assigned to the principal.
* `identity_ext_id`:- (Optional) The external identifier of an existing principal entity. Provide this OR `identity_details`.
* `identity_details`:- (Optional) Details used to resolve an external (AD/LDAP/SAML) principal. Provide this OR `identity_ext_id`.

~> **Note**: `domain_ext_id`, `role_ext_id`, `identity_ext_id`, and `identity_details` are immutable; changing any of them forces the resource to be recreated.

### Identity Details

The `identity_details` block supports the following:

* `identity_provider_ext_id`:- (Required) The external identifier of the identity provider.
* `principal_value`:- (Required) The unique lookup value for the user, such as a UPN, username, or email address.
* `type`:- (Required) The identity type. Valid values: `"USER"`, `"GROUP"`.

## Attribute Reference

The following attributes are exported:

* `ext_id`:- (Computed) A globally unique identifier of an instance that is suitable for external consumption.
* `created_by`:- The username of the user who created the membership.
* `created_time`:- The date and time when the membership was created.
* `principal`:- The resolved principal name associated with this membership.
* `tenant_id`:- A globally unique identifier that represents the tenant that owns this entity.
* `links`:- A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
  * `href`:- The URL at which the entity described by the link can be accessed.
  * `rel`:- A name that identifies the relationship of the link to the object that is returned by the URL.
* `state_info`:- The current lifecycle state of the membership.
  * `state`:- The membership state. Valid values: `"PENDING"`, `"ACTIVE"`, `"ERROR"`, `"DELETING"`.

## Import

This helps to manage existing entities which are not created through Terraform. Tenant membership resources can be imported using the `UUID` (`ext_id` in v4 terms). For example:

```hcl
// create its configuration in the root module. For example:
resource "nutanix_tenant_membership_v2" "import_membership" {}

// execute this command in cli
terraform import nutanix_tenant_membership_v2.import_membership <UUID>
```

See detailed information in [Nutanix Create Domain Membership v4](https://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a2).
