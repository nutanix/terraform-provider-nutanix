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

# pull cluster data
data "nutanix_clusters_v2" "clusters" {}

# pull the desired (non Prism Central) cluster
locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# create a volume group
resource "nutanix_volume_group_v2" "example" {
  name                               = var.volume_group_name
  description                        = "volume group for VM attachments list example"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  created_by                         = "example"
  cluster_reference                  = local.cluster_ext_id
  usage_type                         = "USER"
  is_hidden                          = false
}

# create a VM to attach to the volume group
resource "nutanix_virtual_machine_v2" "example" {
  name              = var.vm_name
  num_sockets       = 1
  memory_size_bytes = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  power_state = "OFF"
}

# attach the VM to the volume group (index on the SCSI bus is optional)
resource "nutanix_volume_group_vm_v2" "example" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  vm_ext_id           = nutanix_virtual_machine_v2.example.id
  index               = 1
}

# List datasource — query all VM attachments for the volume group.
# NOTE: This API has been deprecated.
data "nutanix_volume_group_vms_v2" "attachments" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  depends_on          = [nutanix_volume_group_vm_v2.example]
}

# List datasource with pagination / filter query parameters
data "nutanix_volume_group_vms_v2" "filtered-attachments" {
  volume_group_ext_id = nutanix_volume_group_v2.example.id
  page                = 0
  limit               = 10
  depends_on          = [nutanix_volume_group_vm_v2.example]
}
