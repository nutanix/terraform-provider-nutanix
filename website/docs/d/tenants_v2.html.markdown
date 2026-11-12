---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenants_v2"
sidebar_current: "docs-nutanix-datasource-tenants-v2"
description: |-
  Lists all tenants with pagination support.
---

# nutanix_tenants_v2

Lists all tenants with pagination support.

## Example

```hcl
data "nutanix_tenants_v2" "example" {}
```

## Argument Reference

No arguments are required.

## Attributes Reference

The following attributes are exported:

- `tenants`: List of tenants. Each element exports the same attributes as the `nutanix_tenant_v2` data source.

See detailed information in [Nutanix List Tenants V4](https://developers.nutanix.com/api-reference?namespace=tenancy&version=v4.0.a1#tag/Tenants/operation/listTenants).
