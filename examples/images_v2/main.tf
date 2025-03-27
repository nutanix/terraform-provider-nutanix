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
  port     = 9440
  insecure = true
}

resource "nutanix_images_v2" "img-1" {
  name        = "test-image"
  description = "img desc"
  type        = "ISO_IMAGE"
  source {
    url_source {
      url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
    }
  }
}

data "nutanix_clusters_v2" "clusters" {}

locals {
  clusterExtID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

resource "nutanix_images_v2" "img-2" {
  name        = "test-image"
  description = "img desc"
  type        = "DISK_IMAGE"
  source {
    url_source {
      url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
    }
  }
  cluster_location_ext_ids = [
    local.clusterExtID
  ]
}
