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

#create a virtual machine with minium configuration
resource "nutanix_virtual_machine_v2" "vm-2" {
  name                 = "vm-example-2"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.clusterExtId
  }
  power_state = "OFF"
}

# create RP with Vm Rp
resource "nutanix_recovery_points_v2" "rp-example-1" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2029-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "APPLICATION_CONSISTENT"
  vm_recovery_points {
    vm_ext_id           = nutanix_virtual_machine_v2.vm-1.id
    name                = "vm-recovery-point-1"
    expiration_time     = "2029-09-17T09:20:42Z"
    status              = "COMPLETE"
    recovery_point_type = "APPLICATION_CONSISTENT"
  }
}
# create RP with multiple Vm Rp
resource "nutanix_recovery_points_v2" "rp-example-2" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2029-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "CRASH_CONSISTENT"
  vm_recovery_points {
    vm_ext_id = nutanix_virtual_machine_v2.vm-1.id
  }
  vm_recovery_points {
    vm_ext_id = nutanix_virtual_machine_v2.vm-2.id
  }
}


# create a volume group
resource "nutanix_volume_group_v2" "vg-1" {
  name              = "volume-group-example-1"
  description       = "Test Create Volume group with spec"
  created_by        = "example"
  cluster_reference = local.clusterExtId
}

resource "nutanix_volume_group_v2" "vg-2" {
  name              = "volume-group-example-2"
  description       = "Test Create Volume group with spec"
  created_by        = "example"
  cluster_reference = local.clusterExtId
}

# create RP with VG Rp
resource "nutanix_recovery_points_v2" "rp-example-3" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2029-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "CRASH_CONSISTENT"
  volume_group_recovery_points {
    volume_group_ext_id = nutanix_volume_group_v2.vg-1.id
  }
}

# create RP with multiple VG Rp
resource "nutanix_recovery_points_v2" "rp-example-4" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2029-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "CRASH_CONSISTENT"
  volume_group_recovery_points {
    volume_group_ext_id = nutanix_volume_group_v2.vg-1.id
  }
  volume_group_recovery_points {
    volume_group_ext_id = nutanix_volume_group_v2.vg-2.id
  }
}


# create RP with multiple VG and Vms Rp
resource "nutanix_recovery_points_v2" "rp-example-5" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2029-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "CRASH_CONSISTENT"
  vm_recovery_points {
    vm_ext_id = nutanix_virtual_machine_v2.vm-1.id
  }
  vm_recovery_points {
    vm_ext_id = nutanix_virtual_machine_v2.vm-2.id
  }
  volume_group_recovery_points {
    volume_group_ext_id = nutanix_volume_group_v2.vg-1.id
  }
  volume_group_recovery_points {
    volume_group_ext_id = nutanix_volume_group_v2.vg-2.id
  }
}

# get RP data
data "nutanix_recovery_point_v2" "rp-example-6" {
  ext_id = nutanix_recovery_points_v2.rp-example-1.id
}

# list all RP
data "nutanix_recovery_points_v2" "rp-example-7" {
  depends_on = [nutanix_recovery_points_v2.rp-example-1, nutanix_recovery_points_v2.rp-example-2, nutanix_recovery_points_v2.rp-example-3, nutanix_recovery_points_v2.rp-example-4, nutanix_recovery_points_v2.rp-example-5]
}

# list all RP with filter
data "nutanix_recovery_points_v2" "rp-example-8" {
  filter = "name eq '${nutanix_recovery_points_v2.rp-example-1.name}'"
}

# list all RP with limit
data "nutanix_recovery_points_v2" "rp-example-9" {
  limit      = 3
  depends_on = [nutanix_recovery_points_v2.rp-example-1, nutanix_recovery_points_v2.rp-example-2, nutanix_recovery_points_v2.rp-example-3, nutanix_recovery_points_v2.rp-example-4, nutanix_recovery_points_v2.rp-example-5]

}
