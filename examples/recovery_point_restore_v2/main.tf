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


# create RP with multiple VG and Vms Rp
resource "nutanix_recovery_points_v2" "rp-example" {
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

# restore RP
resource "nutanix_recovery_point_restore_v2" "rp-restore-example" {
  ext_id         = nutanix_recovery_points_v2.rp-example.id
  cluster_ext_id = local.clusterExtId
  vm_recovery_point_restore_overrides {
    vm_recovery_point_ext_id = nutanix_recovery_points_v2.rp-example.vm_recovery_points[0].ext_id
  }
  vm_recovery_point_restore_overrides {
    vm_recovery_point_ext_id = nutanix_recovery_points_v2.rp-example.vm_recovery_points[1].ext_id
  }
  volume_group_recovery_point_restore_overrides {
    volume_group_recovery_point_ext_id = nutanix_recovery_points_v2.rp-example.volume_group_recovery_points[0].ext_id
    volume_group_override_spec {
      name = "vg-1-rp-example-restore"
    }
  }
  volume_group_recovery_point_restore_overrides {
    volume_group_recovery_point_ext_id = nutanix_recovery_points_v2.rp-example.volume_group_recovery_points[1].ext_id
    volume_group_override_spec {
      name = "vg-2-rp-example-restore"
    }
  }
}
