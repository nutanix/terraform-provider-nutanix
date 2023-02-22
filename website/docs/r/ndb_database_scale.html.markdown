---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_database_scale"
sidebar_current: "docs-nutanix-resource-ndb-database-scale"
description: |-
  Scaling the database extends the storage size proportionally across the attached virtual disks or volume groups. Scaling is supported for both single and HA instances.
  This operation submits a request to scale out the database instance in Nutanix database service (NDB).
---

# nutanix_ndb_database_scale

Provides a resource to scale the database instance based on the input parameters. 

## Example Usage

```hcl

    // resource to scale the database

    resource "nutanix_ndb_database_scale" "scale" {
        application_type = "{{ Application Type }}"
        database_uuid = "{{ database_id }}"
        data_storage_size = 1
    }
```

## Argument Reference

* `database_uuid`: (Required) Database id
* `application_type`: (Required) type of instance. eg: postgres_database
* `data_storage_size`: (Required) data area (in GiB) to be added to the existing database.
* `pre_script_cmd`: (Optional) pre script command
* `post_script_cmd`: (Optional) post script command
* `scale_count`: (Optional) scale count helps to scale the same instance with same config


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
