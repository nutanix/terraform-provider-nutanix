---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_perform_inventory_v2"
sidebar_current: "docs-nutanix_lcm_perform_inventory_v2"
description: |-
  Perform an inventory operation to identify/scan entities on the cluster that can be updated through LCM.
---

# nutanix_lcm_perform_inventory_v2

Perform an inventory operation to identify/scan entities on the cluster that can be updated through LCM.

## Example

```hcl

# perform inventory
resource "nutanix_lcm_perform_inventory_v2" "inventory" {
  x_cluster_id = "0005a104-0b0b-4b0b-8005-0b0b0b0b0b0b"
}

# perform software-only inventory
resource "nutanix_lcm_perform_inventory_v2" "software_inventory" {
  x_cluster_id   = "0005a104-0b0b-4b0b-8005-0b0b0b0b0b0b"
  inventory_type = "SOFTWARE"
}

# perform node-specific inventory
resource "nutanix_lcm_perform_inventory_v2" "node_inventory" {
  x_cluster_id   = "0005a104-0b0b-4b0b-8005-0b0b0b0b0b0b"
  inventory_type = "NODE"
  node_list      = ["node-uuid-1", "node-uuid-2"]
}
```

## Argument Reference
The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.
* `inventory_type`: (Optional) The scope of the inventory scan to perform. When omitted, a full inventory is performed.

    | Enum     | Description                                                                                                                                            |
    |----------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
    | FULL     | Full inventory. Scans all nodes in the cluster for both software and firmware entities and checks the LCM repository for available updates.             |
    | SOFTWARE | Software-only inventory. Scans only software entities (e.g. AOS, NCC, LCM framework) without scanning firmware components.                             |
    | NODE     | Node-scoped inventory. Scans only the specific nodes listed in `node_list`. Requires `node_list` to be provided.                                       |
    | RESCAN   | Rescan inventory. Re-checks the LCM repository for available updates without performing a full hardware scan of cluster nodes. Faster than FULL.        |

* `node_list`: (Optional) List of node UUIDs to inventory. Required when `inventory_type` is set to `NODE`; ignored for other inventory types.

See detailed information in [Nutanix LCM Perform Inventory v4](https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.3#tag/Inventory/operation/performInventory)
