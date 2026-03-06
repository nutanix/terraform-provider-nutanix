terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.4.0"
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

# Create a multidomain project
resource "nutanix_project_v2" "example" {
  name        = "multidomain-project-example"
  description = "Example project for multidomain namespace"
}

# Get project by ext_id
data "nutanix_project_v2" "fetch" {
  ext_id = nutanix_project_v2.example.id
}

# List all projects
data "nutanix_projects_v2" "all" {}

# Create a multidomain resource group
resource "nutanix_resource_group_v2" "example" {
  name           = "multidomain-resource-group-example"
  project_ext_id = nutanix_project_v2.example.ext_id
}

# Get resource group by ext_id
data "nutanix_resource_group_v2" "fetch_rg" {
  ext_id = nutanix_resource_group_v2.example.id
}

# List all resource groups
data "nutanix_resource_groups_v2" "all_rg" {}

output "project_ext_id" {
  value = nutanix_project_v2.example.ext_id
}

output "project_name" {
  value = data.nutanix_project_v2.fetch.name
}

output "projects_count" {
  value = length(data.nutanix_projects_v2.all.projects)
}

output "resource_group_ext_id" {
  value = nutanix_resource_group_v2.example.ext_id
}

output "resource_group_name" {
  value = data.nutanix_resource_group_v2.fetch_rg.name
}

output "resource_groups_count" {
  value = length(data.nutanix_resource_groups_v2.all_rg.resource_groups)
}
