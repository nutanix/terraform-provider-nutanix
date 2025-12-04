---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_maintenance_window"
sidebar_current: "docs-nutanix-resource-ndb-maintenance_window"
description: |-
  A maintenance window allows you to set a schedule that is used to automate repeated maintenance tasks such as OS patching and database patching. NDB allows you to create a maintenance window and then associate the maintenance window with a list of database server VMs or an instance. This operation submits a request to create, update and delete maintenance window in Nutanix database service (NDB).
---

# nutanix_ndb_maintenance_window

Provides a resource to create maintenance window based on the input parameters. 

## Example Usage

### resource to create weekly maintenance window
```hcl
    resource nutanix_ndb_maintenance_window acctest-managed {
        name = "test-maintenance"
        description = "desc"
        duration = 3
        recurrence = "WEEKLY"
        day_of_week = "TUESDAY"
        start_time = "17:04:47" 
    }
```

### resource to create monthly maintenance window
```hcl
    resource nutanix_ndb_maintenance_window acctest-managed{
        name = "test-maintenance"
        description = "description"
        duration = 2
        recurrence = "MONTHLY"
        day_of_week = "TUESDAY"
        start_time = "17:04:47" 
        week_of_month = 4
	}
```


## Argument Reference
* `name`: (Required) Name for the maintenance window.
* `description`: (Optional) Description for maintenance window
* `recurrence`: (Required) Supported values [ MONTHLY, WEEKLY ]
* `start_time`: (Required) start time for maintenance window to trigger
* `duration`: (Optional) duration in hours. Default is 2
* `day_of_week`: (Optional) Day of the week to trigger maintenance window. Supports [ MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, SUNDAY ]
* `week_of_month`: (Optional) week of the month. Supports [1, 2, 3, 4] .
* `timezone`: timezone . Default is Asia/Calcutta . 

### a Weekly or Monthly schedule.
* If you select Weekly, select the day and time when the maintenance window triggers.
* If you select Monthly, select a week (1st, 2nd, 3rd, or 4th ), day of the week, and a time when the maintenance window triggers.


## Attributes Reference
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
