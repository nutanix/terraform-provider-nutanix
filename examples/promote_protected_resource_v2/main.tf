terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

############# Example 1: Promote a protected virtual machine on remote site #############

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
  name                 = "tf-example-vm"
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

############# Example 2: Promote a protected volume group on remote cluster #############

# Promote a protected volume group on remote cluster
# This example promotes a protected volume group on a remote cluster.
# Steps:
# 1. Get Unconfigured Node
# 2. fill the variables with the correct values

# This example performs the following steps:
# 1. Create a new cluster
# 2. Register the cluster to PC
# 3. Create a category.
# 4. Create a protection policy.
# 5. Create a volume group.
# 6. Associate the category to the volume group. and wait for 7 minutes to make sure that the vg is protected.
# 7. Promote the protected volume group.

# Note: You can replace volume group with virtual machine

# list Clusters
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}


locals {
  clusterExtId = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
  pcExtId      = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id

  ### Fill these variables with the correct values ############################
  # remote cluster is the cluster that we are going to create
  remoteClusterIP       = "10.xx.xx.xx"
  remoteClusterVIP = "10.xx.xx.xx"
  # local cluster is the PE connected to the PC
  localClusterIP = "10.xx.xx.xx" # PE IP
  localClusterVIP       = "10.xx.xx.xx" # PE VIP
  # username and password to ssh to the unconfigured node, used to reset the cluster password and modify the firewall rules
  remoteClusterUsername = "nutanix"
  remoteClusterPassword = "password-1"

  localClusterUsername = "nutanix"
  localClusterPassword = "password-1"

  # username and password to reset the cluster password
  resetRemoteClusterUsername = "admin"
  resetRemoteClusterPassword = "password-2"

  # Random number to avoid name conflicts
  randomNum = 99

  ##############################################################################

  ## Commands
  # commands to reset the cluster password and modify the firewall rules
  # no need to change these commands, just make sure that you have the correct values for the variables above
  resetClusterPassword = "/home/nutanix/prism/cli/ncli user reset-password user-name=${local.resetRemoteClusterUsername} password=${local.resetRemoteClusterPassword}"

  remoteClusterSSHCommand = "sshpass -p '${local.remoteClusterPassword}' ssh -o StrictHostKeyChecking=no ${local.remoteClusterUsername}@${local.remoteClusterIP}"
  localClusterSSHCommand  = "sshpass -p '${local.localClusterPassword}' ssh -o StrictHostKeyChecking=no ${local.localClusterUsername}@${local.localClusterIP}"

  resetClusterPasswordCommand = "${local.remoteClusterSSHCommand} '${local.resetClusterPassword}'"

  modifyFirewallRulesCommand       = "/usr/local/nutanix/cluster/bin/modify_firewall -f -r"
  modifyLocalClusterFirewallRules  = "${local.localClusterSSHCommand} '${local.modifyFirewallRulesCommand} ${local.remoteClusterIP},${local.remoteClusterVIP} -p 2030,2036,2073,2090,8740 -i eth0'"
  modifyRemoteClusterFirewallRules = "${local.remoteClusterSSHCommand} '${local.modifyFirewallRulesCommand} ${local.localClusterIP},${local.localClusterVIP}   -p 2030,2036,2073,2090,8740 -i eth0'"
}

# check if the nodes is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "example-discover-cluster-node" {
  ext_id       = local.pcExtId
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = local.remoteClusterIP
    }
  }

  ## check if the node is  un configured or not
  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 1
      error_message = "The node ${local.remoteClusterIP} are not unconfigured"
    }
  }

  depends_on = [data.nutanix_clusters_v2.clusters]
}

# create a new cluster
resource "nutanix_cluster_v2" "remote-cluster" {
  name = "tf-example-cluster-${local.randomNum}"
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = local.remoteClusterIP
        }
      }
    }
  }
  config {
    cluster_function = ["AOS"]
    cluster_arch = "X86_64"
    fault_tolerance_state {
      domain_awareness_level = "DISK"
    }
    redundancy_factor = 1
  }
  network {
    external_address {
      ipv4 {
        value = local.remoteClusterVIP
      }
    }
  }

  # Reset the cluster password
  provisioner "local-exec" {
    command    = local.resetClusterPasswordCommand
    on_failure = continue
  }
  # Set lifecycle to ignore changes
  lifecycle {
    ignore_changes = [network.0.smtp_server.0.server.0.password, links, categories, config.0.cluster_function]
  }
  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.example-discover-cluster-node]
}


# register the cluster to pc
resource "nutanix_pc_registration_v2" "node-registration" {
  pc_ext_id = local.pcExtId
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = local.remoteClusterIP
          }
        }
        credentials {
          authentication {
            username = local.resetRemoteClusterUsername
            password = local.resetRemoteClusterPassword
          }
        }
      }
    }
  }
  # Modify the firewall rules on Remote cluster
  provisioner "local-exec" {
    command    = local.modifyRemoteClusterFirewallRules
    when       = create
    on_failure = continue
  }
  depends_on = [nutanix_cluster_v2.remote-cluster]
}

# create a category, protection policy, volume group and associate it to the volume group
# list Clusters
data "nutanix_clusters_v2" "new-cls" {
  filter = "name eq '${nutanix_cluster_v2.remote-cluster.name}'"
  depends_on = [nutanix_pc_registration_v2.node-registration]
}

locals {
  newClusterExtId = data.nutanix_clusters_v2.new-cls.cluster_entities.0.ext_id
}

# Create Category
resource "nutanix_category_v2" "cat" {
  key   = "tf-example-category-pp-promote-vg-${local.randomNum}"
  value = "tf_example_category_pp_promote_vg_${local.randomNum}"

  # Modify the firewall rules on Local cluster
  provisioner "local-exec" {
    command    = local.modifyLocalClusterFirewallRules
    on_failure = continue
  }
  # Delay 5 minutes before destroying the resource to make sure that synced data is deleted
  provisioner "local-exec" {
    command    = "sleep 300"
    when       = destroy
    on_failure = continue
  }
  depends_on = [nutanix_pc_registration_v2.node-registration, nutanix_cluster_v2.remote-cluster]
}

resource "nutanix_protection_policy_v2" "pp" {
  name        = "tf-example-promote-pp-vg-${local.randomNum}"
  description = "create a new protected vg and promote it"

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
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "target"
    is_primary            = false
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.newClusterExtId]
      }
    }
  }

  category_ids = [nutanix_category_v2.cat.id]

}

resource "nutanix_volume_group_v2" "vg" {
  name              = "tf-example-promote-vg-${local.randomNum}"
  description       = "create a new protected vg to be promoted"
  cluster_reference = local.clusterExtId
  lifecycle {
    ignore_changes = [cluster_reference]
  }
  depends_on = [nutanix_protection_policy_v2.pp]
}


resource "nutanix_associate_category_to_volume_group_v2" "vg-cat" {
  ext_id = nutanix_volume_group_v2.vg.id
  categories {
    ext_id = nutanix_category_v2.cat.id
  }
  provisioner "local-exec" {
    # sleep 7 min to wait for the vg to be protected
    command = "sleep 420"
  }
}

resource "nutanix_promote_protected_resource_v2" "promote-vg" {
  ext_id = nutanix_volume_group_v2.vg.id
  depends_on = [nutanix_associate_category_to_volume_group_v2.vg-cat]
}
