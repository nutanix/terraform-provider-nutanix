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


# Create a new Template from done only using vm reference
resource "nutanix_template_v2" "example-1" {
  template_name        = "tf-example-template"
  template_description = "test create template from vm using terraform"
  template_version_spec {
    version_source {
      template_vm_reference {
        ext_id = "<VM_UUID>"
      }
    }
  }
}

# Deploy a template
resource "nutanix_deploy_templates_v2" "test" {
  ext_id            = resource.nutanix_template_v2.example-1.id
  number_of_vms     = 1
  cluster_reference = "<CLUSTER_UUID>"
  override_vm_config_map {
    name                 = "example-tf-template-deploy"
    memory_size_bytes    = 4294967296
    num_sockets          = 2
    num_cores_per_socket = 1
    num_threads_per_core = 1
    nics {
      backing_info {
        is_connected = true
        model        = "VIRTIO"
      }
      network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = " <SUBNET_UUID>"
        }
        vlan_mode                 = "ACCESS"
        should_allow_unknown_macs = false
      }
    }
  }
  depends_on = [
    resource.nutanix_template_v2.example-1
  ]
}
