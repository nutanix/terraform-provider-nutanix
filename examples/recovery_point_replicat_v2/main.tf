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

# replicate RP
resource "nutanix_recovery_point_replicate_v2" "test" {
  ext_id         = nutanix_recovery_points_v2.rp-example.id
  cluster_ext_id = "<cluster-uuid>" # remote cluster uuid
  pc_ext_id      = "<pc-uuid>" # remote pc uuid
  depends_on     = [nutanix_recovery_points_v2.test]
}