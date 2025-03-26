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

# List Prism Central
data "nutanix_clusters_v2" "pc" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  pcExtID      = data.nutanix_clusters_v2.pc.cluster_entities[0].ext_id
}

# Example here, we are enabling Auto Inventory, adding Auto Inventory Schedule and enabling auto upgrade
resource "nutanix_lcm_config_v2" "lcm-configuration-update" {
    x_cluster_id = local.pcExtID
    is_auto_inventory_enabled = true
	auto_inventory_schedule = "16:30"
    has_module_auto_upgrade_enabled = true
}

# Read LCM config.
data "nutanix_lcm_config_v2" "get-lcm-configuration" {
    x_cluster_id = local.pcExtID
    depends_on = [nutanix_lcm_config_v2.lcm-configuration-update]
}

