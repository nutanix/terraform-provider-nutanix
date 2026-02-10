terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
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

resource "nutanix_network_function_v2" "nf" {
  name                    = var.network_function_name
  description             = "Network function managed by Terraform"
  high_availability_mode  = "ACTIVE_PASSIVE"
  failure_handling        = "NO_ACTION"
  traffic_forwarding_mode = "INLINE"

  data_plane_health_check_config {
    failure_threshold = 2
    interval_secs     = 5
    success_threshold = 2
    timeout_secs      = 5
  }

  nic_pairs {
    ingress_nic_reference = var.ingress_nic_reference
    egress_nic_reference  = var.egress_nic_reference
    is_enabled            = true
  }
}

data "nutanix_network_function_v2" "nf" {
  ext_id = nutanix_network_function_v2.nf.ext_id
}

data "nutanix_network_functions_v2" "nfs" {
  filter = "name eq '${var.network_function_name}'"
}

