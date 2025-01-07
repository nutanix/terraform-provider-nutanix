
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
  port     = var.nutanix_port
  insecure = true
}

resource "nutanix_vm_network_device_assign_ip_v2" "test" {
  vm_ext_id = "<VM_EXT_ID>"
  ext_id    = "<NETWORK_DEVICE_EXT_ID>"
  ip_address {
    value = "<IP_ADDRESS>"
  }
}