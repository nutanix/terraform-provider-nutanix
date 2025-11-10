terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.0"
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

#pull all clusters data
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

#create local variable pointing to desired cluster
locals {
  clusters_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

# Create categories
resource "nutanix_category_v2" "cat-1" {
  key         = "environment"
  value       = "production"
  description = "Production environment category"
}

resource "nutanix_category_v2" "cat-2" {
  key         = "department"
  value       = "engineering"
  description = "Engineering department category"
}

# Associate categories with cluster
resource "nutanix_cluster_categories_v2" "cluster-categories" {
  cluster_ext_id = local.clusters_ext_id
  categories = [
    nutanix_category_v2.cat-1.id,
    nutanix_category_v2.cat-2.id
  ]
}
