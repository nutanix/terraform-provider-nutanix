terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.3.2"
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

# List Clusters
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  limit  = 1
}

data "nutanix_subnets_v2" "subnets" {
  filter = "name eq 'vlan.800'"
}


locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
  subnet_ext_id  = data.nutanix_subnets_v2.subnets.subnets[0].ext_id
}

# Create VM with some specific requirements
resource "nutanix_virtual_machine_v2" "vm-example" {
  name              = "vm-example"
  num_sockets       = 2
  memory_size_bytes = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.cluster_ext_id
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
}

# Create Ova from the VM
resource "nutanix_ova_v2" "ov-vm-example" {
  name = "tf-ova-vm-example"
  source {
    ova_vm_source {
      vm_ext_id        = nutanix_virtual_machine_v2.vm-example.id
      disk_file_format = "QCOW2"
    }
  }
}

resource "nutanix_ova_vm_deploy_v2" "vm-from-ova" {
  ext_id = nutanix_ova_v2.ov-vm-example.id
  override_vm_config {
    name              = "${nutanix_virtual_machine_v2.vm-example.name}-from-ova"
    memory_size_bytes = 8 * 1024 * 1024 * 1024 # 8 GiB
    nics {
      backing_info {
        is_connected = true
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = local.subnet_ext_id
        }
        vlan_mode     = "TRUNK"
        trunked_vlans = ["1"]
      }
    }
  }
  cluster_location_ext_id = local.cluster_ext_id
}
