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

# pull the specified virtual machine data
data "nutanix_virtual_machines_v2" "ngt-vm" {
  filter = "name eq '${var.vm_name}'"
}

resource "nutanix_ngt_upgrade_v2" "upgrade-ngt" {
  ext_id = data.nutanix_virtual_machines_v2.ngt-vm.vms.0.ext_id

  reboot_preference {
    schedule_type = "IMMEDIATE"
  }
}
