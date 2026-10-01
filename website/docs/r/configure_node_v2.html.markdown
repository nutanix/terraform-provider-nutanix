---
layout: "nutanix"
page_title: "NUTANIX: nutanix_configure_node_v2"
sidebar_current: "docs-nutanix-resource-configure-node-v2"
description: |-
  Configures the node in the hardware provider with the given settings.
---

# nutanix_configure_node_v2

Configures the node in the hardware provider with the given settings. Exactly one configuration block must be supplied.

## Example

```hcl
resource "nutanix_configure_node_v2" "example" {
  configure_server {
    node_ext_id = "00000000-0000-0000-0000-000000000000"

    server_settings_config {
      network_adaptor_fec_mode = "OFF"
    }
  }
}
```

## Argument Reference

Exactly one of the following configuration blocks must be provided:

* `configure_server` - (Optional) Configure a single server in the hardware provider. See [Configure Server](#configure-server) below.
* `un_configure_server` - (Optional) Remove the configuration from a node. Contains `node_ext_id` (Required).
* `pre_cluster_config` - (Optional) Pre-cluster operations on a set of nodes. Contains `nodes` (Required, list of `ext_id`).
* `post_cluster_config` - (Optional) Post-cluster operations on a set of nodes. Contains `cluster_ext_id` (Required) and `nodes` (Required, list of `ext_id`).

### Configure Server

* `node_ext_id` - (Required) External ID of the node.
* `group_id` - (Optional) Identifier of the group (e.g., organization or tenant) to which the node belongs. Required for Cisco nodes.
* `server_settings_config` - (Optional) Server settings configuration. See [Server Settings Config](#server-settings-config) below.

### Server Settings Config

* `network_adaptor_fec_mode` - (Optional) Forward error correction mode. One of `CL91`, `CL74`, `OFF`.
* `server_identity_pool_ext_id` - (Optional) UUID of the server identity pool to be used for server configuration.

## Attributes Reference

* `ext_id` - Identifier of the configure task tracking the asynchronous operation.

See detailed information in [Nutanix Configure Node V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
