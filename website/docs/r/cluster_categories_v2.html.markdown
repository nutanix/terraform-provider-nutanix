---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_categories_v2"
sidebar_current: "docs-nutanix-resource-nutanix-cluster-categories-v2"
description: |-
  Associate/Disassociate categories to the cluster identified by {clusterExtId}. Creating this resource will associate categories to the cluster, and destroying it will disassociate them.
---

# nutanix_cluster_categories_v2

Associate/Disassociate categories to the cluster identified by {clusterExtId}.

Creating this resource will associate the specified categories to the cluster. Destroying this resource will disassociate the categories from the cluster.

## Example Usage

```hcl
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
  cluster_ext_id = "9aa06bb6-bec6-45a6-9f13-f8685f61efc6"
  categories = [
    nutanix_category_v2.cat-1.id,
    nutanix_category_v2.cat-2.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ext_id`: -(Required) Cluster UUID.
* `categories`: -(Required) Set of category IDs to associate with the cluster.

## Attributes Reference

No attributes are exported.


See detailed information in [Nutanix Cluster - Associate Categories to Cluster V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.1#tag/Clusters/operation/associateCategoriesToCluster).
