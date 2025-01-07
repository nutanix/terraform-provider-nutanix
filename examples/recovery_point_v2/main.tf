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


# create RP with Vm Rp
resource "nutanix_recovery_points_v2" "rp-example-1" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2024-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "APPLICATION_CONSISTENT"
  vm_recovery_points {
    vm_ext_id = "<Vm-uuid-1>"
    name     = "vm-recovery-point-1"
    expiration_time = "2024-09-17T09:20:42Z"
    status = "COMPLETE"
    recovery_point_type = "APPLICATION_CONSISTENT"
  }
}

# create RP with multiple Vm Rp
resource "nutanix_recovery_points_v2" "rp-example-2" {
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
}

# create RP with VG Rp
resource "nutanix_recovery_points_v2" "rp-example-3" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2024-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "CRASH_CONSISTENT"
  volume_group_recovery_points {
    volume_group_ext_id = "<VG-uuid-1>"
  }
}

# create RP with multiple VG Rp
resource "nutanix_recovery_points_v2" "rp-example-4" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2024-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "CRASH_CONSISTENT"
  volume_group_recovery_points {
    volume_group_ext_id = "<VG-uuid-1>"
  }
  volume_group_recovery_points {
    volume_group_ext_id = "<VG-uuid-2>"
  }
}


# create RP with multiple VG and Vms Rp
resource "nutanix_recovery_points_v2" "rp-example-5" {
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

# get RP data
data "nutanix_recovery_point_v2" "rp-example-6" {
  ext_id = nutanix_recovery_points_v2.rp-example.id
}

# list all RP
data "nutanix_recovery_points_v2" "rp-example-7" {}

# list all RP with filter
data "nutanix_recovery_points_v2" "rp-example-8" {
  filter = "name eq 'terraform-test-recovery-point'"
}

# list all RP with limit
data "nutanix_recovery_points_v2" "rp-example-9" {
  limit = 10
}

# vm recovery point details
data "nutanix_vm_recovery_point_info_v2" "rp-example-10" {
  recovery_point_ext_id = nutanix_recovery_points_v2.example-1.ext_id
  ext_id                = nutanix_recovery_points_v2.example-1.vm_recovery_points[0].ext_id
}
