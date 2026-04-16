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

# pull subnet data
data "nutanix_subnets_v2" "vm-subnet" {
  filter = "name eq '${var.subnet_name}'"
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

# Deploy a template
resource "nutanix_deploy_templates_v2" "deploy-example" {
  ext_id            = nutanix_template_v2.example-1.id
  number_of_vms     = 1
  cluster_reference = local.cluster_ext_id
  override_vm_config_map {
    name                 = "example-tf-template-deploy"
    memory_size_bytes    = 4294967296
    num_sockets          = 2
    num_cores_per_socket = 1
    num_threads_per_core = 1
    nics {
      nic_backing_info {
        virtual_ethernet_nic {
          is_connected = true
          model        = "VIRTIO"
        }
      }
      nic_network_info {
        virtual_ethernet_nic_network_info {
          nic_type = "NORMAL_NIC"
          subnet {
            ext_id = data.nutanix_subnets_v2.vm-subnet.subnets[0].ext_id
          }
          vlan_mode                 = "ACCESS"
          should_allow_unknown_macs = false
        }
      }
    }
  }
}
