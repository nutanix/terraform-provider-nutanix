---
layout: "nutanix"
page_title: "NUTANIX: nutanix_manage_alert_v2"
sidebar_current: "docs-nutanix-resource-manage-alert-v2"
description: |-
  Acknowledges or resolves the alert identified by external identifier.
---

# nutanix_manage_alert_v2

Acknowledges or resolves the alert identified by external identifier. This is an action resource — it triggers the manage-alert action on creation and is removed from state on destroy.

## Example Usage

```hcl
resource "nutanix_manage_alert_v2" "acknowledge" {
  ext_id      = "00000000-0000-0000-0000-000000000000"
  action_type = "ACKNOWLEDGE"
}

resource "nutanix_manage_alert_v2" "resolve" {
  ext_id      = "00000000-0000-0000-0000-000000000000"
  action_type = "RESOLVE"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: - (Required) Unique identifier of an alert that can be resolved or acknowledged.
* `action_type`: - (Required) The action to perform on the alert. Valid values are `ACKNOWLEDGE` and `RESOLVE`.

## Attribute Reference

The following attributes are exported:

* `task_ext_id`: A globally unique identifier for the task.

See detailed information in [Nutanix Monitoring v4 Manage Alerts](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.0).
