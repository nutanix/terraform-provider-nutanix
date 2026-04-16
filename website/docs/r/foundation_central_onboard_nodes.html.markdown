---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_onboard_nodes"
sidebar_current: "docs-nutanix-resource-foundation-central-onboard-nodes"
description: |-
  Onboard nodes from a hardware manager (such as Cisco Intersight) into Foundation Central.
---

# nutanix_foundation_central_onboard_nodes

Onboard nodes from a hardware manager (such as Cisco Intersight) into Foundation Central.

## Example Usage

```hcl
resource "nutanix_foundation_central_onboard_nodes" "node" {
  node_serial = "ABC12345D6E"
}
```


## Argument Reference

The following arguments are supported:

* `node_serial`: Serial number of the node to onboard

## Attributes Reference

The following attributes are exported:

* `block_serial`: Block serial number of the node
* `imaged_node_uuid`: UUID of the imaged_node in FC
* `model`: Model of the node
* `node_state`: State of the node (e.g. STATE_ONBOARDED)
* `node_type`: Type of node (e.g. "on-prem")


## Error 

Incase of error (such as node serial not found), an error will be generated. 

## lifecycle

* `Create` : - Resource will trigger onboarding of the node.
* `delete` : - Node will be removed from FC "imaged nodes" but will still be available in the hardware manager for re-onboarding. 
