---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_database_snapshot"
sidebar_current: "docs-nutanix-resource-ndb-database-snapshot"
description: |-
    NDB time machine allows you to capture and replicate snapshots of the source database across multiple clusters (as defined in the DAM policy) at the time and frequency specified in the schedule. 
    This operation submits a request to perform snapshot of the database instance in Nutanix database service (NDB).
---

# nutanix_ndb_database_snapshot

Provides a resource to perform the snapshot for database instance based on the input parameters. 

## Example Usage

```hcl
    // resource to create snapshot with time machine id

    resource "nutanix_ndb_database_snapshot" "name" {
        time_machine_id = "{{ tms_ID }}"
        name = "test-snap"
        remove_schedule_in_days = 1
    }

    // resource to craete snapshot with time machine name

    resource "nutanix_ndb_database_snapshot" "name" {
        time_machine_name = "{{ tms_name }}"
        name = "test-snap"
        remove_schedule_in_days = 1
    }

```

## Argument Reference

* `time_machine_id`: (Optional) Time Machine Id
* `time_machine_name`:(Optional) Time Machine Name
* `name`: (Optional) Snapshot name. Default value is era_manual_snapshot. 
* `remove_schedule_in_days`: (Optional) Removal schedule after which the snapshot should be removed.
* `expiry_date_timezone`: (Optional) Default is set to Asia/Calcutta
* `replicate_to_clusters`: (Optional) snapshots to be replicated to clusters. 


## Attributes Reference

* `id`: name of snapshot
* `description`: description of snapshot
* `properties`: properties 
* `date_created`: created date
* `date_modified`: modified date
* `properties`: properties 
* `tags`: tags
* `snapshot_uuid`: snapshot uuid 
* `nx_cluster_id`: nx cluster id
* `protection_domain_id`: protection domain
* `parent_snapshot_id`: parent snapshot id
* `database_node_id`: database node id
* `app_info_version`: App info version
* `status`: status
* `type`: type
* `applicable_types`: Applicable types
* `snapshot_timestamp`: snapshot timeStamp
* `software_snapshot_id`: software snapshot id
* `software_database_snapshot`: software database snapshot
* `dbserver_storage_metadata_version`: dbserver storage metadata version
* `santised_from_snapshot_id`: sanitized  snapshot id
* `timezone`: timezone
* `processed`: processed
* `database_snapshot`: database snapshot
* `from_timestamp`: from timestamp
* `to_timestamp`: to timestamp
* `dbserver_id`: dbserver id
* `dbserver_name`: dbserver name
* `dbserver_ip`:dbserver ip
* `replicated_snapshots`: replicated snapshots
* `software_snapshot`: software snapshot
* `santised_snapshots`:santised snapshots
* `snapshot_family`: snapshot family
* `snapshot_timestamp_date`: snapshot timestamp date
* `lcm_config`: LCM config
* `parent_snapshot`: parent snapshot
* `snapshot_size`: snapshot size


See detailed information in [NDB Database Snapshot](https://www.nutanix.dev/api_references/ndb/#/ca9e7ed109f08-take-snapshot).