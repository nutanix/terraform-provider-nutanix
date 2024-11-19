terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}



# Remove Nodes from Cluster
resource "nutanix_cluster_remove_node_v2" "cluster_node"{
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  node_uuids = ["00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"]
  should_skip_remove = false
  should_skip_prechecks = true
  extra_params {
    should_skip_upgrade_check = true
    skip_space_check = true
    should_skip_add_check = false
  }
}