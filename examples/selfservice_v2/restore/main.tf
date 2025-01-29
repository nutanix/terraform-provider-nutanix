terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.99.99"
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

data "nutanix_calm_app_snapshots" "snapshots" {
  app_uuid = var.app_uuid
  length = 250
  offset = 0
}

#create local variable pointing to desired recovery point
locals {
	snapshot_uuid = [
	  for snapshot in data.nutanix_calm_app_snapshots.snapshots.entities :
	  snapshot.uuid if snapshot.name == var.snapshot_name
	][0]
}

resource "nutanix_calm_app_restore" "RestoreAction" {
  restore_action_name = var.restore_action_name
  app_uuid = var.app_uuid
  snapshot_uuid = local.snapshot_uuid
}