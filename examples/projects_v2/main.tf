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

# pull cluster data
data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

data "nutanix_storage_containers_v2" "storage_container" {}

# Create a project (project_id defaults to name when omitted)
resource "nutanix_project_v2" "example" {
  name        = "multidomain-project-example"
  description = "Example project for multidomain namespace"
}

# Get project by ext_id
data "nutanix_project_v2" "fetch" {
  ext_id = nutanix_project_v2.example.id
}

# List all projects
data "nutanix_projects_v2" "all" {
  filter = "extId eq '${nutanix_project_v2.example.ext_id}'"
}

# Create a resource group, scoped with project ext_id
resource "nutanix_resource_group_v2" "resource_group_1" {
  name           = "multidomain-resource-group-example"
  project_ext_id = nutanix_project_v2.example.ext_id
  placement_targets {
    cluster_ext_id = local.cluster_ext_id
    storage_containers {
      ext_id = data.nutanix_storage_containers_v2.storage_container.storage_containers[3].ext_id
    }
    storage_containers {
      ext_id = data.nutanix_storage_containers_v2.storage_container.storage_containers[0].ext_id
    }
  }
}

# Get resource group by ext_id
data "nutanix_resource_group_v2" "fetch_rg" {
  ext_id = nutanix_resource_group_v2.resource_group_1.id
}

# List all resource groups, scoped with project ext_id
data "nutanix_resource_groups_v2" "all_rg" {
  filter = "projectExtId eq '${nutanix_project_v2.example.ext_id}'"
}
