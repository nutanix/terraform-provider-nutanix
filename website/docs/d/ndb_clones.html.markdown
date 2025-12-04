---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_clones"
sidebar_current: "docs-nutanix-datasource-ndb-clones"
description: |-
  List all the clone in Nutanix Database Service
---

# nutanix_ndb_clones

List all the clone present in Nutanix Database Service

## Example Usage

```hcl

    data "nutanix_ndb_clones" "clones"{ }

    data "nutanix_ndb_clones" "clones"{
        filters{
            detailed= true
        }
    }

```

## Argument Reference

* `filters`: (Optional) Fetches the clone info based on given params

### filters

* `detailed`: (Optional) Load entities with complete details. Default is false
* `any_status`: (Optional) Get entity(s) if it satisfies query criteria irrespective of status (retrieve even deleted). Default is false
* `load_dbserver_cluster`: (Optional) Load cluster info. Default is false
* `timezone`: (Optional) Default is UTC
* `order_by_dbserver_cluster`: (Optional) Sorted by dbserver cluster. Default is false
* `order_by_dbserver_logical_cluster`: (Optional) Sorted by dbserver logical cluster.  Default is false


## Attribute Reference

* `clones`: List of clones based on filters

### clones

* `id`: cloned id 
* `name`: cloned name
* `description`: cloned description
* `date_created`: date created for clone
* `date_modified`: last modified date for clone
* `tags`: allows you to assign metadata to entities (clones, time machines, databases, and database servers) by using tags.
* `properties`: properties of clone
* `clustered`: clustered or not
* `clone`: clone or not
* `database_name`: database name
* `type`: type 
* `database_cluster_type`: database cluster type
* `status`: status of clone
* `database_status`: database status 
* `dbserver_logical_cluster_id`: dbserver logical cluster id
* `time_machine_id`: time machine id
* `parent_time_machine_id`: parent time machine id
* `time_zone`: time zone
* `info`: cloned info 
* `metric`: Metric of clone
* `parent_database_id`: parent database id
* `parent_source_database_id`: parent source database id
* `lcm_config`: LCM Config
* `time_machine`: Time machine info
* `dbserver_logical_cluster`: dbserver logical cluster 
* `database_nodes`: database nodes associated with database instance 
* `linked_databases`: linked databases within database instance
* `databases`: database for a cloned instance


See detailed information in [NDB Clones](https://www.nutanix.dev/api_references/ndb/#/02b17b417ac8a-get-a-list-of-all-clones).