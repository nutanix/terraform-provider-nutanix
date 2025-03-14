---
layout: "nutanix"
page_title: "NUTANIX: nutanix_lcm_status_v2"
sidebar_current: "docs-nutanix-datasource-lcm-status-v2"
description: |-
  Get the LCM framework status. 
---

# nutanix_lcm_status_v2

Get the LCM framework status. Represents the Status of LCM. Status represents details about a pending or ongoing action in LCM.

## Example

```hcl
# List Prism Central
data "nutanix_clusters_v2" "pc" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}
locals {
  pcExtID      = data.nutanix_clusters_v2.pc.cluster_entities[0].ext_id
}

data "nutanix_lcm_status_v2" "lcm_framework_status" {
  x_cluster_id = local.pcExtID
}
```

# Argument Reference
The following arguments are supported:

* `x_cluster_id`: (Optional) Cluster uuid on which the resource is present or operation is being performed.

See detailed information in [Nutanix LCM Status v4] https://developers.nutanix.com/api-reference?namespace=lifecycle&version=v4.0#tag/Status
