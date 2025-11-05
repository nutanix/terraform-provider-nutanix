terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
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



#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  pc_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]
}

# Add Cluster with minimum configuration.

resource "nutanix_cluster_v2" "example" {
  name = "terraform-example-cluster"
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = var.node_ip
        }
      }
    }
  }
  config {
    cluster_function  = ["AOS"]
    redundancy_factor = 1
    cluster_arch      = "X86_64"
    fault_tolerance_state {
      domain_awareness_level = "DISK"
    }
  }
  # after create a cluster you need to reset the pe ui password
  provisioner "local-exec" {
    command = "sshpass -p '${var.pe_password}' ssh ${var.pe_username}@${var.node_ip} '/home/nutanix/prism/cli/ncli user reset-password user-name=${var.username} password=${var.password}'"
  }
}

# after create a cluster you need to register the cluster with prism central
# to be able to do read, update, delete operations on the cluster and use it
resource "nutanix_pc_registration_v2" "pc1" {
  pc_ext_id = local.pc_ext_id
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = var.node_ip
          }
        }
        credentials {
          authentication {
            username = var.username
            password = var.password
          }
        }
      }
    }
  }
}

# List all  Directory Services.
data "nutanix_clusters_v2" "example" {
  depends_on = [nutanix_pc_registration_v2.pc1]
}

# List Clusters with filter.
data "nutanix_clusters_v2" "example" {
  filter     = "name eq '${nutanix_cluster_v2.example.name}'"
  depends_on = [nutanix_pc_registration_v2.pc1]
}


# Get a Directory Service.
data "nutanix_cluster_v2" "example" {
  ext_id     = nutanix_cluster_v2.example.id
  depends_on = [nutanix_pc_registration_v2.pc1]
}


########################################################################
#### Create 3 Node Cluster then register it with prism central  ########
#### then adding node to the cluster created above              ########
########################################################################
data "nutanix_clusters_v2" "cluster" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}
data "nutanix_clusters_v2" "pc" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  clusterExtId = data.nutanix_clusters_v2.cluster.cluster_entities[0].ext_id
  pcExtId      = data.nutanix_clusters_v2.pc.cluster_entities[0].ext_id
}



############################ cluster with 3 nodes

## check if the nodes is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "cluster-nodes" {
  ext_id       = local.pcExtId
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = var.nodes_ip[0]
    }
  }
  ip_filter_list {
    ipv4 {
      value = var.nodes_ip[1]
    }
  }
  ip_filter_list {
    ipv4 {
      value = var.nodes_ip[2]
    }
  }
  ## check if the 3 nodes are un configured or not
  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 3
      error_message = "The nodes are not unconfigured"
    }
  }
}
resource "nutanix_cluster_v2" "cluster-3nodes" {
  name   = "tf-cluster-3nodes"
  dryrun = false
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = var.nodes_ip[0]
        }
      }
    }
    node_list {
      controller_vm_ip {
        ipv4 {
          value = var.nodes_ip[1]
        }
      }
    }
    node_list {
      controller_vm_ip {
        ipv4 {
          value = var.nodes_ip[2]
        }


      }
      should_skip_host_networking   = false
      should_skip_pre_expand_checks = true

    }
    ## Uncomment the block below after creating and registering a 3-node cluster to Prism Central;
    #  then run 'terraform apply' to add a 4th node. To remove the node later, comment the block again and rerun 'terraform apply'.
    # node_list {
    #   controller_vm_ip {
    #     ipv4 {
    #       value = var.nodes_ip[3]
    #     }
    #   }
    #   should_skip_host_networking   = false
    #   should_skip_pre_expand_checks = true
    # }

  }
  config {
    cluster_function = ["AOS"]
    cluster_arch     = "X86_64"
    fault_tolerance_state {
      domain_awareness_level = "NODE"
    }
  }

  provisioner "local-exec" {
    command = "ssh-keygen -f ~/.ssh/known_hosts -R ${var.nodes_ip[1]};   sshpass -p '${var.pe_password}' ssh -o StrictHostKeyChecking=no ${var.pe_username}@${var.nodes_ip[1]} '/home/nutanix/prism/cli/ncli user reset-password user-name=${var.username} password=${var.password}'"

    on_failure = continue
  }

  lifecycle {
    ignore_changes = [links, categories, config.0.cluster_function]
  }

  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.cluster-nodes]
}


## we need to register on of 3 nodes cluster to pc
resource "nutanix_pc_registration_v2" "nodes-registration" {
  pc_ext_id = local.pcExtId
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = var.nodes_ip[0]
          }
        }
        credentials {
          authentication {
            username = var.username
            password = var.password
          }
        }
      }
    }
  }
  depends_on = [nutanix_cluster_v2.cluster-3nodes]

  provisioner "local-exec" {
    command    = " sleep 5s"
    on_failure = continue
  }
}
