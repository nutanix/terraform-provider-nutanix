terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
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



# Standard create: build a Virtual Switch from scratch with explicit
# host-NIC, bond mode, and IGMP configurations.
resource "nutanix_virtual_switch_v2" "standard_vs" {
  name        = "example-standard-vs"
  description = "A standard Virtual Switch built from scratch"
  bond_mode   = "ACTIVE_BACKUP"
  mtu         = 1500

  clusters {
    ext_id = var.cluster_ext_id
    hosts {
      ext_id        = var.host_ext_id
      host_nics     = [var.host_nic]
      active_uplink = var.host_nic
    }
  }

  igmp_spec {
    is_snooping_enabled = true
    snooping_timeout    = 300
    querier_spec {
      is_querier_enabled = false
    }
  }
}


# Create a VPC Virtual Switch Mapping to permit East-West traffic
# through the newly created standard Virtual Switch.
resource "nutanix_vpc_virtual_switch_mapping_v2" "example_mapping" {
  mappings {
    virtual_switch_uuid      = nutanix_virtual_switch_v2.standard_vs.id
    cluster_uuids            = [var.cluster_ext_id]
    is_all_traffic_permitted = true

    # It is good practice to assign organizational metadata
    metadata {
      project_reference_id = var.project_ext_id
    }
  }
}

# ---------------------------------------------------------
# Data Sources
# ---------------------------------------------------------

# Retrieve a list of all VPC Virtual Switch Mappings in the environment
data "nutanix_vpc_virtual_switch_mappings_v2" "all_mappings" {
  # Add a depends_on to ensure the resource exists before querying (useful for tests)
  depends_on = [nutanix_vpc_virtual_switch_mapping_v2.example_mapping]
}
