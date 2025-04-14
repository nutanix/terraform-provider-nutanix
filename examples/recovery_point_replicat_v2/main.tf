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
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

#create a virtual machine with minium configuration
resource "nutanix_virtual_machine_v2" "vm-1" {
  name                 = "vm-example-1"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.clusterExtId
  }
  power_state = "OFF"
}


# create RP with Vm Rp
resource "nutanix_recovery_points_v2" "rp" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2029-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "APPLICATION_CONSISTENT"
  vm_recovery_points {
    vm_ext_id = nutanix_virtual_machine_v2.vm-1.id
  }
}

# replicate RP
resource "nutanix_recovery_point_replicate_v2" "replicate-rp" {
  ext_id         = nutanix_recovery_points_v2.rp.id
  cluster_ext_id = var.remote_cluster_uuid # remote cluster uuid
  pc_ext_id      = var.remote_pc_uuid      # remote pc uuid
}
