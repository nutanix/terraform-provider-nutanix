terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=1.0.0"
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

# Fetch a list of events with optional filtering
data "nutanix_events_v2" "example" {
  limit = 10
}

output "total_events" {
  value = length(data.nutanix_events_v2.example.events)
}

output "first_event_type" {
  value = length(data.nutanix_events_v2.example.events) > 0 ? data.nutanix_events_v2.example.events[0].event_type : "none"
}
