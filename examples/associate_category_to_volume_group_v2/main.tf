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

#pull cluster data
data "nutanix_clusters_v2" "clusters" {}

#pull desired cluster data from setup
locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

// Create a volume group
resource "nutanix_volume_group_v2" "vg1" {
  name                               = "volume-group-example-001235"
  description                        = "Create Volume group with spec"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  target_name                        = "volumegroup-test-001235"
  created_by                         = "example"
  cluster_reference                  = local.cluster_ext_id
  iscsi_features {
    enabled_authentications = "CHAP"
    target_secret           = "pass.1234567890"
  }

  storage_features {
    flash_mode {
      is_enabled = true
    }
  }
  usage_type = "USER"
  is_hidden  = false

  # ignore changes to target_secret, target secret will not be returned in terraform plan output
  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret
    ]
  }
}



#creating category
resource "nutanix_category_v2" "vg-category" {
  key         = "category_example_key"
  value       = "category_example_value"
  description = "category example to associate with volume group"
}


# Associate categories to volume group
resource "nutanix_associate_category_to_volume_group_v2" "attach_category" {
  ext_id = nutanix_volume_group_v2.vg1.id
  categories {
    ext_id = nutanix_category_v2.vg-category.id
  }
}


# pull associated category data using list categories data source
data "nutanix_categories_v2" "associated_vg" {
  filter     = "extId eq '${nutanix_category_v2.vg-category.id}'"
  expand     = "associations"
  depends_on = [nutanix_associate_category_to_volume_group_v2.attach_category]
}

# pull associated category data fetch category data source
data "nutanix_category_v2" "associated_vg" {
  ext_id     = nutanix_category_v2.vg-category.id
  expand     = "associations"
  depends_on = [nutanix_associate_category_to_volume_group_v2.attach_category]
}
