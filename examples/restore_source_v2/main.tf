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

// using cluster_location
resource "nutanix_restore_source_v2" "cluster-location"{
  location {
    cluster_location {
      config {
        ext_id = var.cluster_ext_id
      }
    }
  }
}

//using object store location
resource "nutanix_restore_source_v2" "object-store-location"{
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

// get the restore source
data "nutanix_restore_source_v2" "restore-source" {
  ext_id = nutanix_restore_source_v2.cluster-location.id
}