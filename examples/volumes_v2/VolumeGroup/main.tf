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
  cluster1 = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_volume_group_v2" "example" {
  name                               = "example-volume-group"
  description                        = "Example Volume Group created by Terraform"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  cluster_reference                  = local.cluster1
  storage_features {
    flash_mode {
      is_enabled = false
    }
  }
  usage_type = "USER"
  is_hidden  = false
}

data "nutanix_volume_group_v2" "example" {
  ext_id     = nutanix_volume_group_v2.example.id
  depends_on = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_groups_v2" "example" {
  depends_on = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_group_stats_v2" "example" {
  ext_id     = nutanix_volume_group_v2.example.id
  start_time = "2024-01-01T00:00:00Z"
  depends_on = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_group_metadata_v2" "example" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  depends_on          = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_group_category_associations_v2" "example" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  depends_on          = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_group_external_iscsi_attachments_v2" "example" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  depends_on          = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_group_vm_attachments_v2" "example" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  depends_on          = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_group_disks_v2" "example" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  depends_on          = [nutanix_volume_group_v2.example]
}
