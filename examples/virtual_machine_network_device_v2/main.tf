
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
  port     = var.nutanix_port
  insecure = true
}

resource "nutanix_vm_network_device_v2" "nic" {
  vm_ext_id = var.vm_uuid
  network_info {
    nic_type = "DIRECT_NIC"
    subnet {
      ext_id = var.subnet_uuid
    }
    ipv4_config {
      should_assign_ip = true
      ip_address {
        value         = "10.51.144.215"
        prefix_length = 32
      }
    }
  }
}