---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_app_custom_action"
sidebar_current: "docs-nutanix_self_service_app"
description: |-
  Triggers custom action execution using it's name in Self Service Application.
---

# nutanix_self_service_app_custom_action

Triggers custom action execution using it's name in Self Service Application.

## Example Usage:

``` hcl
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = "NAME OF BLUEPRINT"
    app_name        = "NAME OF APPLICATION"
    app_description = "DESCRIPTION OF APPLICATION"
}

resource "nutanix_self_service_app_custom_action" "test" {
    app_name        = nutanix_self_service_app_provision.test.app_name
    action_name = "NAME OF ACTION"
}
```

## Argument Reference

The following arguments are supported:

* `app_name`: - (Optional) The name of the application.
* `action_name`: - (Required) The name of the action to run.

## Attribute Reference

The following attributes are exported:

* `app_uuid`: - (Optional) The UUID of the application.
* `runlog_uuid`: - (Computed) The UUID of the runlog associated with the execution of the custom action. This can be used to track the progress or status of the action execution.

See detailed information in [Run action in app](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/Apps/paths/~1apps~1%7Buuid%7D~1actions~1run/post).