---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_snapshot"
sidebar_current: "docs-nutanix-datasource-ndb-snapshot"
description: |-
  Describes a snaphot in Nutanix Database Service
---

# nutanix_ndb_snapshot

Describes the snapshot present in Nutanix Database Service

## Example Usage

```hcl

    data "nutanix_ndb_snapshot" "snaps"{
        snapshot_id = "{{ snapshot_id }}"
        filters {
            load_replicated_child_snapshots = true
        }
    }
```

## Argument Reference

* `snapshot_id`: (Required) Snapshot ID to be given
* `filters`: (Optional) Filters will fetch the snapshot details as per input 

### filters
* `timezone`: (Optional) Default is UTC
* `load_replicated_child_snapshots`: (Optional) load child snapshots. Default is false

## Attribute Reference

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

See detailed information in [NDB Snapshot](https://www.nutanix.dev/api_references/ndb/#/d50fb18097051-get-snapshot-by-value-type).