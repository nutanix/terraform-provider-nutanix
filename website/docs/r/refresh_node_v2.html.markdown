---
layout: "nutanix"
page_title: "NUTANIX: nutanix_refresh_node_v2"
sidebar_current: "docs-nutanix-resource-refresh-node-v2"
description: |-
  Refreshes a node's information by fetching the latest data from the hardware provider.
---

# nutanix_refresh_node_v2

Refreshes a node's information by fetching the latest data from the hardware provider.

## Example

```hcl
resource "nutanix_refresh_node_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `ext_id` - (Required) External ID of the node to refresh.

See detailed information in [Nutanix Refresh Node V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
