---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_time_machine"
sidebar_current: "docs-nutanix-datasource-ndb-time_machine"
description: |-
  Describes a time machine in Nutanix Database Service
---

# nutanix_ndb_time_machine

 Describes a time machine present in Nutanix Database Service

## Example Usage

```hcl

    data "nutanix_ndb_time_machine" "tm"{
        time_machine_id = {{ time_machine_id }} 
    }

```

## Argument Reference

* `time_machine_id`: (Optional) time machine id
* `time_machine_name`:(Optional) time machine name

## Attribute Reference

* `id`: time machine id
* `name`: time machine name
* `description`: time machine description
* `date_created`: date created
* `date_modified`: date modified
* `access_level`: access level to time machines
* `properties`: properties of time machines
* `tags`: tags
* `clustered`: clustered or not
* `clone`: clone time machine or not
* `database_id`: database id 
* `type`: type of time machine
* `category`:  category of time machine
* `status`: status of time machine
* `ea_status`: ea status of time machine
* `scope`: scope
* `sla_id`: sla id
* `schedule_id`: schedule id
* `database`: database info
* `clones`: clone info
* `source_nx_clusters`: source clusters
* `sla_update_in_progress`: sla update in progress
* `metric`: Metric info
* `sla_update_metadata`: sla update metadata
* `sla`: sla info
* `schedule`: schedule info

See detailed information in [NDB Time Machine](https://www.nutanix.dev/api_references/ndb/#/5a3e0962ef442-get-time-machine-by-value-type).