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


resource "nutanix_subnet_v2" "subnet" {
  name              = "tf-example-subnet"
  description       = "terraform test subnet to assign ip"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = 425
  is_external       = false
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "192.168.1.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "192.168.1.1"
      }
      pool_list {
        start_ip {
          value = "192.168.1.20"
        }
        end_ip {
          value = "192.168.1.60"
        }
      }
    }
  }
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "vm-example"
  description          = "create vm example"
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
      bus_type = "IDE"
      index    = 0
    }
  }
  nics {
    network_info {
      nic_type = "DIRECT_NIC"
      subnet {
        ext_id = nutanix_subnet_v2.subnet.id
      }
      vlan_mode = "ACCESS"
    }
  }
  power_state = "OFF"
}


# migrate the network device of the vm, assign ip
resource "nutanix_vm_network_device_migrate_v2" "assign-migrate" {
  vm_ext_id = nutanix_virtual_machine_v2.vm.id
  ext_id    = nutanix_virtual_machine_v2.vm.nics.0.ext_id
  subnet {
    ext_id = nutanix_subnet_v2.subnet.ext_id
  }
  migrate_type = "ASSIGN_IP"
   ip_address {
    value = "192.168.1.55"
  }
}


# migrate the network device of the vm, release ip
resource "nutanix_vm_network_device_migrate_v2" "release-migrate" {
  vm_ext_id = nutanix_virtual_machine_v2.vm.id
  ext_id    = nutanix_virtual_machine_v2.vm.nics.0.ext_id
  subnet {
    ext_id = nutanix_subnet_v2.subnet.ext_id
  }
  migrate_type = "RELEASE_IP"
}
