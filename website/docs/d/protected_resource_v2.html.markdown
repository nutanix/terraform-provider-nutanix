---
layout: "nutanix"
page_title: "NUTANIX: protected_resource_v2"
sidebar_current: "docs-nutanix-datasource-protected-resource-v2"
description: |-
  Get a protected resource


---

# protected_resource_v2

Get the details of the specified protected resource such as the restorable time ranges available on the local Prism Central and the state of replication to the targets specified in the applied protection policies. This applies only if the entity is protected in a minutely or synchronous schedule. Other protection schedules are not served by this endpoint yet, and are considered not protected.


## Example 1: Get Protected Virtual Machine

```hcl

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

```

## Example 2: Get Protected Volume Group

```hcl
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
  depends_on = [nutanix_associate_category_to_volume_group_v2.example]
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.


## Attributes Reference
The following attributes are exported:

* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `links`: 
* `entity_ext_id`: The external identifier of the VM or the volume group associated with the protected resource.
* `ext_id`: - The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.
* `entity_type`: Protected resource entity type. Possible values are: VM, VOLUME_GROUP.
* `source_site_reference`: Details about the data protection site in the Prism Central.
* `site_protection_info`: The data protection details for the protected resource that are relevant to any of the sites in the local Prism Central, like the time ranges available for recovery.
* `replication_states`: Replication related information about the protected resource.
* `consistency_group_ext_id`: External identifier of the Consistency group which the protected resource is part of.
* `category_fq_names`: Category key-value pairs associated with the protected resource at the time of protection. The category key and value are separated by '/'. For example, a category with key 'dept' and value 'hr' will be represented as 'dept/hr'.

### Links
The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Source Site Reference
The source_site_reference attribute supports the following:

* `mgmt_cluster_ext_id`: External identifier of the Prism Central.
* `cluster_ext_id`: External identifier of the cluster.

### Site Protection Info
The site_protection_info attribute supports the following:

* `recovery_info`:  The restorable time range details that can be used to recover the protected resource.
* `location_reference`: Details about the data protection site in the Prism Central.
* `synchronous_replication_role`: Synchronous Replication role related information of the protected resource. Possible values are:
  - `DECOUPLED`: VM is no longer in Synchronous Replication, and all the actions are blocked on VM, except a delete operation.
  - `SECONDARY`: This is the target site for VM in Synchronous Replication.
  - `INDEPENDENT`: VM is no longer in Synchronous Replication, and not replicating to the configured recovery cluster.
  - `PRIMARY`: VM is in Synchronous Replication, and is active on the primary site.

#### Recovery Info
The recovery_info attribute supports the following:

* `restorable_time_ranges`: The restorable time range details that can be used to recover the protected resource.

#### Restorable Time Range
The restorable_time_ranges attribute supports the following:

* `start_time`: UTC date and time in ISO 8601 format representing the time when the restorable time range for the entity starts.
* `end_time`: UTC date and time in ISO 8601 format representing the time when the restorable time range for the entity starts.

#### Location Reference
The location_reference attribute supports the following:

* `mgmt_cluster_ext_id`: External identifier of the Prism Central.
* `cluster_ext_id`: External identifier of the cluster.

### Replication States
The replication_states attribute supports the following:

* `protection_policy_ext_id`: The external identifier of the Protection policy associated with the protected resource.
* `recovery_point_objective_seconds`: The recovery point objective of the schedule in seconds.
* `replication_status`: Status of replication to a specified target site. Possible values are:
    - `IN_SYNC`: The specified recovery point objective is met on the target site and failover can be performed.
    - `SYNCING`: The system is trying to meet the specified recovery point objective for the target site via ongoing replications and failover can't yet be performed.
    - `OUT_OF_SYNC`: The replication schedule is disabled and there are no ongoing replications. Manual action might be needed by the user to meet the recovery point objective.
* `target_site_reference`: Details about the data protection site in the Prism Central.

#### Target Site Reference
The target_site_reference attribute supports the following:

* `mgmt_cluster_ext_id`: External identifier of the Prism Central.
* `cluster_ext_id`: External identifier of the cluster.


See detailed information in [Nutanix Get Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/getProtectedResourceById).

