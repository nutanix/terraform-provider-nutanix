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
  cluster1 = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

##########################
### Resources
##########################


# create a volume group 
resource "nutanix_volume_group_v2" "volume_group_example" {
  name                               = var.volume_group_name
  description                        = "Test Create Volume group with spec"
  should_load_balance_vm_attachments = false
  sharing_status                     = var.volume_group_sharing_status
  target_name                        = "volumegroup-test-001234"
  created_by                         = "example"
  cluster_reference                  = local.cluster1
  iscsi_features {
    enabled_authentications = "CHAP"
    target_secret           = var.volume_group_target_secret
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



# create a volume group with attachement_type , protocol and disks
resource "nutanix_volume_group_v2" "volume_group_example" {
  name                               = var.volume_group_name
  description                        = "Test Create Volume group with spec"
  should_load_balance_vm_attachments = false
  sharing_status                     = var.volume_group_sharing_status
  target_name                        = "volumegroup-test-001234"
  created_by                         = "example"
  cluster_reference                  = local.cluster1
  iscsi_features {
    enabled_authentications = "CHAP"
    target_secret           = var.volume_group_target_secret
  }

  storage_features {
    flash_mode {
      is_enabled = true
    }
  }
  usage_type = "USER"
  attachment_type = "DIRECT"
  protocol = "ISCSI"
  disks {
    disk_size_bytes = 10 * 1024 * 1024 * 1024
    index = 1
    disk_data_source_reference {
      name        = "vg-disk-%[1]s"
      ext_id      = "<storage_container_uuid>"
      entity_type = "STORAGE_CONTAINER"
      uris        = ["uri1","uri2"]
    }
    disk_storage_features {
      flash_mode {
        is_enabled = false
      }
    }
  }
  is_hidden  = false

  # ignore changes to target_secret, target secret will not be returned in terraform plan output
  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret
    ]
  }
}


##########################
### Data Sources
##########################

# pull all volume groups
data "nutanix_volume_groups_v2" "vgs_example" {
  filter = "startswith(name, 'value')"
  limit  = 2
  page   = 0
}

# pull a specific volume group
data "nutanix_volume_group_v2" "vg_example" {
  ext_id = nutanix_volume_group_v2.volume_group_example.id
}
