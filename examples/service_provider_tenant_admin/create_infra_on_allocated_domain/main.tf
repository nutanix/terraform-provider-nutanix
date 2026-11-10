terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

# Service Provider Tenant Admin — create infrastructure on a domain that has
# already been allocated to the tenant (project + image + VM on an existing VLAN).
#
# Note (KB-22890): Calls to Nutanix Central may return 404 when the SP FQDN Host
# header includes an explicit :443 port. This can affect Terraform and other API
# clients. Until a fix ships, add both <host> and <host>:443 to
# ALLOWED_CALLBACK_HOSTS on the Dion service and restart it.
# Details: https://portal.nutanix.com/page/documents/kbs/details?targetId=kA0VO000000EVaD0AW
provider "nutanix" {
  username                 = var.nutanix_username
  password                 = var.nutanix_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  domain_name              = var.domain_name
  tenant_name              = var.tenant_name
  port                     = 443
}

# Project created by the Service Provider domain allocation.
data "nutanix_projects_v2" "tenant_projects" {}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_storage_containers_v2" "storage_containers" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

# Tenant Admin cannot create VLAN subnets; reuse an allocated VLAN by name.
data "nutanix_subnets_v2" "vm_subnet" {
  filter = "name eq '${var.subnet_name}'"
}

locals {
  project_ext_id = data.nutanix_projects_v2.tenant_projects.projects[0].ext_id
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  storage_container_ext_id = data.nutanix_storage_containers_v2.storage_containers.storage_containers[0].ext_id
  subnet_ext_id            = data.nutanix_subnets_v2.vm_subnet.subnets[0].ext_id
}

resource "nutanix_images_v2" "image" {
  name                     = var.image_name
  description              = "Tenant admin image on allocated domain"
  type                     = "DISK_IMAGE"
  project_ext_id           = local.project_ext_id
  cluster_location_ext_ids = [local.cluster_ext_id]

  source {
    url_source {
      url = var.image_url
    }
  }
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = var.vm_name
  description          = "Tenant admin VM on allocated domain"
  num_cores_per_socket = 1
  num_sockets          = 1
  memory_size_bytes    = 4 * 1024 * 1024 * 1024

  cluster {
    ext_id = local.cluster_ext_id
  }

  project {
    ext_id = local.project_ext_id
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

output "project_ext_id" {
  value = local.project_ext_id
}

output "vm_ext_id" {
  value = nutanix_virtual_machine_v2.vm.id
}

output "image_ext_id" {
  value = nutanix_images_v2.image.id
}
