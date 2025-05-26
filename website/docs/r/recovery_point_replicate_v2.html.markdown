---
layout: "nutanix"
page_title: "NUTANIX: nutanix_recovery_point_replicate_v2"
sidebar_current: "docs-nutanix-resource-recovery-point-replicate-v2"
description: |-
  This operation Replicate a recovery point
---

# nutanix_recovery_point_replicate_v2

External identifier of the replicated recovery point can be found in the task completion details under the key

## Example Usage

```hcl
# replicate RP
resource "nutanix_recovery_point_replicate_v2" "rp-replicate" {
  ext_id         = "150a7ed0-9d05-4f35-a060-16dac4c835d0"
  cluster_ext_id = "1cefd0f5-6d38-4c9b-a07c-bdd2db004224"
  pc_ext_id      = "8a938cc5-282b-48c4-81be-de22de145d07"
}

```


## Argument Reference

External identifier of the cluster and the Prism Central where the recovery point is to be replicated. The recovery point identified by `ext_id` is replicated from the current Prism Central to the remote Prism Central with the external identifier `pc_ext_id` . This replication allows the data-protection service to decide on which cluster to perform the automatic replication. However, the user also has the option to choose the cluster identified by `cluster_ext_id` to which the recovery point identified by extId should be replicated. Cross-AZ replication can be performed only by users having legacy roles

The following arguments are supported:

* `ext_id`: -(Required) The external identifier that can be used to retrieve the recovery point using its URL.
* `cluster_ext_id`: -(Required) External identifier of the cluster.
* `pc_ext_id`: -(Required) External identifier of the Prism Central.


## Attribute Reference

The following attributes are exported:

* `ext_id`: - The external identifier that can be used to retrieve the recovery point using its URL.
* `cluster_ext_id`: - External identifier of the cluster.
* `pc_ext_id`: - External identifier of the Prism Central.
* `replicated_rp_ext_id`: - External identifier of replicated recovery point.

See detailed information in [Nutanix Replicate a Recovery Point V4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/RecoveryPoints/operation/replicateRecoveryPoint).
