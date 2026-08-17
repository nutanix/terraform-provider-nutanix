---
layout: "nutanix"
page_title: "NUTANIX: nutanix_scheduled_delete_tenant_v2"
sidebar_current: "docs-nutanix-resource-scheduled-delete-tenant-v2"
description: |-
  Schedule deletion of a tenant by its external identifier after a configurable retention period.
---

# nutanix_scheduled_delete_tenant_v2

Schedule deletion of a tenant by its external identifier after a configurable retention period.

## Example

```hcl
resource "nutanix_scheduled_delete_tenant_v2" "example" {
  ext_id            = "00000000-0000-0000-0000-000000000000"
  retention_seconds = 3600
}
```

## Argument Reference

The following arguments are supported:

- `ext_id`: (Required) External identifier of the tenant.
- `retention_seconds`: (Required) The number of seconds to retain the tenant before it is permanently deleted. A value of 0 deletes the tenant immediately.

See detailed information in [Nutanix Scheduled Delete Tenant V4](https://developers.nutanix.com/).
