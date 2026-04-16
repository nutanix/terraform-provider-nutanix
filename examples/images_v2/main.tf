terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}


#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}
#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}
# create image from vm disk source
data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# create a image from iso source
resource "nutanix_images_v2" "example-1" {
  name = "example-image-1"
  type = "ISO_IMAGE"
  source {
    url_source {
      url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
    }
  }
}



data "nutanix_storage_containers_v2" "sc" {
  filter = "clusterExtId eq '${local.cluster_ext_id}'"
  limit  = 1
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "tf-example-vm-disk"
  description          = "desc vm"
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
      }
    }
  }
  power_state = "OFF"
}

resource "nutanix_images_v2" "image-vm-disk" {
  name        = "tf-example-image-vm-disk"
  description = "desc image"
  type        = "DISK_IMAGE"
  source {
    vm_disk_source {
      ext_id = nutanix_virtual_machine_v2.vm.disks.0.ext_id
    }
  }
  cluster_location_ext_ids = [
    local.cluster_ext_id
  ]
}


# Create image using object lite source
resource "nutanix_images_v2" "object-liteStore-img" {
  name        = "image-object-lite-example"
  description = "Image created from object store"
  type        = "DISK_IMAGE"
  source {
    object_lite_source {
      key = var.lite_source_key
    }
  }
  lifecycle {
    ignore_changes = [
      source
    ]
  }
}


# pull all images
data "nutanix_images_v2" "image-list" {
  depends_on = [nutanix_images_v2.example-1, nutanix_images_v2.image-vm-disk]
}

# pull image with filter
data "nutanix_images_v2" "image-filtered" {
  filter = "name eq '${nutanix_images_v2.example-1.name}'"
}

# pull image with pagination and limit
data "nutanix_images_v2" "images-paginated" {
  page       = 0
  limit      = 13
  depends_on = [nutanix_images_v2.example-1, nutanix_images_v2.image-vm-disk]

}

# get image by id
data "nutanix_image_v2" "image-by-id" {
  ext_id = nutanix_images_v2.example-1.id
}
