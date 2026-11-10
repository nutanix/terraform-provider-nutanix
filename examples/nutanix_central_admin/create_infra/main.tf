terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

# Nutanix Central Admin — create infrastructure on a registered PC domain
# via nutanix_central_endpoint + domain_name (no tenant_name).
provider "nutanix" {
  username                 = var.nutanix_username
  password                 = var.nutanix_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  domain_name              = var.domain_name
  port                     = 443
}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_storage_containers_v2" "storage_containers" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

data "nutanix_subnets_v2" "vm_subnet" {
  filter = "name eq '${var.subnet_name}'"
}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  storage_container_ext_id = data.nutanix_storage_containers_v2.storage_containers.storage_containers[0].ext_id
  subnet_ext_id            = data.nutanix_subnets_v2.vm_subnet.subnets[0].ext_id
}

resource "nutanix_images_v2" "image" {
  name                     = var.image_name
  description              = "NC admin image on registered domain"
  type                     = "DISK_IMAGE"
  cluster_location_ext_ids = [local.cluster_ext_id]

  source {
    url_source {
      url = var.image_url
    }
  }
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = var.vm_name
  description          = "NC admin VM on registered domain"
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
              image_ext_id = nutanix_images_v2.image.id
            }
          }
        }
        storage_container {
          ext_id = local.storage_container_ext_id
        }
      }
    }
  }

  nics {
    nic_backing_info {
      virtual_ethernet_nic {
        is_connected = true
      }
    }
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = local.subnet_ext_id
        }
      }
    }
  }

  power_state = "ON"
}

output "cluster_ext_id" {
  value = local.cluster_ext_id
}

output "vm_ext_id" {
  value = nutanix_virtual_machine_v2.vm.id
}

output "image_ext_id" {
  value = nutanix_images_v2.image.id
}
