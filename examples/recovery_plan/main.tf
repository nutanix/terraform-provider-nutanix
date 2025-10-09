terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.9.5"
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

# Example 1: Create a recovery plan with stage list
resource "nutanix_recovery_plan" "recovery_plan_stage_list" {
  name        = "tf-example-recovery-plan-stage-list"
  description = "This is a sample recovery plan with stage list for demonstration purposes."
  stage_list {
    stage_work {
      recover_entities {
        entity_info_list {
          categories {
            name  = "Environment"
            value = "Dev"
          }
        }
      }
    }
    stage_uuid      = "ab788130-0820-4d07-a1b5-b0ba4d3a4254"
    delay_time_secs = 0
  }
  parameters {}
}


# Example 2: Create a recovery plan with stage list and network mapping


data "nutanix_clusters" "clusters" {}

locals {
  clusterUUID = (data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
  ? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid)
}

resource "nutanix_virtual_machine" "vm1" {
  name         = "test-dou-vm"
  cluster_uuid = local.clusterUUID

  boot_device_order_list = ["DISK", "CDROM"]
  boot_type              = "LEGACY"
  num_vcpus_per_socket   = 1
  num_sockets            = 1
  memory_size_mib        = 186

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }

  categories {
    name  = "Environment"
    value = "Staging"
  }
}

resource "nutanix_recovery_plan" "recovery_plan_with_network" {
  name        = "recovery-plan-example-with-network"
  description = "This is a sample recovery plan with network mapping for demonstration purposes."
  stage_list {
    stage_work {
      recover_entities {
        entity_info_list {
          any_entity_reference_name = nutanix_virtual_machine.vm1.name
          any_entity_reference_kind = "vm"
          any_entity_reference_uuid = nutanix_virtual_machine.vm1.id
        }
      }
    }
    stage_uuid      = "ab788130-0820-4d07-a1b5-b0ba4d3a4254"
    delay_time_secs = 0
  }
  parameters {
    network_mapping_list {
      availability_zone_network_mapping_list {
        availability_zone_url = var.source_az_url
        recovery_network {
          name = "vlan.800"
          subnet_list {
            gateway_ip                  = "10.38.2.129"
            prefix_length               = 24
            external_connectivity_state = "DISABLED"
          }
        }
        test_network {
          name = "vlan.800"
          subnet_list {
            gateway_ip                  = "192.168.0.1"
            prefix_length               = 24
            external_connectivity_state = "DISABLED"
          }
        }
        cluster_reference_list {
          kind = "cluster"
          uuid = var.source_cluster_uuid
        }
      }
      availability_zone_network_mapping_list {
        availability_zone_url = var.target_az_url
        recovery_network {
          name = "vlan.800"
          subnet_list {
            gateway_ip                  = "10.38.4.65"
            prefix_length               = 24
            external_connectivity_state = "DISABLED"
          }
        }
        test_network {
          name = "vlan.800"
          subnet_list {
            gateway_ip                  = "192.168.0.1"
            prefix_length               = 24
            external_connectivity_state = "DISABLED"
          }
        }
        cluster_reference_list {
          kind = "cluster"
          uuid = var.target_cluster_uuid
        }
      }
    }
  }
}


# Get Recovery Plan by UUID
data "nutanix_recovery_plan" "get_recovery_plan" {
  recovery_plan_id = nutanix_recovery_plan.recovery_plan_with_network.id
}

# List Recovery Plans
data "nutanix_recovery_plans" "list_recovery_plans" {
  depends_on = [nutanix_recovery_plan.recovery_plan_stage_list, nutanix_recovery_plan.recovery_plan_with_network]
}
