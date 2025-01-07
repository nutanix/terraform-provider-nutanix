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


# create RP with multiple VG and Vms Rp
resource "nutanix_recovery_points_v2" "rp-example" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2024-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "CRASH_CONSISTENT"
  vm_recovery_points {
    vm_ext_id = "<Vm-uuid-1>"
  }
  vm_recovery_points {
    vm_ext_id = "<Vm-uuid-2>"
  }
  volume_group_recovery_points {
    volume_group_ext_id = "<VG-uuid-1>"
  }
  volume_group_recovery_points {
    volume_group_ext_id = "<VG-uuid-2>"
  }
}

# restore RP
resource "nutanix_recovery_point_restore_v2" "rp-restore-example" {
  ext_id         = nutanix_recovery_points_v2.rp-example.id
  cluster_ext_id = "<cluster-uuid>"
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
  depends_on = [nutanix_recovery_points_v2.test]
}