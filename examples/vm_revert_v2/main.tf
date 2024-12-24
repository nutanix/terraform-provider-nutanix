terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
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


# revert vm from recovery point
resource "nutanix_vm_revert_v2" "example" {
  ext_id = "<VM_UUID>"
  vm_recovery_point_ext_id = "<Vm_Recovery_Point_UUID>"
}