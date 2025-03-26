#############################################################################
# Example main.tf for Nutanix + Terraform
#
# Author: haroon.dweikat@nutanix.com
#
# This script is a quick demo of how to use the following provider objects:
# - providers
#     - terraform-provider-nutanix
# - data sources
#     - nutanix_volume_iscsi_client_v2
#     - nutanix_volume_iscsi_clients_v2
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
### Data Sources
##########################



# List all the iSCSI clients.
data "nutanix_volume_iscsi_clients_v2" "list-iscsi-clients" {}

# list iSCSI clients with a filter.
data "nutanix_volume_iscsi_clients_v2" "list-iscsi-clients-filter" {
  filter = "clusterReference eq '${local.cluster_ext_id}'"
}

# list iSCSI clients with a limit and pagination.
data "nutanix_volume_iscsi_clients_v2" "list-iscsi-clients-limit" {
  page  = 2
  limit = 1
}

# Fetch an iSCSI client details.
data "nutanix_volume_iscsi_client_v2" "v_iscsi_client_example" {
  ext_id = data.nutanix_volume_iscsi_clients_v2.list-iscsi-clients.iscsi_clients.0.ext_id
}
