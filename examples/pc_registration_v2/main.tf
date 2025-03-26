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

data "nutanix_clusters_v2" "clusters" {}

locals {
  pc_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]
}

// DomainManagerRemoteClusterSpec
resource "nutanix_pc_registration_v2 " "pc1" {
  pc_ext_id = local.pc_ext_id
  remote_cluster {
    domain_manager_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = var.cvm_ip
          }
        }
        credentials {
          authentication {
            username = var.username
            password = var.password
          }
        }
      }
      cloud_type = "NUTANIX_HOSTED_CLOUD"
    }
  }
}

// AOSRemoteClusterSpec
resource "nutanix_pc_registration_v2 " "pc2" {
  pc_ext_id = local.pc_ext_id
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = var.cvm_ip
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

// ClusterReference
resource "nutanix_pc_registration_v2 " "pc1" {
  pc_ext_id = local.pc_ext_id
  remote_cluster {
    cluster_reference {
      ext_id = var.cluster_ext_id
    }
  }
}
