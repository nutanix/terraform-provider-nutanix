---
layout: "nutanix"
page_title: "NUTANIX: nutanix_suspend_tenant_v2"
sidebar_current: "docs-nutanix-resource-suspend-tenant-v2"
description: |-
  Suspend a tenant by its external identifier, temporarily blocking all of its users without losing data.
---

# nutanix_suspend_tenant_v2

Suspend a tenant by its external identifier, temporarily blocking all of its users without losing data.

## Example

```hcl
resource "nutanix_suspend_tenant_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) External identifier of the tenant.

See detailed information in [Nutanix Suspend Tenant V4](https://developers.nutanix.com/).
