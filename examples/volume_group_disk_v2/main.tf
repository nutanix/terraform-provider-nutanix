#############################################################################
# Example main.tf for Nutanix + Terraform
#
# Author: haroon.dweikat@nutanix.com
#
# This script is a quick demo of how to use the following provider objects:
# - providers
#     - terraform-provider-nutanix
# - resources
#     - nutanix_volume_group_disk_v2
# - data sources
#     - nutanix_volume_group_disks_v2
#     - nutanix_volume_group_disk_v2
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

##########################
### Resources
##########################


# create a volume group
resource "nutanix_volume_group_v2" "volume_group_example" {
  name                               = "volume-group-example-001234"
  description                        = "Create Volume group with spec example"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  target_name                        = "volumegroup-test-001234"
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

data "nutanix_storage_containers_v2" "sg" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

# create a volume group disk, and attach it to the volume group
resource "nutanix_volume_group_disk_v2" "disk_example" {
  volume_group_ext_id = resource.nutanix_volume_group_v2.volume_group_example.id
  # This Attribute is used to specify the index of the disk in the volume group.
  # its Optional, if not provided, the disk will be added at the end of the volume group.
  # if provided, the disk will be added at the specified index. make sure the index is unique.
  index       = 1
  description = "create volume disk example"
  # disk size in bytes
  disk_size_bytes = 5368709120

  disk_data_source_reference {
    name        = "disk1"
    ext_id      = data.nutanix_storage_containers_v2.sg.storage_containers[0].ext_id
    entity_type = "STORAGE_CONTAINER"
    uris        = ["uri1", "uri2"]
  }

  disk_storage_features {
    flash_mode {
      is_enabled = false
    }
  }

  # ignore changes to disk_data_source_reference, disk data source reference will not be returned in terraform plan output
  lifecycle {
    ignore_changes = [
      disk_data_source_reference
    ]
  }
}


##########################
### Data Sources
##########################

# pull all disks in a volume group
data "nutanix_volume_group_disks_v2" "vg_disks_example" {
  volume_group_ext_id = nutanix_volume_group_v2.volume_group_example.id
  filter              = "startswith(storageContainerId, '${nutanix_volume_group_disk_v2.disk_example.disk_data_source_reference.0.ext_id}')"
  limit               = 2
}

# pull a specific disk in a volume group
data "nutanix_volume_group_disk_v2" "vg_disk_example" {
  ext_id              = nutanix_volume_group_disk_v2.disk_example.id
  volume_group_ext_id = nutanix_volume_group_v2.volume_group_example.id
}
