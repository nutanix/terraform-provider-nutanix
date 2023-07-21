---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_database_restore"
sidebar_current: "docs-nutanix-resource-ndb-database-restore"
description: |-
    Restoring allows you to restore a source instance registered with NDB to a snapshot or point in time supported by the source instance time machine. You can restore an instance by using a snapshot ID, the point-in-time recovery (PITR) timestamp, or the latest snapshot. 
    This operation submits a request to restore the database instance in Nutanix database service (NDB).
---

# nutanix_ndb_database_restore

Provides a resource to restore the database instance based on the input parameters. 

## Example Usage

```hcl
    // resource to  database restore with Point in Time

    resource "nutanix_ndb_database_restore" "name" {
        database_id= "{{ database_id }}"
        user_pitr_timestamp = "2022-12-28 00:54:30"
        time_zone_pitr = "Asia/Calcutta"
    }

    // resource to database restore with snapshot uuid

    resource "nutanix_ndb_database_restore" "name" {
        database_id= "{{ database_id }}"
        snapshot_id= "{{ snapshot id }}"
    }
```

## Argument Reference

* `database_id`: (Required) database id
* `snapshot_id`: (Optional) snapshot id from you want to use for restoring the instance 
* `latest_snapshot`: (Optional) latest snapshot id
* `user_pitr_timestamp`: (Optional) the time to which you want to restore your instance.
* `time_zone_pitr`: (Optional) timezone . Should be used with  `user_pitr_timestamp`
* `restore_version`: (Optional) helps to restore the database with same config. 

## Attributes Reference

* `name`: Name of database instance
* `description`: description of database instance
* `databasetype`: type of database
* `properties`: properties of database created
* `date_created`: date created for db instance
* `date_modified`: date modified for instance
* `tags`: allows you to assign metadata to entities (clones, time machines, databases, and database servers) by using tags.
* `clone`: whether instance is cloned or not
* `database_name`: name of database
* `type`: type of database
* `database_cluster_type`: database cluster type
* `status`: status of instance
* `dbserver_logical_cluster_id`: dbserver logical cluster id
* `time_machine_id`: time machine id of instance 
* `time_zone`: timezone on which instance is created xw
* `info`: info of instance
* `metric`: Stores storage info regarding size, allocatedSize, usedSize and unit of calculation that seems to have been fetched from PRISM.
* `parent_database_id`: parent database id
* `lcm_config`: LCM config of instance
* `time_machine`: Time Machine details of instance
* `dbserver_logical_cluster`: dbserver logical cluster
* `database_nodes`: database nodes associated with database instance 
* `linked_databases`: linked databases within database instance


See detailed information in [NDB Database Restore](https://www.nutanix.dev/api_references/ndb/#/215deb77fba60-restore-database).