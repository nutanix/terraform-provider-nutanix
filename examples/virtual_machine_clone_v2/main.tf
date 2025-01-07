#Here we will create a vm clone 
#the variable "" present in terraform.tfvars file.
#Note - Replace appropriate values of variables in terraform.tfvars file as per setup

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

#create a virtual machine clone
resource "nutanix_vm_clone_v2" "clone-vm" {
  vm_ext_id            = var.vm_uuid
  name                 = "clone-vm-123"
  num_sockets          = "2"
  num_cores_per_socket = "2"
  num_threads_per_core = "2"
  memory_size_bytes    = 4096
  guest_customization {
    config {
      sysprep {
        install_type = "PREPARED"
        sysprep_script {
          unattend_xml {
            value = ""
          }
        }
      }
    }
  }
  boot_config {
    legacy_boot {
      boot_device {
        boot_device_disk {
          disk_address {
            bus_type = "SCSI"
            index    = 1
          }
        }
        boot_device_nic {
          mac_address = ""
        }
      }
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }

  nics {
    ext_id = data.nutanix_subnet.subnet.id
    backing_info {
      model        = "VIRTIO"
      mac_address  = ""
      is_connected = "true"
      num_queues   = 1
    }
    network_info {}
  }

}
