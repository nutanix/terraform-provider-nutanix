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
* `owner_id`: - owner ID
* `properties`: - properties
* `tags`: - tags attached
* `clustered`: - if clustered or not
* `clone`: - if cloned
* `era_created`: - if era created
* `internal`: - if internal database
* `placeholder`: - NA
* `database_name`: - database instance name
* `type`: - database engine type
* `status`: - status of database instance
* `database_status`: - NA
* `dbserver_logical_cluster_id`: - NA
* `time_machine_id`: - time machine ID
* `parent_time_machine_id`: - parent time machine ID
* `time_zone`: - timezone
* `info`: - info regarding disks, vm, storage, etc.
* `group_info`: - group info
* `metadata`: - metadata of database instance
* `metric`: - metrics
* `category`: - category of instance
* `parent_database_id`: - parent database ID
* `parent_source_database_id`: - parent source database ID
* `lcm_config`: - lcm configuration
* `time_machine`: - time machine related config info
* `database_nodes`: - nodes info
* `dbserver_logical_cluster`: - NA
* `linked_databases`: - list of databases created in instance with info
* `databases`: - NA
* `database_group_state_info`: - NA


See detailed information in [Database Instance](https://www.nutanix.dev/api_references/era/#/b3A6MjIyMjI1NDA-get-a-database-using-id).
