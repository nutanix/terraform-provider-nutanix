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

# Run System-Defined Checks on a cluster
resource "nutanix_run_system_defined_checks_v2" "example" {
  cluster_ext_id                              = "00000000-0000-0000-0000-000000000000"
  should_run_all_checks                       = true
  should_send_report_to_configured_recipients = false
}

# Manage cluster-specific configuration for a System-Defined Alert Policy
resource "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id                       = "00000000-0000-0000-0000-000000000000"
}

# Get a list of all System-Defined Alert Policies
data "nutanix_sda_policies_v2" "policies-list" {}

# Get details of a specific System-Defined Alert Policy
data "nutanix_sda_policy_v2" "get-policy" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}

# Get cluster configs for a System-Defined Alert Policy
data "nutanix_sda_cluster_configs_v2" "cluster-configs-list" {
  system_defined_policy_ext_id = "00000000-0000-0000-0000-000000000000"
}

# Get a specific cluster config for a System-Defined Alert Policy
data "nutanix_sda_cluster_config_v2" "get-cluster-config" {
  system_defined_policy_ext_id = "00000000-0000-0000-0000-000000000000"
  ext_id                       = "00000000-0000-0000-0000-000000000000"
}
