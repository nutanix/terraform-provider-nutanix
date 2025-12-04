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


data "nutanix_clusters" "clusters" {}

locals {
  cluster_uuid = (data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
  ? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid)
}

resource "nutanix_protection_rule" "protection_rule_example" {
  name        = "protection-rule-example"
  description = "This is a sample protection rule for demonstration purposes."
  ordered_availability_zone_list {
    availability_zone_url = var.source_az_url
    cluster_uuid          = local.cluster_uuid
  }
  ordered_availability_zone_list {
    availability_zone_url = var.target_az_url
    cluster_uuid          = var.target_cluster_uuid
  }

  availability_zone_connectivity_list {
    source_availability_zone_index      = 0
    destination_availability_zone_index = 1
    snapshot_schedule_list {
      recovery_point_objective_secs = 3600
      snapshot_type                 = "CRASH_CONSISTENT"
      local_snapshot_retention_policy {
        num_snapshots = 4
      }
    }
  }
  availability_zone_connectivity_list {
    source_availability_zone_index      = 1
    destination_availability_zone_index = 0
    snapshot_schedule_list {
      recovery_point_objective_secs = 3600
      snapshot_type                 = "CRASH_CONSISTENT"
      local_snapshot_retention_policy {
        num_snapshots = 4
      }
    }
  }
  category_filter {
    params {
      name   = "Environment"
      values = ["Dev"]
    }
  }
}
