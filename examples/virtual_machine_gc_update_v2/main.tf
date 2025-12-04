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


data "nutanix_subnets_v2" "subnet" {
  filter = "name eq '${var.subnet_name}'"
}

data "nutanix_storage_containers_v2" "ngt-sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
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
          ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
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
    network_info {
      nic_type = "NORMAL_NIC"
      subnet {
        ext_id = data.nutanix_subnets_v2.subnet.subnets[0].ext_id
      }
      vlan_mode = "ACCESS"
    }
  }
  power_state = "OFF"
}

#update a virtual machine guest customization for next boot
resource "nutanix_vm_gc_update_v2" "gc-update" {
  ext_id = nutanix_virtual_machine_v2.vm.id
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


