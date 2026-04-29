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

# Fetch a single event by its external identifier
data "nutanix_event_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}

output "event_type" {
  value = data.nutanix_event_v2.example.event_type
}

output "event_message" {
  value = data.nutanix_event_v2.example.message
}
