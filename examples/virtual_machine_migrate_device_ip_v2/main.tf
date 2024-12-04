
terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
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

# Assign Ip address \
resource "nutanix_vm_network_device_migrate_v2" "example" {
  vm_ext_id = "<VM_UUID>"
  ext_id    = "<NIC_UUID>"
  subnet {
    ext_id = "<SUBNET_UUID>"
  }
  migrate_type = "ASSIGN_IP"
  ip_address {
    value = "<IP_ADDRESS>"
  }
}

# release Ip address
resource "nutanix_vm_network_device_migrate_v2" "example" {
  vm_ext_id = "<VM_UUID>"
  ext_id    = "<NIC_UUID>"
  subnet {
    ext_id = "<SUBNET_UUID>"
  }
  migrate_type = "RELEASE_IP"

}