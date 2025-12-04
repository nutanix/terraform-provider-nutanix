---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_databases"
sidebar_current: "docs-nutanix-datasource-ndb-databases"
description: |-
 List all database instances in Nutanix Database Service
---

# nutanix_ndb_database

List all database instances in Nutanix Database Service

## Example Usage

```hcl
data "nutanix_ndb_databases" "dbs" {}

output "dbs_output" {
 value = data.nutanix_ndb_databases.dbs
}

```

## Attribute Reference

The following attributes are exported:

* `database_instances`: - list of database instances

## database_instances

The following attributes are exported for each database_instances:

* `id`: - id of database instance
* `name`: - name of database instance
* `description`: - description
* `date_created`: - creation date
* `date_modified`: - date modified 
* `properties`: - properties
* `tags`: - tags attached
* `clustered`: - if clustered or not
* `clone`: - if cloned
* `database_name`: - database instance name
* `type`: - database engine type
* `status`: - status of database instance
* `dbserver_logical_cluster_id`: - NA
* `time_machine_id`: - time machine ID
* `time_zone`: - timezone
* `info`: - info regarding disks, vm, storage, etc.
* `metric`: - metrics
* `parent_database_id`: - parent database ID
* `lcm_config`: - lcm configuration
* `time_machine`: - time machine related config info
* `database_nodes`: - nodes info
* `dbserver_logical_cluster`: - NA
* `linked_databases`: - list of databases created in instance with info


See detailed information in [List Database Instances](https://www.nutanix.dev/api_references/ndb/#/1e508756bcdcc-get-all-the-databases).
