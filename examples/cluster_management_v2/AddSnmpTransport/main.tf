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

resource "nutanix_add_snmp_transport_v2" "example" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  port           = 162
  protocol       = "UDP"
}
