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
# List Prism Central
data "nutanix_clusters_v2" "pc" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}
locals {
  pcExtID      = data.nutanix_clusters_v2.pc.cluster_entities[0].ext_id
}

# perform inventory
resource "nutanix_lcm_perform_inventory_v2" "inventory" {
  x_cluster_id = local.pcExtID
}
```

# Argument Reference
The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.

See detailed information in [Nutanix LCM Perform Inventory v4] https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Inventory/operation/performInventory