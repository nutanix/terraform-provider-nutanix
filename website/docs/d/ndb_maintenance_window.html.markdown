---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_maintenance_window"
sidebar_current: "docs-nutanix-datasource-ndb-maintenance-window"
description: |-
 Describes a maintenance window in Nutanix Database Service
---

# nutanix_ndb_maintenance_window

Describes a maintenance window in Nutanix Database Service

## Example Usage

```hcl
    data "nutanix_ndb_maintenance_window" "window"{
        id = "{{ maintenance_window_id }}"
    } 
```

## Argument Reference

The following arguments are supported:

* `id`: (Required) Maintenance window id.

## Attribute Reference

The following attributes are exported:
* `name`: name of maintenance window
* `description`: description of maintenance window
* `schedule`: schedule of maintenance window
* `owner_id`: owner id of maintenance window
* `date_created`: created date of maintenance window
* `date_modified`: modified date of maintenance window
* `access_level`: access level
* `properties`: properties of maintenance window
* `tags`: tags of maintenance window 
* `status`: status of maintennace window
* `next_run_time`: next run time for maintenance window to trigger 
* `entity_task_assoc`: entity task association for maintenance window
* `timezone`: timezone