---
layout: "nutanix"
page_title: "NUTANIX: protected_resource_v2"
sidebar_current: "docs-nutanix-datasource-protected-resource-v2"
description: |-
  Get a protected resource


---

# protected_resource_v2

Get the details of the specified protected resource such as the restorable time ranges available on the local Prism Central and the state of replication to the targets specified in the applied protection policies. This applies only if the entity is protected in a minutely or synchronous schedule. Other protection schedules are not served by this endpoint yet, and are considered not protected.


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

resource "nutanix_virtual_machine_v2" "vm"{
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

data "nutanix_protected_resource_v2" "pr"{
  ext_id = nutanix_virtual_machine_v2.vm.id
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.


## Attributes Reference
The following attributes are exported:



See detailed information in [Nutanix Get Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/getProtectedResourceById).

