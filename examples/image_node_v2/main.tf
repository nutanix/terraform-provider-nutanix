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

# Variant 1 - host_installation: install a patched hypervisor / host OS image on the node.
resource "nutanix_image_node_v2" "host_installation" {
  ext_id = var.node_ext_id

  host_installation {
    patched_image_ext_id = var.patched_image_ext_id
  }
}

# Variant 2 - aos_installation: install AOS on the node.
# resource "nutanix_image_node_v2" "aos_installation" {
#   ext_id = var.node_ext_id
#
#   aos_installation {
#     aos_image_ext_id = var.aos_image_ext_id
#     cvm_memory_gb    = 32
#
#     management_network {
#       vlan_id = 0
#
#       ip {
#         ipv4 {
#           value         = "192.0.2.10"
#           prefix_length = 24
#         }
#       }
#
#       gateway {
#         ipv4 {
#           value         = "192.0.2.1"
#           prefix_length = 24
#         }
#       }
#     }
#   }
# }
