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



# Add Node to Cluster
resource "nutanix_cluster_add_node_v2" "cluster_node" {
  cluster_ext_id = "00000000-0000-0000-0000-000000000000"
  node_params {
    block_list {
      node_list {
        node_uuid                = "00000000-0000-0000-0000-000000000000"
        block_uuid               = "00000000-0000-0000-0000-000000000000"
        node_position            = "<node_position>"
        hypervisor_type          = "XEN"
        is_robo_mixed_hypervisor = true
        hypervisor_hostname      = "<hypervisor_hostname>"
        hypervisor_version       = "9.9.99"
        nos_version              = "9.9.99"
        ipmi_ip {
          ipv4 {
            value = "10.0.0.1"
          }
        }
        digital_certificate_map_list {
          key   = "key"
          value = "value"
        }
        model = "<model>"
      }
      should_skip_host_networking = false
    }
    node_list {
      node_uuid                = "00000000-0000-0000-0000-000000000000"
      block_uuid               = "00000000-0000-0000-0000-000000000000"
      node_position            = "<node_position>"
      hypervisor_type          = "XEN"
      is_robo_mixed_hypervisor = true
      hypervisor_hostname      = "<hypervisor_hostname>"
      hypervisor_version       = "9.9.99"
      nos_version              = "9.9.99"
      ipmi_ip {
        ipv4 {
          value = "10.0.0.1"
        }
      }

    }
    bundle_info {
      name = "<name>"
    }
  }
  config_params {
    should_skip_discovery = false
    should_skip_imaging   = true
    is_nos_compatible     = true
    target_hypervisor     = "<target_hypervisor>"
  }
  should_skip_add_node          = false
  should_skip_pre_expand_checks = true

  remove_node_params {
    extra_params {
      should_skip_upgrade_check = false
      skip_space_check          = false
      should_skip_add_check     = false
    }
    should_skip_remove    = false
    should_skip_prechecks = false
  }

}
