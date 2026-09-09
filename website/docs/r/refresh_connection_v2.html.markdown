---
layout: "nutanix"
page_title: "NUTANIX: nutanix_refresh_connection_v2"
sidebar_current: "docs-nutanix-resource-refresh-connection-v2"
description: |-
  Refreshes nodes or resources from a hardware provider connection.
---

# nutanix_refresh_connection_v2

Refreshes nodes or resources (IP address pools, MAC address pools, and server identity pools) from a hardware provider connection. This is a one-shot, state-changing action resource.

## Example

```hcl
resource "nutanix_refresh_connection_v2" "example" {
  hardware_provider_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id                   = "11111111-1111-1111-1111-111111111111"

  refresh_resources_spec {
    should_refresh_ip_pools              = true
    should_refresh_mac_pools             = true
    should_refresh_server_identity_pools = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `hardware_provider_ext_id`: (Required) External ID of the hardware provider.
* `ext_id`: (Required) External ID of the hardware provider connection to refresh.
* `refresh_resources_spec`: (Optional) Refresh resource pools from the connection. See [Refresh Resources Spec](#refresh-resources-spec) below.
* `discover_nodes_spec`: (Optional) Discover nodes from the connection. See [Discover Nodes Spec](#discover-nodes-spec) below.

### Refresh Resources Spec

* `group_ids`: (Optional) List of group identifiers to scope the refresh to.
* `should_refresh_ip_pools`: (Optional) Whether to refresh IP address pools.
* `should_refresh_mac_pools`: (Optional) Whether to refresh MAC address pools.
* `should_refresh_server_identity_pools`: (Optional) Whether to refresh server identity pools.

### Discover Nodes Spec

* `group_ids`: (Optional) List of group identifiers to scope node discovery to.

See detailed information in [Nutanix Hardware Providers V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
