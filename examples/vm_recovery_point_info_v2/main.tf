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
resource "nutanix_recovery_points_v2" "rp-example" {
  name                = "terraform-test-recovery-point"
  expiration_time     = "2024-09-17T09:20:42Z"
  status              = "COMPLETE"
  recovery_point_type = "APPLICATION_CONSISTENT"
  vm_recovery_points {
    vm_ext_id = "<Vm-uuid-1>"
  }
}


# get VM recovery point info
data "nutanix_vm_recovery_point_info_v2" "example" {
  recovery_point_ext_id = nutanix_recovery_points_v2.rp-example.ext_id
  ext_id                = nutanix_recovery_points_v2.rp-example.vm_recovery_points[0].ext_id
  depends_on            = [nutanix_recovery_points_v2.rp-example]
}