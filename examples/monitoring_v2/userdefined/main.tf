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

resource "nutanix_uda_policy_v2" "example" {
  title       = "example-uda-policy"
  entity_type = "VM"
  description = "Example User-Defined Alert policy"
  is_enabled  = true
  trigger_wait_period = 600

  trigger_conditions {
    condition {
      metric_name = "hypervisor_cpu_usage_ppm"
      operator    = "GREATER_THAN"
      threshold_value {
        int_value = 900000
      }
    }
    condition_type = "STATIC_THRESHOLD"
    severity_level = "CRITICAL"
  }
}

data "nutanix_uda_policy_v2" "example" {
  ext_id = nutanix_uda_policy_v2.example.id
}

data "nutanix_uda_policies_v2" "example" {
  depends_on = [nutanix_uda_policy_v2.example]
}

resource "nutanix_find_conflicting_uda_policies_v2" "example" {
  title       = "conflict-check-policy"
  entity_type = "VM"

  trigger_conditions {
    condition {
      metric_name = "hypervisor_cpu_usage_ppm"
      operator    = "GREATER_THAN"
      threshold_value {
        int_value = 900000
      }
    }
    condition_type = "STATIC_THRESHOLD"
    severity_level = "CRITICAL"
  }
}
