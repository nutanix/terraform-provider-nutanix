terraform {
  required_providers {
    nutanix = {
      source  = "nutanixtemp/nutanix"
      version = "1.99.99"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  insecure = true
  port     = var.nutanix_port
}

data "nutanix_self_service_snapshot_policy_list" "test_snapshot" {
  bp_name = "sample_blueprint"
  length = 250
  offset = 0
}

resource "nutanix_self_service_app_recovery_point" "test_1" {
  app_uuid = var.app_uuid
  action_name = var.snapshot_action_name
  recovery_point_name = var.recovery_point_name
}