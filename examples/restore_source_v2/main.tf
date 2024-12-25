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

// using cluster location
resource "nutanix_restore_source_v2" "example-1" {
  location {
    cluster_location {
      config {
        ext_id = "cluster uuid"
      }
    }
  }
}


// using object store location
resource "nutanix_restore_source_v2" "example-2" {
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
        rpo_in_minutes = 70
      }
    }
  }
}

