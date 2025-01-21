terraform{
    required_providers {
        nutanix = {
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

// Create Protection Policy - Synchronous
resource "nutanix_protection_policy_v2" "pp_sync"{
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

  category_ids = ["<category_ids>"]
}



// Create Protection Policy - Linear Retention
resource "nutanix_protection_policy_v2" "pp_liner"{
  name     = "pp_example_2"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      retention {
        linear_retention {
          local  = 1
          remote = 1
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      retention {
        linear_retention {
          local  = 1
          remote = 1
        }
      }
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

  category_ids = ["<category_ids>"]
}

// Create Protection Policy - Auto Rollup Retention
resource "nutanix_protection_policy_v2" "pp_auto"{
  name     = "pp_example_3"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 20
      start_time = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      start_time = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
        }
      }
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

  category_ids = ["<category_ids>"]
}


// List Protection Policies
data "nutanix_protection_policies_v2" "protectiojn-policies" {}

// with filter
data "nutanix_protection_policies_v2" "pp-filter" {
  filter = "name eq 'pp_example_2'"
}

// with limit
data "nutanix_protection_policies_v2" "pp-limit" {
  limit = 4
}


// get protection policy by ext id
data "nutanix_protection_policy_v2" "pp-ex1" {
  ext_id = nutanix_protection_policy_v2.pp_liner.id
}