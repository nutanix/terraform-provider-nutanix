---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_snapshots"
sidebar_current: "docs-nutanix-datasource-ndb-snapshots"
description: |-
  List all snaphots in Nutanix Database Service
---

# nutanix_ndb_snapshots

List all snapshots present in Nutanix Database Service

## Example Usage

```hcl

    data "nutanix_ndb_snapshots" "snaps"{ }

    data "nutanix_ndb_snapshots" "snaps"{ 
        filters{
            time_machine_id = "{{ time_machine_id }}"
        }
    }
```

## Argument Reference

* `filters`: (Optional) filters help to fetch the snapshots based on input

### filters
* `time_machine_id`: (Optional) Fetches all the snapshots for a given time machine

## Attribute Reference 

* `snapshots`: List of snapshots

### snapshots

* `id`: name of snapshot
* `description`: description of snapshot
* `properties`: properties 
* `owner_id`: owner id 
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
* `metadata`: metadata of snapshot 
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


See detailed information in [NDB Snapshots](https://www.nutanix.dev/api_references/ndb/#/d0b89ff892448-get-list-of-all-snapshots).