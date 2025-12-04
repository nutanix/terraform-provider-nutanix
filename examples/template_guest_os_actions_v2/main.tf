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


data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_virtual_machine_v2" "vm" {
  name                 = "vm-example"
  description          = "create vm example"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.cluster_ext_id
  }
  power_state = "OFF"
}


# Create a new Template from done only using vm reference
resource "nutanix_template_v2" "example-1" {
  template_name        = "tf-example-template"
  template_description = "test create template from vm using terraform"
  template_version_spec {
    version_source {
      template_vm_reference {
        ext_id = nutanix_virtual_machine_v2.vm.id
      }
    }
  }
}


# initiate the template
resource "nutanix_template_guest_os_actions_v2" "example-initiate" {
  ext_id     = nutanix_template_v2.example-1.id
  action     = "initiate"
  version_id = nutanix_template_v2.example-1.template_version_spec.0.ext_id
}

resource "nutanix_template_guest_os_actions_v2" "example-cancel" {
  ext_id     = nutanix_template_v2.example-1.id
  action     = "cancel"
  depends_on = [nutanix_template_guest_os_actions_v2.example-initiate]
}
