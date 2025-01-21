terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "2.0.0"
        }
      nutanix-2 = {
        source = "nutanix/nutanix"
        version = "2.0.0"
      }
    }
}

#defining nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}

# we can add another nutanix provider and setup it with the remote provider configuration
# and specify the provider name in the promote_protected_resource_v2 resource block.

provider "nutanix-2" {
  username = var.nutanix_remote_username
  password = var.nutanix_remote_password
  endpoint = var.nutanix_remote_endpoint
  port = 9440
  insecure = true
}

resource "nutanix_protection_policy_v2" "pp_1"{
  provider = nutanix
  name     = "pp_example_1"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }

  replication_locations {
    domain_manager_ext_id = "domain_manager_ext_id_local"
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = "domain_manager_ext_id_remote"
    label                 = "target"
    is_primary            = false
  }

  category_ids = ["<category_id>"]
}

resource "nutanix_virtual_machine_v2" "vm" {
  provider = nutanix
  name                 = "%[2]s"
  description          = "%[3]s"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
  }
  categories {
    ext_id = "<category_id>"
  }
  power_state = "OFF"
  depends_on = [nutanix_protection_policy_v2.test]
}

# wait some time for the VM to be created to be protected
# you need add delay

resource "nutanix_promote_protected_resource_v2" "pp-vm"{
  provider = nutanix-2 # specify the provider name, to promote the protected resource in the remote site
  ext_id   = nutanix_virtual_machine_v2.vm.id
}

