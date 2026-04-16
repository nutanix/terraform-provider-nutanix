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

#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  cluster1 = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

#creating storage container
resource "nutanix_storage_containers_v2" "example" {
  name                                     = "example_storage_container"
  cluster_ext_id                           = local.cluster1
  logical_advertised_capacity_bytes        = 1073741824000
  logical_explicit_reserved_capacity_bytes = 32
  replication_factor                       = 1
  nfs_whitelist_addresses {
    ipv4 {
      value         = "192.168.15.0"
      prefix_length = 32
    }
  }
  erasure_code                          = "OFF"
  is_inline_ec_enabled                  = false
  has_higher_ec_fault_domain_preference = false
  cache_deduplication                   = "OFF"
  on_disk_dedup                         = "OFF"
  is_compression_enabled                = true
  is_internal                           = false
  is_software_encryption_enabled        = false
}

#output the storage container info
output "storage-container" {
  value = nutanix_storage_containers_v2.example
}


#list all storage containers
data "nutanix_storage_containers_v2" "list-storage-containers" {
  depends_on = [ data.nutanix_storage_container_v2.example ]
}


#pull a storage containers data by ext id
data "nutanix_storage_container_v2" "example" {
  ext_id         = nutanix_storage_containers_v2.example.id
}
