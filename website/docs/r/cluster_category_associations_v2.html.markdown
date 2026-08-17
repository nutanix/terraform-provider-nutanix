---
layout: "nutanix"
page_title: "NUTANIX: nutanix_cluster_category_associations_v2"
sidebar_current: "docs-nutanix-resource-cluster-categories-v2"
description: |-
  Associate and disassociate categories with a cluster identified by `cluster_ext_id`.
---

# nutanix_cluster_category_associations_v2

Associate and disassociate categories with a cluster identified by `cluster_ext_id`.

-> **Note:** The resource manages category associations for the cluster. To change the associations, update the `categories` set.

## Example Usage

```hcl
resource "nutanix_category_v2" "category-1" {
  key         = "environment"
  value       = "production"
  description = "Production environment category"
}

resource "nutanix_category_v2" "category-2" {
  key         = "team"
  value       = "platform"
  description = "Platform team category"
}

data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

resource "nutanix_cluster_category_associations_v2" "cluster_categories" {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
  categories      = [nutanix_category_v2.category-1.id, nutanix_category_v2.category-2.id]
}
```

## Argument Reference

The following arguments are supported:

- `cluster_ext_id`: (Required) The external identifier of the cluster (UUID).
- `categories`: (Required) Set of category external identifiers (UUIDs) to associate with the cluster.

## Import

This resource supports importing an existing cluster category association by cluster external ID.

```bash
terraform import nutanix_cluster_category_associations_v2.cluster_categories <cluster_ext_id>
```

See detailed information in [Nutanix Cluster - Category Associations V4](https://developers.nutanix.com/api-reference?namespace=clustermgmt&version=v4.3#tag/Clusters/operation/associateCategoriesToCluster).