---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_maintenance_task"
sidebar_current: "docs-nutanix-resource-ndb-maintenance_task"
description: |-
  This operation submits a request to create, update and delete maintenance task association with database servers vms in Nutanix database service (NDB).
---

# nutanix_ndb_maintenance_task

Provides a resource to associate a maintenance window with database server VM based on the input parameters. 

## Example Usage

### resource to associated maintenance window with OS_PATCHING
```hcl
    resource "nutanix_ndb_maintenance_task" "name" {
        dbserver_id = [
            "{{ dbserver_vm_id }}"
        ]
        maintenance_window_id = "{{ maintenance_window_id }}"
        tasks{
            task_type = "OS_PATCHING"
        }
    }
```

### resource to associated maintenance window with DB_PATCHING
```hcl
    resource "nutanix_ndb_maintenance_task" "name" {
        dbserver_id = [
            "{{ dbserver_vm_id }}"
        ]
        maintenance_window_id = "{{ maintenance_window_id }}"
        tasks {
            task_type = "DB_PATCHING"
        }
    }
```

### resource to associated maintenance window with pre and post command on each task
```hcl
    resource "nutanix_ndb_maintenance_task" "name" {
        dbserver_id = [
            "{{ dbserver_vm_id }}"
        ]
        maintenance_window_id = "{{ maintenance_window_id }}"
        tasks {
            task_type = "DB_PATCHING"
            pre_command = "{{ pre_command for db patching }}"
            post_command = "{{ post_command for db patching }}"
        }
        tasks{
            task_type = "OS_PATCHING"
            pre_command = "{{ pre_command for os patching}}"
            post_command = "{{ post_command for os patching }}"
        }
    }
```

## Argument Reference

The following arguments are supported:

* `maintenance_window_id`: (Required) maintenance window id which has to be associated
* `dbserver_id`: (Optional) dbserver vm id. Conflicts with "dbserver_cluster"
* `dbserver_cluster`: (Optional) dbserver cluster ids. Conflicts with "dbserver_id"
* `tasks`: (Optional) task input for Operating System Patching or Database Patching or both

### tasks
* `task_type`: (Required) type of task. Supports [ "OS_PATCHING", "DB_PATCHING" ]
* `pre_command`: (Optional) command that you want to run before patching the OS/DB
* `post_command`: (Optional) command that you want to run after patching the OS/DB

## Attributes Reference

The following attributes are exported:

* `entity_task_association`: Entity Task Association  List.


### entity_task_association
* `id`: id of maintenance window
* `name`: name of of maintenance window
* `description`: description of maintenance window
* `owner_id`: owner id of task
* `date_created`: created date of task
* `date_modified`: modified date of task
* `access_level`: access level of tasks
* `properties`: properties of task
* `tags`: tags of task
* `maintenance_window_id`: maintenance window id
* `maintenance_window_owner_id`: maintenance window owner id
* `entity_id`: entity id
* `entity_type`: type of the entity. i.e. DBSERVER
* `status`: status of task
* `task_type`: type of the task. OS or DB 
* `payload`: list of pre post commands of OS or DB task


### payload
* `pre_post_command`: Pre Post command of Task 
* `pre_post_command.pre_command`: pre command of task
* `pre_post_command.post_command`: post command of task