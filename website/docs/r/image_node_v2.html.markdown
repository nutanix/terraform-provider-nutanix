---
layout: "nutanix"
page_title: "NUTANIX: nutanix_image_node_v2"
sidebar_current: "docs-nutanix-resource-image-node-v2"
description: |-
  Installs AOS or hypervisor or host OS on the node.
---

# nutanix_image_node_v2

Installs AOS or hypervisor or host OS on the node. Exactly one configuration block must be supplied.

## Example

```hcl
resource "nutanix_image_node_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"

  host_installation {
    patched_image_ext_id = "33333333-3333-3333-3333-333333333333"
  }
}
```

## Argument Reference

* `ext_id` - (Required) External ID of the node.

Exactly one of the following configuration blocks must be provided:

* `host_installation` - (Optional) Install a patched hypervisor or host OS image. See [Host Installation](#host-installation) below.
* `aos_installation` - (Optional) Install AOS on the node. See [AOS Installation](#aos-installation) below.

### Host Installation

* `patched_image_ext_id` - (Required) External ID of a patched hypervisor or host OS image.
* `patched_image_url` - (Optional) URL of the patched image to be installed on the node. If the URL is provided, nodes download the image from the given URL instead of downloading it from Foundation Central.

### AOS Installation

* `aos_image_ext_id` - (Required) External ID of the AOS image managed by Foundation Central.
* `cvm_memory_gb` - (Optional) Memory allocation for the CVM in GB.
* `management_network` - (Optional) Management network configuration. Contains `gateway`, `ip` (each with `ipv4` / `ipv6` `value` and `prefix_length`) and `vlan_id`.

See detailed information in [Nutanix Image Node V4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3).
