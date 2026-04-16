#############################################################################
# Example main.tf for Nutanix + Terraform
#
# Author: haroon.dweikat@nutanix.com
#
# This script is a quick demo of how to use the following provider objects:
# - providers
#     - terraform-provider-nutanix
# - resources
#     - nutanix_volume_group_vm_v2
# - data sources
#     - nutanix_volume_group_vms_v2
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
data "nutanix_clusters_v2" "clusters" {
}

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

resource "nutanix_virtual_machine_v2" "vm" {
  name              = "tf-vg-vm-example"
  is_agent_vm       = false
  num_sockets       = 1
  memory_size_bytes = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }

  power_state = "OFF"
}

# attach the volume group to a VM
resource "nutanix_volume_group_vm_v2" "vg_vm_example" {
  volume_group_ext_id = nutanix_volume_group_v2.vg-example-1.id
  vm_ext_id           = nutanix_virtual_machine_v2.vm.id
}
