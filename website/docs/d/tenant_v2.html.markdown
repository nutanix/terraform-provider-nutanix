---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_v2"
sidebar_current: "docs-nutanix-datasource-tenant-v2"
description: |-
  Retrieves a specific tenant by its external identifier.
---

# nutanix_tenant_v2

Retrieves a specific tenant by its external identifier.

## Example

```hcl
data "nutanix_tenant_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) External identifier of the tenant.

## Attributes Reference

The following attributes are exported:

- `name`: The name of the tenant.
- `description`: A brief description of the tenant.
- `url`: The URL of the tenant.
- `added_by`: The identifier of the user who created the tenant.
- `domains`: List of external identifiers of domains associated with the tenant.
- `admin_info`: The local admin details of the tenant.
- `aggregated_quota`: The aggregated resource quota for the tenant.
- `created_time`: The timestamp when the tenant was created.
- `state_info`: The state details of the tenant.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `links`: A HATEOAS style link for the response.

### admin_info

- `username`: The username of the user.
- `first_name`: The first name of the user.
- `last_name`: The last name of the user.
- `email`: The email address of the user.

### aggregated_quota

- `ext_id`: The external identifier of the quota for the placement target.
- `memory`: The amount of memory allocated as quota, in bytes.
- `storage`: The amount of storage allocated as quota, in bytes.
- `vcpu`: The number of vCPUs allocated as quota.

### state_info

- `state`: The state of the tenant.

See detailed information in [Nutanix Get Tenant V4](https://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a1#tag/Tenants/operation/getTenantById).
