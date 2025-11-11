terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.0"
    }
  }
}

# Defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

# Get all AOS clusters (excluding Prism Central)
data "nutanix_clusters_v2" "aos_clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  # Get the first AOS cluster ext_id for association
  cluster_ext_ids = [
    for cluster in data.nutanix_clusters_v2.aos_clusters.cluster_entities :
    cluster.ext_id
  ]
}

# Create a cluster profile
resource "nutanix_cluster_profile_v2" "example" {
  name              = "terraform-example-cluster-profile"
  description       = "Example Cluster Profile created via Terraform"
  allowed_overrides = ["NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"]

  name_server_ip_list {
    ipv4 {
      value = "8.8.8.8"
    }
  }

  ntp_server_ip_list {
    ipv4 {
      value = "10.40.64.15"
    }
    fqdn {
      value = "ntp.example.com"
    }
  }

  pulse_status {
    is_enabled          = true
    pii_scrubbing_level = "DEFAULT"
  }
}

# Associate cluster profile with a single cluster (dryrun mode - for testing)
# Uncomment to test association without actually applying it
resource "nutanix_cluster_profile_association_v2" "dryrun_example" {
  ext_id   = nutanix_cluster_profile_v2.example.id
  clusters = [local.cluster_ext_ids[0]]
  dryrun   = true
}

# Associate cluster profile with multiple clusters
resource "nutanix_cluster_profile_association_v2" "multi_cluster_example" {
  ext_id   = nutanix_cluster_profile_v2.example.id
  clusters = length(local.cluster_ext_ids) > 0 ? local.cluster_ext_ids : []
  dryrun   = false
}

# Associate cluster profile with specific clusters by name (optional - uncomment if needed)
data "nutanix_clusters_v2" "named_clusters" {
  filter = "name eq '${var.cluster_name}'"
}

resource "nutanix_cluster_profile_association_v2" "named_cluster_example" {
  ext_id = nutanix_cluster_profile_v2.example.id
  clusters = [
    for cluster in data.nutanix_clusters_v2.named_clusters.cluster_entities :
    cluster.ext_id
  ]
  dryrun = false
}

# Output examples
output "cluster_profile_id" {
  value       = nutanix_cluster_profile_v2.example.id
  description = "The ID of the created cluster profile"
}

output "associated_clusters" {
  value       = local.cluster_ext_ids
  description = "List of cluster IDs associated with the profile"
}

