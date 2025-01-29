terraform {
  required_providers {
    nutanix = {
      source  = "nutanixtemp/nutanix"
      version = "1.99.99"
    }
  }
}

provider "nutanix" {
  username = "admin"
  password = "Nutanix.123"
  endpoint = "10.101.176.123"
  insecure = true
  port     = 9440
}

data "nutanix_calm_snapshot_policy_list" "test_snapshot" {
  bp_name = "bp2"
  length = 250
  offset = 0
}

resource "nutanix_calm_app_recovery_point" "test_1" {
  app_uuid = "59e00130-b7a0-40ef-b6e3-25ae608648fd"
  action_name = "Snapshot_test_1"
  recovery_point_name = "snap0"
}