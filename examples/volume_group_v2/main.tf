#############################################################################
# Example main.tf for Nutanix + Terraform
#
# Author: haroon.dweikat@nutanix.com
#
# This script is a quick demo of how to use the following provider objects:
# - providers
#     - terraform-provider-nutanix
# - resources
#     - nutanix_volume_group_v2
# - data sources
#     - nutanix_volume_groups_v2
#     - nutanix_volume_group_v2
# - script Variables
#     - clusterid's for targeting clusters within prism central
#
# Feel free to reuse, comment, and contribute, so that others may learn.
#
#############################################################################
### Define Provider Info for terraform-provider-nutanix
### This is where you define the credentials for ** Prism Central **
###
### NOTE:
###   While it may be possible to use Prism Element directly, Nutanix's
###   provider is not structured or tested for this. Using Prism Central will
###   give the broadest capabilities across the board


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

#pull cluster data
data "nutanix_clusters_v2" "clusters" {}

#pull desired cluster data from setup
locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}


data "nutanix_storage_containers_v2" "sg" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

##########################
### Resources
##########################


# create a volume group
resource "nutanix_volume_group_v2" "vg-example-1" {
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



# create a volume group with attachment_type , protocol and disks
resource "nutanix_volume_group_v2" "vg-example-2" {
  name                               = "volume-group-example-001236"
  description                        = "Create Volume group with spec"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  target_name                        = "volumegroup-example-001236"
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
  usage_type      = "USER"
  attachment_type = "DIRECT"
  protocol        = "ISCSI"
  disks {
    disk_size_bytes = 10 * 1024 * 1024 * 1024
    index           = 1
    disk_data_source_reference {
      name        = "vg-disk-example"
      ext_id      = data.nutanix_storage_containers_v2.sg.storage_containers[0].ext_id
      entity_type = "STORAGE_CONTAINER"
      uris        = ["uri1", "uri2"]
    }
    disk_storage_features {
      flash_mode {
        is_enabled = false
      }
    }
  }
  is_hidden = false

  # ignore changes to target_secret, target secret will not be returned in terraform plan output
  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret
    ]
  }
}

# create a volume group with minimum required fields
resource "nutanix_volume_group_v2" "vg-example-3" {
  name              = "volume-group-example-001237"
  cluster_reference = local.cluster_ext_id
}

##########################
### Data Sources
##########################

# list all volume groups
data "nutanix_volume_groups_v2" "volume-groups" {
  depends_on = [nutanix_volume_group_v2.vg-example-1, nutanix_volume_group_v2.vg-example-2]
}
# pull all volume groups
data "nutanix_volume_groups_v2" "filter-volume-groups" {
  filter = "name eq '${nutanix_volume_group_v2.vg-example-1.name}'"
  limit  = 1
  page   = 0
}

# pull a specific volume group
data "nutanix_volume_group_v2" "vg_example" {
  ext_id = nutanix_volume_group_v2.vg-example-1.id
}
