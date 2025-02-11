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

// Create a volume group
resource "nutanix_volume_group_v2" "example"{
  name                               = "test_volume_group"
  description                        = "Test Volume group with min spec and no Auth"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  target_name                        = "volumegroup-test-0"
  created_by                         = "Test"
  cluster_reference                  = "<Cluster uuid>"
  iscsi_features {
    enabled_authentications = "CHAP"
    target_secret           = "1234567891011"
  }

  storage_features {
    flash_mode {
      is_enabled = true
    }
  }
  usage_type = "USER"
  is_hidden  = false

  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret
    ]
  }
}


# List categories
data "nutanix_categories_v2" "categories"{}

# Associate categories to volume group
resource "nutanix_associate_category_to_volume_group_v2" "example"{
  ext_id = nutanix_volume_group_v2.example.id
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.0.ext_id
  }
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.1.ext_id
  }
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.2.ext_id
  }
}