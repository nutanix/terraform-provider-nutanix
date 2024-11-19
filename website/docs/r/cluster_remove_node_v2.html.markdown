---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_remove_node_v2"
sidebar_current: "docs-nutanix-resource-nutanix-cluster-remove-node-v2"
description: |-
  Removes nodes from cluster identified by {extId}.

---

# nutanix_cluster_add_node_v2

Removes nodes from cluster identified by {extId}.


## Example Usage

```hcl
resource "nutanix_cluster_remove_node_v2" "cluster_node"{
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  node_uuids = ["00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"]
  should_skip_remove = false
  should_skip_prechecks = true
  extra_params {
    should_skip_upgrade_check = true
    skip_space_check = true
    should_skip_add_check = false
  }
}
```


## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: -(Required) Cluster UUID.
* `should_skip_prechecks`: -(Optional) Indicates if prechecks can be skipped for node removal.
* `should_skip_remove`: -(Optional) Indicates if node removal can be skipped.
* `node_uuids`: -(Required) List of node UUIDs to be removed.
* `extra_params`: -(Optional) Extra parameters for node addition.

### Extra Params
The extra_params block supports the following:

* `should_skip_upgrade_check`: -(Optional) Indicates if upgrade check needs to be skipped or not.
* `skip_space_check`: -(Optional) Indicates if space check needs to be skipped or not.
* `should_skip_add_check`: -(Optional) Indicates if add check needs to be skipped or not.


See detailed information in [Nutanix Cluster - Remove Nodes From Cluster](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.0.b2#tag/Clusters/operation/removeNode).

