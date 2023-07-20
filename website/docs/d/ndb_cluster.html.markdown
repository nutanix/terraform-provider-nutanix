---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_cluster"
sidebar_current: "docs-nutanix-datasource-ndb-cluster"
description: |-
 Describes a cluster in Nutanix Database Service
---

# nutanix_ndb_cluster

Describes a cluster in Nutanix Database Service

## Example Usage

```hcl
data "nutanix_ndb_cluster" "c1" {
  cluster_name = "<era-cluster-name>"
}

output "cluster" {
 value = data.nutanix_ndb_cluster.c1
}

```

## Argument Reference

The following arguments are supported:

* `cluster_id`: ID of cluster
* `cluster_name`: name of cluster

* `cluster_name` and `cluster_id` are mutually exclusive.

## Attribute Reference

The following attributes are exported:

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


See detailed information in [NDB Cluster](https://www.nutanix.dev/api_references/ndb/#/b4f28e2f7b6a9-get-a-cluster-by-id).
