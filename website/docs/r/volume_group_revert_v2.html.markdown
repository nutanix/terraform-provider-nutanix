---
layout: "nutanix"
page_title: "NUTANIX: nutanix_volume_group_revert_v2"
sidebar_current: "docs-nutanix-resource-volume-group-revert-v2"
description: |-
  Reverts a Volume Group from a recovery point.
---

# nutanix_volume_group_revert_v2

Reverts a Volume Group identified by Volume Group external identifier. This API performs an in-place restore from a specified Volume Group recovery point.

## Example Usage

```hcl
resource "nutanix_volume_group_revert_v2" "example" {
  ext_id                              = "d09aeec9-5bb7-4bfd-9717-a051178f6e7c"
  volume_group_recovery_point_ext_id  = "a12bcedf-1234-5678-9012-abcdef123456"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a Volume Group.
* `volume_group_recovery_point_ext_id`: -(Required) The external identifier of the Volume Group recovery point. This is a mandatory field.
