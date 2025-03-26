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
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}


resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "vm-example"
  description          = "create vm example"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster_ext_id
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 0
    }
  }
  power_state = "OFF"
}

resource "nutanix_recovery_points_v2" "rp" {
  name                = "rp-example"
  expiration_time     = "2025-03-23T20:00:30Z"
  status              = "COMPLETE"
  recovery_point_type = "APPLICATION_CONSISTENT"
  vm_recovery_points {
    vm_ext_id = nutanix_virtual_machine_v2.vm.id
  }
}

resource "nutanix_recovery_point_restore_v2" "restore-rp" {
  ext_id         = nutanix_recovery_points_v2.rp.id
  cluster_ext_id = local.cluster_ext_id
  vm_recovery_point_restore_overrides {
    vm_recovery_point_ext_id = nutanix_recovery_points_v2.rp.vm_recovery_points[0].ext_id
  }
  depends_on = [nutanix_virtual_machine_v2.vm]
}

resource "nutanix_vm_revert_v2" "revert-vm" {
  ext_id                   = nutanix_virtual_machine_v2.vm.id
  vm_recovery_point_ext_id = nutanix_recovery_points_v2.rp.vm_recovery_points[0].ext_id
  depends_on               = [nutanix_recovery_point_restore_v2.restore-rp]
}
