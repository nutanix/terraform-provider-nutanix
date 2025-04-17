terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.2.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

# Read available recovery points in application
data "nutanix_self_service_app_snapshots" "snapshots" {
  app_uuid = var.app_uuid
  length = 250
  offset = 0
}

# Create local variable pointing to desired recovery point
locals {
	snapshot_uuid = [
	  for snapshot in data.nutanix_self_service_app_snapshots.snapshots.entities :
	  snapshot.uuid if snapshot.name == var.snapshot_name
	][0]
}

// Restore from a recovery point
resource "nutanix_self_service_app_restore" "RestoreAction" {
  restore_action_name = var.restore_action_name
  app_uuid = var.app_uuid
  snapshot_uuid = local.snapshot_uuid
}