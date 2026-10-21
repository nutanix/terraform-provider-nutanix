---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_v2"
sidebar_current: "docs-nutanix-resource-tenant-v2"
description: |-
  Creates a new tenant.
---

# nutanix_tenant_v2

Creates a new tenant.

> **WriteOnly credentials:** `admin_info.password` is write-only. It is never stored in Terraform state. Re-supply it in configuration on every apply. Do not use `lifecycle.ignore_changes` on this attribute. See the [WriteOnly Credentials](/docs/providers/nutanix/guides/writeonly_credentials.html) guide for upgrade caveats.

## Example

```hcl
resource "nutanix_tenant_v2" "example" {
  name        = "example-tenant"
  description = "Tenant created via Terraform"
  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = "Nutanix.123"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name`: (Optional) The name of the tenant. It must start and end with a letter (including Unicode letters with combining marks) or a digit, and can contain hyphens in between.
- `description`: (Optional) A brief description of the tenant.
- `url`: (Optional) The URL of the tenant.
- `added_by`: (Optional) The identifier of the user who created the tenant.
- `domains`: (Optional) List of external identifiers of domains associated with the tenant.
- `admin_info`: (Optional) The local admin details of the tenant.
- `aggregated_quota`: (Optional) The aggregated resource quota for the tenant.

### admin_info

- `username`: (Optional) The username of the user.
- `first_name`: (Optional) The first name of the user.
- `last_name`: (Optional) The last name of the user.
- `email`: (Optional) The email address of the user.
- `password`: (Optional, WriteOnly) The password of the user.

### aggregated_quota

- `ext_id`: (Optional) The external identifier of the quota for the placement target.
- `memory`: (Optional) The amount of memory allocated as quota, in bytes.
- `storage`: (Optional) The amount of storage allocated as quota, in bytes.
- `vcpu`: (Optional) The number of vCPUs allocated as quota.

## Attributes Reference

The following attributes are exported:

- `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
- `created_time`: The timestamp when the tenant was created.
- `state_info`: The state details of the tenant.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.

### state_info

- `state`: The state of the tenant.

See detailed information in [Nutanix Create Tenant V4](https://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a1#tag/Tenants/operation/createTenant).
