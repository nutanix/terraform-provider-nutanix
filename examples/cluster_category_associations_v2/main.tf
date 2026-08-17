#############################################################################
# Example main.tf for Nutanix + Terraform
#
# This example demonstrates how to associate Prism Central categories
# to a cluster using `nutanix_cluster_category_associations_v2`.
#############################################################################

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.0"
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

# Pick a target cluster. In real usage, filter to your desired cluster(s).
data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

# Create two categories (key/value).
resource "nutanix_category_v2" "category-1" {
  description = "Example category 1"
  key          = "example-category-key-1"
  value        = "example-category-value-1"
}

resource "nutanix_category_v2" "category-2" {
  description = "Example category 2"
  key          = "example-category-key-2"
  value        = "example-category-value-2"
}

# Associate the categories to the cluster.
resource "nutanix_cluster_category_associations_v2" "cluster_categories" {
  cluster_ext_id = local.cluster_ext_id
  categories      = [nutanix_category_v2.category-1.id, nutanix_category_v2.category-2.id]
}

/*
Update example:

To remove `category-2` from the association, change `categories` to:

categories = [nutanix_category_v2.category-1.id]

Import example:

terraform import nutanix_cluster_category_associations_v2.cluster_categories <cluster_ext_id>
*/

