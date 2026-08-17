---
layout: "nutanix"
page_title: "NUTANIX: nutanix_restore_tenant_v2"
sidebar_current: "docs-nutanix-resource-restore-tenant-v2"
description: |-
  Restore a suspended or pending-deletion tenant by its external identifier back to the active state.
---

# nutanix_restore_tenant_v2

Restore a suspended or pending-deletion tenant by its external identifier back to the active state.

## Example

```hcl
resource "nutanix_restore_tenant_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) External identifier of the tenant.

See detailed information in [Nutanix Restore Tenant V4](https://developers.nutanix.com/).
