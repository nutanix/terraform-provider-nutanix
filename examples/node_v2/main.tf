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

# Register and onboard a node for management by Foundation Central.
resource "nutanix_node_v2" "example" {
  manufacturer = var.node_manufacturer
  model        = var.node_model

  identifiers {
    type  = "SERIAL_NUMBER"
    value = var.node_serial
  }
}

# Fetch the node by its external ID (singular datasource).
data "nutanix_node_v2" "get_node" {
  ext_id = nutanix_node_v2.example.ext_id
}

# List all nodes managed by Foundation Central (plural datasource).
data "nutanix_nodes_v2" "list_nodes" {}

# Filtered list example.
data "nutanix_nodes_v2" "filtered_nodes" {
  limit  = 10
  filter = "manufacturer eq '${var.node_manufacturer}'"
}
