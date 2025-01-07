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
  cluster0 = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

#create a virtual machine with minium configuration
resource "nutanix_virtual_machine_v2" "example-1" {
  name                 = "vm-example"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
  }
  power_state = "ON"
}

# create virtual machine with disk
resource "nutanix_virtual_machine_v2" "example-2" {
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
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
          ext_id = "<storage_container_uuid>"
        }
      }
    }
  }
  power_state = "ON"
}

# create virtual machine with disk data source
resource "nutanix_virtual_machine_v2" "example-3" {
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
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
          ext_id = "<storage_container_uuid>"
        }
        data_source {
          reference {
            vm_disk_reference {
              disk_address {
                bus_type = "SCSI"
                index    = 0
              }
              vm_reference {
                ext_id = resource.nutanix_virtual_machine_v2.example-2.id
              }
            }
          }
        }
      }
    }
  }
  power_state = "ON"
  depends_on  = [resource.nutanix_virtual_machine_v2.testWithDisk]
  lifecycle {
    ignore_changes = [
      disks.0.backing_info.0.vm_disk.0.data_source
    ]
  }
}

# create virtual machine with nics
resource "nutanix_virtual_machine_v2" "example-4" {
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
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
          ext_id = "<storage_container_uuid>"
        }
      }
    }
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = "<subnet_uuid>"
      }
      vlan_mode = "ACCESS"
    }
  }
  power_state = "ON"
}

# create virtual machine with legacy boot config 
resource "nutanix_virtual_machine_v2" "example-5" {
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
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
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
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
              image_ext_id = "<image_uuid>"
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
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
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
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster0
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
          ext_id = "<storage_container_uuid>"
        }
      }
    }
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = "<subnet_uuid>"
      }
      vlan_mode = "ACCESS"
    }
  }
  guest_customization {
    config {
      cloud_init {
        cloud_init_script {
          user_data {
            value = "echo 'Hello, World!'"
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
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 2
  cluster {
    ext_id = local.cluster0
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }
  categories {
    ext_id = "<category_1_uuid>"
  }
  categories {
    ext_id = "<category_2_uuid>"
  }
  categories {
    ext_id = "<category_3_uuid>"
  }
  power_state = "ON"
}


# create virtual machine with gpus 
resource "nutanix_virtual_machine_v2" "example-10" {
  name                = "example-vm"
  um_cores_per_socket = 1
  num_sockets         = 2
  cluster {
    ext_id = local.cluster0
  }

  gpus {
    device_id = "<device_id>"
    mode      = "PASSTHROUGH"
    vendor    = "<vendor>"
  }
  power_state = "ON"
}

# create virtual machine with serial port
resource "nutanix_virtual_machine_v2" "example-11" {
  name                 = "example-vm"
  num_cores_per_socket = 1
  num_sockets          = 2
  cluster {
    ext_id = local.cluster0
  }
  nics {
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
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



# list all virtual machines
data "nutanix_virtual_machines_v2" "vms" {}

# list vms with filter
data "nutanix_virtual_machines_v2" "filtered_vms" {
  filter = "name eq 'example-vm'"
}

# get vm by id
data "nutanix_virtual_machine_v2" "vm" {
  ext_id = nutanix_virtual_machine_v2.example-1.id
}
