---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_tms_capability"
sidebar_current: "docs-nutanix-datasource-ndb-tms-capability"
description: |-
  Describes a time machine in Nutanix Database Service
---

# nutanix_ndb_tms_capability

 Describes a time machine present in Nutanix Database Service

## Example Usage

```hcl

    data "nutanix_ndb_tms_capability" "tms_cap"{
        time_machine_id = {{ timeMachine_ID }}
    }

```

## Argument Reference

* `time_machine_id`: (Required) Time machine Id

## Attribute Reference

* `output_time_zone`: output time zone
* `type`: type of tms
* `nx_cluster_id`: cluster id where time machine is present
* `source`: source of time machine
* `nx_cluster_association_type`: cluster association 
* `sla_id`: SLA id
* `overall_continuous_range_end_time`: continuous range end time info
* `last_continuous_snapshot_time`: last continuous snapshot time
* `log_catchup_start_time`: log catchup start time
* `heal_with_reset_capability`: heal with reset capability
* `database_ids`: database ids
* `log_time_info`: log time info
* `capability`: capability info
* `capability_reset_time`: capability reset time
* `last_db_log`: last db log info
* `last_continuous_snapshot`: last continuous snapshot info

See detailed information in [NDB Time Machine Capability](https://www.nutanix.dev/api_references/ndb/#/011b39e32bdc5-get-capability-of-given-time-machine) .
