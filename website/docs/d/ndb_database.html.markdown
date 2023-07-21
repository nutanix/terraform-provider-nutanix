---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_database"
sidebar_current: "docs-nutanix-datasource-ndb-database"
description: |-
 Describes a database instance in Nutanix Database Service
---

# nutanix_ndb_database

Describes a database instance in Nutanix Database Service

## Example Usage

```hcl
data "nutanix_ndb_database" "db1" {
 database_id = "<sample-id>"
}

output "db1_output" {
 value = data.nutanix_ndb_database.db1
}

```

## Argument Reference

The following arguments are supported:

* `database_id`: ID of database instance

## Attribute Reference

The following attributes are exported:

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


See detailed information in [Database Instance](https://www.nutanix.dev/api_references/ndb/#/7ea718d287345-get-the-database-by-value-type).
