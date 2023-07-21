---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_clusters"
sidebar_current: "docs-nutanix-datasource-ndb-clusters"
description: |-
 List all clusters in Nutanix Database Service
---

# nutanix_ndb_clusters

List all clusters in Nutanix Database Service

## Example Usage

```hcl
data "nutanix_ndb_clusters" "clusters" {
}

output "clusters_op" {
 value = data.nutanix_ndb_clusters.clusters
}

```

## Attribute Reference

The following arguments are exported:

* `clusters`: list of clusters

## clusters

The following attributes are exported for each cluster:

* `id`: - id of cluster
* `name`: - name of cluster
* `unique_name`: - unique name of cluster
* `ip_addresses`: - IP address
* `fqdns`: - fqdn
* `nx_cluster_uuid`: - nutanix cluster uuid
* `description`: - description
* `cloud_type`: - cloud type
* `date_created`: - creation date
* `date_modified`: - date modified
* `version`: - version
* `owner_id`: - owner UUID
* `status`: - current status
* `hypervisor_type`: - hypervisor type
* `hypervisor_version`: - hypervisor version
* `properties`: - list of properties
* `reference_count`: - NA
* `username`: - username 
* `password`: - password
* `cloud_info`: - cloud info
* `resource_config`: - resource related consumption info
* `management_server_info`: - NA
* `entity_counts`: - no. of entities related
* `healthy`: - if healthy status

See detailed information in [NDB Clusters](https://www.nutanix.dev/api_references/ndb/#/b00cac8329db1-get-a-list-of-all-clusters).
