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

data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}


# pull storage container data
data "nutanix_storage_containers_v2" "sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

# pull subnet data
data "nutanix_subnets_v2" "vm-subnet" {
  filter = "name eq '${var.subnet_name}'"
}

# pull image data
data "nutanix_images_v2" "vm-image" {
  filter = "name eq '${var.image_name}'"
  limit  = 1
}

#pull all categories data
data "nutanix_categories_v2" "categories-list" {}


#create a virtual machine with minium configuration
resource "nutanix_virtual_machine_v2" "example-1" {
  name                 = "vm-example-1"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster_ext_id
  }
  power_state = "ON"
}

# create virtual machine with disk
resource "nutanix_virtual_machine_v2" "example-2" {
  name                 = "example-vm-2"
  num_cores_per_socket = 1
  num_sockets          = 1
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
        disk_size_bytes = "1073741824"
        storage_container {
          ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
        }
      }
    }
  }
  power_state = "ON"
}

# create virtual machine with disk data source
resource "nutanix_virtual_machine_v2" "example-3" {
  name                 = "example-vm-3"
  num_cores_per_socket = 1
  num_sockets          = 1
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
        disk_size_bytes = 1073741824
        storage_container {
          ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
        }
        data_source {
          reference {
            vm_disk_reference {
              disk_address {
                bus_type = "SCSI"
                index    = 0
              }
              vm_reference {
                ext_id = nutanix_virtual_machine_v2.example-2.id
              }
            }
          }
        }
      }
    }
  }
  power_state = "ON"
  lifecycle {
    ignore_changes = [
      disks.0.backing_info.0.vm_disk.0.data_source
    ]
  }
}

# create virtual machine with nics
resource "nutanix_virtual_machine_v2" "example-4" {
  name                 = "example-vm-4"
  num_cores_per_socket = 1
  num_sockets          = 1
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
        disk_size_bytes = "1073741824"
        storage_container {
          ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
        }
      }
    }
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.vm-subnet.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }
  power_state = "ON"
}

# create virtual machine with legacy boot config
resource "nutanix_virtual_machine_v2" "example-5" {
  name                 = "example-vm-5"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster_ext_id
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  power_state = "ON"
}


# create virtual machine with legacy boot device
resource "nutanix_virtual_machine_v2" "example-6" {
  name                 = "example-vm-6"
  num_cores_per_socket = 1
  num_sockets          = 1
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
              image_ext_id = data.nutanix_images_v2.vm-image.images[0].ext_id
            }
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
            index    = 0
          }
        }
      }
    }
  }
  power_state = "ON"
}


# create virtual machine with cdrom
resource "nutanix_virtual_machine_v2" "example-7" {
  name                 = "example-vm-7"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster_ext_id
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  cd_roms {
    disk_address {
      bus_type = "SATA"
      index    = 0
    }
  }
}

# create virtual machine with guest customization
resource "nutanix_virtual_machine_v2" "example-8" {
  name                 = "example-vm-8"
  num_cores_per_socket = 1
  num_sockets          = 1
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
        disk_size_bytes = "1073741824"
        storage_container {
          ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
        }
      }
    }
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.vm-subnet.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }
  guest_customization {
    config {
      cloud_init {
        cloud_init_script {
          user_data {
            value = base64encode("echo 'Hello, World!'")
          }
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
      guest_customization, cd_roms
    ]
  }
}


# create virtual machine with categories
resource "nutanix_virtual_machine_v2" "example-9" {
  name                 = "example-vm-9"
  num_cores_per_socket = 1
  num_sockets          = 2
  cluster {
    ext_id = local.cluster_ext_id
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.vm-subnet.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }
  categories {
    ext_id = data.nutanix_categories_v2.categories-list.categories[0].ext_id
  }
  categories {
    ext_id = data.nutanix_categories_v2.categories-list.categories[1].ext_id
  }
  categories {
    ext_id = data.nutanix_categories_v2.categories-list.categories[2].ext_id
  }
  power_state = "ON"
}


# create virtual machine with serial port
resource "nutanix_virtual_machine_v2" "example-10" {
  name                 = "example-vm-10"
  num_cores_per_socket = 1
  num_sockets          = 2
  cluster {
    ext_id = local.cluster_ext_id
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.vm-subnet.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }
  serial_ports {
    index        = 2
    is_connected = true
  }
  power_state = "ON"
}


# create virtual machine with gpus
resource "nutanix_virtual_machine_v2" "example-11" {
  name                 = "example-vm-11"
  num_cores_per_socket = 1
  num_sockets          = 2
  cluster {
    ext_id = local.cluster_ext_id
  }

  gpus {
    device_id = "<device_id>"
    mode      = "PASSTHROUGH"
    vendor    = "<vendor>"
  }
  power_state = "ON"
}


# create virtual machine with gest customization sysprep script
resource "nutanix_virtual_machine_v2" "example-12" {
  name                 = "example-vm-12"
  num_cores_per_socket = 1
  num_sockets          = 1
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
        disk_size_bytes = "1073741824"
        storage_container {
          ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
        }
      }
    }
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.vm-subnet.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }
  guest_customization {
    config {
      sysprep {
        install_type = "PREPARED"
        sysprep_script {
          unattend_xml {
            value = base64encode(file(var.unattend_xml_path)) # encoded unattend_xml file value or base64 encoded string value
          }
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
      guest_customization, cd_roms
    ]
  }
}

resource "nutanix_virtual_machine_v2" "example-13" {
  name                 = "example-13"
  description          = "vm example with uefi boot"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster_ext_id
  }
  boot_config {
    uefi_boot {
      boot_order = ["NETWORK", "DISK", "CDROM", ]
    }
  }
  power_state = "OFF"
}


# list all virtual machines
data "nutanix_virtual_machines_v2" "vms" {}

# list vms with filter
data "nutanix_virtual_machines_v2" "filtered_vms" {
  filter = "name eq '${nutanix_virtual_machine_v2.example-1.name}'"
}

# list vms with limit and pagination
data "nutanix_virtual_machines_v2" "paginated_vms" {
  page  = 2
  limit = 4
}

# get vm by id
data "nutanix_virtual_machine_v2" "vm" {
  ext_id = nutanix_virtual_machine_v2.example-1.id
}
