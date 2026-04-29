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

# Data source: List all cluster configs for a System-Defined Alert Policy
data "nutanix_sda_cluster_configs_v2" "all" {
  system_defined_policy_ext_id = var.system_defined_policy_ext_id
}

# Data source: Get a specific cluster config
data "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = var.system_defined_policy_ext_id
  ext_id                       = var.cluster_ext_id
}

# Resource: Manage a cluster config for a System-Defined Alert Policy
resource "nutanix_sda_cluster_config_v2" "example" {
  system_defined_policy_ext_id = var.system_defined_policy_ext_id
  ext_id                       = var.cluster_ext_id
  is_enabled                   = true
}

# Output the list of cluster configs
output "cluster_configs" {
  value = data.nutanix_sda_cluster_configs_v2.all.cluster_configs
}

# Output the specific cluster config
output "cluster_config" {
  value = data.nutanix_sda_cluster_config_v2.example
}
