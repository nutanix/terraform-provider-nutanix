terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.0.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  insecure = var.nutanix_insecure
  port     = var.nutanix_port
}

# Fetch list of all audits
data "nutanix_audits_v2" "all_audits" {}

# Output the first audit ext_id if available
output "first_audit_ext_id" {
  value = length(data.nutanix_audits_v2.all_audits.audits) > 0 ? data.nutanix_audits_v2.all_audits.audits[0].ext_id : ""
}

# Fetch a specific audit by ext_id
# Replace the ext_id below with an actual audit ext_id from your environment
data "nutanix_audit_v2" "specific_audit" {
  ext_id = length(data.nutanix_audits_v2.all_audits.audits) > 0 ? data.nutanix_audits_v2.all_audits.audits[0].ext_id : ""
}

# Output audit details
output "audit_details" {
  value = {
    ext_id             = data.nutanix_audit_v2.specific_audit.ext_id
    audit_type         = data.nutanix_audit_v2.specific_audit.audit_type
    service_name       = data.nutanix_audit_v2.specific_audit.service_name
    creation_time      = data.nutanix_audit_v2.specific_audit.creation_time
    operation_type     = data.nutanix_audit_v2.specific_audit.operation_type
    status             = data.nutanix_audit_v2.specific_audit.status
    message            = data.nutanix_audit_v2.specific_audit.message
  }
}
