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
  insecure = true
  port     = var.nutanix_port
}

// Provision an application
resource "nutanix_self_service_app_provision" "test" {
    bp_name         = var.blueprint_name
    app_name        = var.app_name
    app_description = var.app_description
}

// Execute custom action with a name
resource "nutanix_self_service_app_custom_action" "test" {
    app_name        = var.app_name
    action_name = var.action_name
}