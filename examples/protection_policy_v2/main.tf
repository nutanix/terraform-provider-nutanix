terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "2.1.0"
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


# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

# list Clusters
data "nutanix_clusters_v2" "clusters" {}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][
  0
  ]
}
# Synchronous Protection Policy
# Create Category
resource "nutanix_category_v2" "synchronous-pp-category" {
  key = "category-synchronous-protection-policy"
  value = "category_synchronous_protection_policy"
}

resource "nutanix_protection_policy_v2" "synchronous-protection-policy"{
  name        = "synchronous_protection_policy"

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
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17" # Remote Domain Manager UUID
    label                 = "target"
    is_primary            = false
  }

  category_ids = [nutanix_category_v2.synchronous-pp-category.id]
}


# Linear Retention Protection Policy
# Create Category
resource "nutanix_category_v2" "linear-retention-pp-category" {
  key = "category-linear-retention-protection-policy"
  value = "category_linear_retention_protection_policy"
}

resource "nutanix_protection_policy_v2" "linear-retention-protection-policy" {
  name = "linear-retention-protection-policy"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds = 7200
      recovery_point_type                   = "CRASH_CONSISTENT"
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
      recovery_point_objective_time_seconds = 7200
      recovery_point_type                   = "CRASH_CONSISTENT"
      retention {
        linear_retention {
          local  = 1
          remote = 1
        }
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17" # Remote Domain Manager UUID
    label      = "target"
    is_primary = false
  }

  category_ids = [nutanix_category_v2.linear-retention-pp-category.id]
}


# Auto Rollup Retention Protection Policy
# Create Category
resource "nutanix_category_v2" "auto-rollup-pp-category" {
  key = "category-auto-rollup-retention-protection-policy"
  value = "category_auto_rollup_retention_protection_policy"
}
# Create Auto Rollup Retention Protection Policy
resource "nutanix_protection_policy_v2" "auto-rollup-retention-protection-policy" {
  name = "auto_rollup_retention_protection_policy"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 20
      start_time                                    = "18h:10m"
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
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 30
      start_time                                    = "18h:10m"
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
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17" # Remote Domain Manager UUID
    label      = "target"
    is_primary = false
  }

  category_ids = [nutanix_category_v2.auto-rollup-pp-category.id]
}


// List Protection Policies
data "nutanix_protection_policies_v2" "protection-policies" {}

// with filter
data "nutanix_protection_policies_v2" "pps-filter" {
  filter = "name eq 'auto_rollup_retention_protection_policy'"
}

// with limit
data "nutanix_protection_policies_v2" "pp-limit" {
  limit = 4
}


// get protection policy by ext id
data "nutanix_protection_policy_v2" "pp-ex1" {
  ext_id = nutanix_protection_policy_v2.auto-rollup-retention-protection-policy.id
}
