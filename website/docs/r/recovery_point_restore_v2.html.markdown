---
layout: "nutanix"
page_title: "NUTANIX: nutanix_recovery_point_replicate_v2"
sidebar_current: "docs-nutanix-resource-recovery-point-replicate-v2"
description: |-
  This operation Restore a recovery point identified by {extId}.
---

# nutanix_recovery_point_replicate_v2
This operation Restore a recovery point identified by {extId}.
A comma separated list of the created VM and volume group external identifiers can be found in the task completion details under the keys `vm_ext_ids` and `volume_group_ext_ids` respectively.

## Example Usage

```hcl
  # restore RP 
resource "nutanix_recovery_point_restore_v2" "example"{
  ext_id         = "<recovery_point_external_id>"
  cluster_ext_id = "<cluster_external_id>"
  vm_recovery_point_restore_overrides {
    vm_recovery_point_ext_id = "<vm_recovery_point_uuid>"
  }
  volume_group_recovery_point_restore_overrides {
    volume_group_recovery_point_ext_id = "<volume_group_recovery_point_uuid>"
    volume_group_override_spec {
      name = "<volume_group_restored_name>"
    }
  }
}

```


## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier that can be used to retrieve the recovery point using its URL.
* `cluster_ext_id`: -(Required) External identifier of the cluster.
* `vm_recovery_point_restore_overrides`: -(Optional) List of specifications to restore a specific VM recovery point(s) that are a part of the top-level recovery point. A specific VM recovery point can be selected for restore by specifying its external identifier along with override specification (if any).
* `volume_group_recovery_point_restore_overrides`: -(Optional) List of specifications to restore a specific volume group recovery point(s) that are a part of the top-level recovery point. A specific volume group recovery point can be selected for restore by specifying its external identifier along with override specification (if any).
  

### vm_recovery_point_restore_overrides 

* `vm_recovery_point_ext_id`: -(Required) External identifier of a VM recovery point, that is a part of the top-level recovery point.

### volume_group_recovery_point_restore_overrides

* `volume_group_recovery_point_ext_id`: -(Required) External identifier of a volume group recovery point, that is a part of the top-level recovery point.
* `volume_group_override_spec`: -(Optional) Protected resource/recovery point restore that overrides the volume group configuration. The specified properties will be overridden for the restored volume group.

#### volume_group_override_spec

* `name`: -(Optional) The name of the restored volume group.


## Attribute Reference

The following attributes are exported:

* `ext_id`: - The external identifier that can be used to retrieve the recovery point using its URL.
* `cluster_ext_id`: - External identifier of the cluster.
* `vm_recovery_point_restore_overrides`: - List of specifications to restore a specific VM recovery point(s) that are a part of the top-level recovery point. A specific VM recovery point can be selected for restore by specifying its external identifier along with override specification (if any).
* `volume_group_recovery_point_restore_overrides`: - List of specifications to restore a specific volume group recovery point(s) that are a part of the top-level recovery point. A specific volume group recovery point can be selected for restore by specifying its external identifier along with override specification (if any).
* `vm_ext_ids`: - List of external identifiers of the created(restored) VMs.
* `volume_group_ext_ids`: - List of external identifiers of the created(restored) volume groups.

See detailed information in [Nutanix Recovery Point V4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0).