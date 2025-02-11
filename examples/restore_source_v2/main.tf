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
  alise    = "nutanix-pe"
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
  clusterExtID       = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][
  0
  ]
}

resource "nutanix_backup_target_v2" "cluster-location"{
  domain_manager_ext_id = local.domainManagerExtID
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtID
      }
    }
  }
}

resource "nutanix_restore_source_v2" "cluster-location"{
  provider = nutanix.pe
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtID
      }
    }
  }
  depends_on = [nutanix_backup_target_v2.cluster-location]
}

// get the restore source
data "nutanix_restore_source_v2" "restore-source" {
  ext_id = nutanix_restore_source_v2.cluster-location.id
}