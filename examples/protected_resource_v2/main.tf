terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "2.1.0"
        }
    }
}

#defining nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}

#################### Example 1: Get Protected Resource - VM ####################
# Create a protected virtual machine on remote site and get
# This example demonstrates how to get a protected virtual machine details
# steps:
# 1. Define the provider for the remote site
# 2. List domain Managers, Clusters for the local and remote sites
# 3. Create a category and a protection policy, on the local site
# 4. Create a virtual machine and associate it with the protection policy, on local site
# 5. Get the protected virtual machine details

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
  key         = "tf-example-category-pp-restore-vm"
  value       = "tf_example_category_pp_restore_vm"
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
  name                 = "tf-vm-example-restore"
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

data "nutanix_protected_resource_v2" "protected-vm" {
  ext_id = nutanix_virtual_machine_v2.vm.id
}



#################### Example 2: Get Protected Resource - Volume Group ####################
# Create a protected volume group on remote site and get details of the protected volume group
# This example demonstrates how to get a protected volume group .
# steps:
# 1. Define the provider for the remote site
# 2. List domain Managers, Clusters for the local and remote sites
# 3. Create a category and a protection policy, on the local site
# 4. Create a volume group and associate it with the category on the local site
# 5. Get the protected volume group details


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
  key         = "tf-example-category-pp-restore-vg"
  value       = "tf_example_category_pp_restore_vg"
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
  name                 = "tf-vg-example-restore"
  description          = "volume group for restore"
  cluster_reference                  = local.clusterExtId
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

data "nutanix_protected_resource_v2" "protected-vg" {
  ext_id = nutanix_volume_group_v2.vg.id
}
