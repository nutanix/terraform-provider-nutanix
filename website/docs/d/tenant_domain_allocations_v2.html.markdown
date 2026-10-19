---
layout: "nutanix"
page_title: "NUTANIX: nutanix_tenant_domain_allocations_v2"
sidebar_current: "docs-nutanix-datasource-tenant-domain-allocations-v2"
description: |-
  Lists all domain allocations for a specific tenant.
---

# nutanix_tenant_domain_allocations_v2

Lists all domain allocations for a specific tenant.

## Example

```hcl
data "nutanix_tenant_domain_allocations_v2" "example" {
  tenant_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `tenant_ext_id`: (Required) External identifier of the tenant.

## Attributes Reference

The following attributes are exported:

- `domain_allocations`: List of domain allocations. Each element exports the same attributes as the `nutanix_tenant_domain_allocation_v2` data source.

See detailed information in [Nutanix List Tenant Domain Allocations V4](https://developers.nutanix.com/).
