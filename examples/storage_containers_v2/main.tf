terraform{
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = "1.3.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}

#pull all clusters data
data "nutanix_clusters" "clusters"{}

#create local variable pointing to desired cluster
locals {
  cluster1 = [
    for cluster in data.nutanix_clusters.clusters.entities :
    cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
  ][0]
}

#creating storage container
resource "nutanix_storage_containers_v2" "example" {
  name = "<storage_container_name>"
  cluster_ext_id = local.cluster1
  logical_advertised_capacity_bytes = <logical_advertised_capacity_bytes>
  logical_explicit_reserved_capacity_bytes = <logical_explicit_reserved_capacity_bytes>
  replication_factor = <replication_factor>
  nfs_whitelist_addresses {
    ipv4  {
      value = "<nfs_whitelist_addresses>"
      prefix_length ="<prefix_length>"
    }
  }
  erasure_code = "OFF"
  is_inline_ec_enabled = false
  has_higher_ec_fault_domain_preference = false
  cache_deduplication = "OFF"
  on_disk_dedup = "OFF"
  is_compression_enabled = true
  is_internal = false
  is_software_encryption_enabled = false
}

#output the storage container info
output "subnet" {
  value   = nutanix_storage_containers_v2.example
}


#pull all storage containers data in the system
data "nutanix_storage_containers_v2" "example"{}


#pull a storage containers data by ext id
data "nutanix_storage_container_v2" "example"{
  cluster_ext_id = local.cluster1
  storage_container_ext_id = nutanix_storage_containers_v2.example.id
}
