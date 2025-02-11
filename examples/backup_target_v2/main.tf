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

// using cluster_location
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



//using object store location
resource "nutanix_backup_target_v2" "object-store-location"{
  domain_manager_ext_id = local.domainManagerExtID
  location {
    object_store_location {
      provider_config {
        bucket_name = var.bucket_name
        region      = var.region
        credentials {
          access_key_id     = var.access_key_id
          secret_access_key = var.secret_access_key
        }
      }
      backup_policy {
        rpo_in_minutes = 120
      }
    }
  }
  lifecycle {
    ignore_changes = [
      location[0].object_store_location[0].provider_config[0].credentials
    ]
  }
}


// list backup targets
data "nutanix_backup_targets_v2" "backup-targets" {
  domain_manager_ext_id = local.domainManagerExtID
}

// get backup target
data "nutanix_backup_target_v2" "backup-target" {
  domain_manager_ext_id = local.domainManagerExtID
  ext_id = nutanix_backup_target_v2.cluster-location.id
}


