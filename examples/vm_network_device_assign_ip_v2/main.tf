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


data "nutanix_storage_containers_v2" "sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

resource "nutanix_subnet_v2" "vm-subnet" {
  name              = "example-subnet"
  description       = "terraform test subnet to assign ip"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = 576
  is_external       = false
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "192.168.11.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "192.168.11.1"
      }
      pool_list {
        start_ip {
          value = "192.168.11.2"
        }
        end_ip {
          value = "192.168.11.22"
        }
      }
    }
  }
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "example-vm"
  description          = "create vm to be cloned"
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
        disk_size_bytes = "1073741824"
        storage_container {
          ext_id = data.nutanix_storage_containers_v2.sc.storage_containers[0].ext_id
        }
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
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = nutanix_subnet_v2.vm-subnet.id
        }
        vlan_mode = "ACCESS"
      }
    }
  }
  power_state = "OFF"
}


resource "nutanix_vm_network_device_migrate_v2" "assign_ip" {
  vm_ext_id = nutanix_virtual_machine_v2.vm.id
  ext_id    = nutanix_virtual_machine_v2.vm.nics.0.ext_id
  subnet {
    ext_id = nutanix_subnet_v2.vm-subnet.id
  }
  migrate_type = "ASSIGN_IP"
  ip_address {
    value = "192.168.11.14"
  }
}


# release ip from the vm
resource "nutanix_vm_network_device_migrate_v2" "release_ip" {
  vm_ext_id = nutanix_virtual_machine_v2.vm.id
  ext_id    = nutanix_virtual_machine_v2.vm.nics.0.ext_id
  subnet {
    ext_id = nutanix_subnet_v2.vm-subnet.id
  }
  migrate_type = "RELEASE_IP"
}
