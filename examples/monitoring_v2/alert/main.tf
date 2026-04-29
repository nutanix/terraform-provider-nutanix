variable "nutanix_endpoint" {}
variable "nutanix_username" {}
variable "nutanix_password" {}
variable "nutanix_port" {}
variable "nutanix_insecure" {}

provider "nutanix" {
  endpoint = var.nutanix_endpoint
  username = var.nutanix_username
  password = var.nutanix_password
  port     = var.nutanix_port
  insecure = var.nutanix_insecure
}

# Update alert email configuration
resource "nutanix_alert_email_configuration_v2" "example" {
  is_enabled              = true
  is_email_digest_enabled = true
  email_contact_list      = ["admin@example.com"]
}

# Manage alert - acknowledge
resource "nutanix_manage_alert_v2" "example" {
  ext_id      = "00000000-0000-0000-0000-000000000000"
  action_type = "ACKNOWLEDGE"
}

# Get a single alert by ID
data "nutanix_alert_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}

# List all alerts
data "nutanix_alerts_v2" "example" {}
