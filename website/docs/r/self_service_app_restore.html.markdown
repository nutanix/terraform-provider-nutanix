---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_app_restore"
sidebar_current: "docs-nutanix_self_service_app_restore"
description: |-
  Run restore action in application to restore from recovery point.
---

# nutanix_self_service_app_restore

Run restore action in application to restore from recovery point.

## Example Usage

```hcl
resource "nutanix_self_service_app_recovery_point" "test" {
    app_name = "NAME OF APPLICATION"
    action_name = "SNAPSHOT ACTION NAME"
    recovery_point_name = "RECOVERY POINT NAME"
}

# Read available recovery points in app
data "nutanix_self_service_app_snapshots" "snapshots" {
    app_name = "NAME OF APPLICATION"
    length = 250
    offset = 0
    depends_on = [nutanix_self_service_app_recovery_point.test]
}

locals {
    snapshot_uuid = [
    for snapshot in data.nutanix_self_service_app_snapshots.snapshots.entities :
    snapshot.uuid if snapshot.name == "SNAPSHOT ACTION NAME"
    ][0]
}

# Restore from recovery point
resource "nutanix_self_service_app_restore" "test" {
    restore_action_name = "RESTORE ACTION NAME"
    app_name =  "NAME OF APPLICATION"
    snapshot_uuid = local.snapshot_uuid
}
```

## Argument Reference

The following arguments are supported:

* `app_name`: - (Optional) The name of the application
* `app_uuid`: - (Optional) The UUID of the application.
* `snapshot_uuid`: - (Required) The UUID of the snapshot to which the application will be restored.

Both (`app_name` and `app_uuid`) are optional but atleast one of them should be provided for resource to work.

## Attribute Reference

* `restore_action_name`: - (Required) The name of the restore action to be performed.
* `state`: - (Computed) This will be set after the restore action has been processed.


See detailed information in [Run restore action in app](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Apps/paths/~1apps~1%7Buuid%7D~1actions~1%7Baction_uuid%7D~1run/post).