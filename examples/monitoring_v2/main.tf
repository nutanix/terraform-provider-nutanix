terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1.0"
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

# List audits
data "nutanix_audits_v2" "audits" {}

# List audits with limit
data "nutanix_audits_v2" "audits_limited" {
  limit = 10
}

# List audits with filter
data "nutanix_audits_v2" "audits_filtered" {
  filter = "serviceName eq 'Nutanix'"
}

# Get audit by ext_id
data "nutanix_audit_v2" "audit" {
  ext_id = data.nutanix_audits_v2.audits.audits.0.ext_id
}

# Output the audit details
output "audit_type" {
  value = data.nutanix_audit_v2.audit.audit_type
}

output "audit_status" {
  value = data.nutanix_audit_v2.audit.status
}

output "audit_operation_type" {
  value = data.nutanix_audit_v2.audit.operation_type
}

output "total_audits" {
  value = length(data.nutanix_audits_v2.audits.audits)
}
