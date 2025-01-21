---
layout: "nutanix"
page_title: "NUTANIX: restore_protected_resource_v2"
sidebar_current: "docs-nutanix-resource-restore-protected-resource-v2"
description: |-
  Restore the specified protected resource from its state at the given timestamp on the given cluster. This is only relevant if the entity is protected in a minutely schedule at the given timestamp.



---

# restore_protected_resource_v2

Restore the specified protected resource from its state at the given timestamp on the given cluster. This is only relevant if the entity is protected in a minutely schedule at the given timestamp.


## Example

```hcl


resource "nutanix_protection_policy_v2" "pp_1"{
  name     = "pp_example_1"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }

  replication_locations {
    domain_manager_ext_id = "domain_manager_ext_id_local"
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = "domain_manager_ext_id_remote"
    label                 = "target"
    is_primary            = false
  }

  category_ids = ["<category_ids>"]
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "%[2]s"
  description          = "%[3]s"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
  }
  categories {
    ext_id = local.category1
  }
  power_state = "OFF"
  depends_on = [nutanix_protection_policy_v2.test]
}

# wait some time for the VM to be created to be protected
# you need add delay 

resource "nutanix_restore_protected_resource_v2" "rp-vm" {
  ext_id = nutanix_virtual_machine_v2.vm.id
  cluster_ext_id = "<cluster_ext_id>"
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.
* `cluster_ext_id`: -(Required) The external identifier of the cluster on which the entity has valid restorable time ranges. The restored entity will be created on the same cluster.
* `restore_time`: -(Optional) UTC date and time in ISO 8601 format representing the time from when the state of the entity should be restored. This needs to be a valid time within the restorable time range(s) for the protected resource.


## Attributes Reference
The following attributes are exported:



See detailed information in [Nutanix Restore Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/restoreProtectedResourcen   ).

