---
layout: "nutanix"
page_title: "NUTANIX: nutanix_promote_protected_resource_v2"
sidebar_current: "docs-nutanix-resource-promote-protected-resource-v2"
description: |-
  Promotes the specified synced entity at the target site. This is only relevant if the synced entity is protected in a synchronous schedule.

---

# nutanix_promote_protected_resource_v2

Promotes the specified synced entity at the target site. This is only relevant if the synced entity is protected in a synchronous schedule.


## Example: Promote a protected virtual machine on remote site

```hcl

# Promote a protected virtual machine on remote site
# This example promotes a protected virtual machine on a remote site.
# Steps:
# 1. Define the provider for the remote site
# 2. List domain Managers, Clusters for the local and remote sites
# 3. Create a category and a protection policy, on the local site
# 4. Create a virtual machine and associate it with the protection policy, on local site
# 5. Promote the protected virtual machine on the remote site

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
resource "nutanix_category_v2" "cat" {
  provider = nutanix
  key      = "tf-synchronous-pp"
  value    = "tf_synchronous_pp"
}

resource "nutanix_protection_policy_v2" "sync-pp" {
  provider    = nutanix
  name        = "tf-sync-pp"
  description = "create sync pp for vm"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }

  replication_locations {
    domain_manager_ext_id = local.localPcExtId
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = local.remotePcExtId
    label                 = "target"
    is_primary            = false
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.remoteClusterExtId]
      }
    }
  }

  category_ids = [nutanix_category_v2.cat.id]
}

resource "nutanix_virtual_machine_v2" "vm" {
  provider             = nutanix
  name                 = "tf-test-vm"
  description          = "create a new protected vm and get it"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.clusterExtId
  }
  categories {
    ext_id = nutanix_category_v2.cat.id
  }
  power_state = "OFF"
  depends_on = [nutanix_protection_policy_v2.sync-pp]
  provisioner "local-exec" {
    # sleep 5 min to wait for the vm to be protected
    command = "sleep 300"
    when    = create
  }
}

resource "nutanix_promote_protected_resource_v2" "promote-vm" {
  provider = nutanix.remote
  ext_id   = nutanix_virtual_machine_v2.vm.id
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: -(Required) The external identifier of a protected VM or volume group that can be used to retrieve the protected resource.


See detailed information in [Nutanix Promote Protected Resource v4](https://developers.nutanix.com/api-reference?namespace=dataprotection&version=v4.0#tag/ProtectedResources/operation/promoteProtectedResource).
