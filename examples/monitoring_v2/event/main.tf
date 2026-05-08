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

data "nutanix_events_v2" "events-list" {}

data "nutanix_events_v2" "filtered-events" {
  limit = 2
}

data "nutanix_event_v2" "get-event" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
