---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_sla"
sidebar_current: "docs-nutanix-datasource-ndb-sla"
description: |-
 Describes a SLA in Nutanix Database Service
---

# nutanix_ndb_sla

Describes a SLA in Nutanix Database Service

## Example Usage

```hcl
data "nutanix_ndb_sla" "sla1" {
 sla_name = "test-sla"
}

output "sla" {
 value = data.nutanix_ndb_sla.sla1
}

```

## Argument Reference

The following arguments are supported:

* `sla_id`: SLA ID for query
* `sla_name`: SLA Name for query

* `sla_id` and `sla_name` are mutually exclusive.

## Attribute Reference

The following attributes are exported:

* `id`: - id of sla
* `name`: - sla name
* `unique_name`: - unique name
* `description`: - description of sla
* `owner_id`: - owner ID
* `system_sla`: - if system sla
* `date_created`: - creation date
* `date_modified`: - last modified
* `continuous_retention`: - continuous retention of logs limit
* `daily_retention`: - Daily snapshots retention limit
* `weekly_retention`: - weeky snapshots retention limit
* `monthly_retention`: - Monthly snapshots retention limit
* `quartely_retention`: - Daily snapshots retention limit
* `yearly_retention`: - Yearly snapshots retention limit
* `reference_count`: - Reference count
* `pitr_enabled`: - If point in time recovery enabled
* `current_active_frequency`: - Current active frequency



See detailed information in [Nutanix Database Service SLA](https://www.nutanix.dev/api_references/ndb/#/eabbbab3f7eff-get-sla-by-name).