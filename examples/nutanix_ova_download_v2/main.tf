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

# List Prism Central
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
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


# Create Ova from Url
resource "nutanix_ova_v2" "ov-url-example" {
  name = "tf-ova-url-example"
  source {
    ova_url_source {
      url              = var.ova_url
      should_allow_insecure_url = true
    }
  }
}


// Create a new OVA using the object_lite_source
resource "nutanix_ova_v2" "ova-object-lite"{
  name = "tf-example-ova-object-lite-example"
  source {
    object_lite_source {
      key = var.ova_key
    }
  }
}

resource "nutanix_ova_download_v2" "download-ova" {
  ova_ext_id = nutanix_ova_v2.ov-vm-example.id
}

output "ova_file_path" {
  value = nutanix_ova_download_v2.download-ova.ova_file_path
}
