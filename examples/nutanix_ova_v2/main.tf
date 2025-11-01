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
      url                       = var.ova_url
      should_allow_insecure_url = true
    }
  }
  cluster_location_ext_ids = [local.cluster_ext_id]
}


// Create a new OVA using the object_lite_source
resource "nutanix_ova_v2" "ova-object-lite" {
  name = "tf-example-ova-object-lite-example"
  source {
    object_lite_source {
      key = var.ova_key
    }
  }
  cluster_location_ext_ids = [local.cluster_ext_id]
}


// list all Ovas
data "nutanix_ova_v2" "ovas-list" {}

// limit, page and ordered by ovas
data "nutanix_ova_v2" "ovas-limit" {
  limit      = 20
  page       = 1
  order_by   = "createTime"
  depends_on = [nutanix_ova_v2.ov-vm-example]
}


// filtered Ovas examples
data "nutanix_ova_v2" "filtered-name" {
  filter     = "name eq 'tf-ova-vm-example'"
  depends_on = [nutanix_ova_v2.ov-vm-example]
}

data "nutanix_ova_v2" "filtered-disk-format" {
  filter     = "diskFormat eq Vmm.Content.OvaDiskFormat'QCOW2'"
  depends_on = [nutanix_ova_v2.ov-vm-example]
}

data "nutanix_ova_v2" "filtered-parent-vm" {
  filter     = "parentVm eq 'LinuxServer_VM'"
  depends_on = [nutanix_ova_v2.ov-vm-example]
}

data "nutanix_ova_v2" "filtered-size" {
  filter     = "sizeBytes eq ${4 * 1024 * 1024 * 1024}"
  depends_on = [nutanix_ova_v2.ov-vm-example]
}

// Get ova details by ext id
data "nutanix_ova_v2" "ova" {
  ext_id = nutanix_ova_v2.ov-vm-example.id
}
