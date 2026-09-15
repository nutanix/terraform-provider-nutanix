terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

#####################################################################
# NIC Profiles
#
# A NIC Profile is a Prism Central construct used to enable advanced
# host NIC capabilities such as SR-IOV, datapath (network) offload and
# PCIe pass-through on supported NIC families (e.g. NVIDIA/Mellanox
# ConnectX / BlueField). A profile selects NICs by `nic_family`
# (vendor_id:device_id) and associates one or more Host NICs to it.
#####################################################################

# Example 1: SR-IOV NIC Profile in Ethernet mode.
# Associates two Host NICs from the given NIC family.
resource "nutanix_nic_profile_v2" "sriov" {
  name             = "tf-nic-profile-sriov"
  description      = "SR-IOV NIC profile managed by Terraform."
  nic_family       = var.nic_family
  host_nic_ext_ids = var.host_nic_ext_ids

  capability_config {
    capability_type = "SRIOV"
  }

  operating_mode = "ETHERNET"
}

# Example 2: Datapath (Network) offload NIC Profile in Ethernet mode.
# Useful for BlueField / DPU based offload. No Host NICs are associated
# at create time; they can be added later by populating host_nic_ext_ids.
resource "nutanix_nic_profile_v2" "dp_offload" {
  name        = "tf-nic-profile-dp-offload"
  description = "Datapath offload NIC profile managed by Terraform."
  nic_family  = var.nic_family

  capability_config {
    capability_type = "DP_OFFLOAD"
  }

  operating_mode = "ETHERNET"
}

# Example 3: SR-IOV NIC Profile in InfiniBand mode (East-West traffic).
resource "nutanix_nic_profile_v2" "infiniband" {
  name        = "tf-nic-profile-infiniband"
  description = "InfiniBand SR-IOV NIC profile managed by Terraform."
  nic_family  = var.nic_family

  capability_config {
    capability_type = "SRIOV"
  }

  operating_mode = "INFINIBAND"
}

#####################################################################
# Data Sources
#####################################################################

# Fetch a single NIC Profile by its ext_id (UUID).
data "nutanix_nic_profile_v2" "by_id" {
  ext_id = nutanix_nic_profile_v2.sriov.id
}

# List all NIC Profiles in Prism Central.
data "nutanix_nic_profiles_v2" "all" {}

#####################################################################
# Outputs
#####################################################################

output "sriov_nic_profile_id" {
  description = "ext_id (UUID) of the created SR-IOV NIC profile."
  value       = nutanix_nic_profile_v2.sriov.id
}

output "sriov_nic_profile_host_nic_references" {
  description = "Host NICs currently associated with the SR-IOV NIC profile."
  value       = data.nutanix_nic_profile_v2.by_id.host_nic_references
}

output "all_nic_profiles" {
  description = "All NIC profiles available in Prism Central."
  value       = data.nutanix_nic_profiles_v2.all.nic_profiles
}
