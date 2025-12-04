---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_slas"
sidebar_current: "docs-nutanix-datasource-ndb-slas"
description: |-
 Lists all SLAs in Nutanix Database Service
---

# nutanix_ndb_slas

Lists all SLAs in Nutanix Database Service

## Example Usage

```hcl
data "nutanix_ndb_slas" "slas" {}

output "sla" {
 value = data.nutanix_ndb_slas.slas
}

```

## Attribute Reference

The following attributes are exported:

* `slas`: - list of slas

### slas

Each sla in list exports following attributes:

* `id`: - ID of sla
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

See detailed information in [Nutanix Database Service SLAs](https://www.nutanix.dev/api_references/ndb/#/fbcfae008ec2a-get-all-sl-as).
