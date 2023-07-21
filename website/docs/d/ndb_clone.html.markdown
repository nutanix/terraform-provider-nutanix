---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_clone"
sidebar_current: "docs-nutanix-datasource-ndb-clone"
description: |-
  Describes a clone in Nutanix Database Service
---

# nutanix_ndb_clone

Describes the clone present in Nutanix Database Service

## Example Usage

```hcl

    data "nutanix_ndb_clone" "name" {
        clone_name = "test-inst-tf-check"
    } 

    data "nutanix_ndb_clone" "name" {
        clone_name = "test-inst-tf-check"
        filters{
            detailed= true
        }
    }

```

## Argument Reference

* `clone_id`: (Optional) Clone id
* `clone_name`: (Optional) Clone Name
* `filters`: (Optional) Fetches info based on filter

### filters
* `detailed`: (Optional) Load entities with complete details. Default is false
* `any_status`: (Optional) Get entity(s) if it satisfies query criteria irrespective of status (retrieve even deleted). Default is false
* `load_dbserver_cluster`:(Optional) Load cluster info. Default is false
* `timezone`:(Optional) Default is UTC


## Attribute Reference

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

See detailed information in [NDB Clone](https://www.nutanix.dev/api_references/ndb/#/fd37879a2d8c0-get-clone-by-value-type).