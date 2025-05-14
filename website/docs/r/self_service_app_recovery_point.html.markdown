---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_app_recovery_point"
sidebar_current: "docs-nutanix_self_service_app_recovery_point"
description: |-
  Run snapshot action in application to create recovery point.
---

# nutanix_self_service_app_recovery_point

Run snapshot action in application to create recovery point.

## Example Usage

```hcl
resource "nutanix_self_service_app_recovery_point" "test" {
    app_name = "NAME OF APPLICATION"
    action_name = "SNAPSHOT ACTION NAME"
    recovery_point_name = "RECOVERY POINT NAME"
}
```

## Argument Reference

The following arguments are supported:

* `app_name`: - (Optional) The name of the application
* `app_uuid`: - (Required) The UUID of the application.

Both (`app_name` and `app_uuid`) are optional but atleast one of them should be provided for resource to work.

## Attribute Reference

* `action_name`: - (Required) The name of the snapshot action to trigger.
* `recovery_point_name`: - (Required) The name of recovery point.


See detailed information in [Run snapshot action in app](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Apps/paths/~1apps~1%7Buuid%7D~1actions~1%7Baction_uuid%7D~1run/post).