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
