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

# pull cluster data
data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# pull image data, provide the image name that supports NGT
data "nutanix_images_v2" "ngt-image" {
  filter = "name eq '${var.image_name}'"
  limit  = 1
}

# pull storage container data
data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

# pull subnet data
data "nutanix_subnets_v2" "subnet" {
  filter = "name eq '${var.subnet_name}'"
}

# create a VM to install NGT
resource "nutanix_virtual_machine_v2" "ngt-vm" {
  name                 = "vm-example-ngt"
  description          = "vm to test ngt installation"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_bytes    = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }

  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        data_source {
          reference {
            image_reference {
              image_ext_id = data.nutanix_images_v2.ngt-image.images[0].ext_id
            }
          }
        }
        disk_size_bytes = 20 * 1024 * 1024 * 1024
      }
    }
  }

  cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 0
    }
  }

  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }

  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  power_state = "ON"

  lifecycle {
    ignore_changes = [guest_tools]
  }

  depends_on = [data.nutanix_clusters_v2.clusters, data.nutanix_images_v2.ngt-image, data.nutanix_storage_containers_v2.ngt-sc]
}

# install NGT on the VM
resource "nutanix_ngt_installation_v2" "install-ngt" {
  ext_id = nutanix_virtual_machine_v2.ngt-vm.id
  credential {
    username = var.username
    password = var.password
  }
  reboot_preference {
    schedule_type = "SKIP"
  }
  capablities = ["VSS_SNAPSHOT", "SELF_SERVICE_RESTORE"]
}


# shutdown the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmShuts" {
  ext_id     = nutanix_virtual_machine_v2.ngt-vm.id
  action     = "shutdown"
  depends_on = [nutanix_ngt_installation_v2.install-ngt]
}

# restart the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmReboot" {
  ext_id     = nutanix_virtual_machine_v2.ngt-vm.id
  action     = "reboot"
  depends_on = [nutanix_ngt_installation_v2.install-ngt]

}

# guest-shutdown the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmGuestShuts" {
  ext_id     = nutanix_virtual_machine_v2.ngt-vm.id
  action     = "guest_shutdown"
  depends_on = [nutanix_ngt_installation_v2.install-ngt]

}

# guest-restart the virtual machine
resource "nutanix_vm_shutdown_action_v2" "vmGuestReboot" {
  ext_id     = nutanix_virtual_machine_v2.ngt-vm.id
  action     = "guest_reboot"
  depends_on = [nutanix_ngt_installation_v2.install-ngt]

}
