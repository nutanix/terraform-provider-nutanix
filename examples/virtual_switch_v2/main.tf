terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
    }
  }
}

# Configure the Nutanix Provider
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Fetch cluster and host information for virtual switch creation
data "nutanix_clusters_v2" "test" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}
data "nutanix_hosts_v2" "test" {}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.test.cluster_entities[0].ext_id
  host_ext_id    = data.nutanix_hosts_v2.test.host_entities[0].ext_id
}


resource "nutanix_project_v2" "test" {
  name        = "tf-test-project"
  project_id  = "tf-test-project"
  description = "Test project"
}

data "nutanix_storage_containers_v2" "test" {
  filter = "name eq 'SelfServiceContainer'"
}

resource "nutanix_resource_group_v2" "test" {
  name           = "tf-test-resource-group"
  project_ext_id = nutanix_project_v2.test.ext_id
  placement_targets {
    cluster_ext_id = local.cluster_ext_id
    storage_containers {
      ext_id = data.nutanix_storage_containers_v2.test.storage_containers[0].ext_id
    }
  }
}

# Standard create: build a Virtual Switch from scratch with explicit
# host-NIC, bond mode, and IGMP configurations.
#
# Sharing is now handled natively by this resource via shared_with_projects.
# Listed projects are shared on create/update and unshared on removal/delete.
resource "nutanix_virtual_switch_v2" "standard_vs" {
  name        = "example-standard-vs"
  description = "A standard Virtual Switch built from scratch"
  bond_mode   = "NONE"
  mtu         = 1500

  clusters {
    ext_id = local.cluster_ext_id
    hosts {
      ext_id        = local.host_ext_id
    }
  }

  igmp_spec {
    is_snooping_enabled = true
    snooping_timeout    = 300
  }

  # Share this Virtual Switch with one or more projects.
  shared_with_projects = [nutanix_project_v2.test.ext_id]
}


# Migrate-from-bridge: convert a pre-existing OVS bridge on the host into
# a Virtual Switch. Setting clusters[].existing_bridge_name triggers the migration.
# Note: Ensure the bridge was created out-of-band beforehand.
resource "nutanix_virtual_switch_v2" "migrated_vs" {
  name        = "example-migrated-vs"
  description = "Virtual Switch built from an existing OVS bridge"

  clusters {
    ext_id               = var.cluster_ext_id
    existing_bridge_name = var.existing_bridge_name
    hosts {
      ext_id = var.host_ext_id
    }
  }
}

# ---------------------------------------------------------
# Data Sources
# ---------------------------------------------------------

# Read detailed information about the newly created standard Virtual Switch
data "nutanix_virtual_switch_v2" "get_standard_vs" {
  ext_id = nutanix_virtual_switch_v2.standard_vs.id
}

# Retrieve a list of all Virtual Switches in the environment
data "nutanix_virtual_switches_v2" "list_all_vs" {}

# Retrieve node schedulable statuses
data "nutanix_node_schedulable_statuses_v2" "statuses" {}
