---
layout: "nutanix"
page_title: "NUTANIX: restore_protected_resource_v2"
sidebar_current: "docs-nutanix-resource-restore-protected-resource-v2"
description: |-
  Restore the specified protected resource from its state at the given timestamp on the given cluster. This is only relevant if the entity is protected in a minutely schedule at the given timestamp.



---

# restore_protected_resource_v2

Restore the specified protected resource from its state at the given timestamp on the given cluster. This is only relevant if the entity is protected in a minutely schedule at the given timestamp.


## Example 1: Restore Virtual Machine

```hcl
# Restore a protected virtual machine on remote site
# This example demonstrates how to restore a protected virtual machine on remote site.
# steps:
# 1. Define the provider for the remote site
# 2. List domain Managers, Clusters for the local and remote sites
# 3. Create a category and a protection policy, on the local site
# 4. Create a virtual machine and associate it with the protection policy, on local site
# 5. Restore the virtual machine on the remote site


# define another alias for the provider, this time for the remote PC
provider "nutanix" {
  alias    = "remote"
  username = var.nutanix_remote_username
  password = var.nutanix_remote_password
  endpoint = var.nutanix_remote_endpoint
  insecure = true
  port     = 9440
}


# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {
  provider = nutanix
}

# list Clusters
data "nutanix_clusters_v2" "clusters" {
  provider = nutanix
}

# remote pc list
data "nutanix_pcs_v2" "pcs-list-remote" {
  provider = nutanix.remote
}

# remote cluster list
data "nutanix_clusters_v2" "clusters-remote" {
  provider = nutanix.remote
}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][
  0
  ]
  localPcExtId       = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
  remoteClusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters-remote.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][
  0
  ]
  remotePcExtId = data.nutanix_pcs_v2.pcs-list-remote.pcs[0].ext_id
}

# Create Category
resource "nutanix_category_v2" "example" {
  provider    = nutanix
  key         = "tf-test-category-pp-restore-vm"
  value       = "tf_test_category_pp_restore_vm"
  description = "category for protection policy and protected vm"
}

resource "nutanix_protection_policy_v2" "pp-vm" {
  provider    = nutanix
  name        = "pp_example_1"
  description = "protection policy for restore vm"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = local.localPcExtId
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.remotePcExtId
    label                 = "target"
    is_primary            = false
  }

  category_ids = [nutanix_category_v2.example.id]
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "tf-vm-exmaple-restore"
  description          = "virtual machine for restore"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.clusterExtId
  }
  categories {
    ext_id = nutanix_category_v2.example.id
  }
  power_state = "OFF"
  depends_on = [nutanix_protection_policy_v2.pp-vm]
  provisioner "local-exec" {
    # sleep 5 min to wait for the vm to be protected
    command = "sleep 300"
    when    = create
  }
}

resource "nutanix_restore_protected_resource_v2" "rp-vm" {
  provider       = nutanix.remote
  ext_id         = nutanix_virtual_machine_v2.vm.id
  cluster_ext_id = local.remoteClusterExtId
}

```

## Example 2: Restore Volume Group

```hcl
# Restore a protected volume group on remote site
# This example demonstrates how to restore a protected volume group on remote site.
# steps:
# 1. Define the provider for the remote site
# 2. List domain Managers, Clusters for the local and remote sites
# 3. Create a category and a protection policy, on the local site
# 4. Create a volume group and associate it with the category on the local site
# 5. Restore the volume group on the remote site


# define another alias for the provider, this time for the remote PC
provider "nutanix" {
  alias    = "remote"
  username = var.nutanix_remote_username
  password = var.nutanix_remote_password
  endpoint = var.nutanix_remote_endpoint
  insecure = true
  port     = 9440
}


# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {
  provider = nutanix
}

# list Clusters
data "nutanix_clusters_v2" "clusters" {
  provider = nutanix
}

# remote pc list
data "nutanix_pcs_v2" "pcs-list-remote" {
  provider = nutanix.remote
}

# remote cluster list
data "nutanix_clusters_v2" "clusters-remote" {
  provider = nutanix.remote
}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][
  0
  ]
  localPcExtId       = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
  remoteClusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters-remote.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][
  0
  ]
  remotePcExtId = data.nutanix_pcs_v2.pcs-list-remote.pcs[0].ext_id
}

# Create Category
resource "nutanix_category_v2" "example" {
  provider    = nutanix
  key         = "tf-test-category-pp-restore-vg"
  value       = "tf_test_category_pp_restore_vg"
  description = "category for protection policy and protected vg"
}

resource "nutanix_protection_policy_v2" "pp-vg" {
  provider    = nutanix
  name        = "pp_example_1"
  description = "protection policy for restore vg"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = local.localPcExtId
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.remotePcExtId
    label                 = "target"
    is_primary            = false
  }

  category_ids = [nutanix_category_v2.example.id]
}

resource "nutanix_volume_group_v2" "vg" {
  name              = "tf-vg-exmaple-restore"
  description       = "volume group for restore"
  cluster_reference = local.clusterExtId
  depends_on = [nutanix_protection_policy_v2.pp-vg]
}

resource "nutanix_associate_category_to_volume_group_v2" "example" {
  ext_id = nutanix_volume_group_v2.vg.id
  categories {
    ext_id = nutanix_category_v2.example.id
  }
  provisioner "local-exec" {
    # sleep 7 min to wait for the vg to be protected
    command = "sleep 420"
  }
}

resource "nutanix_restore_protected_resource_v2" "rp-vg" {
  provider       = nutanix.remote
  ext_id         = nutanix_volume_group_v2.vg.id
  cluster_ext_id = local.remoteClusterExtId
}
```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.
* `cluster_ext_id`: -(Required) The external identifier of the cluster on which the entity has valid restorable time ranges. The restored entity will be created on the same cluster.
* `restore_time`: -(Optional) UTC date and time in ISO 8601 format representing the time from when the state of the entity should be restored. This needs to be a valid time within the restorable time range(s) for the protected resource.


See detailed information in [Nutanix Restore Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/restoreProtectedResourcen).

