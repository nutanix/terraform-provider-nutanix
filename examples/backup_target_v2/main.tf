terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1"
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
resource "nutanix_backup_target_v2" "example-1" {
  domain_manager_ext_id = "<domain_manager_uuid>"
  location {
    cluster_location {
      config {
        ext_id = "cluster uuid"
      }
    }
  }
}

// using object_store_location
resource "nutanix_backup_target_v2" "example-2" {
  domain_manager_ext_id = "<domain_manager_uuid>"
  location {
    object_store_location {
      provider_config {
        bucket_name = "bucket name"
        region      = "region"
        credentials {
          access_key_id     = "id"
          secret_access_key = "key"
        }
      }
      backup_policy {
        rpo_in_minutes = 0
      }
    }
  }
}

// list backup targets
data "nutanix_backup_targets_v2" "backup-targets" {
  domain_manager_ext_id = "<domain_manager_uuid>"
}

// get backup target
data "nutanix_backup_target_v2" "backup-target" {
  domain_manager_ext_id = "<domain_manager_uuid>"
  ext_id = "<backup_target_uuid>"
}


