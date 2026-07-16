#############################################################################
# Example main.tf for Nutanix + Terraform
#
# This script is a quick demo of how to use the following provider objects:
# - providers
#     - terraform-provider-nutanix
# - data sources
#     - nutanix_volume_disk_stats_v2
#
# Feel free to reuse, comment, and contribute, so that others may learn.
#############################################################################

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

##########################
### Data Sources
##########################

# Query the Volume Disk stats identified by {diskExtId} in the Volume Group
# identified by {volumeGroupExtId} over the [start_time, end_time] window.
data "nutanix_volume_disk_stats_v2" "example" {
  volume_group_ext_id = var.volume_group_ext_id
  ext_id              = var.volume_disk_ext_id
  start_time          = var.start_time
  end_time            = var.end_time
  sampling_interval   = var.sampling_interval
  stat_type           = var.stat_type
}
