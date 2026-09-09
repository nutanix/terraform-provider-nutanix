terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Variant 1 - configure_server: configure a single server in the hardware provider.
resource "nutanix_configure_node_v2" "configure_server" {
  configure_server {
    node_ext_id = var.node_ext_id
    group_id    = var.group_id

    server_settings_config {
      network_adaptor_fec_mode    = "OFF"
      server_identity_pool_ext_id = var.server_identity_pool_ext_id
    }
  }
}

# Variant 2 - un_configure_server: remove the configuration from a node.
# resource "nutanix_configure_node_v2" "un_configure_server" {
#   un_configure_server {
#     node_ext_id = var.node_ext_id
#   }
# }

# Variant 3 - pre_cluster_config: pre-cluster operations on a set of nodes.
# resource "nutanix_configure_node_v2" "pre_cluster" {
#   pre_cluster_config {
#     nodes {
#       ext_id = var.node_ext_id
#     }
#   }
# }

# Variant 4 - post_cluster_config: post-cluster operations on a set of nodes.
# resource "nutanix_configure_node_v2" "post_cluster" {
#   post_cluster_config {
#     cluster_ext_id = var.cluster_ext_id
#     nodes {
#       ext_id = var.node_ext_id
#     }
#   }
# }
