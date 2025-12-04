terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

data "nutanix_clusters_v2" "clusters" {}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}


# List all the hosts
data "nutanix_hosts_v2" "hosts" {
  cluster_ext_id = local.clusterExtId
}

# List all the hosts with filter
data "nutanix_hosts_v2" "hosts" {
  cluster_ext_id = local.clusterExtId
  filter         = "name eq 'host-1'"
}

# Get the host details
data "nutanix_host_v2" "host" {
  cluster_ext_id = local.clusterExtId
  ext_id         = data.nutanix_hosts_v2.hosts.host_entities[0].ext_id
}
