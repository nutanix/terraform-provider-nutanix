---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_log_catchups"
sidebar_current: "docs-nutanix-resource-ndb-log-catchups"
description: |-
    A log catch-up operation copies transaction logs from the source database based on a specified schedule. The schedule can be provided during database registration or provisioning or can be modified later. 
    This operation submits a request to perform log catchups of the database instance in Nutanix database service (NDB).
---

# nutanix_ndb_log_catchups

Provides a resource to perform the log cactup for database instance based on the input parameters. 

## Example Usage

```hcl
    resource "nutanix_ndb_log_catchups" "name" {
        time_machine_id = "{{ timeMachineID }}"
    }

    resource "nutanix_ndb_log_catchups" "name" {
        database_id = "{{ DatabaseID }}"
    }
```

## Argument Reference

* `time_machine_id`: (Optional) Time machine id of 
* `database_id`: (Optional)
* `for_restore`: (Optional) Logs to Backup. The database may contain additional logs. Backup any remaining logs before restore or they will be lost.
* `log_catchup_version`: (Optional) it helps to perform same operation with same config.


See detailed information in [NDB Log Catchups](https://www.nutanix.dev/api_references/ndb/#/6100cd9959e52-start-log-catchup-for-given-time-machine) .