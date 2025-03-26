#############################################################################
# Example main.tf for Nutanix + Terraform
#
# Author: haroon.dweikat@nutanix.com
#
# This script is a quick demo of how to use the following provider objects:
# 1 - configure provider for PC
# 2 - configure provider for PE
# 3 - get PC ExtID and Cluster ExtID
# 4 - create a backup target using PC Provider
# 5 - create a restore source using PE Provider
# 6 - get the list of restorable domain managers using the restore source ExtID and PE Provider
#
# Feel free to reuse, comment, and contribute, so that others may learn.
#####################################################################################
terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1.0"
    }
  }
}

#defining nutanix configuration for PC
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_pc_endpoint
  port     = 9440
  insecure = true
}

#defining nutanix configuration for PE
provider "nutanix" {
  alias    = "pe"
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_pe_endpoint
  port     = 9440
  insecure = true
}

data "nutanix_clusters_v2" "clusters" {}

locals {
  domainManagerExtID = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
    ][
    0
  ]
  clusterExtID = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
    ][
    0
  ]
}

resource "nutanix_pc_backup_target_v2" "cluster-location" {
  domain_manager_ext_id = local.domainManagerExtID
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtID
      }
    }
  }
}

resource "nutanix_pc_restore_source_v2" "cluster-location" {
  provider = nutanix.pe
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtID
      }
    }
  }
  depends_on = [nutanix_pc_backup_target_v2.cluster-location]
}

// Get the list of restorable domain managers
data "nutanix_restorable_pcs_v2" "restorable_pcs" {
  provider              = nutanix.pe
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.ext_id
}
