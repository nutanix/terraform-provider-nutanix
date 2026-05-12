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
  name              = "example-volume-group"
  description       = "Example Volume Group"
  cluster_reference = local.cluster1
  sharing_status    = "NOT_SHARED"
  storage_features {
    flash_mode {
      is_enabled = false
    }
  }
  usage_type = "USER"
  is_hidden  = false
}

data "nutanix_volume_group_v2" "get-volume-group" {
  ext_id = nutanix_volume_group_v2.example.id
}

data "nutanix_volume_groups_v2" "list-volume-groups" {
  depends_on = [nutanix_volume_group_v2.example]
}

data "nutanix_volume_group_stats_v2" "get-stats" {
  ext_id     = nutanix_volume_group_v2.example.id
  start_time = "2024-01-01T00:00:00Z"
  end_time   = "2024-01-02T00:00:00Z"
}

data "nutanix_volume_group_metadata_v2" "get-metadata" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
}

data "nutanix_volume_group_category_associations_v2" "get-categories" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
}

data "nutanix_volume_group_iscsi_attachments_v2" "get-iscsi" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
}

data "nutanix_volume_group_vm_attachments_v2" "get-vm-attachments" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
}

data "nutanix_volume_group_disks_v2" "get-disks" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
}
