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

// using cluster location
resource "nutanix_pc_restore_source_v2" "cluster-location" {
  provider = nutanix.pe
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtID
      }
    }
  }
}

# wait some time until the restore point is created
# keep reading the backup target until the last_sync_time is updated
data "nutanix_pc_backup_target_v2" "targets" {
  domain_manager_ext_id = local.domainManagerExtId
  ext_id                = nutanix_pc_backup_target_v2.cluster-location.id
}

# Get the restorable pcs using the restore source extID
data "nutanix_restorable_pcs_v2" "restorable-pcs" {
  provider              = nutanix.pe
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.ext_id
}

# after the restore point is created, you can list restore points
data "nutanix_pc_restore_points_v2" "list-restore-points" {
  provider                         = nutanix.pe
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.restorable-pcs.restorable_pcs.0.ext_id
  restore_source_ext_id            = nutanix_pc_restore_source_v2.cluster-location.id
}

# get restore point details
data "nutanix_pc_restore_point_v2" "get-restore-point" {
  provider                         = nutanix.pe
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
  restore_source_ext_id            = nutanix_pc_restore_source_v2.cluster-location.id
  ext_id                           = data.nutanix_pc_restore_points_v2.list-restore-points.restore_points.0.ext_id
}
