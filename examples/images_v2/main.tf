terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
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

# create image from vm disk source
resource "nutanix_images_v2" "example-2" {
  name = "example-image-2"
  type = "DISK_IMAGE"
  source {
    vm_disk_source {
      ext_id = resource.nutanix_virtual_machine_v2.test.disks.0.ext_id
    }
  }
  cluster_location_ext_ids = [
    local.cluster0
  ]
  depends_on = [nutanix_virtual_machine_v2.test]
}
