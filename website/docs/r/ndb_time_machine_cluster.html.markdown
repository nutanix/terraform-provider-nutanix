---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_tms_cluster"
sidebar_current: "docs-nutanix-resource-ndb-tms-cluster"
description: |-
  NDB multi-cluster allows you to manage time machine data availability across all the registered Nutanix clusters in NDB. This operation submits a request to add, update and delete clusters in time machine data availability for Nutanix database service (NDB).
---

# nutanix_ndb_tms_cluster

Provides a resource to manage time machine data availability across all the registered Nutanix clusters in NDB.

## Example Usage

```hcl
    resource "nutanix_ndb_tms_cluster" "cls" {
        time_machine_id = "{{ tms_id }}"
        nx_cluster_id = "{{ cluster_id }}"
        sla_id = "{{ sla_id }}"
    }
```

## Argument Reference

The following arguments are supported:

* `time_machine_id`: (Required) time machine id 
* `nx_cluster_id`: (Required) Nutanix cluster id on the associated registered clusters.
* `sla_id`: (Required) SLA id for the associated cluster.

* `type`:  (Optional) Default value is "OTHER"

## Attributes Reference

The following attributes are exported:

* `status`: status of the cluster associated with time machine
* `schedule_id`: schedule id of the data associated with time machine
* `owner_id`: owner id 
* `source_clusters`: source clusters in time machines
* `log_drive_status`: log drive status of time machine
* `date_created`: created date of time machine associated with cluster
* `date_modified`: modified date of time machine associated with cluster
* `log_drive_id`: log drive id
* `description`: description of nutanix cluster associated with time machine
* `source`: source is present or not